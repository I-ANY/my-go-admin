package starPortal

import (
	"fmt"
	"log"
	"os"
	"testing"
)

var (
	clientId, _     = os.LookupEnv("CLIENT_ID")
	clientSecret, _ = os.LookupEnv("CLIENT_SECRET")
	apiToken, _     = os.LookupEnv("API_TOKEN")
)

func TestStarPortalClient_GetUsers(t *testing.T) {
	c, _ := NewStarPortalClient("https://star.xaidc.com", apiToken)
	users, err := c.GetUsers(false)
	if err != nil {
		log.Fatal(err)
	}
	for _, user := range users {
		fmt.Println(user.Name, user.DeptIDList)
	}
}
func TestClient(t *testing.T) {
	c, _ := NewStarPortalClient("https://star-portal.xaidc.com:8443", apiToken)

	users, err := c.GetUserToken(
		"66493ca0445345648515adfabb6e10d0",
		clientId,
		clientSecret,
		"http://localhost:5173/login?redirect=/dashboard",
	)
	if err != nil {
		log.Fatalf("%+v", err)
	}

	c2, _ := NewStarPortalClient("https://star-portal.xaidc.com:8443", users.AccessToken)
	userInfo, err := c2.GetUserInfo()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	fmt.Println(userInfo)

	fmt.Println(users)
}
