package flushduty

type EventReqPayload struct {
	TitleRule   *string           `json:"title_rule"`
	EventStatus *string           `json:"event_status"`
	AlertKey    *string           `json:"alert_key"`
	Description *string           `json:"description"`
	Labels      map[string]string `json:"labels"`
	Images      []*EventReqImage  `json:"images"`
}
type EventReqImage struct {
	Alt  *string `json:"alt"`
	Src  *string `json:"src"`
	Href *string `json:"href"`
}
type EventResp struct {
	RequestId *string         `json:"request_id"`
	Error     *EventRespError `json:"error"`
	Data      *EventRespData  `json:"data"`
}
type EventRespError struct {
	Code    *int    `json:"code"`
	Message *string `json:"message"`
}
type EventRespData struct {
	AlertKey *string `json:"alert_key"`
}
