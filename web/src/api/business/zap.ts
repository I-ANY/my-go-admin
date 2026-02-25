import { defHttp } from '@/utils/http/axios';

export enum Api {
  GetDeviceList = '/v1/business/zap/device/list',
  ExportDeviceList = '/v1/business/zap/device/export',
  GetQualityEvents = '/v1/business/zap/quality/event/list',
  // GetQualityEventsRealtime = '/v1/business/zap/quality/event/list/realtime',
  // CSCM平台
  GetQualityEventsRealtime = '/v1/flask-api/proxy/zap/quality/event/list/realtime',
  GetNodeStatusRealtime = '/v1/flask-api/proxy/zap/device/status/realtime',
  // 业务统一管理平台
  GetQualityEventsRealtimeBytedance = '/v1/business/zap/bytedance/quality/event/list/realtime',
  GetNodeStatusRealtimeBytedance = '/v1/business/zap/bytedance/device/status/realtime',
  GetDeliveryDevice = '/v1/business/zap/delivery/device/list',
  UpdateDeliveryDevice = '/v1/business/zap/delivery/device/update',
  UpdateDeliveryDeviceStatus = '/v1/business/zap/delivery/device/update/status',
  GetZapAreaList = '/v1/business/zap/delivery/area/info',
  DeliveryDeviceImport = '/v1/business/zap/delivery/device/import',
  DeliveryDeviceCommit = '/v1/business/zap/delivery/device/submit',
  DeliveryDeviceSubmittingRefresh = '/v1/business/zap/delivery/device/refresh',
  DeviceTempOffline = '/v1/business/zap/device/tempOffline',
  DeviceForeverOffline = '/v1/business/zap/device/foreverOffline',
  DeviceOnline = '/v1/business/zap/device/online',
  DelZapPlatformToken = '/v1/business/zap/device/token/del',
  SyncEcdnData = '/v1/business/zap/delivery/device/sync/ecdn',
  DeviceUpdate = '/v1/business/zap/device/update',
  GetBwChartData = '/v1/business/zap/device/bw/stats',
}

// 获取设备列表
export const GetDeviceList = (params: Recordable) => {
  return defHttp.get({ url: Api.GetDeviceList, params });
};
// 设备信息导出
export const ExportDeviceList = (params: Recordable) => {
  return defHttp.get({ url: Api.ExportDeviceList, params });
};
// 获取质量事件列表
export const GetQualityEvents = (params: Recordable) => {
  return defHttp.get({ url: Api.GetQualityEvents, params });
};
// 实时获取质量事件列表（后端直接调的业务平台接口）
export const GetQualityEventsRealtime = (params: Recordable) => {
  return defHttp.get({ url: Api.GetQualityEventsRealtime, params });
};
// 节点状态实时查询
export const GetNodeStatusRealtime = (params: Recordable) => {
  return defHttp.get({ url: Api.GetNodeStatusRealtime, params });
};

// 实时获取质量事件列表（后端直接调的业务平台接口）
export const GetQualityEventsRealtimeBytedance = (params: Recordable) => {
  return defHttp.get({ url: Api.GetQualityEventsRealtimeBytedance, params });
};
// 节点状态实时查询
export const GetNodeStatusRealtimeBytedance = (params: Recordable) => {
  return defHttp.get({ url: Api.GetNodeStatusRealtimeBytedance, params });
};

// 交付设备信息列表
export const GetDeliveryDevice = (params: Recordable) => {
  return defHttp.get({ url: Api.GetDeliveryDevice, params });
};
// 编辑修改设备信息
export const UpdateDeliveryDevice = (params: Recordable) => {
  return defHttp.post({ url: Api.UpdateDeliveryDevice, params });
};
// 编辑修改设备状态为关闭
export const UpdateDeliveryDeviceStatus = (params: Recordable) => {
  return defHttp.put({ url: Api.UpdateDeliveryDeviceStatus, params });
};
// 获取地域信息
export const GetZapAreaList = () => {
  return defHttp.get({ url: Api.GetZapAreaList });
};
// 编辑修改设备信息
export const DeliveryDeviceImport = (params: Recordable) => {
  return defHttp.post({ url: Api.DeliveryDeviceImport, params });
};
// 提交设备
export const DeliveryDeviceCommit = (params: Recordable) => {
  return defHttp.post({ url: Api.DeliveryDeviceCommit, params });
};
// 设备提交状态更新
export const DeliveryDeviceSubmittingRefresh = (params: Recordable) => {
  return defHttp.put({ url: Api.DeliveryDeviceSubmittingRefresh, params });
};
// 设备临时下线
export const DeviceTempOffline = (params: Recordable) => {
  return defHttp.put({ url: Api.DeviceTempOffline, params });
};
// 设备永久下线
export const DeviceForeverOffline = (params: Recordable) => {
  return defHttp.put({ url: Api.DeviceForeverOffline, params });
};
// 设备上线
export const DeviceOnline = (params: Recordable) => {
  return defHttp.put({ url: Api.DeviceOnline, params });
};

// 设备上线
export const DeviceUpdate = (params: Recordable) => {
  return defHttp.put({ url: Api.DeviceUpdate, params });
};

// 业务平台删除设备token，重新刷新
export const DelZapPlatformToken = () => {
  return defHttp.put({ url: Api.DelZapPlatformToken });
};

// 业务平台删除设备token，重新刷新
export const SyncEcdnData = (params: Recordable) => {
  return defHttp.post({ url: Api.SyncEcdnData, data: params });
};

// 获取带宽统计数据
export const GetBwChartData = (params: Recordable) => {
  return defHttp.get({ url: Api.GetBwChartData, params });
};
