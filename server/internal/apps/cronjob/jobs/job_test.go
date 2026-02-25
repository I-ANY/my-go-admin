package jobs

import (
	"biz-auto-api/internal/apps/cronjob/jobs/crontab"
	"biz-auto-api/pkg/clickhouse"
	"biz-auto-api/pkg/config"
	"biz-auto-api/pkg/config/types"
	"biz-auto-api/pkg/consts"
	db "biz-auto-api/pkg/db"
	"biz-auto-api/pkg/grpc/clients"
	"biz-auto-api/pkg/logger"
	"biz-auto-api/pkg/queue"
	"biz-auto-api/pkg/redis"
	"biz-auto-api/pkg/snowid"
	"github.com/sirupsen/logrus"
	"log"
	"testing"
	"time"
)

var (
	args *crontab.JobExecArgs
)

func Init() {
	config.SetupCronjobConfig("../../../../config/config-cronjob-dev.yaml")
	var c = config.CronjobConfig
	// 2.加载logger
	logger.Setup(c.Log.Level)
	// 3.初始化数据库连接
	db.Setup(
		c.Mysql.Host,
		c.Mysql.Database,
		c.Mysql.Username,
		c.Mysql.Password,
		c.Mysql.Port,
		c.Mysql.MaxIdleCon,
		c.Mysql.MaxOpenCon,
	)
	redis.Setup(
		c.Redis.Username,
		c.Redis.Password,
		c.Redis.Host,
		c.Redis.DB,
		c.Redis.Port,
		c.Redis.Timeout,
		c.Redis.PoolSize,
	)
	clickhouse.Setup(
		c.Clickhouse.Host,
		c.Clickhouse.Database,
		c.Clickhouse.Username,
		c.Clickhouse.Password,
		c.Clickhouse.Port,
		c.Clickhouse.MaxIdleCon,
		c.Clickhouse.MaxOpenCon,
		false,
	)
	snowid.Setup("test", redis.GetClient())
	args = &crontab.JobExecArgs{
		Args:   "",
		DB:     db.GetDB(),
		Log:    crontab.NewPersistLogger(logger.GetLogger().WithField("type", "test"), 0, db.GetDB()),
		Config: *c,
		Redis:  redis.GetClient(),
		CK:     clickhouse.GetCK(),
	}
	err := clients.SetupVpnGrpcClient(c.VpnClient.Address, time.Duration(c.VpnClient.ConnTimeoutSecond)*time.Second)
	InitSpeedLimitQueue(c, logger.GetLogger())
	if err != nil {
		log.Fatalf("%+v", err)
	}
}

func InitSpeedLimitQueue(c *types.CronjobConfig, log *logrus.Logger) {
	queueName := consts.QueneName_SpeedLimitTask
	queueLogger := queue.NewLogger(logger.GetLogger().WithField("type", "queue"))
	q, err := queue.NewQueue(c.Redis.Username,
		c.Redis.Password,
		c.Redis.Host,
		queueName,
		int64(c.Redis.Port),
		int64(c.Redis.DB),
		3*24*60, //结果保留3天
		nil,
		queueLogger,
	)
	if err != nil {
		log.Fatalf("%+v", err)
	}
	err = queue.Queue.AddQueue(queueName, q)
	if err != nil {
		log.Fatalf("%+v", err)
	}
}
func TestSyncStarPortalUserJob(t *testing.T) {
	Init()
	syncStarPortalUserJob := NewSyncStarPortalUserJob()
	err := syncStarPortalUserJob.Exec(args)
	if err != nil {
		log.Fatalf("%+v", err)
	}

}

func TestNewCollectTingContainerInfoJob(t *testing.T) {
	Init()
	j := NewCollectTingContainerInfoJob()
	err := j.Exec(args)
	if err != nil {
		log.Fatalf("%+v", err)
	}
}

