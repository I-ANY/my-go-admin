package casbin

import (
	db "biz-auto-api/pkg/db"
	"biz-auto-api/pkg/logger"
	"fmt"
	"os"
	"testing"
)

var (
	dbHost, _     = os.LookupEnv("DB_HOST")
	dbUsername, _ = os.LookupEnv("DB_USERNAME")
	dbPassword, _ = os.LookupEnv("DB_PASSWORD")
	dbName, _     = os.LookupEnv("DB_NAME")
)

func Test_RBAC(t *testing.T) {
	logger.Setup("error")
	db.Setup(dbHost, dbName, dbUsername, dbPassword, 3306, 10, 10)
	Setup(db.GetDB())
	e := GetEnforcer()
	policy, err := e.Enforce("1", "/api/v1/system/userPermCode", "get")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(policy)
}
