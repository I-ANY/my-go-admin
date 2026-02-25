import { defHttp } from '@/utils/http/axios';

export enum Api {
  UpdateZPCookie = '/v1/business/zp/cookie/update',
  GetZPCookie = '/v1/business/zp/cookie/info',
  GetZPDeviceList = '/v1/business/zp/device/list',
  ExportDeviceList = '/v1/business/zp/device/export',
  SyncZPDevice = '/v1/business/zp/sync/device',
  GetZPTrafficDetailList = '/v1/business/zp/traffic/detail/list',
  GetZPTraffic95Gather = '/v1/business/zp/traffic/95/gather',
  CreateZPTraffic95 = '/v1/business/zp/traffic/95/create',
  UpdateZPTraffic95 = '/v1/business/zp/traffic/95/update',
  GetZPTraffic95List = '/v1/business/zp/traffic/95/list',
  ExportTrafficDetailList = '/v1/business/zp/traffic/detail/export',
  DeliveryDeviceList = '/v1/business/zp/delivery/device/list',
  DeliveryDeviceEdit = '/v1/business/zp/delivery/device/edit',
  DeliveryDeviceSubmit = '/v1/business/zp/delivery/device/submit',
  DeliveryDeviceStatusRefresh = '/v1/business/zp/delivery/device/refresh',
  DeliveryDeviceBatchEdit = '/v1/business/zp/delivery/device/batch/edit',
  DeliveryDeviceSyncEcdnBw = '/v1/business/zp/delivery/device/sync/ecdnBw',
  GetBwStats = '/v1/business/zp/device/bw/stats/list',
  GetBwStatsDiffDetail = '/v1/business/zp/device/bw/stats/diff/detail',
}

// 获取cookie
export const GetZPCookie = () => {
  return defHttp.get({ url: Api.GetZPCookie });
};
// 更新cookie
export const UpdateZPCookie = (params: Recordable) => {
  return defHttp.post({ url: Api.UpdateZPCookie, params });
};
// 获取设备列表
export const GetDeviceList = (params: Recordable) => {
  return defHttp.get({ url: Api.GetZPDeviceList, params });
};
// 同步设备
export const SyncZPDevice = () => {
  return defHttp.post({ url: Api.SyncZPDevice });
};
// 获取流量数据
export const GetZPTrafficDetailList = (params: Recordable) => {
  return defHttp.get({ url: Api.GetZPTrafficDetailList, params });
};

// 获取流量数95值
export const GetZPTraffic95Gather = (params: Recordable) => {
  return defHttp.get({ url: Api.GetZPTraffic95Gather, params });
};

// 创建95流量数据
export const CreateZPTraffic95 = (data: Recordable) => {
  return defHttp.post({ url: Api.CreateZPTraffic95, data });
};

// 获取流量数据
export const GetZPTraffic95List = (params: Recordable) => {
  return defHttp.get({ url: Api.GetZPTraffic95List, params });
};

// 创建95流量数据
export const UpdateZPTraffic95 = (data: Recordable) => {
  return defHttp.post({ url: Api.UpdateZPTraffic95, data });
};

// 查询交付设备
export const DeliveryDeviceList = (params: Recordable) => {
  return defHttp.get({ url: Api.DeliveryDeviceList, params });
};

// 编辑交付设备信息
export const DeliveryDeviceEdit = (params: Recordable) => {
  return defHttp.put({ url: Api.DeliveryDeviceEdit, params });
};
// 交付设备
export const DeliveryDeviceSubmit = (params: Recordable) => {
  return defHttp.post({ url: Api.DeliveryDeviceSubmit, params });
};
// 刷新交付设备交付状态
export const DeliveryDeviceStatusRefresh = (params: Recordable) => {
  return defHttp.post({ url: Api.DeliveryDeviceStatusRefresh, params });
};
// 批量编辑运营商与带宽
export const DeliveryDeviceBatchEdit = (params: Recordable) => {
  return defHttp.put({ url: Api.DeliveryDeviceBatchEdit, params });
};
// 手动同步ECDN带宽
export const DeliveryDeviceSyncEcdnBw = (params: Recordable) => {
  return defHttp.put({ url: Api.DeliveryDeviceSyncEcdnBw, params });
};

// 获取带宽统计数据
export const GetBwStats = (params: Recordable) => {
  return defHttp.get({ url: Api.GetBwStats, params });
};

// 获取带宽统计数据差异详情
export const GetBwStatsDiffDetail = (params: Recordable) => {
  return defHttp.get({ url: Api.GetBwStatsDiffDetail, params });
};
