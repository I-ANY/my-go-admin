package la_traffic

type CommonResponse[T any] struct {
	Code int64  `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}

type DeviceTrafficItem struct {
	Timestamp     string `json:"timestamp"`
	Downloadspeed int64  `json:"downloadspeed"`
	Uploadspeed   int64  `json:"uploadspeed"`
}
