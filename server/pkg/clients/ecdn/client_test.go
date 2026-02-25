package ecdn

import (
	"fmt"
	"testing"
)

var (
	//key = os.Getenv("ECDN_KEY")
	key = "PJzsUvRtSUOa3JrPpogPpUBOeakUX7vBLyN1LiVxfaGE1kYebtUw93EigqPOhJID"
)

func TestEcdnClient_GetServers(t *testing.T) {
	c, _ := NewEcdnClient("https://cdn.cjjd10.com", key)
	//servers, err := c.GetServerByBusiness("LE")
	servers, err := c.GetServers([]string{}, "", "", "")
	if err != nil {
		panic(err)
	}
	fmt.Println(servers)
}

func TestEcdnClient_DifIsp(t *testing.T) {
	c, _ := NewEcdnClient("https://cdn.cjjd10.com", key)
	req := &DifIspReq{
		Carrier:    1,
		FrankIDS:   "AHALY1L00MAED44BF",
		Note:       "测试",
		Phone:      "15270822181",
		Provincial: "安徽",
		Remind:     true,
	}
	data, err := c.DifIsp(req)
	if err != nil {
		panic(err)
	}
	fmt.Println(data)
}

func TestEcdnClient_GetServerByBusinesses(t *testing.T) {
	c, _ := NewEcdnClient("https://cdn.cjjd10.com", key)
	servers, err := c.GetServerByBusinesses([]string{"AACDN", "QACDN", "LE"})
	if err != nil {
		panic(err)
	}
	fmt.Println(servers)
}
