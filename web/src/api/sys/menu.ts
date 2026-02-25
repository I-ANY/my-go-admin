import { defHttp } from '@/utils/http/axios';
import { getMenuListResultModel } from './model/menuModel';

enum Api {
  GetMenuList = '/v1/system/userMenu',
  MenuRecord = '/v1/system/menu',
}

/**
 * @description: Get user menu based on id
 */

export const getMenuList = () => {
  return defHttp.get<getMenuListResultModel>({ url: Api.GetMenuList });
};

export const updateMenu = (id: number, data) => {
  return defHttp.put({ url: Api.MenuRecord + '/' + id, data });
};

export const addMenu = (data) => {
  return defHttp.post({ url: Api.MenuRecord, data });
};

export const deleteMenu = (id: number) => {
  return defHttp.delete({ url: Api.MenuRecord + '/' + id });
};
