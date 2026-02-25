import { defHttp } from '@/utils/http/axios';

export enum Api {
  GetDeviceList = '/v1/business/k/device/list',
  GetDeviceMacsList = '/v1/business/k/device/macs/list',
  GetBusinessList = '/v1/business/k/device/business/list',
  GetDevTypeList = '/v1/business/k/device/deviceType/list',
  GetAreaList = '/v1/business/k/device/area/list',
  UpdateDeviceMac = '/v1/flask-api/proxy/tencent/mac/replace',
  EditDeviceMac = '/v1/flask-api/proxy/tencent/mac/edit',
  GetDeliveryInfoList = '/v1/business/k/deliveryInfo/list',
  UpdateDeliveryInfo = '/v1/business/k/deliveryInfo',
  BatchUpdateDeliveryInfo = '/v1/business/k/deliveryInfo/update/batch',
  GetDeliveryOptions = '/v1/business/k/delivery/options',
  GetMatchingDemandList = '/v1/business/k/matchingDemand/list',
  // GetDemandList = '/v1/business/k/demand/list',
  GetDemandList = '/v1/flask-api/proxy/tencent/demand/list',
  DemandOccupy = '/v1/flask-api/proxy/tencent/demand/occupy',
  DemandOccupyTask = '/v1/flask-api/proxy/tencent/demand/occupy_task',
  ExportDemandOccupy = '/v1/flask-api/proxy/tencent/demand/occupy/export',
  DeliveryInfoDelivery = '/v1/business/k/deliveryInfo/delivery',
  GetHddDeviceStatusList = '/v1/business/k/hdd/capacity/list',
  GetHddPartitionList = '/v1/business/k/hdd/partition/list',
  GetHddSummaryDailyPeak = '/v1/business/k/hdd/dailyPeak/summary',
  GetHddDailyPeakList = '/v1/business/k/hdd/dailyPeak/list',
  ExportHddDailyPeak = '/v1/business/k/hdd/dailyPeak/export',
  GetHddDailyPeakDetailList = '/v1/business/k/hdd/dailyPeak/detail/list',
  ExportHddDailyPeakDetail = '/v1/business/k/hdd/dailyPeak/detail/export',
  GetHdd5minSummaryList = '/v1/business/k/hdd/5min/summary/list',
  ExportHdd5minSummary = '/v1/business/k/hdd/5min/summary/export',
  ExportHddDeviceStatus = '/v1/business/k/hdd/capacity/export',
  GetHddDiffList = '/v1/business/k/hdd/diff/list',
  ConfirmHddDiff = '/v1/business/k/hdd/diff/confirm',
  ExportHddDiff = '/v1/business/k/hdd/diff/export',
  GetDeliveryInfoMacList = '/v1/business/k/deliveryInfoMac/list',
  BatchUpdateDeliveryStatus = '/v1/business/k/deliveryInfo/status/batch',
  BatchBindDemandCheck = '/v1/business/k/deliveryInfo/batchBind/check',
  GetBatchBindDemandList = '/v1/business/k/batchMatchingDemand/list',
  BatchBindDemand = '/v1/business/k/deliveryInfo/demand/batch',
  ExportDeviceMacs = '/v1/flask-api/proxy/tencent/mac/export',
  GetMacReplaceHistory = '/v1/flask-api/proxy/tencent/mac/replace_history',
  ExportTraffic = '/v1/business/k/traffic/export',
  GetTrafficList = '/v1/business/k/traffic/list',
  DeliveryDifIsp = '/v1/business/k/delivery/difisp',
  MacDisable = '/v1/flask-api/proxy/tencent/mac/disable',
  GetMacHistory = '/v1/flask-api/proxy/tencent/mac/history',
  GetDailyPeakList = '/v1/business/k/traffic/dailyPeak/list',
  ExportDailyPeak = '/v1/business/k/traffic/dailyPeak/export',
  GetFlowStats = '/v1/business/k/traffic/5min/list',
  ExportTrafficBill = '/v1/business/k/traffic/bill/export',
  LockDemand = '/v1/flask-api/proxy/tencent/demand/lock',
  BusinessJoin = '/v1/flask-api/proxy/tencent/delivery/business_join',
  GetQualityAbnormal = '/v1/flask-api/proxy/tencent/quality/abnormal',
  GetDemandDetectList = '/v1/business/k/demand/detect/list',
  ExportDemandDetect = '/v1/business/k/demand/detect/export',
  ExportDemandList = '/v1/flask-api/proxy/tencent/demand/export',
}

