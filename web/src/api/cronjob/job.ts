import { defHttp } from '@/utils/http/axios';

enum Api {
  GetJobList = '/v1/cronjob/job/list',
  Job = '/v1/cronjob/job',
  JobExecLog = '/v1/cronjob/job/exec',
}

export const getJobList = (params: Recordable) => {
  return defHttp.get({ url: Api.GetJobList, params });
};
export const deleteJob = (id: number | string) => defHttp.delete({ url: Api.Job + '/' + id });
export const getJob = (id: number | string) => defHttp.get({ url: Api.Job + '/' + id });

export const updateJob = (id: number | string, data) =>
  defHttp.put({ url: Api.Job + '/' + id, data });

export const createJob = (data: Recordable) => {
  return defHttp.post({ url: Api.Job, data });
};

export const execJob = (id: number | string) => {
  return defHttp.post({ url: Api.Job + '/' + id + '/exec' });
};

export const getJobExecRecordListApi = (params: Recordable) => {
  return defHttp.get({ url: Api.Job + '/exec-record/list', params });
};

export const getExecRecordLog = (id: string | number, params) => {
  return defHttp.get({ url: Api.JobExecLog + '/' + id + '/logs', params });
};
