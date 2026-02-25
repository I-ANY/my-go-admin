package overprovisioning

import (
	"biz-auto-api/internal/models"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Logger interface {
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Debugf(format string, args ...interface{})
}

type Executor struct {
	db     *gorm.DB
	ecdnDB *gorm.DB
	log    Logger
	config *Config
}

func NewExecutor(db, ecdnDB *gorm.DB, log Logger, config *Config) *Executor {
	return &Executor{
		db:     db,
		ecdnDB: ecdnDB,
		log:    log,
		config: config,
	}
}

// Execute 执行超配检测
func (e *Executor) Execute(arg ExecArg) error {
	switch arg.ExecType {
	case "kvm":
		// KM 宿主机清单生成（不做超配判定）
		if err := e.syncKMHostOverProvisioning(); err != nil {
			return err
		}
	case "host":
		// 默认执行业务超配检测（物理机）
		if err := e.syncBusinessOverProvisioningHostname(); err != nil {
			return err
		}
	case "all":
		// 同时执行KM 宿主机清单生成（不做超配判定）和业务超配检测（物理机）
		// 先执行业务超配检测，再执行KM 宿主机清单生成
		if err := e.syncBusinessOverProvisioningHostname(); err != nil {
			return err
		}
		if err := e.syncKMHostOverProvisioning(); err != nil {
			return err
		}
	default:
		return errors.New("invalid exec type")
	}

	e.log.Infof("超配检测完成")
	return nil
}

// getAllBizWithRules 获取存在启用规则的业务小类列表
func (e *Executor) getAllBizWithRules() ([]BizWithRules, error) {
	var bizList []BizWithRules

	// 1. 先取有"专属规则"的业务，并按业务维度去重，通过业务小类/大类表获取业务大类名称
	err := e.db.
		Table("business_over_provisioning_rule_biz AS rb").
		Joins("JOIN business_over_provisioning_rule AS r ON r.id = rb.rule_id").
		Joins("JOIN business_subcategory AS bs ON bs.id = rb.business_id").
		Joins("LEFT JOIN business_category AS bc ON bs.category_id = bc.id").
		Select("DISTINCT rb.business_id, rb.business_name, bc.name AS business_category").
		Where("rb.business_id != 0 AND r.rule_status = 0 AND r.deleted_at IS NULL AND bs.status = 1").
		Scan(&bizList).Error
	if err != nil {
		return nil, err
	}

	// 2. 判断是否存在至少一条"全业务通用"规则（business_id=0）
	var globalRuleCount int64
	if err := e.db.
		Model(&models.BusinessOverProvisioningRule{}).
		Joins("JOIN business_over_provisioning_rule_biz AS rb ON rb.rule_id = business_over_provisioning_rule.id").
		Where("rb.business_id = 0 AND business_over_provisioning_rule.deleted_at IS NULL AND business_over_provisioning_rule.rule_status = 0").
		Count(&globalRuleCount).Error; err != nil {
		return nil, err
	}

	// 3. 如果存在全业务规则，则需要把所有"已启用业务小类"也加入检测范围
	if globalRuleCount > 0 {
		var allBiz []BizWithRules
		if err := e.db.
			Model(&models.BusinessSubcategory{}).
			Joins("LEFT JOIN business_category AS bc ON business_subcategory.category_id = bc.id").
			Where("business_subcategory.status = 1").
			Select("business_subcategory.id AS business_id, business_subcategory.name AS business_name, bc.name AS business_category").
			Scan(&allBiz).Error; err != nil {
			return nil, err
		}

		// 合并已有的 bizList，按 BusinessId 去重
		exists := make(map[int64]struct{}, len(bizList))
		for _, b := range bizList {
			exists[b.BusinessId] = struct{}{}
		}
		for _, b := range allBiz {
			if _, ok := exists[b.BusinessId]; ok {
				continue
			}
			// 复制一份，避免引用问题
			biz := b
			bizList = append(bizList, biz)
		}
	}

	return bizList, nil
}

// getBizActiveRulesCore 查询某业务下所有生效的超配规则
// businessId 和 businessName 两者只需提供一个，另一个传 nil
func (e *Executor) getBizActiveRulesCore(businessId *int64, businessName *string) ([]models.BusinessOverProvisioningRule, error) {
	var rules []models.BusinessOverProvisioningRule
	query := e.db.
		Model(&models.BusinessOverProvisioningRule{}).
		Joins("JOIN business_over_provisioning_rule_biz AS rb ON rb.rule_id = business_over_provisioning_rule.id").
		Where("business_over_provisioning_rule.deleted_at IS NULL AND business_over_provisioning_rule.rule_status = 0")

	// 根据传入参数决定查询条件
	if businessId != nil {
		query = query.Where("(rb.business_id = ? OR rb.business_id = 0)", *businessId)
	} else if businessName != nil {
		query = query.Where("(rb.business_name = ? OR rb.business_id = 0)", *businessName)
	} else {
		return nil, errors.New("either businessId or businessName must be provided")
	}

	err := query.Order("business_over_provisioning_rule.id DESC").Find(&rules).Error
	return rules, err
}

// syncBusinessOverProvisioningHostname 物理机超配检测
func (e *Executor) syncBusinessOverProvisioningHostname() error {
	e.log.Infof("sync_business_over_provisioning_hostname")
	bizList, err := e.getAllBizWithRules()
	if err != nil {
		return err
	}
	if len(bizList) <= 0 {
		e.log.Infof("没有维护业务规则")
		return nil
	}

	// 1. 清理当日物理机超配检测记录
	collectDate := time.Now().Format("2006-01-02")
	if err := e.db.
		Where("collect_date = ? and machine_type = 1", collectDate).
		Delete(&models.BusinessOverProvisioningRecord{}).Error; err != nil {
		e.log.Errorf("清理历史检测记录失败: %+v", err)
		return err
	}

	// 2. 执行物理机超配检测
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 3)
	for _, biz := range bizList {
		e.log.Infof("biz with rules: %+v", biz)
		wg.Add(1)
		semaphore <- struct{}{}
		b := biz
		go func() {
			e.log.Infof("business: %+v", b)
			defer wg.Done()
			if err := e.businessOverProvisioningInspection(&b); err != nil {
				e.log.Errorf("syncBusinessOverProvisioningHostnamebyBsi error: %+v", err)
				return
			}
			defer func() { <-semaphore }()
		}()
	}
	wg.Wait()
	close(semaphore)

	return nil
}