export const GetDeviceList = (params: Recordable) => {
  return defHttp.get({ url: Api.GetDeviceList, params });
};
export const GetDeviceMacsList = (params: Recordable) => {
  return defHttp.get({ url: Api.GetDeviceMacsList, params });
};

export const GetMacReplaceHistory = (params: Recordable) => {
  return defHttp.get({ url: Api.GetMacReplaceHistory, params });
};
export const UpdateDeviceMac = (data: Recordable) => {
  return defHttp.post({ url: Api.UpdateDeviceMac, data, timeout: 5 * 60 * 1000 });
};

export const EditDeviceMac = (data: Recordable) => {
  return defHttp.post({ url: Api.EditDeviceMac, data });
};

export const GetBusinessList = (params: Recordable) => {
  return defHttp.get({ url: Api.GetBusinessList, params });
};
export const GetDevTypeList = () => {
  return defHttp.get({ url: Api.GetDevTypeList });
};
export const GetAreaList = (data: Recordable) => {
  return defHttp.get({ url: Api.GetAreaList, data });
};
export const getDeliveryInfoList = (params: Recordable) => {
  return defHttp.get({ url: Api.GetDeliveryInfoList, params });
};
export const getDeliveryOptions = (params: Recordable) => {
  return defHttp.get({ url: Api.GetDeliveryOptions, params });
};

export const getMatchingDemandList = (params: Recordable) => {
  return defHttp.get({ url: Api.GetMatchingDemandList, params });
};

export const getDemandList = (params: Recordable) => {
  return defHttp.get({ url: Api.GetDemandList, params });
};

export const DemandOccupy = (data: Recordable) => {
  return defHttp.post({ url: Api.DemandOccupy, data }, { isReturnNativeResponse: true });
};

export const DemandOccupyList = (params: Recordable) => {
  return defHttp.get({ url: Api.DemandOccupy, params });
};

export const updateDeliveryInfo = (id: number, data: Recordable) => {
  return defHttp.put({ url: Api.UpdateDeliveryInfo + '/' + id, data });
};

export const batchUpdateDeliveryInfo = (data: {
  ids: number[];
  bizType: string;
  updateData: Recordable;
}) => {
  return defHttp.post({ url: Api.BatchUpdateDeliveryInfo, data });
};

export const updateDeliveryInfoDemand = (id: number, data: Recordable) => {
  return defHttp.put({ url: Api.UpdateDeliveryInfo + '/' + id + '/demand', data });
};

export const deliveryInfoDelivery = (data: Recordable) => {
  return defHttp.post({ url: Api.DeliveryInfoDelivery, data, timeout: 2 * 60 * 1000 });
};
export const getRealDeliveryBw = (params: Recordable) => {
  return defHttp.get({ url: `${Api.UpdateDeliveryInfo}/realDeliveryBw`, params });
};

export const GetHddDeviceStatusList = (data: Recordable) => {
  return defHttp.post({ url: Api.GetHddDeviceStatusList, data, timeout: 5 * 60 * 1000 });
};
export const GetHddPartitionList = (params: Recordable) => {
  return defHttp.get({ url: Api.GetHddPartitionList, params });
};

export const getHddDailyPeakList = (params: Recordable) => {
  return defHttp.get({ url: Api.GetHddDailyPeakList, params });
};

