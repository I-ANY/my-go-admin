import { defHttp } from '@/utils/http/axios';
import { DictTypeItem } from './model/dict';

enum Api {
  GetAllDict = '/v1/system/dict/all',
}
export const GetAllDict = () => {
  return defHttp.get<DictTypeItem[]>({ url: Api.GetAllDict });
};
