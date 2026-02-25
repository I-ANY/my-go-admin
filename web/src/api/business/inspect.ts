import { defHttp } from '@/utils/http/axios';

export enum Api {
  GetTaskList = '/v1/business/inspect/task/list',
  GetResultList = '/v1/business/inspect/result/list',
  GetResultField = '/v1/business/inspect/result/field',
  GetInspectServerList = '/v1/business/inspect/server/list',
  InspectServer = '/v1/business/inspect/server',
  ExportInspectResult = '/v1/business/inspect/result/export',
  GetResultFieldName = '/v1/business/inspect/result/field/name',
  GetResultSummaryStatus = '/v1/business/inspect/result/summary/status',
  GetResultSummaryField = '/v1/business/inspect/result/summary/field',
  GetTaskNames = '/v1/business/inspect/task/names',
}

export const GetTaskList = (params: Recordable) => {
  return defHttp.get({ url: Api.GetTaskList, params });
};

export const GetResultList = (params: Recordable) =>
  defHttp.get({ url: Api.GetResultList, params });

export const GetResultField = (id: number) => defHttp.get({ url: Api.GetResultField + '/' + id });

export const GetInspectServerList = (params: Recordable) =>
  defHttp.get({ url: Api.GetInspectServerList, params });

export const InspectServer = (data: Recordable) => defHttp.post({ url: Api.InspectServer, data });
export const GetResultFieldName = (params: Recordable) =>
  defHttp.get({ url: Api.GetResultFieldName, params });
export const GetResultSummaryStatus = (params: Recordable) =>
  defHttp.get({ url: Api.GetResultSummaryStatus, params });

export const GetResultSummaryField = (params: Recordable) =>
  defHttp.get({ url: Api.GetResultSummaryField, params });

export const GetTaskNames = () => defHttp.get({ url: Api.GetTaskNames });