// businessOverProvisioningInspection 针对单个业务小类执行主机超配检测
func (e *Executor) businessOverProvisioningInspection(biz *BizWithRules) error {
	if biz == nil || biz.BusinessName == nil {
		e.log.Infof("business info is nil, skip")
		return nil
	}

	businessName := *biz.BusinessName
	e.log.Infof("执行业务[%s]超配检测", businessName)

	// 创建 Prometheus 客户端
	prom := NewProm(e.config.PromURL, e.config.PromAuth)
	hosts := NewHosts(prom, e.ecdnDB)

	promHosts, err := hosts.GetBizHosts(businessName)
	if err != nil {
		e.log.Errorf("获取prom信息错误: %+v", err)
		return err
	}

	// 查询该业务下生效的超配规则
	rules, err := e.getBizActiveRulesCore(&biz.BusinessId, nil)
	if err != nil {
		e.log.Errorf("获取业务规则错误: %+v", err)
		return err
	}
	if len(rules) == 0 {
		e.log.Infof("业务[%s]无生效规则，跳过", businessName)
		return nil
	}

	data, err := hosts.GetECdnHost(businessName)
	if err != nil {
		e.log.Errorf("获取ECDN主机信息错误: %+v", err)
		return err
	}

	// 创建检测器
	detector := NewOverProvisioningDetector()
	collectDate := time.Now().Format("2006-01-02")

	// 解析业务组
	businessGroup := "" // 按需求暂时保留为空
	businessCategory := ""
	if biz.BusinessCategory != nil {
		businessCategory = *biz.BusinessCategory
	}

	// 设备外循环，规则内循环
	for id, hardware := range data {
		// 跳过虚拟机
		promHost, ok := promHosts[id]
		if !ok || isKvm(promHost.Hostname) {
			continue
		}

		// 构造Prometheus主机信息
		promInfo := &PrometheusHostInfo{
			ID:              promHost.ID,
			Hostname:        promHost.Hostname,
			BwPlan:          promHost.BwPlan,
			BwCount:         promHost.BwCount,
			BwSingle:        promHost.BwSingle,
			Day95:           promHost.Day95,
			Evening95:       promHost.Evening95,
			Interprovincial: promHost.Interprovincial,
			Kvm:             promHost.Kvm,
			Location:        promHost.Location,
			Owner:           promHost.Owner,
			IP:              promHost.IP,
			CactiNotes:      promHost.CactiNotes,
		}

		// 首个匹配规则模式：找到第一条匹配且检测为超配的规则后立即返回
		var matchedRule *models.BusinessOverProvisioningRule
		var detectionResult *DetectionResult

		// 遍历规则，找到第一条匹配且检测为超配的规则
		for i := range rules {
			rule := &rules[i]

			// 先判断规则是否适用于当前主机
			if !detector.IsMatchedRule(rule, hardware, promInfo) {
				continue
			}

			// 规则适用，执行超配检测
			result, err := detector.DetectOverProvisioning(rule, hardware, promInfo)
			if err != nil {
				e.log.Errorf("主机[%s]规则[%s]检测失败: %+v", promHost.Hostname, rule.RuleName, err)
				continue
			}

			if result.IsOverProvisioned {
				matchedRule = rule
				detectionResult = result
				break
			}
		}

		// 只保存有超配项的检测结果
		if matchedRule != nil && detectionResult != nil && detectionResult.IsOverProvisioned {
			standardStr := formatStandardValues(detectionResult.StandardValues)
			err := e.saveDetectionResult(businessGroup, businessName, businessCategory, promHost, matchedRule.RuleName, standardStr, detectionResult, collectDate)
			if err != nil {
				e.log.Errorf("保存主机[%s]超配检测结果失败: %+v", promHost.Hostname, err)
				continue
			}
			e.log.Infof("已保存主机[%s]超配检测结果，超配项: %s", promHost.Hostname, strings.Join(detectionResult.OverProvisionedItems, "|"))
		}
	}

	return nil
}

