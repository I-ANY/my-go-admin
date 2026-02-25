import { defHttp } from '@/utils/http/axios';

export enum Api {
  GetTrafficList = '/v1/business/b/traffic/list',
  GetTrafficCheckResultList = '/v1/business/b/traffic/checkResult/list',
  ExportTraffic = '/v1/business/b/traffic/export',
  GetTrafficOptions = '/v1/business/b/traffic/options',
}

export const GetTrafficList = (params: Recordable) => {
  return defHttp.get({ url: Api.GetTrafficList, params });
};

export const GetTrafficCheckResultList = (params: Recordable) => {
  return defHttp.get({ url: Api.GetTrafficCheckResultList, params });
};

export const GetTrafficOptions = (params: Recordable) => {
  return defHttp.get({ url: Api.GetTrafficOptions, params });
};
