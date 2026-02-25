package switch_collector

import (
	"fmt"
	"github.com/gosnmp/gosnmp"
	"os"
	"testing"
	"time"
)

var (
	privateKey = os.Getenv("PRIVATE_KEY")
)

//
//func Test_GetPortIds(t *testing.T) {
//	err := cli.BulkWalk(oidIfIndex, func(dataUnit gosnmp.SnmpPDU) error {
//		name := dataUnit.Name
//		value := dataUnit.Value
//		fmt.Println("name", name, "value", value)
//		return nil
//	})
//	if err != nil {
//		panic(err)
//	}
//}
//func Test_GetPortName(t *testing.T) {
//	err := cli.BulkWalk(oidIfName, func(dataUnit gosnmp.SnmpPDU) error {
//		name := dataUnit.Name
//		value := string(dataUnit.Value.([]uint8))
//		fmt.Println("name", name, "value", value)
//		return nil
//	})
//	if err != nil {
//		panic(err)
//	}
//}
//func Test_GetPortStatus(t *testing.T) {
//	err := cli.BulkWalk(oidIfOperStatus, func(dataUnit gosnmp.SnmpPDU) error {
//		name := dataUnit.Name
//		value := dataUnit.Value
//		fmt.Println("name", name, "value", value)
//		return nil
//	})
//	if err != nil {
//		panic(err)
//	}
//}
//func Test_GetPortAlias(t *testing.T) {
//	err := cli.BulkWalk(oidIfAlias, func(dataUnit gosnmp.SnmpPDU) error {
//		oid := dataUnit.Name
//		value := string(dataUnit.Value.([]uint8))
//		fmt.Println("oid", oid, "value", value)
//		return nil
//	})
//	if err != nil {
//		panic(err)
//	}
//}
//func Test_GetPortDesc(t *testing.T) {
//	err := cli.BulkWalk(oidIfDescr, func(dataUnit gosnmp.SnmpPDU) error {
//		oid := dataUnit.Name
//		value := string(dataUnit.Value.([]uint8))
//		fmt.Println("oid", oid, "value", value)
//		return nil
//	})
//	if err != nil {
//		panic(err)
//	}
//}
//
//func Test_GetPortIn(t *testing.T) {
//	err := cli.BulkWalk(oidIfIn, func(dataUnit gosnmp.SnmpPDU) error {
//		oid := dataUnit.Name
//		value := dataUnit.Value
//		fmt.Println("oid", oid, "value", value)
//		return nil
//	})
//	if err != nil {
//		panic(err)
//	}
//}
//
//func Test_GetPortOut(t *testing.T) {
//	err := cli.BulkWalk(oidIfOut, func(dataUnit gosnmp.SnmpPDU) error {
//		oid := dataUnit.Name
//		value := dataUnit.Value
//		fmt.Println("oid", oid, "value", value)
//		return nil
//	})
//	if err != nil {
//		panic(err)
//	}
//}
//func Test_GetAll(t *testing.T) {
//	result, err := cli.GetBulk([]string{oidIfName, oidIfAlias, oidIfOperStatus}, 0, 2)
//	if err != nil {
//		panic(err)
//	}
//	for _, dataUnit := range result.Variables {
//		oid := dataUnit.Name
//		value := dataUnit.Value
//		fmt.Println("oid", oid, "value", value)
//	}
//
//}

func TestSwitchCollector_GetSwitchIfInfo(t *testing.T) {
	sn := &gosnmp.GoSNMP{
		Target: "10.254.120.100",
		//Target:             "111.47.233.1",
		Port:               161,
		Transport:          "udp",
		Community:          "mf_stone",
		Version:            gosnmp.Version2c,
		Timeout:            time.Second * 5,
		Retries:            3,
		MaxOids:            gosnmp.MaxOids * 10,
		ExponentialTimeout: true,
	}
	collector, err := New(sn)
	if err != nil {
		panic(err)
	}
	err = collector.Connect()
	if err != nil {
		panic(err)
	}
	defer collector.Close()
	infos, err := collector.GetSwitchIfBaseInfo()
	if err != nil {
		panic(err)
	}
	for _, info := range infos {
		fmt.Printf("%#v\n", info)
	}
	fmt.Println(len(infos))
}
func TestSwitchCollector_GetIndex(t *testing.T) {
	sn := &gosnmp.GoSNMP{
		//Target:             "183.131.183.33",
		Target:             "111.47.233.1",
		Port:               161,
		Transport:          "udp",
		Community:          "mf_stone",
		Version:            gosnmp.Version2c,
		Timeout:            time.Second,
		Retries:            5,
		MaxOids:            gosnmp.MaxOids * 10,
		ExponentialTimeout: true,
	}
	collector, err := New(sn)
	if err != nil {
		panic(err)
	}
	err = collector.Connect()
	if err != nil {
		panic(err)
	}
	defer collector.Close()
	trafficItems, err := collector.CollectIfTraffic(time.Second*10, 3)
	if err != nil {
		panic(err)
	}
	for _, trafficItem := range trafficItems {
		fmt.Printf("%#v\n", trafficItem)
	}
}

