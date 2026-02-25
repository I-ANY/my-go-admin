package dscp

import (
	"biz-auto-api/internal/models"
	"biz-auto-api/pkg/clients/operator"
	"biz-auto-api/pkg/logger"
	"biz-auto-api/pkg/tools"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/avast/retry-go/v5"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"sync"
	"time"
)

func IssueDscp2Servers(servers []*models.NetworkServerDscp, config *IssueDscpConfig, db *gorm.DB, userId int64, log logger.Logger) []*UpdatePriorityResultItem {
	var (
		wg       sync.WaitGroup
		results  = make([]*UpdatePriorityResultItem, 0)
		resultCh = make(chan *UpdatePriorityResultItem, len(servers))
	)
	wg.Add(len(servers))
	for _, server := range servers {
		go func(server *models.NetworkServerDscp) {
			defer wg.Done()
			resultCh <- IssueDscp2Server(server, config, log)
		}(server)
	}
	wg.Wait()
	close(resultCh)
	var modifyRecords = make([]*models.NetworkDscpModifyRecord, 0)
	for result := range resultCh {
		results = append(results, result)
		modifyRecord := &models.NetworkDscpModifyRecord{
			CreatedBy: &userId,
			Hostname:  result.Hostname,
			Payload:   result.Payload,
			ServerID:  result.ServerID,
			Status:    result.Status,
			Message:   result.Message,
		}
		modifyRecords = append(modifyRecords, modifyRecord)
	}
	err := db.CreateInBatches(modifyRecords, 500).Error
	if err != nil {
		log.Errorf("%+v", err)
	}
	return results
}
func IssueDscp2Server(server *models.NetworkServerDscp, config *IssueDscpConfig, log logger.Logger) (result *UpdatePriorityResultItem) {
	var (
		err error
	)
	result = &UpdatePriorityResultItem{}
	if server.Hostname == nil {
		result.Status = tools.ToPointer(int64(models.ModifyDscpStatusFail))
		result.Message = tools.ToPointer("设备名称为空")
		return result
	}
	result.Hostname = server.Hostname
	result.ServerID = server.ServerID
	payload, err := json.Marshal(config)
	if err != nil {
		err = errors.Wrap(err, "marshal data failed")
		result.Message = tools.ToPointer(err.Error())
		result.Status = tools.ToPointer(int64(models.ModifyDscpStatusFail))
		log.Errorf("%+v", err)
		return result
	}
	payloadStr := fmt.Sprintf("start %v", string(payload))
	result.Payload = tools.ToPointer(payloadStr)
	log.Infof("issue dscp to server %s, payload=%v", tools.ToValue(server.Hostname), payloadStr)
	argsEncode := base64.StdEncoding.EncodeToString([]byte(payloadStr))
	retrier := retry.New(
		retry.Attempts(3),
		retry.Delay(time.Second),
		retry.DelayType(retry.BackOffDelay),
		retry.MaxDelay(time.Second*5),
		retry.LastErrorOnly(true),
	)
	// 未部署的主机名，使用frp-dev
	if len(*server.Hostname) >= 8 && (*server.Hostname)[7:8] == "_" {
		err = retrier.Do(func() error {
			_, err = operator.ExecScriptWithDomain(operator.FrpAddress_Dev, "dscp_judge.py", *server.Hostname, "general", argsEncode, time.Second*20)
			return err
		})
	} else {
		err = retrier.Do(func() error {
			_, err = operator.ExecScript("dscp_judge.py", *server.Hostname, "general", argsEncode, time.Second*20)
			return err
		})
	}

	// 判断结果
	if err != nil {
		result.Message = tools.ToPointer(err.Error())
		result.Status = tools.ToPointer(int64(models.ModifyDscpStatusFail))
		log.Errorf("%+v", err)
		return result
	} else {
		result.Message = tools.ToPointer("操作成功")
		result.Status = tools.ToPointer(int64(models.ModifyDscpStatusSuccess))
		return result
	}
	//return result
}