// KM 宿主机相关类型和函数

type kmHostInfo struct {
	Hostname       string
	Node           string
	Sn             string
	BusinessServer *models.BusinessServer
	Hardware       *models.Hardware
	ActualBizNames []string // 实际业务名（去重后，含完整后缀）
	MainBizNames   []string // 主业务名（去重后，仅用于构造混跑名称）
	SubmitBwGbps   float64  // 提交带宽汇总（Gbps）：宿主机下所有VM的 bwcount*bwsingle 之和
}

func isKvm(hostname string) bool {
	return len(hostname) > 15 && (string(hostname[9]) == "V" || string(hostname[9]) == "U")
}

// syncKMHostOverProvisioning KM 宿主机相关任务：
// 1）混跑宿主机（多业务）：生成 KM 宿主机混跑路由清单（结果标记为正常，仅记录当前配置）
// 2）单业务宿主机：按该业务规则执行超配检测，保存超配结果
func (e *Executor) syncKMHostOverProvisioning() error {
	e.log.Infof("start KM host inspection")
	collectDate := time.Now().Format("2006-01-02")

	// 清理当日该业务的旧检测记录，避免重复数据
	if err := e.db.
		Where("collect_date = ? and machine_type = 3", collectDate).
		Delete(&models.BusinessOverProvisioningRecord{}).Error; err != nil {
		e.log.Errorf("清理历史检测记录失败: %+v", err)
		return err
	}

	// 1. 构建 KM 宿主机聚合信息（不涉及写库）
	kmHosts, err := e.buildKMHostInfos()
	if err != nil {
		return err
	}
	if len(kmHosts) == 0 {
		e.log.Infof("未找到有效的KM宿主机，跳过KM检测")
		return nil
	}

	// 2. 分类处理 KM 宿主机
	for parent, info := range kmHosts {
		if len(info.ActualBizNames) == 1 {
			// 2.1 单业务宿主机：按该业务规则进行超配检测
			bizName := info.ActualBizNames[0]
			if err := e.detectSingleBizKMHost(parent, bizName, info, collectDate); err != nil {
				e.log.Errorf("KM宿主机[%s]业务[%s]超配检测失败: %+v", parent, bizName, err)
			}
		}
		// 混跑宿主机（len > 1）在 saveKMHostInventory 中批量处理
		// 无业务的宿主机（len == 0）跳过
	}

	// 3. 生成 KM 宿主机混跑清单（仅记录当前配置，结果标记为正常）
	if err := e.saveKMHostInventory(kmHosts, collectDate); err != nil {
		return err
	}

	return nil
}

