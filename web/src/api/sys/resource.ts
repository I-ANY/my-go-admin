import { defHttp } from '@/utils/http/axios';

const Api = {
  getResourceList: '/v1/system/resource/list',
  getResourceTableList: '/v1/system/resource/table/list',
  getResourceTableField: '/v1/system/resource/table/field',
  resource: (id: number) => `/v1/system/resource/${id}`,
  addResource: '/v1/system/resource',
  getResourceViewColumns: '/v1/system/resource/view/table-columns',
  getResourceViewSearchFormSchemas: '/v1/system/resource/view/search-form-schemas',
  getResourceDetailList: '/v1/system/resource/detail/list',
  getRoleResourceInfo: '/v1/system/role/resource/info',
  getRoleResourceDetailList: (roleId: number) => `/v1/system/role/${roleId}/resource/detail/list`,
  updateRoleResource: (roleId: number) => `/v1/system/role/${roleId}/resource`,
  getRoleAuthedResource: () => `/v1/system/role/authed/resource`,
  getBusinessResource: () => `/v1/system/resource/subcategory`,
  roleResourceAuth: (id: number) => `/v1/system/role/${id}/resource/auth`,
};

export const getResourceList = (params: Recordable) => {
  return defHttp.get({ url: Api.getResourceList, params });
};
export const getResourceTableList = (params: Recordable) => {
  return defHttp.get({ url: Api.getResourceTableList, params });
};
export const getResourceTableField = (params: Recordable) => {
  return defHttp.get({ url: Api.getResourceTableField, params });
};

export const addResource = (data: Recordable) => {
  return defHttp.post({ url: Api.addResource, data });
};
export const deleteResource = (id: number) => {
  return defHttp.delete({ url: Api.resource(id) });
};
export const updateResource = (id: number, data: Recordable) => {
  return defHttp.put({ url: Api.resource(id), data });
};

export const getResourceViewColumns = (params: Recordable) => {
  return defHttp.get({ url: Api.getResourceViewColumns, params });
};
export const getResourceViewSearchFormSchemas = (params: Recordable) => {
  return defHttp.get({ url: Api.getResourceViewSearchFormSchemas, params });
};
export const getResourceDetailList = (params: Recordable) => {
  return defHttp.get({ url: Api.getResourceDetailList, params });
};
export const getRoleResourceInfo = (params: Recordable) => {
  return defHttp.get({ url: Api.getRoleResourceInfo, params });
};
export const getRoleResourceDetailList = (roleId: number, params: Recordable) => {
  return defHttp.get({ url: Api.getRoleResourceDetailList(roleId), params });
};
export const updateRoleResource = (roleId: number, data: Recordable) => {
  return defHttp.post({ url: Api.updateRoleResource(roleId), data });
};

export const getRoleAuthedResource = (params: Recordable) => {
  return defHttp.get({ url: Api.getRoleAuthedResource(), params });
};

export const getBusinessResource = (params: Recordable) => {
  return defHttp.get({ url: Api.getBusinessResource(), params });
};

export const roleResourceAuth = (id: number, data: Recordable) => {
  return defHttp.post({ url: Api.roleResourceAuth(id), data });
};