func TestSwitchCollector_GetSwitchIfFullInfo(t *testing.T) {
	sn := &gosnmp.GoSNMP{
		//Target: "111.47.233.1",
		//Target: "183.131.183.33",
		//Target: "10.254.120.100", //测试机器
		Target: "10.10.195.2",
		//Target:    "10.238.75.2",
		Port:               161,
		Transport:          "udp",
		Community:          "mf_stone",
		Version:            gosnmp.Version2c,
		Timeout:            time.Duration(2) * time.Second,
		Retries:            3,
		MaxRepetitions:     10,
		MaxOids:            gosnmp.MaxOids * 20,
		ExponentialTimeout: true,
	}
	collector, err := New(sn)
	if err != nil {
		panic(err)
	}
	err = collector.Connect()
	if err != nil {
		panic(err)
	}
	infos, err := collector.GetSwitchIfFullInfo(time.Second*14, 3)
	if err != nil {
		panic(err)
	}
	defer collector.Close()
	for _, info := range infos {
		fmt.Printf("%#v, speed:%#v\n", info, info.Speed)
	}
	var uplinkSpeedOut float64 = 0
	var otherSpeedIn float64 = 0
	for _, info := range infos {
		if info.IsUplink && info.Speed != nil {
			uplinkSpeedOut += info.Speed.Outbps
		}
		if info.IsSpeedLimit && info.Speed != nil {
			otherSpeedIn += info.Speed.Inbps
		}
	}
	fmt.Println("total uplink speed: ", uplinkSpeedOut)
	fmt.Println("total other speed: ", otherSpeedIn)
}

func TestSwitchCollector_SSH(t *testing.T) {
	client, err := NewSwitchSSHClientWithPrivateKey("iaas", "10.254.120.100", 22, privateKey)
	if err != nil {
		panic(err)
	}
	defer client.Close()
	session, err := client.NewSession()
	if err != nil {
		panic(err)
	}
	defer session.Close()
	// 限速
	cmds := []string{
		"system-view",
		"qos car 17M cir 17000 kbps",
		"int 10GE1/0/48",
		"qos car inbound 17M",
		"commit",
		"quit",
		"quit",
		"save",
		"y",
	}
	// 解除限速
	//cmds := []string{
	//	"system-view",
	//	"int 10GE1/0/48",
	//	"undo qos car inbound",
	//	"commit",
	//	"quit",
	//	"quit",
	//	"save",
	//	"y",
	//}
	for _, c := range cmds {
		result, err := session.Exec(c, time.Second*10, false, false)
		if err != nil {
			fmt.Printf("%v", string(result))
			fmt.Printf("%+v", err)
			return
		}
		fmt.Print(string(result))
	}
}

func Test_QosTemplate(t *testing.T) {
	client, err := NewSwitchSSHClientWithPrivateKey("iaas", "10.254.120.100", 22, privateKey)
	if err != nil {
		panic(err)
	}
	defer client.Close()
	err = client.InitSession()
	if err != nil {
		panic(err)
	}
	template, err := client.GetQosTemplate()
	if err != nil {
		panic(err)
	}
	for _, t := range template {
		fmt.Printf("%#v\n", t)
	}
}
func Test_GetInterfaceInfo(t *testing.T) {
	client, err := NewSwitchSSHClientWithPrivateKey("iaas", "10.254.120.100", 22, privateKey)
	if err != nil {
		panic(err)
	}
	defer client.Close()
	err = client.InitSession()
	if err != nil {
		panic(err)
	}
	template, err := client.LimitIfInSpeed("10GE1/0/35", 1000)
	if err != nil {
		fmt.Println(string(template))
		panic(err)
	}
	fmt.Println(string(template))
}

func Test_GetInterfaceInfo1(t *testing.T) {
	client, err := NewSwitchSSHClientWithPrivateKey("iaas", "10.254.120.100", 22, privateKey)
	if err != nil {
		panic(err)
	}
	defer client.Close()
	session, err := client.NewSession()
	if err != nil {
		panic(err)
	}
	defer session.Close()
	template, err := client.UnlimitedInfInSpeed("10GE1/0/35")
	if err != nil {
		fmt.Println(string(template))
		panic(err)
	}
	fmt.Println(string(template))
}

