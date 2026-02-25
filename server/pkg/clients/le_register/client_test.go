package le_register

import (
	"log"
	"testing"
)

func TestAuth(t *testing.T) {

	client := NewClient("https://node.71edge.com", "mingfuyunmi", "7z4&oPaoh@!ih7G6Zv~B", nil)
	if err := client.Auth(); err != nil {
		panic(err)
	}

	sms := []string{
		"28b960d5a9b7c89f5ee0ad4da37c48cc973deeadd1c287c28866883c93dae54268b2d7ee",
		"8de5c15c09affc320bb9e863857f9c41c5aa8f59fbf31b74e28092d01f88a32c68b2d7ee",
	}
	//if err := client.Register(sms); err != nil {
	//	log.Println("register error:", err)
	//	panic(err)
	//}

	result, registerText, err := client.Query(sms)
	if err != nil {
		log.Println("query error:", err)
		panic(err)
	}
	log.Println(result, registerText)
}
