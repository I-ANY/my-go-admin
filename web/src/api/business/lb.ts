import { defHttp } from '@/utils/http/axios';

export enum Api {
  GetBillList = '/v1/business/lb/bill/list',
  GetBillSummaryList = '/v1/business/lb/bill/listsummary',
  ExportBill = '/v1/business/lb/bill/export',
  ExportBillSummary = '/v1/business/lb/bill/listsummaryexport',
}

export const GetBillList = (params: Recordable) => {
  return defHttp.get({ url: Api.GetBillList, params });
};

export const GetBillSummaryList = (params: Recordable) => {
  return defHttp.get({ url: Api.GetBillSummaryList, params });
};
