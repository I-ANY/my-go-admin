package switch_collector

import (
	"biz-auto-api/pkg/tools"
	"fmt"
	"github.com/gosnmp/gosnmp"
	"github.com/pkg/errors"
	"sort"
	"strconv"
	"strings"
	"time"
)

type SwitchCollector struct {
	*gosnmp.GoSNMP
}

func New(snmp *gosnmp.GoSNMP) (*SwitchCollector, error) {
	if snmp == nil {
		return nil, errors.New("snmp client is nil")
	}
	return &SwitchCollector{
		GoSNMP: snmp,
	}, nil
}

func (c *SwitchCollector) GetSwitchIfBusiness() ([]string, error) {
	allTags := make([]string, 0)
	err := c.BulkWalk(oidIfAlias, func(dataUnit gosnmp.SnmpPDU) error {
		name := dataUnit.Name
		value := dataUnit.Value
		if strings.HasPrefix(name, oidIfAlias) && !strings.Contains(strings.ReplaceAll(name, oidIfAlias+".", ""), ".") {
			tagBS, ok := value.([]uint8)
			if !ok {
				return errors.New("value type error")
			}
			tagStr := string(tagBS)
			tagStr = strings.Trim(tagStr, " ")
			tags := strings.Split(tagStr, " ")

			// 限速端口
			if len(tags) >= 2 {
				if strings.HasPrefix(tags[0], speedLimitTagPrefix) && len(tags[1]) > 0 {
					allTags = append(allTags, tags[1])
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, errors.WithMessage(err, "get switch if business failed")
	}
	allTags = tools.RemoveDuplication(allTags)
	allTags = tools.SliceFilter(allTags, func(tag string) bool {
		return tag != ""
	})
	return allTags, nil
}

func (c *SwitchCollector) GetSwitchIfBaseInfo() ([]*IfInfo, error) {
	var ifInfoMapping = make(map[int64]*IfInfo)
	snmpResults, err := c.BulkWalkAllOids([]string{oidIfOperStatus, oidIfAlias, oidIfName})
	if err != nil {
		return nil, errors.Wrap(err, "get if oper status and alias failed")
	}
	// 获取名称
	for _, r := range snmpResults {
		for _, dataUnit := range r.Results {
			name := dataUnit.Name
			value := dataUnit.Value
			if c.OidBelongTo(name, oidIfName) {
				index, err := c.GetIndex(oidIfName, name)
				if err != nil {
					return nil, err
				}
				value, ok := value.([]uint8)
				if !ok {
					return nil, errors.New("value type error")
				}
				if ifInfo, exist := ifInfoMapping[index]; exist {
					ifInfo.Name = string(value)
				} else {
					ifInfoMapping[index] = &IfInfo{
						Index: index,
						Name:  string(value),
					}
				}
			}
		}
	}
	// 获取其他信息
	for _, snmpResult := range snmpResults {
		for _, dataUnit := range snmpResult.Results {
			name := dataUnit.Name
			value := dataUnit.Value
			switch {
			case c.OidBelongTo(name, oidIfOperStatus): // 端口状态数据
				index, err := c.GetIndex(oidIfOperStatus, name)
				if err != nil {
					return nil, err
				}
				ifInfo, exist := ifInfoMapping[index]
				if !exist {
					ifInfo = &IfInfo{
						Index: index,
					}
				}
				statusId := gosnmp.ToBigInt(value).Int64()
				ifInfo.Status = IfStatusMapping[statusId]
				ifInfoMapping[index] = ifInfo
			case c.OidBelongTo(name, oidIfDescr):
				index, err := c.GetIndex(oidIfDescr, name)
				if err != nil {
					return nil, err
				}
				ifInfo, exist := ifInfoMapping[index]
				if !exist {
					ifInfo = &IfInfo{
						Index: index,
					}
				}
				n, ok := value.([]uint8)
				if !ok {
					return nil, errors.New("value type error")
				}
				ifInfo.Name = string(n)
				ifInfoMapping[index] = ifInfo
			case c.OidBelongTo(name, oidIfAlias): // 端口别名数据
				index, err := c.GetIndex(oidIfAlias, name)
				if err != nil {
					return nil, err
				}
				ifInfo, exist := ifInfoMapping[index]
				if !exist {
					ifInfo = &IfInfo{
						Index: index,
					}
				}
				tagBS, ok := value.([]uint8)
				if !ok {
					return nil, errors.New("value type error")
				}
				tagStr := string(tagBS)
				ifInfo.Alias = tagStr
				tagStr = strings.Trim(tagStr, " ")
				tags := strings.Split(tagStr, " ")
				for _, tag := range tags {
					tag = strings.Trim(tag, " ")
					if len(tag) > 0 {
						ifInfo.Tags = append(ifInfo.Tags, tag)
					}
					if tag == uplinkTag {
						ifInfo.IsUplink = true
					}
					if tag == downlinkTag {
						ifInfo.IsDownlink = true
					}

				}
				// 限速端口
				if len(ifInfo.Tags) >= 1 {
					if strings.HasPrefix(ifInfo.Tags[0], speedLimitTagPrefix) {
						ifInfo.IsSpeedLimit = true
					}
				}
				// 限速业务名称
				if len(ifInfo.Tags) >= 2 {
					ifInfo.PeakBusiness = ifInfo.Tags[1]
				}
				ifInfoMapping[index] = ifInfo
			}
		}
	}
	//}
	var ifInfos = make([]*IfInfo, 0, len(ifInfoMapping))
	for _, ifInfo := range ifInfoMapping {
		ifInfos = append(ifInfos, ifInfo)
	}
	return ifInfos, nil
}
func (c *SwitchCollector) GetIndex(prefix, name string) (int64, error) {
	s := strings.ReplaceAll(name, prefix, "")
	idStr := strings.ReplaceAll(s, ".", "")
	idStr = strings.ReplaceAll(idStr, " ", "")
	index, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, errors.WithStack(err)
	}
	return index, nil
}

func (c *SwitchCollector) CollectIfTrafficCounter() ([]*IfTrafficCounter, error) {
	var trafficMap = make(map[int64]*IfTrafficCounter)
	snmpResults, err := c.BulkGetOidsByPage([]string{oidIfIn, oidIfOut}, 50)
	//snmpResults, err := c.GetBulkWithRetry([]string{oidIfIn, oidIfOut}, 70)
	if err != nil {
		return nil, err
	}
	for _, snmpResult := range snmpResults {
		for _, dataUnit := range snmpResult.Results {
			name := dataUnit.Name
			value := dataUnit.Value
			switch {
			case c.OidBelongTo(name, oidIfIn):
				index, err := c.GetIndex(oidIfIn, name)
				if err != nil {
					return nil, err
				}
				if trafficCounter, exist := trafficMap[index]; exist {
					trafficCounter.InByte = gosnmp.ToBigInt(value).Uint64()
					trafficCounter.InTime = snmpResult.FinishTime
				} else {
					trafficMap[index] = &IfTrafficCounter{
						Index:   index,
						InByte:  gosnmp.ToBigInt(value).Uint64(),
						InTime:  snmpResult.FinishTime,
						OutByte: 0,
					}
				}
			case c.OidBelongTo(name, oidIfOut):
				index, err := c.GetIndex(oidIfOut, name)
				if err != nil {
					return nil, err
				}
				if trafficCounter, exist := trafficMap[index]; exist {
					trafficCounter.OutByte = gosnmp.ToBigInt(value).Uint64()
					trafficCounter.OutTime = snmpResult.FinishTime
				} else {
					trafficMap[index] = &IfTrafficCounter{
						Index:   index,
						InByte:  0,
						OutByte: gosnmp.ToBigInt(value).Uint64(),
						OutTime: snmpResult.FinishTime,
					}
				}
			}
		}
	}
	var trafficCounters = make([]*IfTrafficCounter, 0, len(trafficMap))
	for _, trafficCounter := range trafficMap {
		trafficCounters = append(trafficCounters, trafficCounter)
	}
	return trafficCounters, nil
}

func (c *SwitchCollector) CollectIfTraffic(interval time.Duration, times int) ([]*IfTrafficSpeed, error) {
	if times < 3 {
		return nil, errors.New("times must be at least 3 to calculate speed")
	}
	if interval < time.Second {
		return nil, errors.New("interval must be at least 1 second")
	}
	// 按时间点存储采样数据
	samples := make([][]*IfTrafficCounter, times)
	for i := 0; i < times; i++ {
		counters, err := c.CollectIfTrafficCounter()
		if err != nil {
			return nil, err
		}
		samples[i] = counters // 每次采样的数据独立存储
		if i < times-1 {
			time.Sleep(interval)
		}
	}
	trafficCounterMap := make(map[int64][]*IfTrafficCounter)
	for _, sample := range samples { // 遍历每次采样
		for _, counter := range sample { // 遍历采样中的每个接口
			trafficCounterMap[counter.Index] = append(trafficCounterMap[counter.Index], counter)
		}
	}
	return CalculateAverageSpeed(trafficCounterMap), nil
}
func CalculateAverageSpeed(trafficCounterMap map[int64][]*IfTrafficCounter) []*IfTrafficSpeed {
	var trafficSpeeds = make([]*IfTrafficSpeed, 0, len(trafficCounterMap))
	for _, trafficCounters := range trafficCounterMap {
		// 计算IN
		// 按in时间递增排序
		sort.Slice(trafficCounters, func(i, j int) bool {
			return trafficCounters[i].InTime.Before(trafficCounters[j].InTime)
		})
		var inSpeeds = make([]float64, 0, len(trafficCounters)-1)
		for i := 1; i < len(trafficCounters); i++ {
			lastPoint := trafficCounters[i-1]
			currentPoint := trafficCounters[i]
			timeDiff := float64(currentPoint.InTime.UnixMilli()-lastPoint.InTime.UnixMilli()) / 1000
			// 计算总数据差值
			inDiff := currentPoint.InByte - lastPoint.InByte
			if timeDiff <= 0 {
				continue
			}
			// 防止计数器回绕和异常值（检查 bytes）
			if inDiff >= 0 && inDiff <= ReasonableByteDiffLimit {
				inSpeed := (float64(inDiff) / timeDiff) * 8 // 转换为 bps
				inSpeeds = append(inSpeeds, inSpeed)
			}
		}
		// 计算OUT
		// 按out时间递增排序
		sort.Slice(trafficCounters, func(i, j int) bool {
			return trafficCounters[i].OutTime.Before(trafficCounters[j].OutTime)
		})
		var outSpeeds = make([]float64, 0, len(trafficCounters)-1)
		for i := 1; i < len(trafficCounters); i++ {
			lastPoint := trafficCounters[i-1]
			currentPoint := trafficCounters[i]
			timeDiff := float64(currentPoint.OutTime.UnixMilli()-lastPoint.OutTime.UnixMilli()) / 1000
			// 计算总数据差值
			outDiff := currentPoint.OutByte - lastPoint.OutByte
			if timeDiff <= 0 {
				continue
			}
			if outDiff >= 0 && outDiff <= ReasonableByteDiffLimit {
				outSpeed := (float64(outDiff) / timeDiff) * 8 // 转换为 bps
				outSpeeds = append(outSpeeds, outSpeed)
			}
		}

		var inSpeedAvg, outSpeedAvg float64
		if len(inSpeeds) > 0 {
			inSpeedAvg = Avg(inSpeeds)
		}
		if len(outSpeeds) > 0 {
			outSpeedAvg = Avg(outSpeeds)
		}
		trafficSpeeds = append(trafficSpeeds, &IfTrafficSpeed{
			Index:  trafficCounters[0].Index,
			Inbps:  inSpeedAvg,
			Outbps: outSpeedAvg,
			T:      trafficCounters[len(trafficCounters)-1].OutTime, // 具体实时就区出口的最后一条
		})
	}
	return trafficSpeeds
}
func Avg(nums []float64) float64 {
	if len(nums) == 0 {
		return 0
	}
	var sum float64
	for _, num := range nums {
		sum += num
	}
	return sum / float64(len(nums))
}

func (c *SwitchCollector) GetSwitchIfFullInfo(interval time.Duration, times int) ([]*IfInfo, error) {
	ifBaseInfo, err := c.GetSwitchIfBaseInfo()
	if err != nil {
		return nil, err
	}
	if len(ifBaseInfo) == 0 {
		return ifBaseInfo, nil // 如果没有接口，直接返回
	}
	// 获取的数据量需要乘以BulkWalkRedundancyFactor，防止获取数据不全
	ifTrafficSpeeds, err := c.CollectIfTraffic(interval, times)
	if err != nil {
		return nil, err
	}
	ifTrafficMap := make(map[int64]*IfTrafficSpeed)
	for _, ifTrafficSpeed := range ifTrafficSpeeds {
		ifTrafficMap[ifTrafficSpeed.Index] = ifTrafficSpeed
	}
	for _, ifInfo := range ifBaseInfo {
		if ifTrafficSpeed, exist := ifTrafficMap[ifInfo.Index]; exist {
			ifInfo.Speed = ifTrafficSpeed
		}
	}
	return ifBaseInfo, nil
}

func (c *SwitchCollector) GetOids(oidPrefixs []string, ifIndexs []int64) []string {
	var oids = make([]string, 0)
	for _, oid := range oidPrefixs {
		for _, index := range ifIndexs {
			oids = append(oids, fmt.Sprintf("%s.%d", oid, index))
		}
	}
	return oids
}
func (c *SwitchCollector) UseGetBulk(oids []string, maxRepetitions uint32) ([]*SnmpResult, error) {
	var startTime = time.Now()
	result, err := c.GetBulk(oids, 0, maxRepetitions)
	if err != nil {
		return nil, errors.Wrap(err, "get bulk failed")
	}
	var snmpResTime = time.Now()
	var snmpResult = &SnmpResult{
		Results:    result.Variables,
		FinishTime: snmpResTime,
		StartTime:  startTime,
	}
	return []*SnmpResult{snmpResult}, nil
}

func (c *SwitchCollector) UseGetBulkOneByOne(oids []string, maxRepetitions uint32) ([]*SnmpResult, error) {
	var results = make([]*SnmpResult, 0, len(oids))
	for _, oid := range oids {
		result, err := c.UseGetBulk([]string{oid}, maxRepetitions)
		if err != nil {
			return nil, err
		}
		if len(result) == 0 {
			continue
		}
		results = append(results, result...)
	}
	return results, nil
}

func (c *SwitchCollector) BulkGetOidsByPage(oids []string, pageSize uint32) ([]*SnmpResult, error) {
	var results = make([]*SnmpResult, 0)
	for _, oid := range oids {
		result, err := c.BulkGetOidByPage(oid, pageSize)
		if err != nil {
			return nil, err
		}
		results = append(results, result...)
	}
	return results, nil
}

func (c *SwitchCollector) BulkGetOidByPage(oid string, pageSize uint32) ([]*SnmpResult, error) {
	var (
		results = make([]*SnmpResult, 0)
	)
	lastOid := oid
	for {
		startAt := time.Now()
		result, err := c.GetBulk([]string{lastOid}, 0, pageSize)
		if err != nil {
			return nil, errors.Wrap(err, "get bulk failed")
		}
		finishedAt := time.Now()
		if len(result.Variables) == 0 {
			break
		}
		snmpPDUS := make([]gosnmp.SnmpPDU, 0, len(result.Variables))
		for _, pdu := range result.Variables {
			// 当前这个name 属于oid的
			if c.OidBelongTo(pdu.Name, oid) {
				snmpPDUS = append(snmpPDUS, pdu)
			} else {
				// 没有多余的数据了
				break
			}
		}
		r := &SnmpResult{
			Results:    snmpPDUS,
			StartTime:  startAt,
			FinishTime: finishedAt,
		}
		results = append(results, r)
		if len(snmpPDUS) < int(pageSize) {
			break
		}
		// 记录最后一个id，下次查询的时候带过去
		lastOid = result.Variables[len(result.Variables)-1].Name
	}
	return results, nil
}

func (c *SwitchCollector) OidBelongTo(oidName, parentOid string) bool {
	return strings.HasPrefix(oidName, parentOid) && !strings.Contains(strings.ReplaceAll(oidName, parentOid+".", ""), ".")
}
func (c *SwitchCollector) GetBulkWithRetry(oids []string, maxRepetitions uint32) ([]*SnmpResult, error) {
	result, err := c.UseGetBulkOneByOne(oids, maxRepetitions)
	if err == nil {
		return result, nil
	}
	result, err = c.UseGetBulk(oids, maxRepetitions)
	if err != nil {
		return nil, err
	}
	return result, nil
}
func (c *SwitchCollector) UseGet(oids []string, indexs []int64) ([]*SnmpResult, error) {
	var allOids = make([]string, 0, len(oids)*len(indexs))
	for _, index := range indexs {
		for _, oid := range oids {
			allOids = append(allOids, oid+"."+strconv.FormatInt(index, 10))
		}
	}
	startAt := time.Now()
	result, err := c.Get(allOids)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	finishedAt := time.Now()
	var snmpResult = &SnmpResult{
		Results:    result.Variables,
		StartTime:  startAt,
		FinishTime: finishedAt,
	}
	return []*SnmpResult{snmpResult}, nil
}

// BulkWalkAllOids 注意这个方法查询流量数据有缓存
func (c *SwitchCollector) BulkWalkAllOids(oids []string) ([]*SnmpResult, error) {
	var results = make([]*SnmpResult, 0)
	for _, oid := range oids {
		startAt := time.Now()
		result, err := c.BulkWalkAll(oid)
		if err != nil {
			return nil, err
		}
		finishedAt := time.Now()
		results = append(results, &SnmpResult{
			Results:    result,
			StartTime:  startAt,
			FinishTime: finishedAt,
		})
	}
	return results, nil
}
