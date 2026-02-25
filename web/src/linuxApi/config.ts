import { linuxHttp } from '@/utils/http/axios';

export enum Api {
  configList = '/api/configs',
  createEncryptionConfiguration = '/api/configs',
  updateEncryptionConfiguration = '/api/configs',
  deleteEncryptionConfiguration = '/api/configs',
  taskBuildCenterList = '/api/build',
  getAllBasenameList = '/api/configs/basenames',
  createTaskBuildCenter = '/api/build',
  getBasenameList = '/api/build/basenames',
}

// 获取配置列表
export const getEncryptionConfigurationList = (params: Recordable) => {
  return linuxHttp.get({ url: Api.configList, params });
};

// 获取所有 basename 列表
export const getAllBasenameList = () => {
  return linuxHttp.get({ url: Api.getAllBasenameList });
};

// 根据 basename 获取所有 os_type 列表
export const getAllOsTypeList = (basename: string) => {
  return linuxHttp.get({ url:`/api/configs/${basename}/os-types`});
};

// 根据供应商名称ios_sign查询基础配置名称列表
export const getBasenameListByIsoSign = (iso_sign: string) => {
  return linuxHttp.get({ url: iso_sign ? `${Api.getBasenameList}?iso_sign=${iso_sign}` : Api.getBasenameList });
};

// 根据 basename 和 os_type 获取所有 tar_version 列表
export const getAllTarVersionList = (basename: string, os_type: string) => {
  return linuxHttp.get({ url: `/api/configs/${basename}/${os_type}/versions` });
};

// 根据iso_sign,basename获取操作系统列表
export const getAllOsTypeListByIsoSignAndBasename = (iso_sign: string, basename: string) => {
  return linuxHttp.get({ url: '/api/build/os-types', params: { iso_sign, basename } });
};

// 新增配置
export const createEncryptionConfiguration = (data: Recordable) => {
  return linuxHttp.post({ url: Api.createEncryptionConfiguration, data });
};

// 更新配置
export const updateEncryptionConfiguration = (params: Recordable, data: Recordable) => {
  return linuxHttp.put({ url: Api.updateEncryptionConfiguration + '/' + params.basename + '/' + params.os_type + '/' + params.version, data });
};

// 删除配置
export const deleteEncryptionConfiguration = (data: Recordable) => {
  return linuxHttp.delete({ url: Api.deleteEncryptionConfiguration + '/' + data.basename + '/' + data.os_type + '/' + data.tar_version });
};

// 获取任务构建中心列表
export const getTaskBuildCenterList = (params: Recordable) => {
  return linuxHttp.get({ url: Api.taskBuildCenterList, params });
};

// 新增任务
export const createTaskBuildCenter = (data: Recordable) => {
  return linuxHttp.post({ url: Api.taskBuildCenterList, data });
};

// 获取任务详情
export const getTaskBuildCenterDetail = (params: Recordable) => {
  return linuxHttp.get({ url: `/api/build/${params.task_id}/log` });
};
