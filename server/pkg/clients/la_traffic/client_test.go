package la_traffic

import (
	"fmt"
	"testing"
)

func TestClient_GetDeviceTraffic(t *testing.T) {
	client := NewClient("https://dcache.iqiyi.com", "", "")
	trafficData, err := client.GetDeviceTraffic("06AF6EB2532D2D46", 1749430668, 1749437867)
	if err != nil {
		panic(err)
	}
	fmt.Println(trafficData)
}
