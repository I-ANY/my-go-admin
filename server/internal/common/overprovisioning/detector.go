package overprovisioning

import (
	"biz-auto-api/internal/models"
	"fmt"
	"math"
	"strconv"
	"strings"
)

// OverProvisioningDetector 超配检测器
type OverProvisioningDetector struct{}

// NewOverProvisioningDetector 创建超配检测器
func NewOverProvisioningDetector() *OverProvisioningDetector {
	return &OverProvisioningDetector{}
}

// DetectOverProvisioning 执行超配检测
func (d *OverProvisioningDetector) DetectOverProvisioning(
	rule *models.BusinessOverProvisioningRule,
	hardware *models.Hardware,
	promInfo *PrometheusHostInfo,
) (*DetectionResult, error) {

	result := &DetectionResult{
		IsOverProvisioned:    false,
		OverProvisionedItems: []string{},
		StandardValues:       make(map[string]float64),
		CurrentValues:        make(map[string]float64),
	}

	// 解析当前硬件配置
	currentConfig, err := d.parseHardwareConfig(hardware)
	if err != nil {
		return nil, fmt.Errorf("解析硬件配置失败: %v", err)
	}

	// 获取带宽信息
	bandwidthInfo, err := d.parseBandwidthInfo(promInfo)
	if err != nil {
		return nil, fmt.Errorf("解析带宽信息失败: %v", err)
	}

	// 合并当前配置
	for k, v := range bandwidthInfo {
		currentConfig[k] = v
	}

	result.CurrentValues = currentConfig

	// 应用规则配置
	standardConfig, err := d.applyRuleConfig(rule, currentConfig)
	if err != nil {
		return nil, fmt.Errorf("应用规则配置失败: %v", err)
	}

	result.StandardValues = standardConfig

	// 执行检测（传入规则以使用配置开关）
	overProvisionedItems, details := d.performDetection(rule, currentConfig, standardConfig)

	result.IsOverProvisioned = len(overProvisionedItems) > 0
	result.OverProvisionedItems = overProvisionedItems
	result.Details = details

	return result, nil
}

// parseHardwareConfig 解析硬件配置
func (d *OverProvisioningDetector) parseHardwareConfig(hardware *models.Hardware) (map[string]float64, error) {
	config := make(map[string]float64)

	if hardware == nil {
		// 硬件信息为空时直接返回空配置，后续检测会自动跳过缺失项
		return config, nil
	}

	// 解析内存
	if hardware.Memory != "" {
		memStr := strings.ReplaceAll(strings.ReplaceAll(hardware.Memory, "G", ""), "M", "")
		if mem, err := strconv.ParseFloat(memStr, 64); err == nil {
			if strings.Contains(hardware.Memory, "M") {
				mem = mem / 1024 // 转换为GB
			}
			config[Memory] = mem
		}
	}

	// 解析系统盘
	if hardware.DiskDetail.SysDisk.Volume != "" {
		sysDisk := d.convertToTB(hardware.DiskDetail.SysDisk.Volume)
		config[SysDisk] = sysDisk * 1000
	}

	// 解析数据盘
	var ssdSize, hddSize, totalDataDisk float64
	for _, disk := range hardware.DiskDetail.DataDisk.DataDetail {
		diskSize := d.convertToTB(disk.Volume)
		totalDataDisk += diskSize

		if disk.DiskType == "hdd" {
			hddSize += diskSize
		} else {
			ssdSize += diskSize
		}
	}

	config[SSD] = ssdSize
	config[HDD] = hddSize
	config[TotalDataDisk] = totalDataDisk

	return config, nil
}