export const getHddSummaryDailyPeak = (params: Recordable) => {
  return defHttp.get({ url: Api.GetHddSummaryDailyPeak, params });
};

export const GetDeliveryInfoMacList = (params: Recordable) => {
  return defHttp.get({ url: Api.GetDeliveryInfoMacList, params });
};
export const UpdateDeliveryInfoStatus = (id: number, data: Recordable) => {
  return defHttp.put({
    url: `${Api.UpdateDeliveryInfo}/${id}/status`,
    data,
  });
};

export const BatchUpdateDeliveryStatus = (data: Recordable) => {
  return defHttp.put({
    url: Api.BatchUpdateDeliveryStatus,
    data,
  });
};

export const BatchBindDemandCheck = (params: Recordable) => {
  return defHttp.get({
    url: Api.BatchBindDemandCheck,
    params,
  });
};

export const GetBatchBindDemandList = (params: Recordable) => {
  return defHttp.get({
    url: Api.GetBatchBindDemandList,
    params,
  });
};

export const BatchBindDemand = (data: Recordable) => {
  return defHttp.put({
    url: Api.BatchBindDemand,
    data,
  });
};

export const GetTrafficList = (params: Recordable) => {
  return defHttp.get({
    url: Api.GetTrafficList,
    params,
    timeout: 2 * 60 * 1000,
  });
};

export const DeliveryDifIsp = (data: Recordable) => {
  return defHttp.post({
    url: Api.DeliveryDifIsp,
    data,
  });
};

export const GetMacDisable = (params: Recordable) => {
  return defHttp.get({
    url: Api.MacDisable,
    params,
  });
};

export const GetMacHistory = (params: Recordable) => {
  return defHttp.get({
    url: Api.GetMacHistory,
    params,
  });
};

export const GetFlowStats = (params: Recordable) => {
  return defHttp.get({ url: Api.GetFlowStats, params });
};
export const GetDailyPeakList = (params: Recordable) => {
  return defHttp.get({ url: Api.GetDailyPeakList, params });
};

export const ExportDailyPeak = (data: Recordable) => {
  return defHttp.post({ url: Api.ExportDailyPeak, data });
};

export const ExportTrafficBill = (data: Recordable) => {
  return defHttp.post({ url: Api.ExportTrafficBill, data });
};

export const lockDemand = (data: { demand_id: string; is_locked: boolean; locked_bw: number }) => {
  return defHttp.post({ url: Api.LockDemand, data }, { isReturnNativeResponse: true });
};

export const GetHddDiffList = (params: Recordable) => {
  return defHttp.get({ url: Api.GetHddDiffList, params });
};

export const ConfirmHddDiff = (data: Recordable) => {
  return defHttp.post({ url: Api.ConfirmHddDiff, data });
};

export const BusinessJoin = (data: Recordable) => {
  return defHttp.post({ url: Api.BusinessJoin, data }, { isReturnNativeResponse: true });
};

export const GetQualityAbnormal = (params: Recordable) => {
  return defHttp.get({ url: Api.GetQualityAbnormal, params });
};

export const GetDemandDetectList = (params: Recordable) => {
  return defHttp.get({ url: Api.GetDemandDetectList, params });
};

export const ExportDemandDetect = (data: Recordable) => {
  return defHttp.post({ url: Api.ExportDemandDetect, data });
};

// 提交需求占用任务
export const DemandOccupyTask = (data: Recordable) => {
  return defHttp.post({ url: Api.DemandOccupyTask, data });
};

// 获取需求占用任务列表
export const DemandOccupyTaskList = (params: Recordable) => {
  return defHttp.get({ url: Api.DemandOccupyTask, params });
};

// 获取需求缺口统计列表（开放接口）
export const GetDemandGapList = () => {
  return defHttp.get({ url: '/v1/business/k/demand/gap/list' });
};
