package flushduty

import (
	"biz-auto-api/pkg/tools"
	"encoding/json"
	"testing"
)

func TestEventClient(t *testing.T) {
	client := NewEventClient("https://api.flashcat.cloud", "be1e212fb9f956b2a7d1d9152bb4e73a906")
	payload := &EventReqPayload{
		TitleRule:   tools.ToPointer("Test1"),
		EventStatus: tools.ToPointer(EventStatus_Critical),
		//AlertKey:    "test",
		Description: tools.ToPointer("testtttttttttttttttttttttttttttttttttttttttttttttttttttt\nsss\nsss\aaaa"),
		Labels:      map[string]string{"test": "test", "check": "111122233311"},
	}
	res, err := client.SendAlert(payload)
	if err != nil {
		panic(err)
	}
	jsonBytes, _ := json.Marshal(res)
	t.Logf("response: %s", jsonBytes)
	//	KuwjddpFpwwxAcQdQzphue
	//err := client.RecoverAlert("KuwjddpFpwwxAcQdQzphue")
	//if err != nil {
	//	panic(err)
	//}

}
