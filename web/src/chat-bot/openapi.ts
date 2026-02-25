import { defHttp } from '@/utils/http/axios';

enum Api {
  GetDifiApi = '/v1/business/openaiapi/difi/api',
}

export const GetDifiApi = (data: Recordable) => {
  return defHttp.post({ url: Api.GetDifiApi, data });
};
