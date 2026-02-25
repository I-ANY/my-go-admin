import {
  AccountParams,
  DeptListItem,
  MenuParams,
  RoleParams,
  RolePageParams,
  MenuListGetResultModel,
  DeptListGetResultModel,
  AccountListGetResultModel,
  RolePageListGetResultModel,
  RoleListGetResultModel,
} from './model/systemModel';
import { defHttp } from '@/utils/http/axios';

enum Api {
  AccountList = '/v1/system/user/list',
  Account = '/v1/system/user',
  IsAccountExist = '/system/accountExist',
  // DeptList = '/system/getDeptList',
  DeptList = '/v1/system/dept/tree',
  setRoleStatus = '/system/setRoleStatus',
  MenuList = '/system/getMenuList',
  GetAllMenu = '/v1/system/menu/tree',
  RolePageList = '/v1/system/role/list',
  GetAllRoleList = '/v1/system/getAllRoleList',
  RoleRecord = '/v1/system/role',
  OperLogList = '/v1/system/operaLog/list',
}

export const getAccountList = (params: AccountParams) =>
  defHttp.get<AccountListGetResultModel>({ url: Api.AccountList, params });
export const addUser = (data: Recordable) => defHttp.post({ url: Api.Account, data });
export const updateUser = (id: number | string, data: Recordable) =>
  defHttp.put({ url: Api.Account + '/' + id, data });
export const deleteUser = (id: number) => defHttp.delete({ url: Api.Account + '/' + id });

export const getDeptList = (params?: DeptListItem) =>
  defHttp.get<DeptListGetResultModel>({ url: Api.DeptList, params });

export const addDept = (data: Recordable) => defHttp.post({ url: Api.DeptList, data });
export const updateDept = (id: number, data: Recordable) =>
  defHttp.put({ url: Api.DeptList + '/' + id, data });
export const deleteDept = (id: number) => defHttp.delete({ url: Api.DeptList + '/' + id });

export const getMenuList = (params?: MenuParams) =>
  defHttp.get<MenuListGetResultModel>({ url: Api.MenuList, params });

export const getRoleListByPage = (params?: RolePageParams) =>
  defHttp.get<RolePageListGetResultModel>({ url: Api.RolePageList, params });

export const addRole = (data: Recordable) => defHttp.post({ url: Api.RoleRecord, data });
export const updateRole = (id: number | string, data: Recordable) =>
  defHttp.put({ url: Api.RoleRecord + '/' + id, data });

export const deleteRole = (id: number) => defHttp.delete({ url: Api.RoleRecord + '/' + id });

export const getAllRoleList = (params?: RoleParams) =>
  defHttp.get<RoleListGetResultModel>({ url: Api.GetAllRoleList, params });

export const setRoleStatus = (id: number, status: string) =>
  defHttp.post({ url: Api.setRoleStatus, params: { id, status } });

export const isAccountExist = (account: string) =>
  defHttp.post({ url: Api.IsAccountExist, params: { account } }, { errorMessageMode: 'none' });

export const getAllMenus = (params?: MenuParams) =>
  defHttp.get<MenuListGetResultModel>({ url: Api.GetAllMenu, params });

export const getOperaLogList = (params: any) => defHttp.get({ url: Api.OperLogList, params });