func Test_GetInterfaceInfo2(t *testing.T) {
	decimal := FloorbpsToNearest100MbpsDecimal(-110)
	fmt.Println(decimal)
}

func Test_MustLimitIfInSpeed(t *testing.T) {
	client, err := NewSwitchSSHClientWithPrivateKey("iaas", "10.254.120.100", 22, privateKey)
	if err != nil {
		panic(err)
	}
	defer client.Close()
	err = client.InitSession()
	if err != nil {
		panic(err)
	}
	for {
		out, err := client.MustLimitIfInSpeed("10GE1/0/48", 10000000001)
		if err != nil {
			fmt.Printf("%+v", err)
			fmt.Println(string(out))
		}
		fmt.Println(string(out))
		out, err = client.MustUnlimitedIfInSpeed("10GE1/0/48")
		if err != nil {
			fmt.Printf("%+v", err)
			fmt.Println(string(out))
		}
		fmt.Println(string(out))
		time.Sleep(time.Second * 5)
	}
}

func Test_MustUnlimitedIfInSpeed(t *testing.T) {
	client, err := NewSwitchSSHClientWithPrivateKey("iaas", "10.254.120.100", 22, privateKey)
	if err != nil {
		panic(err)
	}
	defer client.Close()
	err = client.InitSession()
	if err != nil {
		panic(err)
	}
	for {
		c, err := client.GetIfConfigs()
		if err != nil {
			fmt.Printf("%+v", err)
		}
		fmt.Printf("%#v", c)
		time.Sleep(time.Second * 5)
	}
}

func Test_shllRun(t *testing.T) {
	fmt.Println(FloorbpsToNearest100MbpsDecimal(201 * 1000 * 1000))
}

func Test_test1(t *testing.T) {

	sn := &gosnmp.GoSNMP{
		//Target: "111.47.233.1",
		//Target: "183.131.183.33",
		//Target: "10.254.120.100", //测试机器
		Target: "10.254.120.100",
		//Target:    "10.238.75.2",
		Port:               161,
		Transport:          "udp",
		Community:          "mf_stone",
		Version:            gosnmp.Version2c,
		Timeout:            time.Duration(2) * time.Second,
		Retries:            3,
		MaxOids:            gosnmp.MaxOids * 10,
		ExponentialTimeout: true,
		MaxRepetitions:     10,
	}
	collector, err := New(sn)
	if err != nil {
		panic(err)
	}
	err = collector.Connect()
	if err != nil {
		panic(err)
	}
	business, err := collector.GetSwitchIfBusiness()
	if err != nil {
		panic(err)
	}
	fmt.Println(business)
}

func Test_test2(t *testing.T) {
	sn := &gosnmp.GoSNMP{
		//Target: "111.47.233.1",
		//Target: "183.131.183.33",
		Target: "10.238.75.2",
		//Target: "123.180.181.1",
		//Target: "10.254.120.100", //测试机器
		Port:               161,
		Transport:          "udp",
		Community:          "mf_stone",
		Version:            gosnmp.Version2c,
		Timeout:            time.Duration(5) * time.Second,
		Retries:            3,
		MaxOids:            gosnmp.MaxOids * 10,
		ExponentialTimeout: false,
		MaxRepetitions:     10,
	}
	collector, err := New(sn)
	if err != nil {
		panic(err)
	}
	err = collector.Connect()
	if err != nil {
		panic(err)
	}
	oids, err := collector.CollectIfTraffic(time.Second*1, 3)
	if err != nil {
		panic(err)
	}
	fmt.Println(len(oids))
}

func Test_GetBulkByPage(t *testing.T) {
	sn := &gosnmp.GoSNMP{
		//Target: "111.47.233.1",
		//Target: "183.131.183.33",
		Target: "10.16.0.3",
		//Target: "123.180.181.1",
		//Target: "10.254.120.100", //测试机器
		Port:               161,
		Transport:          "udp",
		Community:          "mf_stone",
		Version:            gosnmp.Version2c,
		Timeout:            time.Duration(5) * time.Second,
		Retries:            3,
		MaxOids:            gosnmp.MaxOids * 10,
		ExponentialTimeout: false,
		MaxRepetitions:     10,
	}
	collector, err := New(sn)
	if err != nil {
		panic(err)
	}
	err = collector.Connect()
	if err != nil {
		panic(err)
	}
	page, err := collector.BulkGetOidsByPage([]string{oidIfIn, oidIfOut}, 45)
	if err != nil {
		panic(err)
	}
	fmt.Println(page)
}