// parseBandwidthInfo 解析带宽信息
func (d *OverProvisioningDetector) parseBandwidthInfo(promInfo *PrometheusHostInfo) (map[string]float64, error) {
	bandwidth := make(map[string]float64)

	if promInfo == nil {
		return bandwidth, nil
	}

	// 解析规划带宽
	if promInfo.BwPlan != "" {
		if bw, err := strconv.ParseFloat(promInfo.BwPlan, 64); err == nil {
			bandwidth[PlanBw] = bw / 1000 / 1000 / 1000 // 转换为Gbps
		}
	}

	// 解析提交带宽
	if promInfo.BwCount != "" && promInfo.BwSingle != "" {
		if count, err1 := strconv.ParseFloat(promInfo.BwCount, 64); err1 == nil {
			if single, err2 := strconv.ParseFloat(promInfo.BwSingle, 64); err2 == nil {
				submitBw := (count * single) / 1000 / 1000 / 1000 // 转换为Gbps
				bandwidth[SubmitBw] = submitBw
			}
		}
	}

	// 解析95值
	if promInfo.Day95 != "" {
		if day95, err := strconv.ParseFloat(promInfo.Day95, 64); err == nil {
			bandwidth[Day95] = day95 / 1000 / 1000 / 1000 // 转换为Gbps
		}
	}

	if promInfo.Evening95 != "" {
		if evening95, err := strconv.ParseFloat(promInfo.Evening95, 64); err == nil {
			bandwidth[Evening95] = evening95 / 1000 / 1000 / 1000 // 转换为Gbps
		}
	}

	return bandwidth, nil
}

// applyRuleConfig 应用规则配置，生成标准值配置
// 职责：根据规则配置生成各项检测的标准值，不执行实际检测
func (d *OverProvisioningDetector) applyRuleConfig(
	rule *models.BusinessOverProvisioningRule,
	currentConfig map[string]float64,
) (map[string]float64, error) {
	standardConfig := make(map[string]float64)

	// ==================== 第一步：设置基础标配（来自规则字段） ====================
	standardConfig[Memory] = float64(rule.Mem)
	standardConfig[SSD] = float64(rule.SsdSize)
	standardConfig[HDD] = float64(rule.HddSize)
	// 业务盘标配 = SSD + HDD
	standardConfig[TotalDataDisk] = float64(rule.SsdSize) + float64(rule.HddSize)

	// ==================== 第二步：磁盘类型前置校验（PlanConf中的磁盘类型条件） ====================
	// 用于判断规则是否适用于当前设备，不适用则直接返回错误
	if err := d.validateDiskTypeRequirement(rule.PlanConf, currentConfig); err != nil {
		return nil, err
	}

	// ==================== 第三步：设置存储带宽比标准值（独立检测项） ====================
	// 存储带宽比是独立检测项，有自己的标准值，不影响SSD/HDD/业务盘的标准值
	if rule.EnableStorageBwRatioCheck && strings.TrimSpace(rule.StorageBwRatio) != "" {
		if ratio := d.parseStorageBwRatio(rule.StorageBwRatio); ratio > 0 {
			standardConfig[StorageBwRatio] = ratio
		}
	}

	// ==================== 第四步：PlanConf 中的大小类规则最终覆盖 ====================
	// 这些规则具有最高优先级，会直接覆盖前面计算的标准值
	d.applyPlanConfSizeRules(rule.PlanConf, standardConfig)

	return standardConfig, nil
}

// validateDiskTypeRequirement 验证磁盘类型要求（前置校验）
// 如果规则要求特定磁盘类型，但设备不满足，则规则不适用
func (d *OverProvisioningDetector) validateDiskTypeRequirement(
	planConfs []models.PlanConf,
	currentConfig map[string]float64,
) error {
	for _, planConf := range planConfs {
		if planConf.CondType != models.CronTypeDiskType {
			continue
		}

		diskType := strings.ToUpper(strings.TrimSpace(planConf.Cond.Value))
		switch diskType {
		case "SSD":
			// 规则要求有 SSD，但当前配置中 SSD 容量 <= 0，则规则不适用
			if currentConfig[SSD] <= 0 {
				return fmt.Errorf("磁盘类型不匹配: 规则要求SSD, 实际SSD容量<=0")
			}
		case "HDD":
			// 规则要求有 HDD，但当前配置中 HDD 容量 <= 0，则规则不适用
			if currentConfig[HDD] <= 0 {
				return fmt.Errorf("磁盘类型不匹配: 规则要求HDD, 实际HDD容量<=0")
			}
		default:
			// 其他取值暂不强制校验
		}
	}
	return nil
}