// buildKMHostInfos 构建 KM 宿主机聚合信息（不涉及写库）
func (e *Executor) buildKMHostInfos() (map[string]*kmHostInfo, error) {
	// 1. 查询 business_server 中所有业务名为 KM 且非休眠状态的设备（KM 宿主机）
	var kmServers []models.BusinessServer
	if err := e.db.
		Model(&models.BusinessServer{}).
		Where("business = ?", "KM").
		Where("status != ?", models.ServerStatusDormant).
		Find(&kmServers).Error; err != nil {
		e.log.Errorf("查询KM宿主机失败: %+v", err)
		return nil, err
	}
	if len(kmServers) == 0 {
		return map[string]*kmHostInfo{}, nil
	}

	// 构建 KM 宿主机映射：hostname -> kmHostInfo（包含 BusinessServer）
	kmHosts := make(map[string]*kmHostInfo)
	for i := range kmServers {
		bs := &kmServers[i]
		if bs.Hostname == nil || *bs.Hostname == "" {
			continue
		}
		host := *bs.Hostname
		info := &kmHostInfo{
			Hostname:       host,
			BusinessServer: bs,
		}
		if bs.Owner != nil {
			info.Node = *bs.Owner
		}
		if bs.Sn != nil {
			info.Sn = *bs.Sn
		}
		kmHosts[host] = info
	}
	if len(kmHosts) == 0 {
		return map[string]*kmHostInfo{}, nil
	}

	// 2. 从 Prometheus 获取所有 KVM 设备（parent!=""），并按 parent 过滤到 KM 宿主机
	prom := NewProm(e.config.PromURL, e.config.PromAuth)
	query := `mfy_hosts_ecdn_info{origin="自建", parent!=""}`
	bytes, err := prom.Fetch(query, 30*time.Second)
	if err != nil {
		return nil, err
	}

	var hosts HostsResp
	if err := json.Unmarshal(bytes, &hosts); err != nil {
		return nil, err
	}

	// parent -> 实际业务名/主业务名去重
	for _, r := range hosts.Data.Result {
		parent := r.Metric.Parent
		if parent == "" {
			continue
		}
		info, ok := kmHosts[parent]
		if !ok {
			continue
		}
		actualBiz := r.Metric.Business
		if actualBiz == "" {
			continue
		}

		// 实际业务名去重
		if info.ActualBizNames == nil {
			info.ActualBizNames = make([]string, 0)
		}
		hasActual := false
		for _, b := range info.ActualBizNames {
			if b == actualBiz {
				hasActual = true
				break
			}
		}
		if !hasActual {
			info.ActualBizNames = append(info.ActualBizNames, actualBiz)
		}

		// 主业务名（取 '_' 前缀）去重，仅用于构造混跑名称
		mainName := actualBiz
		if idx := strings.Index(actualBiz, "_"); idx > 0 {
			mainName = actualBiz[:idx]
		}
		if mainName == "" {
			continue
		}
		if info.MainBizNames == nil {
			info.MainBizNames = make([]string, 0)
		}
		hasMain := false
		for _, b := range info.MainBizNames {
			if b == mainName {
				hasMain = true
				break
			}
		}
		if !hasMain {
			info.MainBizNames = append(info.MainBizNames, mainName)
		}

		// 累加提交带宽：submit_bw = bwcount * bwsingle / 1e9 (Gbps)
		if r.Metric.BwCount != "" && r.Metric.BwSingle != "" {
			if cnt, err1 := strconv.ParseFloat(r.Metric.BwCount, 64); err1 == nil {
				if single, err2 := strconv.ParseFloat(r.Metric.BwSingle, 64); err2 == nil {
					info.SubmitBwGbps += (cnt * single) / 1000 / 1000 / 1000
				}
			}
		}
	}

	// 2.5. 批量预加载所有 KM 宿主机的 ServersMore 硬件信息，减少数据库查询次数
	serverIDSet := make(map[string]struct{})
	for _, info := range kmHosts {
		bs := info.BusinessServer
		if bs != nil && bs.ServerID != nil && *bs.ServerID != "" {
			serverIDSet[*bs.ServerID] = struct{}{}
		}
	}
	hardwareCache := make(map[string]*models.Hardware)
	if len(serverIDSet) > 0 {
		ids := make([]string, 0, len(serverIDSet))
		for id := range serverIDSet {
			ids = append(ids, id)
		}

		// 分批查询，每批 500 条
		batchSize := 500
		totalBatches := (len(ids) + batchSize - 1) / batchSize
		for i := 0; i < len(ids); i += batchSize {
			end := i + batchSize
			if end > len(ids) {
				end = len(ids)
			}
			batchIDs := ids[i:end]
			batchNum := (i / batchSize) + 1

			var serversMore []models.ServersMore
			if err := e.ecdnDB.
				Model(&models.ServersMore{}).
				Where("id IN ?", batchIDs).
				Find(&serversMore).Error; err != nil {
				e.log.Errorf("批量预加载KM宿主机ServersMore失败(批次 %d/%d): %+v", batchNum, totalBatches, err)
				// 出错时不直接返回，继续处理下一批
				continue
			}

			for j := range serversMore {
				sm := &serversMore[j]
				hardwareCache[sm.ID] = &sm.Hardware
			}
		}
		e.log.Infof("批量预加载了%d条KM宿主机硬件信息(共%d批)", len(hardwareCache), totalBatches)
	}

	// 3. 为每个 KM 宿主机补全硬件信息
	for parent, info := range kmHosts {
		hw, err := e.getKMHostHardware(parent, info.BusinessServer, hardwareCache)
		if err != nil {
			e.log.Errorf("获取KM宿主机[%s]硬件信息失败: %+v", parent, err)
			continue
		}
		info.Hardware = hw
	}

	return kmHosts, nil
}

