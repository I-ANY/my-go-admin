import { defHttp } from '@/utils/http/axios';
import { getApiListModel } from './model/api';
import { BasicPageParams } from '../model/baseModel';

enum Api {
  GetApiList = '/v1/system/api/list',
}

export const getApiList = (params: BasicPageParams) => {
  return defHttp.get<getApiListModel>({ url: Api.GetApiList, params });
};
