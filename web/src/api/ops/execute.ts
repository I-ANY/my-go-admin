import { defHttp } from '@/utils/http/axios';

// 获取业务名称
export const GetBizName = (params: Recordable) => {
  return defHttp.get({ url: '/v1/execute/biz/name', params });
};

// 获取设备列表
export const GetDeviceList = (params: Recordable) => {
  return defHttp.get({ url: '/v1/execute/device/list', params });
};

// 获取执行任务列表
export const GetExecuteTaskRecord = (params: Recordable) => {
  return defHttp.get({ url: '/v1/ops/execute/tasks', params });
};
export const GetExecuteTaskRecordUser = (params: Recordable) => {
  return defHttp.get({ url: '/v1/ops/execute/tasks/user', params });
};

// 获取执行任务详情
export const GetExecuteRecordDetail = (params: Recordable) => {
  return defHttp.get({ url: '/v1/ops/execute/tasks/result', params });
};

// 功能选择树
export const GetOperationsTreeData = (params: Recordable) => {
  return defHttp.post(
    { url: '/v1/ops/script-configs/tree', data: params },
    { isTransformResponse: false },
  );
};

// 新增操作组
export const AddOperation = (params: Recordable) => {
  return defHttp.post({ url: '/v1/ops/script-configs', params });
};
// 修改操作组
export const EditOperation = (id: number, params: Recordable) => {
  return defHttp.put({ url: '/v1/ops/script-configs/' + id, params });
};

// 删除操作组
export const DelOperation = (id: number) => {
  return defHttp.delete({ url: '/v1/ops/script-configs/' + id });
};

// 获取操作组信息
export const GetScriptTasks = (params: Recordable) => {
  return defHttp.get({ url: '/v1/ops/script-configs', params });
};

// 拷贝功能或操作组
export const CopyScriptConfig = (params: Recordable) => {
  return defHttp.post({ url: '/v1/ops/script-configs/copy', params });
};

// 新增脚本任务
export const AddScriptTask = (params: Recordable) => {
  params.type = 'function';
  return defHttp.post({ url: '/v1/ops/script-configs', params });
};

// 修改脚本任务
export const EditScriptTask = (params: Recordable) => {
  return defHttp.put({ url: '/v1/ops/script-configs', params });
};
// 删除脚本任务
export const DelScriptTask = (id: number) => {
  return defHttp.delete({ url: '/v1/ops/script-configs/' + id });
};

// 执行
export const ScriptExecute = (params: Recordable) => {
  return defHttp.post({ url: '/v1/ops/execute/script', params });
};

// 执行结果导出
export const ExportExeResultDetail = (params: Recordable) => {
  return defHttp.post({ url: '/v1/ops/execute/result/export', params });
};
