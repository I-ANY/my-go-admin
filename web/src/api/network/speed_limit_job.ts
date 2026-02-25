import { defHttp } from '@/utils/http/axios';

export const Api = {
  GetSpeedLimitJobList: () => '/v1/network/speed-limit/job/list',
  GetEcdnRoomList: () => '/v1/network/speed-limit/ecdn-room/list',
  ImportPeakSchema: () => `/v1/network/speed-limit/import/peak-schema`,
  GetRoomSwitches: () => `/v1/network/speed-limit/room/switches`,
  GetSwitchBusinessTags: () => '/v1/network/speed-limit/switch/business-tags',
  CreateSpeedLimitJob: () => '/v1/network/speed-limit/job',
  UpdateSpeedLimitJob: (id: number) => `/v1/network/speed-limit/job/${id}`,
  GetSpeedLimitJobExecRecordList: (id: number) => `/v1/network/speed-limit/${id}/exec-record/list`,
};

export const GetSpeedLimitJobList = (params: Recordable) => {
  return defHttp.get({ url: Api.GetSpeedLimitJobList(), params });
};
export const GetEcdnRoomList = (params: Recordable) => {
  return defHttp.get({ url: Api.GetEcdnRoomList(), params });
};
export const ImportPeakSchema = (params: Recordable) => {
  return defHttp.get({ url: Api.ImportPeakSchema(), params });
};
export const GetRoomSwitches = (params: Recordable) => {
  return defHttp.get({ url: Api.GetRoomSwitches(), params });
};
export const GetSwitchBusinessTags = (params: Recordable) => {
  return defHttp.get({ url: Api.GetSwitchBusinessTags(), params });
};
export const CreateSpeedLimitJob = (data: Recordable) => {
  return defHttp.post({ url: Api.CreateSpeedLimitJob(), data });
};
export const UpdateSpeedLimitJob = (id: number, data: Recordable) => {
  return defHttp.put({ url: Api.UpdateSpeedLimitJob(id), data });
};

export const GetSpeedLimitJobExecRecordList = (id: number, params: Recordable) => {
  return defHttp.get({ url: Api.GetSpeedLimitJobExecRecordList(id), params });
};
