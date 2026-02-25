import { defHttp } from '@/utils/http/axios';

enum Api {
  Category = '/v1/price/category',
  Mode = '/v1/price/mode',
  Record = '/v1/price/record',
  Region = '/v1/price/region',
  Zone = '/v1/price/zone',
  Biz = '/v1/price/biz',
  Owners = '/v1/price/owners',
  Business = '/v1/price/record/category',
  Webhook = '/v1/price/webhook',
  Users = '/v1/price/users',
  Node = '/v1/price/nodeRecord',
}

// category
export const getCategory = (params: Recordable) => {
  return defHttp.get({ url: Api.Category, params });
};
export const createCategory = (data: Recordable) => {
  return defHttp.post({ url: Api.Category, data });
};
export const deleteCategory = (id: number | string) =>
  defHttp.delete({ url: Api.Category + '/' + id });
export const updateCategory = (data: Recordable) => defHttp.put({ url: Api.Category, data });

// mode
export const getMode = (params: Recordable) => {
  return defHttp.get({ url: Api.Mode, params });
};
export const createMode = (data: Recordable) => {
  return defHttp.post({ url: Api.Mode, data });
};
export const deleteMode = (id: number | string) => defHttp.delete({ url: Api.Mode + '/' + id });
export const updateMode = (data: Recordable) => defHttp.put({ url: Api.Mode, data });

// record
export const getRecord = (params: Recordable) => {
  return defHttp.get({ url: Api.Record, params });
};
export const createRecord = (data: Recordable) => {
  return defHttp.post({ url: Api.Record, data });
};
export const deleteRecords = (params: Recordable) => defHttp.delete({ url: Api.Record, params });
export const deleteRecord = (id: number | string) => defHttp.delete({ url: Api.Record + '/' + id });
export const updateRecord = (data: Recordable) => defHttp.put({ url: Api.Record, data });
export const updateBatchNumPointRecord = (data: Recordable) =>
  defHttp.put({ url: Api.Record + '/batch/point', data });
export const updateBatchNumRecord = (data: Recordable) =>
  defHttp.put({ url: Api.Record + '/batch/num', data });
export const getAllRecord = (params: Recordable) => {
  return defHttp.get({ url: Api.Record + '/all', params });
};
export const getOwnerRecord = (params: Recordable) => {
  return defHttp.get({ url: Api.Record + '/owner', params });
};

// Region
export const getRegion = (params: Recordable) => {
  return defHttp.get({ url: Api.Region, params });
};

// Zone
export const getZone = (params: Recordable) => {
  return defHttp.get({ url: Api.Zone, params });
};

// Biz
export const getBiz = (params: Recordable) => {
  return defHttp.get({ url: Api.Biz, params });
};

// update biz
export const updateBiz = (data: Recordable) => defHttp.put({ url: Api.Biz, data });

// update biz notify
export const updateNotify = (data: Recordable) => defHttp.put({ url: Api.Biz + '/notify', data });

// Biz
export const getNotRelBiz = (params: Recordable) => {
  return defHttp.get({ url: Api.Biz + '/residual', params });
};

// Hosts
export const getHosts = (params: Recordable) => {
  return defHttp.get({ url: Api.Category + '/host', params });
};

// 导入业务组
export const importCategory = (data: Recordable) => {
  return defHttp.post({ url: Api.Category + '/import', data });
};

// 导入单价
export const importRecord = (data: Recordable) => {
  return defHttp.post({ url: Api.Record + '/import', data });
};

// owners
export const getOwners = () => {
  return defHttp.get({ url: Api.Owners });
};

export const getBusiness = (params: Recordable) => {
  return defHttp.get({ url: Api.Business, params });
};

// export
export const exportRecord = (data: Recordable) => {
  return defHttp.post({ url: Api.Record + '/export', data });
};

// query webhook
export const getWebhooks = () => {
  return defHttp.get({ url: Api.Webhook });
};

// add webhook
export const createWebhook = (data: Recordable) => {
  return defHttp.post({ url: Api.Webhook, data });
};

// del webhook
export const deleteWebhook = (id: number | string) =>
  defHttp.delete({ url: Api.Webhook + '/' + id });

// update webhook
export const updateWebhook = (data: Recordable) => defHttp.put({ url: Api.Webhook, data });

// query user
export const getUsers = () => {
  return defHttp.get({ url: Api.Users });
};

export const getNodesByCategory = (params: Recordable) => {
  return defHttp.get({ url: Api.Node + '/category', params });
};

export const getNodePrice = (params: Recordable) => {
  return defHttp.get({ url: Api.Node, params });
};

export const updateNodes = (data: Recordable) => {
  return defHttp.put({ url: Api.Node + '/batch/num', data });
};

export const exportNodeRecord = (data: Recordable) => {
  return defHttp.post({ url: Api.Node + '/export', data });
};

export const updateNode = (data: Recordable) => {
  return defHttp.put({ url: Api.Node, data });
};

export const importNodeRecord = (data: Recordable) => {
  return defHttp.post({ url: Api.Node + '/import', data });
};

export const getRecordOption = (params: Recordable) => {
  return defHttp.get({ url: Api.Record + '/option', params });
};

export const mockProfit = (data: Recordable) => {
  return defHttp.post({ url: Api.Record + '/mock/profit', data });
};
