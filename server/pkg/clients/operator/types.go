package operator

type RemoteExecResponse struct {
	Code     int    `json:"code"`
	Hostname string `json:"hostname"`
	Result   string `json:"result"`
}
