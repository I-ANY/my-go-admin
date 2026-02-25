import { defHttp } from '@/utils/http/axios';

export const Api = {
  GetCategoryList: '/v1/business/biz/category',
  GetCategoryListAll: '/v1/business/biz/category/all',
  CreateCategory: '/v1/business/biz/category/create',
  UpdateCategoryList: '/v1/business/biz/category/update',
  DelCategory: '/v1/business/biz/category/del',
  // CreateSubcategoryList = '/v1/business/biz/subcategory/create',
  GetSubcategoryList: '/v1/business/biz/subcategory',
  GetSubcategoryListAll: '/v1/business/biz/subcategory/all',
  GetSubcategoryFilterList: '/v1/business/biz/subcategory/filter',
  UpdateSubcategory: '/v1/business/biz/subcategory/update',
  GetAuthedBiz: () => '/v1/business/biz/authed-biz',
  GetAuthedServer: () => '/v1/business/cm/server',
  GetServerHistoryInfo: () => '/v1/business/cm/server/history/list',
  // DelSubcategory = '/v1/business/biz/subcategory/del',
  // BusinessGroup (CRUD)
  BusinessGroupList: '/v1/business/biz/group',
  BusinessGroupOptions: '/v1/business/biz/group/options',
  CreateBusinessGroup: '/v1/business/biz/group',
  UpdateBusinessGroup: (id: number) => '/v1/business/biz/group/' + id,
  DeleteBusinessGroup: (id: number) => '/v1/business/biz/group/' + id,
};

export const CreateCategory = (params: Recordable) => {
  return defHttp.post({ url: Api.CreateCategory, params });
};

export const GetCategoryList = (params: Recordable) => {
  return defHttp.get({ url: Api.GetCategoryList, params });
};

export const GetCategoryListAll = (params: Recordable) => {
  return defHttp.get({ url: Api.GetCategoryListAll, params });
};

export const UpdateCategory = (id: number, params: Recordable) => {
  return defHttp.post({ url: Api.UpdateCategoryList + '/' + id, params });
};

export const DelCategory = (id: string | number) => {
  return defHttp.delete({ url: Api.DelCategory + '/' + id });
};

// export const CreateSubcategoryList = (params: Recordable) => {
//   return defHttp.post({ url: Api.CreateSubcategoryList, params });
// };

export const GetSubcategoryList = (params: Recordable) => {
  return defHttp.get({ url: Api.GetSubcategoryList, params });
};

export const GetSubcategoryListAll = (params: Recordable) => {
  return defHttp.get({ url: Api.GetSubcategoryListAll, params });
};

export const GetSubcategoryFilterList = (params: Recordable) => {
  return defHttp.get({ url: Api.GetSubcategoryFilterList, params });
};

export const UpdateSubcategory = (params: Recordable) => {
  return defHttp.post({ url: Api.UpdateSubcategory, params });
};

// export const DelSubcategory = (id: string | number) => {
//   return defHttp.delete({ url: Api.DelSubcategory + '/' + id });
// };

export const GetAuthedBiz = (params: Recordable) => {
  return defHttp.get({ url: Api.GetAuthedBiz(), params });
};

export const GetAuthedServer = (params: Recordable) => {
  return defHttp.get({ url: Api.GetAuthedServer(), params });
};

// BusinessGroup (CRUD)
export const GetBusinessGroupList = (params: Recordable) => {
  return defHttp.get({ url: Api.BusinessGroupList, params });
};

export const GetBusinessGroupOptions = () => {
  return defHttp.get({ url: Api.BusinessGroupOptions });
};

export const CreateBusinessGroup = (params: Recordable) => {
  return defHttp.post({ url: Api.CreateBusinessGroup, params });
};

export const UpdateBusinessGroup = (id: number, params: Recordable) => {
  return defHttp.put({ url: Api.UpdateBusinessGroup(id), params });
};

export const DeleteBusinessGroup = (id: number) => {
  return defHttp.delete({ url: Api.DeleteBusinessGroup(id) });
};

export const GetServerHistoryInfo = (params: Recordable) => {
  return defHttp.get({ url: Api.GetServerHistoryInfo(), params });
};
