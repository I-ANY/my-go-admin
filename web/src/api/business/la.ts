import { defHttp } from '@/utils/http/axios';

export const Api = {
  GetTrafficList: '/v1/business/la/traffic/list',
  GetBusiness: '/v1/business/la/business',
  ExportTraffic: '/v1/business/la/traffic/export',
  GetDeviceList: '/v1/business/la/device/list',
  GetDeviceHostHistoryList: '/v1/business/la/device/host-history/list',
  GetSLAInfoList: '/v1/business/la/sla/info/list',
  GetAbnormalHostList: '/v1/business/la/sla/abnormal-host/list',
  ExportSLAInfo: () => '/v1/business/la/sla/detail/export',
  ExportDeviceInfo: () => '/v1/business/la/device/info/export',
};

export const GetTrafficList = (params: Recordable) => {
  return defHttp.get({
    url: Api.GetTrafficList,
    params,
  });
};
export const GetBusiness = (params: Recordable) => {
  return defHttp.get({
    url: Api.GetBusiness,
    params,
  });
};

export const GetDeviceList = (params: Recordable) => {
  return defHttp.get({
    url: Api.GetDeviceList,
    params,
  });
};

export const GetDeviceHostHistoryList = (params: Recordable) => {
  return defHttp.get({
    url: Api.GetDeviceHostHistoryList,
    params,
  });
};
export const GetSLAInfoList = (params: Recordable) => {
  return defHttp.get({
    url: Api.GetSLAInfoList,
    params,
  });
};
export const GetAbnormalHostList = (params: Recordable) => {
  return defHttp.get({
    url: Api.GetAbnormalHostList,
    params,
  });
};