func TestNewSyncEcdnServerJob(t *testing.T) {
	Init()
	args.Args = `{"useRedis": false}`
	j := NewSyncEcdnServerJob()
	err := j.Exec(args)
	if err != nil {
		log.Fatalf("%+v", err)
	}
}
func TestNewSyncEcdnBusinessJob(t *testing.T) {
	Init()
	j := NewSyncBusinessTypeJob()
	err := j.Exec(args)
	if err != nil {
		log.Fatalf("%+v", err)
	}
}
func TestNewSyncKDeliveryInfoJob(t *testing.T) {
	Init()
	j := NewSyncKDeliveryInfoJob()
	err := j.Exec(args)
	if err != nil {
		log.Fatalf("%+v", err)
	}
}
func TestNewNetworkSpeedLimitEnqueueJob(t *testing.T) {
	Init()
	j := NewNetworkSpeedLimitEnqueueJob()
	err := j.Exec(args)
	if err != nil {
		log.Fatalf("%+v", err)
	}
}
func TestNewBizPrioritySummaryJob(t *testing.T) {
	Init()
	j := NewBizPrioritySummaryJob()
	args.Args = `{"syncSuitableServer": true, "calBizPrioritySummary": true}`
	err := j.Exec(args)
	if err != nil {
		log.Fatalf("%+v", err)
	}
}
func TestNewSyncKToDeliveryServerInfo(t *testing.T) {
	Init()
	j := NewSyncKToDeliveryServerInfo()
	err := j.Exec(args)
	if err != nil {
		log.Fatalf("%+v", err)
	}
}
func TestNewServerDscpLockJob(t *testing.T) {
	Init()
	j := NewServerDscpLockJob()
	err := j.Exec(args)
	if err != nil {
		log.Fatalf("%+v", err)
	}
}
func TestNewBilibiliInfo(t *testing.T) {
	Init()
	j := NewSyncBilibiliBillJob()
	err := j.Exec(args)
	if err != nil {
		log.Fatalf("%+v", err)
	}
}

func TestPreprocessKHDDDataJob_Exec(t *testing.T) {
	Init()
	j := NewPreprocessKHDDDataJob()
	args.Args = `{"startDate": "2025-12-15", "endDate": "2025-12-23"}`
	err := j.Exec(args)
	if err != nil {
		log.Fatalf("%+v", err)
	}
}
func Test_RooomSpeedLimitJob_Exec(t *testing.T) {
	Init()
	j := NewRoomSpeedLimitJob()
	args.Args = `{
    "collectInterval": 16,
    "collectTimes": 3,
    "rpcClientTimeout": 15,
    "snmpTimeout":10,
    "roomSuffixes": ["0","1","2","3","4","5","6","7","8","9"],
    "unlimitDiffGBStep":15,
    "execErrorAlertIntervalMinute": 15,
    "integrationKey": "be1e212fb9f956b2a7d1d9152bb4e73a906"
}`
	err := j.Exec(args)
	if err != nil {
		log.Fatalf("%+v", err)
	}
}
func TestNewCalKHddDifferenceJob(t *testing.T) {
	Init()
	j := NewCalKHddDifferenceJob()
	args.Args = `{"startDate": "2025-07-16", "endDate": "2025-07-17"}`
	err := j.Exec(args)
	if err != nil {
		log.Fatalf("%+v", err)
	}
}
func TestNewBilibiliSumarayInfo(t *testing.T) {
	Init()
	j := NewSyncBilibiliBillSunmaryJob()
	err := j.Exec(args)
	if err != nil {
		log.Fatalf("%+v", err)
	}
}

func TestNewSyncBilibiliHostnameInfoJob(t *testing.T) {
	Init()
	j := NewSyncBilibiliHostnameInfoJob()
	err := j.Exec(args)
	if err != nil {
		log.Fatalf("%+v", err)
	}
}

func TestNewCheckBTrafficJob(t *testing.T) {
	Init()
	args.Args = `{"startTime": "2024-12-23 12:23:23", "endTime": "2024-12-23 19:19:23"}`
	args.Args = `{"beforeStartMinute":70, "beforeEndMinute":30 ,"CMCCSpeedDiffPercent":0.05, "CUCCSpeedDiffPercent":0.05,"CTCCSpeedDiffPercent":0.05}`
	j := NewCheckBTrafficJob()
	err := j.Exec(args)
	if err != nil {
		log.Fatalf("%+v", err)
	}
}

