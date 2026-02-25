package jobs

import "biz-auto-api/internal/apps/cronjob/jobs/crontab"

func RegisterJobs() {
	crontab.JobList = map[string]crontab.JobExec{
		"ExampleJob":                 NewExampleJob(),
		"SyncZpTrafficDetailJob":     NewSyncZpTrafficDetailJob(),
		"SyncStarPortalUser":         NewSyncStarPortalUserJob(),
		"CollectTingContainerInfo":   NewCollectTingContainerInfoJob(),
		"SyncPriceBizJob":            NewSyncPriceBizJob(),
		"SyncEcdnServerJob":          NewSyncEcdnServerJob(),
		"SyncKDeliveryInfoJob":       NewSyncKDeliveryInfoJob(),
		"SyncKToDeliveryServerInfo":  NewSyncKToDeliveryServerInfo(),
		"NotifyBizJob":               NewNotifyBizJob(),
		"SyncBilibiliBillJob":        NewSyncBilibiliBillJob(),
		"SaGUIDCheckJob":             NewSaGUIDCheckJob(),
		"SaGuidAlertJob":             NewSaGUIDAlertJob(),
		"SaUtilizationAlertJob":      NewSaUtilizationAlertJob(),
		"SyncBilibiliBillSummaryJob": NewSyncBilibiliBillSunmaryJob(),
		"SyncDailyPeakTrafficJob":    NewSyncDailyPeakTrafficJob(),
		"SyncBusinessTypeJob":        NewSyncBusinessTypeJob(),
		"SyncEcdnDBServerJob":        NewSyncEcdnDBServerJob(),
		"SyncEcdnTrafficJob":         NewSyncEcdnTrafficJob(),
		"CalcTUSTrafficJob":          NewCalcTUSTrafficJob(),
		"SaResourceGapAlertJob":      NewSaResourceGapAlertJob(),
		"LeRegisterJob":              NewLeRegisterJob(),
		"SyncZPDeviceJob":            NewSyncZPDeviceJob(),
		"SyncZPTrafficJob":           NewSyncZPTrafficJob(),
		"SyncZPTrafficAddJob":        NewSyncZPTrafficAddJob(),
		"ZPUnusualDeviceAlertJob":    NewZPUnusualDeviceAlertJob(),
		"ZpCookieUpdateAlertJob":     NewZpCookieUpdateAlertJob(),
		"ZP95AutoPublishJob":         NewZP95AutoPublishJob(),
		"DeviceDataBackupJob":        NewDataBackupJob(),
		"SyncEcdnUTSOffineServerJob": NewSyncEcdnUTSOffineServerJob(),
		"ZapDeliveryDeviceScanJob":   NewZapDeliveryDeviceScanJob(),
		"ZapCookieUpdateAlertJob":    NewZapCookieUpdateAlertJob(),
		"TencentDemandGapJob":        NewTencentDemandGapJob(),
		"DetectTencentDemandJob":     NewDetectTencentDemandJob(),
		"ZapExportDeviceJob":         NewZapExportDeviceJob(),
		"ZapISPDiffAlertJob":         NewZAPISPDiffAlertJob(),
		"ZPBwStatsJob":               NewZPBwStatsJob(),
	}
	crontab.JobList["FixHddPartitionDataJob"] = NewFixHddPartitionDataJob()
	crontab.JobList["CheckInspectStatusJob"] = NewCheckInspectStatusJob()
	crontab.JobList["SyncLaServerToCacheJob"] = NewSyncLaServerToCacheJob()
	crontab.JobList["ClearHistoryData"] = NewClearHistoryDataJob()
	crontab.JobList["CalLaSLAJob"] = NewCalLaSLAJob()
	crontab.JobList["RoomSpeedLimit"] = NewRoomSpeedLimitJob()
	crontab.JobList["PreprocessKHDDDataJob"] = NewPreprocessKHDDDataJob()
	crontab.JobList["CalKHddDifferenceJob"] = NewCalKHddDifferenceJob()
	crontab.JobList["NewSyncNodeCurRatioJob"] = NewSyncNodeCurRatioJob()
	crontab.JobList["SyncPeakRoomJob"] = NewSyncPeakRoomJob()
	crontab.JobList["CalPeakRoom95"] = NewCalPeakRoom95()
	crontab.JobList["ServerDscpLockJob"] = NewServerDscpLockJob()
	crontab.JobList["Fix95RoomTrafficData"] = NewFix95RoomTrafficData()
	crontab.JobList["CachePeakRoomSwitch"] = NewCachePeakRoomSwitchJob()
	crontab.JobList["BizPrioritySummary"] = NewBizPrioritySummaryJob()
	crontab.JobList["SyncLaDeviceJob"] = NewSyncLaDeviceJob()
	crontab.JobList["SyncLaTrafficJob"] = NewSyncLaTrafficJob()
	crontab.JobList["InspectServer"] = NewInspectServerJob()
	crontab.JobList["CheckBTraffic"] = NewCheckBTrafficJob()
	crontab.JobList["SyncKTrafficDataToCKJob"] = NewSyncKTrafficDataToCKJob()
	crontab.JobList["NetworkSpeedLimitEnqueueJob"] = NewNetworkSpeedLimitEnqueueJob()
	crontab.JobList["OptimizeCKTableJob"] = NewOptimizeCKTableJob()
	crontab.JobList["NewCalNodeRedundancy"] = NewCalNodeRedundancyJob()
	crontab.JobList["SyncNodeJob"] = NewSyncNodeJob()
	crontab.JobList["CompositeInspect"] = NewCalculateJob()
	crontab.JobList["CompositeInspectV2"] = NewCalculateJobV2()
	crontab.JobList["SyncBiliHostnameInfoJob"] = NewSyncBilibiliHostnameInfoJob()
	crontab.JobList["BizOperationInspect"] = NewCalcCompBizJob()
	crontab.JobList["SyncNodeConfigJob"] = NewSyncNodeConfigJob()
	crontab.JobList["SyncAssessmentRuleJob"] = NewSyncAssessmentRuleJob()
	crontab.JobList["SyncBusinessOverProvisioningJob"] = NewBusinessOverProvisioningJob()
	crontab.JobList["OpsTaskScannerJob"] = NewOpsTaskScannerJob()
}