// applyPlanConfSizeRules 应用 PlanConf 中的大小类规则（最终覆盖）
// 这些规则会直接覆盖基础标配值
func (d *OverProvisioningDetector) applyPlanConfSizeRules(
	planConfs []models.PlanConf,
	standardConfig map[string]float64,
) {
	for _, planConf := range planConfs {
		switch planConf.CondType {
		case models.CronTypeSystemDiskSize:
			d.applySizeRule(planConf, standardConfig, SysDisk)
		case models.CronTypeSsdSize:
			d.applySizeRule(planConf, standardConfig, SSD)
		case models.CronTypeHddSize:
			d.applySizeRule(planConf, standardConfig, HDD)
		case models.CronTypeTotalDataDiskSize:
			d.applySizeRule(planConf, standardConfig, TotalDataDisk)
		case models.CronTypeTotalMemorySize:
			d.applySizeRule(planConf, standardConfig, Memory)
		}
	}
}

// performDetection 执行检测（根据规则配置开关决定检测项）
// 根据规则开关，逐一检测各项是否超配，返回超配项列表和详情
func (d *OverProvisioningDetector) performDetection(
	rule *models.BusinessOverProvisioningRule,
	currentConfig, standardConfig map[string]float64,
) ([]string, string) {

	var overProvisionedItems []string
	var details []string

	// ==================== 硬件资源检测 ====================
	// 1. 内存检测
	if rule.EnableMemoryCheck {
		d.detectItem(Memory, currentConfig, standardConfig, &overProvisionedItems, &details)
	}

	// 2. SSD 检测
	if rule.EnableSsdCheck {
		d.detectItem(SSD, currentConfig, standardConfig, &overProvisionedItems, &details)
	}

	// 3. HDD 检测
	if rule.EnableHddCheck {
		d.detectItem(HDD, currentConfig, standardConfig, &overProvisionedItems, &details)
	}

	// 4. 业务盘总容量检测
	if rule.EnableDataDiskCheck {
		d.detectItem(TotalDataDisk, currentConfig, standardConfig, &overProvisionedItems, &details)
	}

	// 5. 系统盘检测
	if rule.EnableSysDiskCheck {
		d.detectItem(SysDisk, currentConfig, standardConfig, &overProvisionedItems, &details)
	}

	// ==================== 存储带宽比检测（独立检测项） ====================
	// 存储带宽比是独立检测项，有自己的检测逻辑，与SSD/HDD/业务盘分开检测
	if rule.EnableStorageBwRatioCheck {
		d.detectStorageBwRatio(rule, currentConfig, standardConfig, &overProvisionedItems, &details)
	}

	return overProvisionedItems, strings.Join(details, "; ")
}

// detectItem 检测单个项目是否超配（通用方法）
func (d *OverProvisioningDetector) detectItem(
	itemKey string,
	currentConfig, standardConfig map[string]float64,
	overProvisionedItems *[]string,
	details *[]string,
) {
	current, ok1 := currentConfig[itemKey]
	standard, ok2 := standardConfig[itemKey]
	if !ok1 || !ok2 || standard < 0 {
		return
	}

	if current > standard {
		*overProvisionedItems = append(*overProvisionedItems, d.translateItemName(itemKey))
		*details = append(*details, d.formatOverProvisioningDetail(itemKey, current, standard))
	}
}

// parseStorageBwRatio 解析存储带宽比规则值
// 格式如 "2:1" 表示 2TB/Gbps，返回比值（如 2.0）
func (d *OverProvisioningDetector) parseStorageBwRatio(ratioStr string) float64 {
	parts := strings.Split(ratioStr, ":")
	if len(parts) != 2 {
		return 0
	}
	num, err1 := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	den, err2 := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	if err1 != nil || err2 != nil || den == 0 {
		return 0
	}
	// 返回比值：存储容量(TB) / 带宽(Gbps)
	return num / den
}

