import { defHttp } from '@/utils/http/axios';

export enum Api {
  GetTsDetail = '/v1/business/ts/uts',
  GetTsHibernationHost = '/v1/business/ts/hibernationHost/list',
}

export const GetTsDetail = (params: Recordable) => {
  return defHttp.get({ url: Api.GetTsDetail, params });
};

export const GetTsHibernationHost = (params: Recordable) => {
  return defHttp.get({ url: Api.GetTsHibernationHost, params });
};