func TestNewSyncPeakRoomJob(t *testing.T) {
	Init()
	j := NewSyncPeakRoomJob()
	err := j.Exec(args)
	if err != nil {
		log.Fatalf("%+v", err)
	}
}
func TestNewCalPeakRoom95(t *testing.T) {
	Init()
	//"startTime": "2025-07-30 00:00:00",
	//	"endTime": "2025-07-31 23:59:59",
	args.Args = `
{
    "startTime": "2025-10-27 00:00:00",
    "endTime": "2025-10-27 16:20:00",
    "syncTraffic": true,
    "cal95": true,
    "datasource": "ck"
}
`

	//, "syncTrafficWithCache": true,
	//    "syncRoomId":[24160]
	//    "syncRoomId": [22785,4571,22562]
	j := NewCalPeakRoom95()
	err := j.Exec(args)
	if err != nil {
		log.Fatalf("%+v", err)
	}
}
func TestNewCalNodeRedundancyJob(t *testing.T) {
	Init()
	args.Args = `
{
    "business": [
        "ZP_LIVE_C1_M",
        "ZP_LIVE_C2_M",
        "ZP_LIVE_C3_M",
        "ZP_LIVE_C4_M",
        "ZP_LIVE_C5_M",
        "ZP_LIVE_D1_M"
    ],
    "startTime": "2025-12-16 00:00:00",
    "endTime": "2025-12-17 23:59:59"
}
`
	//"relativeTimeMinutes": 20,
	//	"batchTimeMinutes": 5
	//,
	//"startTime": "2025-12-09 00:00:00",
	//	"endTime": "2025-12-09 23:59:59"
	j := NewCalNodeRedundancyJob()
	err := j.Exec(args)
	if err != nil {
		log.Fatalf("%+v", err)
	}
}
func TestCachePeakRoomSwitch_Exec(t *testing.T) {
	Init()
	j := NewCachePeakRoomSwitchJob()
	err := j.Exec(args)
	if err != nil {
		log.Fatalf("%+v", err)
	}
}

func TestNewSyncNodeCurRatioJob(t *testing.T) {
	Init()
	j := NewSyncNodeCurRatioJob()
	err := j.Exec(args)
	if err != nil {
		log.Fatalf("%+v", err)
	}
}
func TestFix95RoomTrafficData_Exec(t *testing.T) {
	Init()
	j := NewFix95RoomTrafficData()
	//"startTime": "2025-08-07 11:00:00",
	//	"endTime": "2025-08-07 12:00:00"
	args.Args = `
{
	"relativeBeginMinute": 60,
	"relativeEndMinute": 30
}
`
	err := j.Exec(args)
	if err != nil {
		log.Fatalf("%+v", err)
	}
}

//	func TestNewTestAddKTrafficData(t *testing.T) {
//		Init()
//		j := NewAddTestKTrafficDataJob()
//		err := j.Exec(args)
//		if err != nil {
//			log.Fatalf("%+v", err)
//		}
//	}
func TestClearHistoryDataJob_Exec(t *testing.T) {
	Init()
	j := NewClearHistoryDataJob()
	args.Args = `
{
	"enableDeleteKHddData": false,
	"enableDeleteInspectData": false,
	"inspectRetainDays": 1000,
	"enableDeleteOperaLogData": false,
    "operaLogRetainDays":15,
    "enableDeleteJobData":false,
    "jobRetainDays": 30
}
`
	err := j.Exec(args)
	if err != nil {
		log.Fatalf("%+v", err)
	}
}
func TestNewSyncKTrafficDataToCKJob(t *testing.T) {
	Init()
	args.Args = `{"startId":0,"endId":7,"batchSize":2}`
	j := NewSyncKTrafficDataToCKJob()
	err := j.Exec(args)
	if err != nil {
		log.Fatalf("%+v", err)
	}
}

func TestNewOptimizeCKTableJob(t *testing.T) {
	Init()
	//args.Args = `{"kTrafficPartitions":["20250509","20250510","20250508"]}`
	j := NewOptimizeCKTableJob()
	err := j.Exec(args)
	if err != nil {
		log.Fatalf("%+v", err)
	}
}

func TestNewCheckInspectStatusJob(t *testing.T) {
	Init()
	args.Args = `{ "timeoutSecond": 600 }`
	j := NewCheckInspectStatusJob()
	err := j.Exec(args)
	if err != nil {
		log.Fatalf("%+v", err)
	}
}