// detectStorageBwRatio 检测存储带宽比（独立检测项）
// 实际存储带宽比 = 实际存储容量(TB) / 实际带宽(Gbps)
// 如果实际比值 > 规则要求的比值，则超配
func (d *OverProvisioningDetector) detectStorageBwRatio(
	rule *models.BusinessOverProvisioningRule,
	currentConfig map[string]float64,
	standardConfig map[string]float64,
	overProvisionedItems *[]string,
	details *[]string,
) error {
	// 1. 获取规则要求的存储带宽比值
	standardRatio, ok := standardConfig[StorageBwRatio]
	if !ok || standardRatio <= 0 {
		return fmt.Errorf("规则未配置存储带宽比或配置无效")
	}

	// 2. 根据 StorageBwRatioDiskType 选择实际存储容量（TB）
	var actualStorageTB float64
	switch rule.StorageBwRatioDiskType {
	case 0: // 仅 SSD
		actualStorageTB = currentConfig[SSD]
	case 1: // 仅 HDD
		actualStorageTB = currentConfig[HDD]
	default: // SSD+HDD (默认，diskType=2)
		actualStorageTB = currentConfig[TotalDataDisk]
	}

	// 3. 获取实际带宽（Gbps），优先使用提交带宽，没有则使用规划带宽
	actualBandwidth := currentConfig[SubmitBw]
	if actualBandwidth == 0 {
		actualBandwidth = currentConfig[PlanBw]
	}

	// 4. 如果存储容量或带宽为0，无法计算比值，跳过检测
	if actualStorageTB <= 0 || actualBandwidth <= 0 {
		return fmt.Errorf("存储容量或带宽为0，无法计算存储带宽比")
	}

	// 5. 计算实际存储带宽比 = 存储容量(TB) / 带宽(Gbps)
	actualRatio := actualStorageTB / actualBandwidth

	// 6. 保存实际比值到 currentConfig，用于输出
	currentConfig[StorageBwRatio] = actualRatio

	// 7. 检测：如果实际比值 > 规则要求的比值，则超配
	if actualRatio > standardRatio {
		*overProvisionedItems = append(*overProvisionedItems, d.translateItemName(StorageBwRatio))
		*details = append(*details, d.formatOverProvisioningDetail(StorageBwRatio, actualRatio, standardRatio))
	}

	return nil
}

// 应用大小规则
func (d *OverProvisioningDetector) applySizeRule(
	planConf models.PlanConf,
	standardConfig map[string]float64,
	sizeKey string,
) {
	if value, err := strconv.ParseFloat(planConf.Cond.Value, 64); err == nil {
		standardConfig[sizeKey] = value
	}
}

// evaluateCondition 评估条件
func (d *OverProvisioningDetector) evaluateCondition(cond models.CondDetail, currentValue float64, compareValue ...float64) bool {
	threshold, err := strconv.ParseFloat(cond.Value, 64)
	if err != nil {
		return false
	}

	if len(compareValue) > 0 {
		threshold = compareValue[0]
	}

	switch cond.Cond {
	case ">":
		return currentValue > threshold
	case ">=":
		return currentValue >= threshold
	case "=":
		return math.Abs(currentValue-threshold) < 0.001
	case "<":
		return currentValue < threshold
	case "<=":
		return currentValue <= threshold
	case "!=":
		return math.Abs(currentValue-threshold) >= 0.001
	default:
		return false
	}
}

// convertToTB 转换为TB
func (d *OverProvisioningDetector) convertToTB(volume string) float64 {
	sizeType := "G"
	if strings.Contains(strings.TrimSpace(volume), "T") {
		sizeType = "T"
	}
	sizeStr := strings.ReplaceAll(volume, sizeType, "")
	size, _ := strconv.ParseFloat(sizeStr, 64)
	if sizeType == "T" {
		return size
	}
	return size / 1000
}

// translateItemName 转换项目名称
func (d *OverProvisioningDetector) translateItemName(item string) string {
	nameMap := map[string]string{
		Memory:         "内存",
		SSD:            "SSD",
		HDD:            "HDD",
		TotalDataDisk:  "业务盘",
		StorageBwRatio: "存储带宽比",
	}
	if name, exists := nameMap[item]; exists {
		return name
	}
	return item
}

// formatOverProvisioningDetail 格式化超配详情
func (d *OverProvisioningDetector) formatOverProvisioningDetail(item string, current, standard float64) string {
	itemName := d.translateItemName(item)
	return fmt.Sprintf("%s超配：当前%.1f，标配%.1f", itemName, current, standard)
}

// IsMatchedRule 检查规则是否匹配（无论是否超配）
func (d *OverProvisioningDetector) IsMatchedRule(
	rule *models.BusinessOverProvisioningRule,
	hardware *models.Hardware,
	promInfo *PrometheusHostInfo,
) bool {
	// 1. 跨省调度判断：规则要求跨省调度，但设备不是跨省，则不匹配
	if rule.IsProvinceScheduling == 1 {
		if promInfo == nil || !strings.Contains(promInfo.CactiNotes, "跨省-") {
			return false
		}
	}

	// 2. 解析当前配置
	currentConfig, err := d.parseHardwareConfig(hardware)
	if err != nil {
		return false
	}

	bandwidthInfo, err := d.parseBandwidthInfo(promInfo)
	if err != nil {
		return false
	}

	for k, v := range bandwidthInfo {
		currentConfig[k] = v
	}

	// 3. 评估所有 PlanConf 条件（AND 关系）
	return d.evaluateAllConditions(rule.PlanConf, currentConfig)
}

