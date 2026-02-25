import { defHttp } from '@/utils/http/axios';

export const Api = {
  GetServerPriorityList: () => '/v1/network/priority/server/list',
  GetBizPriorityList: () => '/v1/network/priority/biz/list',
  GetBizOwnerPriority: () => `/v1/network/priority/owner/biz`,
  UpdateBizPriority: () => `/v1/network/priority/biz`,
  UpdateServerPriority: () => `/v1/network/priority/server`,
  ExportBizPriority: () => '/v1/network/priority/biz/export',
  ExportServerPriority: () => '/v1/network/priority/server/export',
};

export const GetServerPriorityList = (params: Recordable) => {
  return defHttp.get({ url: Api.GetServerPriorityList(), params });
};

export const GetBizPriorityList = (params: Recordable) => {
  return defHttp.get({ url: Api.GetBizPriorityList(), params });
};

export const GetBizOwnerPriority = (params: Recordable) => {
  return defHttp.get({ url: Api.GetBizOwnerPriority(), params });
};

export const UpdateBizPriority = (data: Recordable) => {
  return defHttp.post({ url: Api.UpdateBizPriority(), data, timeout: 120 * 1000 });
};

export const UpdateServerPriority = (data: Recordable) => {
  return defHttp.post({ url: Api.UpdateServerPriority(), data, timeout: 120 * 1000 });
};