// getKMHostHardware 根据 KM 宿主机名和 BusinessServer 获取硬件信息
// hardwareCache: 预加载的硬件信息缓存，key 为 serverID
func (e *Executor) getKMHostHardware(parent string, bs *models.BusinessServer, hardwareCache map[string]*models.Hardware) (*models.Hardware, error) {
	// 1. 优先使用预加载缓存（通过 ServerID 命中）
	if bs != nil && bs.ServerID != nil && *bs.ServerID != "" && hardwareCache != nil {
		if hw, ok := hardwareCache[*bs.ServerID]; ok && hw != nil {
			return hw, nil
		}
	}

	// 2. 通过 BusinessServer 的 ServerID 单次查询（降级方案，如果缓存未命中）
	// if bs != nil && bs.ServerID != nil && *bs.ServerID != "" {
	// 	var sm models.ServersMore
	// 	if err := e.ecdnDB.
	// 		Model(&models.ServersMore{}).
	// 		Where("id = ?", *bs.ServerID).
	// 		First(&sm).Error; err == nil {
	// 		return &sm.Hardware, nil
	// 	}
	// }

	// 3. 通过 Prometheus -> ServersMore 获取（按 hostname 精确匹配）
	prom := NewProm(e.config.PromURL, e.config.PromAuth)
	query := fmt.Sprintf(`mfy_hosts_ecdn_info{origin="自建",hostname="%s"}`, parent)
	bytes, err := prom.Fetch(query, 30*time.Second)
	if err != nil {
		return nil, err
	}

	var hostResp HostsResp
	if err := json.Unmarshal(bytes, &hostResp); err != nil {
		return nil, err
	}
	var serverID string
	for _, r := range hostResp.Data.Result {
		if r.Metric.Hostname == parent {
			serverID = r.Metric.ID
			break
		}
	}
	if serverID == "" {
		return nil, fmt.Errorf("未在Prometheus中找到KM宿主机[%s]对应的server_id", parent)
	}

	var sm models.ServersMore
	if err := e.ecdnDB.
		Model(&models.ServersMore{}).
		Where("id = ?", serverID).
		First(&sm).Error; err != nil {
		return nil, err
	}
	return &sm.Hardware, nil
}

