import { defHttp } from '@/utils/http/axios';

export enum Api {
  // 任务配置管理
  TaskConfig = '/v1/ops/agent/task-configs',

  // 任务执行
  TaskDetail = '/v1/ops/agent/task',
  GetTaskList = '/v1/ops/agent/tasks',
}

// 参数选项定义
export interface ParamOption {
  label: string;
  value: string;
}

// 步骤参数定义
export interface StepParam {
  name: string;
  type: 'select' | 'input';
  required: boolean;
  default?: any;
  help?: string; // input类型的帮助文本
  options?: ParamOption[]; // select类型的选项列表
}

// 任务配置相关接口
export interface TaskStep {
  name: string;
  description: string;
  script_path?: string;
  script_name?: string;
  script_url?: string;
  timeout: number;
  required: boolean;
  on_failure: string;
  params?: StepParam[];
}

export interface TaskConfigParams {
  business: string;
  name: string;
  taskType: string;
  description?: string;
  params: Record<string, any>;
  steps: TaskStep[];
}

export interface TaskConfigItem {
  id: number;
  business: string;
  name: string;
  taskType: string;
  description: string;
  params: Record<string, any>;
  steps: TaskStep[];
  createBy: number;
  updateBy: number;
  createdAt: string;
  updatedAt: string;
}

export interface TaskConfigListParams {
  pageIndex: number;
  pageSize: number;
  business?: string;
  taskType?: string;
  name?: string;
}

// 任务执行相关接口
export interface TaskSubmitParams {
  task_name: string;
  task_type: string;
  server_id: string;
  execute_mode: 'sync' | 'async';
  timeout?: number;
  max_retries?: number;
  params: Record<string, any>;
  steps?: TaskStep[];
  callback_url?: string;
  user_data?: Record<string, any>;
}

export interface TaskStepResult {
  name: string;
  status: string;
  start_time?: string;
  end_time?: string;
  output: string;
  exit_code: number;
  error?: string;
}

export interface TaskExecItem {
  task_id: string;
  task_name: string;
  task_type: string;
  server_id: string;
  execute_mode: 'sync' | 'async';
  status: string;
  start_time?: string;
  end_time?: string;
  timeout: number;
  max_retries: number;
  params: Record<string, any>;
  steps_result: TaskStepResult[];
  created_at: string;
  updated_at: string;
}

export interface TaskListParams {
  pageIndex: number;
  pageSize: number;
  task_id?: string;
  task_type?: string;
  server_id?: string;
  execute_mode?: 'sync' | 'async';
  status?: string;
}

// API调用函数
export const getTaskConfigList = (params: TaskConfigListParams) =>
  defHttp.get<{ items: TaskConfigItem[]; total: number }>({
    url: Api.TaskConfig,
    params,
  });

export const createTaskConfig = (data: TaskConfigParams) =>
  defHttp.post<TaskConfigItem>({
    url: Api.TaskConfig,
    data,
  });

export const getTaskConfig = (id: number) =>
  defHttp.get<TaskConfigItem>({
    url: `${Api.TaskConfig}/${id}`,
  });

export const updateTaskConfig = (id: number, data: TaskConfigParams) =>
  defHttp.put<TaskConfigItem>({
    url: `${Api.TaskConfig}/${id}`,
    data,
  });

export const deleteTaskConfig = (id: number) =>
  defHttp.delete({
    url: `${Api.TaskConfig}/${id}`,
  });

export const submitTask = (data: TaskSubmitParams) =>
  defHttp.post<{ task_id: string; status: string; message: string }>({
    url: Api.TaskDetail,
    data,
  });

export const getTaskStatus = (taskId: string) =>
  defHttp.get<TaskExecItem>({
    url: `${Api.TaskDetail}/${taskId}/status`,
  });

export const getTaskList = (params: TaskListParams) =>
  defHttp.get<{ items: TaskExecItem[]; total: number }>({
    url: Api.GetTaskList,
    params,
  });