// IsMatchedRuleForKM 检查规则是否匹配（用于 KM 宿主机，无 promInfo）
// KM 宿主机没有 Prometheus 信息，无法判断跨省调度和带宽相关条件
func (d *OverProvisioningDetector) IsMatchedRuleForKM(
	rule *models.BusinessOverProvisioningRule,
	hardware *models.Hardware,
	promInfo *PrometheusHostInfo,
) bool {
	// 1. 跨省调度判断：KM 宿主机无法判断跨省调度，如果规则要求跨省则不匹配
	if rule.IsProvinceScheduling == 1 {
		return false
	}

	// 2. 解析当前配置（仅硬件信息）
	currentConfig, err := d.parseHardwareConfig(hardware)
	if err != nil {
		return false
	}

	// KM 宿主机没有 promInfo，无法获取带宽信息
	// 但如果规则中有带宽相关条件，需要特殊处理

	// 3. 评估所有 PlanConf 条件（AND 关系）
	return d.evaluateAllConditionsForKM(rule.PlanConf, currentConfig)
}

// evaluateAllConditionsForKM 评估所有条件组合（KM 宿主机版本）
// 对于带宽相关条件，如果 currentConfig 中没有对应值，则跳过该条件
func (d *OverProvisioningDetector) evaluateAllConditionsForKM(
	planConfs []models.PlanConf,
	currentConfig map[string]float64,
) bool {
	// 如果没有条件，则默认匹配
	if len(planConfs) == 0 {
		return true
	}

	// 当前实现：所有条件采用 AND 关系（标准配置 cond_type=0 不参与条件判断）
	for _, planConf := range planConfs {
		if planConf.CondType == models.CronTypeStandard {
			continue
		}

		// 对于带宽相关条件，如果 currentConfig 中没有值，跳过该条件
		if planConf.CondType == models.CronTypeSubmitBandwidth ||
			planConf.CondType == models.CronTypePlanBandwidth {
			if _, ok := currentConfig[SubmitBw]; !ok {
				continue
			}
		}

		matched := d.evaluateConditionValue(planConf, currentConfig)
		if !matched {
			return false
		}
	}

	return true
}

// evaluateAllConditions 评估所有条件组合
func (d *OverProvisioningDetector) evaluateAllConditions(
	planConfs []models.PlanConf,
	currentConfig map[string]float64,
) bool {
	// 如果没有条件，则默认匹配
	if len(planConfs) == 0 {
		return true
	}

	// 当前实现：所有条件采用 AND 关系（标准配置 cond_type=0 不参与条件判断）
	for _, planConf := range planConfs {
		if planConf.CondType == models.CronTypeStandard {
			continue
		}

		matched := d.evaluateConditionValue(planConf, currentConfig)
		if !matched {
			return false
		}
	}

	return true
}

// evaluateConditionValue 评估单个条件值
func (d *OverProvisioningDetector) evaluateConditionValue(
	planConf models.PlanConf,
	currentConfig map[string]float64,
) bool {
	var currentValue float64

	switch planConf.CondType {
	case models.CronTypeSubmitBandwidth:
		currentValue = currentConfig[SubmitBw]
	case models.CronTypePlanBandwidth:
		currentValue = currentConfig[PlanBw]
	case models.CronTypeStorageBwRatio:
		currentValue = currentConfig[TotalDataDisk]
	case models.CronTypeSystemDiskSize:
		currentValue = currentConfig[SysDisk]
	case models.CronTypeSsdSize:
		currentValue = currentConfig[SSD]
	case models.CronTypeHddSize:
		currentValue = currentConfig[HDD]
	case models.CronTypeTotalDataDiskSize:
		currentValue = currentConfig[TotalDataDisk]
	case models.CronTypeTotalMemorySize:
		currentValue = currentConfig[Memory]
	case models.CronTypeDiskType:
		currentValue = currentConfig[DiskType]
	default:
		return true // 其他类型暂时默认匹配
	}

	return d.evaluateCondition(planConf.Cond, currentValue)
}