// detectSingleBizKMHost 对单业务的 KM 宿主机按业务规则执行超配检测
func (e *Executor) detectSingleBizKMHost(
	parent string,
	bizName string,
	info *kmHostInfo,
	collectDate string,
) error {
	// 1. 查询该业务的所有生效规则
	rules, err := e.getBizActiveRulesCore(nil, &bizName)
	if err != nil {
		return fmt.Errorf("查询业务规则失败: %w", err)
	}
	if len(rules) == 0 {
		e.log.Infof("业务[%s]无生效规则，跳过KM宿主机[%s]检测", bizName, parent)
		return nil
	}

	// 2. 创建检测器
	detector := NewOverProvisioningDetector()

	// 3. 构造包含提交带宽的 Prometheus 信息（用于 DetectOverProvisioning）
	promInfo := &PrometheusHostInfo{}
	if info.SubmitBwGbps > 0 {
		promInfo.BwCount = fmt.Sprintf("%.0f", info.SubmitBwGbps) // 数量（Gbps）
		promInfo.BwSingle = "1000000000"                          // 1Gbps = 1e9 bps
	}

	// 4. 遍历规则，找到第一条匹配且超配的规则
	for i := range rules {
		rule := &rules[i]

		// 4.1 判断规则是否适用（不匹配则跳过）
		if !detector.IsMatchedRuleForKM(rule, info.Hardware, promInfo) {
			continue
		}

		// 4.2 使用 DetectOverProvisioning 执行检测（统一使用标准检测方法）
		result, err := detector.DetectOverProvisioning(rule, info.Hardware, promInfo)
		if err != nil {
			// 规则不适用（如磁盘类型不匹配），跳过
			e.log.Debugf("KM宿主机[%s]规则[%s]检测失败: %+v", parent, rule.RuleName, err)
			continue
		}

		// 4.3 如果超配，保存结果并退出
		if result.IsOverProvisioned {
			return e.saveKMHostDetectionResultFromResult(parent, bizName, info, rule, result, collectDate)
		}
	}

	return nil
}

// saveKMHostDetectionResultFromResult 保存KM宿主机超配检测结果（从 DetectionResult）
func (e *Executor) saveKMHostDetectionResultFromResult(
	parent, bizName string,
	info *kmHostInfo,
	rule *models.BusinessOverProvisioningRule,
	result *DetectionResult,
	collectDate string,
) error {
	currentCfgStr := formatCurrentValues(result.CurrentValues)
	standardCfgStr := formatStandardValues(result.StandardValues)

	record := &models.BusinessOverProvisioningRecord{
		CollectDate:          collectDate,
		RuleName:             rule.RuleName,
		BusinessGroup:        "",
		BusinessCategory:     "",
		BusinessName:         bizName,
		Node:                 info.Node,
		Hostname:             parent,
		Sn:                   info.Sn,
		HostMachine:          "",                                         // 宿主机本身，无需再填宿主机名
		MachineType:          int64(3),                                   // 宿主机
		Result:               int64(models.OverProvisioningResultFailed), // 超配
		OverProvisioningItem: getOverProvisioningItems(result.OverProvisionedItems),
		Standard:             standardCfgStr,
		CurrentConfiguration: currentCfgStr,
	}
	return e.db.Create(record).Error
}

