import { defHttp } from '@/utils/http/axios';
import {
  LoginParams,
  LoginResultModel,
  GetUserInfoModel,
  starPortalLoginParams,
} from './model/userModel';

import { ErrorMessageMode } from '#/axios';

enum Api {
  Login = '/v1/system/login',
  GetStarPortalLoginUrl = '/v1/system/star-portal/login-url',
  StarPortalLogin = '/v1/system/star-portal/login',
  RefreshToken = '/v1/system/refreshToken',
  Logout = '/v1/system/logout',
  GetUserInfo = '/v1/system/userInfo',
  GetPermCode = '/v1/system/userPermCode',
  TestRetry = '/testRetry',
}

/**
 * @description: user login api
 */
export function loginApi(params: LoginParams, mode: ErrorMessageMode = 'modal') {
  return defHttp.post<LoginResultModel>(
    {
      url: Api.Login,
      params,
    },
    {
      errorMessageMode: mode,
    },
  );
}
export function refreshTokenApi(mode: ErrorMessageMode = 'message') {
  return defHttp.post<LoginResultModel>(
    {
      url: Api.RefreshToken,
    },
    {
      errorMessageMode: mode,
    },
  );
}
/**
 * @description: getUserInfo
 */
export function getUserInfo() {
  return defHttp.get<GetUserInfoModel>({ url: Api.GetUserInfo }, { errorMessageMode: 'none' });
}

export function getPermCode() {
  return defHttp.get<string[]>({ url: Api.GetPermCode });
}

export function doLogout() {
  return defHttp.post({ url: Api.Logout });
}

export function testRetry() {
  return defHttp.get(
    { url: Api.TestRetry },
    {
      retryRequest: {
        isOpenRetry: true,
        count: 5,
        waitTime: 1000,
      },
    },
  );
}

export function getStarPortalLoginUrl(params: Recordable) {
  return defHttp.get({ url: Api.GetStarPortalLoginUrl, params });
}
export function starPortalLogin(params: starPortalLoginParams, mode: ErrorMessageMode = 'modal') {
  return defHttp.post<LoginResultModel>(
    {
      url: Api.StarPortalLogin,
      params,
    },
    {
      errorMessageMode: mode,
    },
  );
}
