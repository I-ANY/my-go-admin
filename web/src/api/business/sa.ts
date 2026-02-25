import { defHttp } from '@/utils/http/axios';

export enum Api {
  GetGuidList = '/v1/business/sa/guid/list',
  GuidCheck = '/v1/business/sa/guid/check',
  GetResourceGapList = '/v1/business/sa/resourceGap/list',
  GetResourceGapTypeList = '/v1/business/sa/resourceGap/type/list',
  GetResourceGapProvinceList = '/v1/business/sa/resourceGap/province/list',

  GetUtilizationList = '/v1/business/sa/utilization/list',
  GetLimitList = '/v1/business/sa/limit/list',
  ExportUtilizationList = '/v1/business/sa/utilization/list/export',
}

// 获取GUID列表
export const GetGuidList = (params: Recordable) => {
  return defHttp.get({ url: Api.GetGuidList, params });
};

// Guid检测
export const GuidCheck = () => {
  return defHttp.post({ url: Api.GuidCheck });
};

export const GetResourceGapList = (params: Recordable) => {
  return defHttp.get({ url: Api.GetResourceGapList, params });
};

export const GetResourceGapTypeList = (params: Recordable) => {
  return defHttp.get({ url: Api.GetResourceGapTypeList, params });
};

export const GetResourceGapProvinceList = (params: Recordable) => {
  return defHttp.get({ url: Api.GetResourceGapProvinceList, params });
};

export const GetUtilizationList = (params: Recordable) => {
  return defHttp.get({ url: Api.GetUtilizationList, params });
};

export const GetLimitList = (params: Recordable) => {
  return defHttp.get({ url: Api.GetLimitList, params });
};
