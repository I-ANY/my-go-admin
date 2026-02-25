package flask_api

type DeliveryPayload struct {
	BizType  string        `json:"biz_type"`
	Provider string        `json:"provider"`
	DemandId string        `json:"demand_id"`
	Data     []*MacPayload `json:"data"`
	FrankID  string        `json:"frank_id"`
}

type MacPayload struct {
	Mac                  string `json:"mac"`
	DeliveryId           int64  `json:"delivery_id"`
	DeviceType           int64  `json:"device_type"`
	Province             string `json:"province"`
	City                 string `json:"city"`
	Isp                  string `json:"isp"`
	DownloadBw           int64  `json:"download_bw"`
	UploadBw             int64  `json:"upload_bw"`
	IsProvinceScheduling *int64 `json:"is_province_scheduling"`
	NetworkType          string `json:"network_type"`
	Ip                   string `json:"ip"`
	IsCoverDiffIsp       *int64 `json:"is_cover_diff_isp"`
	CoverDiffIsp         string `json:"cover_diff_isp"`
}

type DeliveryResp struct {
	Code int64  `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		TaskId int64 `json:"task_id"`
	} `json:"data"`
}
