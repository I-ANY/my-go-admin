import { defHttp } from '@/utils/http/axios';

enum Api {
  GetContainerInfoList = '/v1/business/ting/container/info/list',
  GetContainerInfo = '/v1/business/ting/container/info',
  GetContainerDetailList = '/v1/business/ting/container/detail/list',
}
export const GetContainerInfoList = (params: Recordable) => {
  return defHttp.get({ url: Api.GetContainerInfoList, params });
};
export const GetContainerInfo = (id: number) => {
  return defHttp.get({ url: Api.GetContainerInfo + '/' + id });
};
export const GetContainerDetailList = (params: Recordable) => {
  return defHttp.get({ url: Api.GetContainerDetailList, params });
};