// saveKMHostInventory 生成 KM 宿主机混跑清单（仅记录当前配置，结果标记为正常）
func (e *Executor) saveKMHostInventory(kmHosts map[string]*kmHostInfo, collectDate string) error {
	var records []models.BusinessOverProvisioningRecord
	for parent, info := range kmHosts {
		// 仅为"混跑"宿主机生成清单记录：实际业务名去重后数量必须 > 1
		if len(info.ActualBizNames) <= 1 {
			continue
		}

		// KM宿主机业务名生成规则：混跑_主业务1_主业务2_...
		mixedBizName := ""
		if len(info.MainBizNames) > 0 {
			mixedBizName = "混跑_" + strings.Join(info.MainBizNames, "_")
		} else {
			mixedBizName = "混跑_" + strings.Join(info.ActualBizNames, "_")
		}

		currentConfig := buildKMHostCurrentConfig(info.Hardware, info.SubmitBwGbps)
		currentCfgStr := formatCurrentValues(currentConfig)

		record := models.BusinessOverProvisioningRecord{
			CollectDate:          collectDate,
			RuleName:             "",
			BusinessGroup:        "",
			BusinessCategory:     "",
			BusinessName:         mixedBizName,
			Node:                 info.Node,
			Hostname:             parent,
			Sn:                   info.Sn,
			HostMachine:          "",                                         // 宿主机本身，无需再填宿主机名
			MachineType:          int64(3),                                   // 宿主机
			Result:               int64(models.OverProvisioningResultNormal), // 结果标记为正常
			SubmittedBandwidth:   0,
			OverProvisioningItem: "",
			Standard:             "",
			CurrentConfiguration: currentCfgStr,
		}

		records = append(records, record)
	}

	if len(records) == 0 {
		return nil
	}
	if err := e.db.Create(&records).Error; err != nil {
		e.log.Errorf("批量保存KM宿主机混跑记录失败: %+v", err)
		return err
	}
	e.log.Infof("已批量保存%d条KM宿主机混跑记录", len(records))
	return nil
}

// buildKMHostCurrentConfig 构造KM宿主机当前硬件配置（仅关注内存/系统盘/SSD/HDD/业务盘）
func buildKMHostCurrentConfig(hardware *models.Hardware, submitBwGbps float64) map[string]float64 {
	config := make(map[string]float64)
	if hardware == nil {
		return config
	}

	// 内存（GB）
	if hardware.Memory != "" {
		memStr := strings.ReplaceAll(strings.ReplaceAll(hardware.Memory, "G", ""), "M", "")
		if mem, err := strconv.ParseFloat(memStr, 64); err == nil {
			if strings.Contains(hardware.Memory, "M") {
				mem = mem / 1024 // M 转 GB
			}
			config["内存"] = mem
		}
	}

	// 系统盘（GB，按 T->1000G 计算）
	if hardware.DiskDetail.SysDisk.Volume != "" {
		sizeType := "G"
		vol := strings.TrimSpace(hardware.DiskDetail.SysDisk.Volume)
		if strings.Contains(vol, "T") {
			sizeType = "T"
		}
		sizeStr := strings.ReplaceAll(vol, sizeType, "")
		if size, err := strconv.ParseFloat(sizeStr, 64); err == nil {
			if sizeType == "T" {
				config["系统盘"] = size * 1000
			} else {
				config["系统盘"] = size
			}
		}
	}

	// 数据盘：SSD/HDD/业务盘（容量单位：TB）
	var ssdSize, hddSize, totalDataDisk float64
	for _, disk := range hardware.DiskDetail.DataDisk.DataDetail {
		vol := strings.TrimSpace(disk.Volume)
		sizeType := "G"
		if strings.Contains(vol, "T") {
			sizeType = "T"
		}
		sizeStr := strings.ReplaceAll(vol, sizeType, "")
		size, _ := strconv.ParseFloat(sizeStr, 64)
		if sizeType == "G" {
			size = size / 1000
		}

		totalDataDisk += size
		if disk.DiskType == "hdd" {
			hddSize += size
		} else {
			ssdSize += size
		}
	}

	if ssdSize > 0 {
		config["SSD"] = ssdSize
	}
	if hddSize > 0 {
		config["HDD"] = hddSize
	}
	if totalDataDisk > 0 {
		config["业务盘"] = totalDataDisk
	}
	// 提交带宽（Gbps，宿主机下所有VM提交带宽累加）
	if submitBwGbps > 0 {
		config["提交带宽"] = submitBwGbps
	}
	return config
}

