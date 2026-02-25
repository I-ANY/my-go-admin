import { defHttp } from '@/utils/http/axios';

export enum Api {
  GetZPExtneralTrafficDetail = '/v1/business/zp/extneral/traffic/detail',
}
// 获取外部流量明细
export const GetZPExtneralTrafficDetail = (params: Recordable) => {
  return defHttp.get({ url: Api.GetZPExtneralTrafficDetail, params });
};
