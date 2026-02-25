import { defHttp } from '@/utils/http/axios';

export const Api = {
  getPeakRoomList: () => '/v1/network/room/peak/list',
  getPeakBusiness: () => '/v1/network/room/peak/business',
  updatePeakRoom: (id: number) => `/v1/network/room/peak/${id}`,
  updatePeakRoomSpeedLimitConfig: (id: number) => `/v1/network/room/${id}/speed-limit/config`,
  getPeakRoomTag: (id: number) => `/v1/network/room/${id}/tags`,
  getSpeedLimitRecordList: () => '/v1/network/room/speed-limit/record/list',
  getRoomSwitch: (id: number) => `/v1/network/room/${id}/switches`,
  getSinglePort95: (id: number) => `/v1/network/room/peak/${id}/port`,
};

export const getPeakRoomList = (params: Recordable) => {
  return defHttp.get({ url: Api.getPeakRoomList(), params });
};
export const getPeakBusiness = (params: Recordable) => {
  return defHttp.get({ url: Api.getPeakBusiness(), params });
};
export const updatePeakRoom = (id: number, data: Recordable) => {
  return defHttp.put({ url: Api.updatePeakRoom(id), data });
};
export const updatePeakRoomSpeedLimitConfig = (id: number, data: Recordable) => {
  return defHttp.put({ url: Api.updatePeakRoomSpeedLimitConfig(id), data });
};

export const getPeakRoomTag = (id: number) => {
  return defHttp.get({ url: Api.getPeakRoomTag(id) });
};

export const getSpeedLimitRecordList = (params: Recordable) => {
  return defHttp.get({ url: Api.getSpeedLimitRecordList(), params });
};
export const getRoomSwitch = (id: number, params: Recordable) => {
  return defHttp.get({ url: Api.getRoomSwitch(id), params });
};
export const getSinglePort95 = (id: number, params: Recordable) => {
  return defHttp.get({ url: Api.getSinglePort95(id), params });
};