// saveDetectionResult 保存超配检测结果
func (e *Executor) saveDetectionResult(
	businessGroup string,
	businessName string,
	businessCategory string,
	promHost *PromHost,
	ruleName string,
	standard string,
	result *DetectionResult,
	collectDate string,
) error {
	kvm, err := strconv.ParseInt(promHost.Kvm, 10, 64)
	if err != nil {
		return err
	}
	machineType := getMachineType(kvm)

	// 只保存超配的检测结果，状态固定为超配
	record := &models.BusinessOverProvisioningRecord{
		CollectDate:          collectDate,
		RuleName:             ruleName,
		BusinessGroup:        businessGroup,
		BusinessCategory:     businessCategory,
		BusinessName:         businessName,
		Node:                 promHost.Owner,
		Hostname:             promHost.Hostname,
		Sn:                   promHost.Sn,
		MachineType:          int64(machineType),
		Result:               int64(models.OverProvisioningResultFailed), // 固定为超配状态
		OverProvisioningItem: getOverProvisioningItems(result.OverProvisionedItems),
		Standard:             standard,
		CurrentConfiguration: formatCurrentValues(result.CurrentValues),
	}

	// 保存到数据库
	return e.db.Create(record).Error
}

// getMachineType 获取机型类型
func getMachineType(kvm int64) int64 {
	if kvm == 1 {
		return 2
	}
	return 1
}

// getOverProvisioningItems 获取超配项目
func getOverProvisioningItems(items []string) string {
	if len(items) == 0 {
		return ""
	}
	result := ""
	for i, item := range items {
		if i > 0 {
			result += "|"
		}
		result += item
	}
	return result
}

// formatStandardValues 格式化标准值
func formatStandardValues(values map[string]float64) string {
	// Standard 输出顺序（存在才输）：
	// 1. 内存 2. SSD 3. HDD 4. 业务盘 5. 提交带宽 6. 规划带宽 7. 95值 8. 晚高峰95值
	order := []string{
		"内存",
		"SSD",
		"HDD",
		"业务盘",
		"提交带宽",
		"规划带宽",
		"95值",
		"晚高峰95值",
	}
	return formatValuesWithOrder(values, order)
}

// formatCurrentValues 格式化当前值
func formatCurrentValues(values map[string]float64) string {
	// CurrentConfiguration 输出顺序（存在才输）：
	// 1. 内存 2. 系统盘 3. SSD 4. HDD 5. 业务盘 6. 提交带宽 7. 规划带宽 8. 95值 9. 晚高峰95值
	order := []string{
		"内存",
		"系统盘",
		"SSD",
		"HDD",
		"业务盘",
		"提交带宽",
		"规划带宽",
		"95值",
		"晚高峰95值",
	}
	return formatValuesWithOrder(values, order)
}

// formatValuesWithOrder 按指定顺序格式化值
func formatValuesWithOrder(values map[string]float64, order []string) string {
	var parts []string
	for _, key := range order {
		value, ok := values[key]
		if !ok {
			continue
		}

		// day95 / evening95 以 Gbps 保留两位小数，其它保留一位
		format := "%.1f"
		if key == "95值" || key == "晚高峰95值" {
			format = "%.2f"
		}

		// 不同指标的单位说明
		unit := ""
		switch key {
		case "内存", "系统盘":
			unit = "GB"
		case "SSD", "HDD", "业务盘":
			unit = "TB"
		case "提交带宽", "规划带宽", "95值", "晚高峰95值":
			unit = "Gbps"
		}

		parts = append(parts, fmt.Sprintf("%s:"+format+"%s", key, value, unit))
	}
	return strings.Join(parts, "|")
}
