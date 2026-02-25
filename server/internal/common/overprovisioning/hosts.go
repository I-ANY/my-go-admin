package overprovisioning

import (
	"biz-auto-api/internal/models"

	"gorm.io/gorm"
)

type Hosts struct {
	prom   *Prometheus
	ecdnDB *gorm.DB
}

func NewHosts(prom *Prometheus, ecdnDB *gorm.DB) *Hosts {
	return &Hosts{
		prom:   prom,
		ecdnDB: ecdnDB,
	}
}

func (h *Hosts) GetBizHosts(bizName string) (map[string]*PromHost, error) {
	lines, err := h.prom.GetHostsByBiz(bizName)
	if err != nil {
		return nil, err
	}
	promHosts := make(map[string]*PromHost, len(lines.Data.Result))
	for _, line := range lines.Data.Result {
		promHosts[line.Metric.ID] = &PromHost{
			ID:              line.Metric.ID,
			Hostname:        line.Metric.Hostname,
			BwPlan:          line.Metric.BwPlan,
			BwCount:         line.Metric.BwCount,
			BwSingle:        line.Metric.BwSingle,
			CactiNotes:      line.Metric.CactiNotes,
			Day95:           line.Metric.Day95,
			Evening95:       line.Metric.Evening95,
			Owner:           line.Metric.Owner,
			Sn:              line.Metric.Sn,
			Interprovincial: line.Metric.Interprovincial,
			Kvm:             line.Metric.Kvm,
			Location:        line.Metric.Location,
			IP:              line.Metric.IP,
		}
	}
	return promHosts, nil
}

func (h *Hosts) GetPromHostID(biz string) ([]string, error) {
	hosts, err := h.prom.GetHostsByBiz(biz)
	if err != nil {
		return nil, err
	}
	var hostID []string
	for _, host := range hosts.Data.Result {
		hostID = append(hostID, host.Metric.ID)
	}
	return hostID, nil
}

func (h *Hosts) GetECdnHost(biz string) (map[string]*models.Hardware, error) {
	hostIDs, err := h.GetPromHostID(biz)
	if err != nil {
		return nil, err
	}
	result := make(map[string]*models.Hardware)

	queryDb := func(id []string) ([]*models.ServersMore, error) {
		var data []*models.ServersMore
		if err := h.ecdnDB.Model(&models.ServersMore{}).Where("id in ?", id).Find(&data).Error; err != nil {
			return nil, err
		}
		return data, nil
	}

	// 分批查询ecdn数据
	var queryID []string
	batchNum := 80

	for i, hostID := range hostIDs {
		queryID = append(queryID, hostID)
		// 当达到批量大小 或者 已到最后一个元素
		if len(queryID) >= batchNum || i == len(hostIDs)-1 {
			data, err := queryDb(queryID)
			if err != nil {
				queryID = queryID[:0]
				continue
			}

			for _, item := range data {
				result[item.ID] = &item.Hardware
			}
			queryID = queryID[:0]
		}
	}
	return result, nil
}