func TestNewInspectServerJob(t *testing.T) {
	Init()
	args.Args = `[
  {
    "categories":["KP2"],
    "subcategories":["KPC_M"],
    "idcType":1,
    "excludeServerTags":["离线","裁撤","关停","下架","异常","让量","测试","瓶颈","故障"]
  },
  {
    "categories":["KPC"],
    "subcategories":["KPC_M"],
    "idcType":1,
    "excludeServerTags":["离线","裁撤","关停","下架","异常","让量","测试","瓶颈","故障"]
  }
]`
	args.Args = `[
  {
   "excludeServerTags":["离线","裁撤","关停","下架","异常","让量","测试","瓶颈","故障"],
    "categories":["KPC"],
    "idcType":2
  }
]`
	j := NewInspectServerJob()
	err := j.Exec(args)
	if err != nil {
		log.Fatalf("%+v", err)
	}
}

func TestNewSyncLaDeviceJob(t *testing.T) {
	Init()
	j := NewSyncLaDeviceJob()
	err := j.Exec(args)
	if err != nil {
		log.Fatalf("%+v", err)
	}
}
func TestCalcCompBizJob2(t *testing.T) {
	Init()
	job := NewCalcCompBizJob()
	err := job.Exec(args)
	if err != nil {
		log.Fatalf("%+v", err)
	}
}

func TestSyncLaServerToCacheJob(t *testing.T) {
	Init()
	j := NewSyncLaServerToCacheJob()
	err := j.Exec(args)
	if err != nil {
		log.Fatalf("%+v", err)
	}
}
func TestNewCalLaSLAJob(t *testing.T) {
	Init()
	j := NewCalLaSLAJob()
	args.Args = `
{
"startTime": "2025-10-14 10:00:00",
"endTime": "2025-10-16 11:00:00",
"enableAlert": true
}
`
	err := j.Exec(args)
	if err != nil {
		log.Fatalf("%+v", err)
	}
}

func TestNewFixHddPartitionDataJob(t *testing.T) {
	Init()
	args.Args = `
{
"startTime": "2025-06-01 00:00:00",
"endTime": "2025-06-30 23:59:59"
}
`
	j := NewFixHddPartitionDataJob()
	err := j.Exec(args)
	if err != nil {
		log.Fatalf("%+v", err)
	}
}
func TestNewSyncLaTrafficJob(t *testing.T) {
	Init()
	//args.Args = `{"startTime": "2025-06-25 00:00:00", "endTime": "2025-06-25 23:59:59"}`
	args.Args = `{"startTime": "2025-09-01 14:00:00", "endTime": "2025-09-01 14:10:00", "onlyEmpty": false}`
	//args.Args = `{"beforeStartMinute": 120, "beforeEndMinute": 60, "onlyEmpty": true}`
	j := NewSyncLaTrafficJob()
	err := j.Exec(args)
	if err != nil {
		log.Fatalf("%+v", err)
	}
}

func TestNewSyncEcdnDBServerJob(t *testing.T) {
	Init()
	args.Args = `{"startTime": "2025-08-14 00:00:00", "endTime": "2025-08-18 23:59:59"}`
	//args.Args = `{"startTime": "", "endTime": ""}`
	j := NewSyncEcdnDBServerJob()
	err := j.Exec(args)
	if err != nil {
		log.Fatalf("%+v", err)
	}
}
func TestNewCalcTUSTrafficJob(t *testing.T) {
	Init()
	//args.Args = `{"startTime": "", "endTime": ""}`
	j := NewCalcTUSTrafficJob()
	err := j.Exec(args)
	if err != nil {
		log.Fatalf("%+v", err)
	}
}

func TestNewLeRegisterJob(t *testing.T) {
	Init()
	//args.Args = `{"startTime": "", "endTime": ""}`
	j := NewLeRegisterJob()
	err := j.Exec(args)
	if err != nil {
		log.Fatalf("%+v", err)
	}
}

func TestNewSyncEcdnUTSOffineServerJob(t *testing.T) {
	Init()
	args.Args = `{"startTime": "2025-09-23 15:44:47", "endTime": "2025-09-25 15:44:47"}`
	j := NewSyncEcdnUTSOffineServerJob()
	err := j.Exec(args)
	if err != nil {
		log.Fatalf("%+v", err)
	}
}

func TestNewBusinessOverProvisioningJob(t *testing.T) {
	Init()
	args.Args = `{"execType": "sync_business_over_provisioning_hostname"}`
	j := NewBusinessOverProvisioningJob()
	err := j.Exec(args)
	if err != nil {
		log.Fatalf("%+v", err)
	}
}
