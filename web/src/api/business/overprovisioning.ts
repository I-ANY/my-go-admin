import { defHttp } from '@/utils/http/axios';

export enum Api {
  // Rule Management
  OverProvisioningRule = '/v1/business/over-provisioning/rule',
  // UpdateRule = '/v1/business/over-provisioning/rule',
  // DeleteRule = '/v1/business/over-provisioning/rule',
  ExecuteDetection = '/v1/business/over-provisioning/execute',

  // Detection Records
  GetRecordList = '/v1/business/over-provisioning/list',
  RecordListExport = '/v1/business/over-provisioning/export',

  // 无规则业务
  NoRuleBusiness = '/v1/business/over-provisioning/no-rule-business',
  NoRuleBusinessExport = '/v1/business/over-provisioning/no-rule-business/export',
  UpdateNoRuleBusinessRemark = '/v1/business/over-provisioning/no-rule-business/update',
}

// --- Rule Management ---

export const getRuleList = (params: Recordable) => {
  return defHttp.get({ url: Api.OverProvisioningRule, params });
};

export const createRule = (data: Recordable) => {
  return defHttp.post({
    url: Api.OverProvisioningRule,
    data,
    headers: {
      'Content-Type': 'application/json',
    },
  });
};

export const updateRule = (id: number | string, data: Recordable) => {
  return defHttp.put({
    url: `${Api.OverProvisioningRule}/${id}`,
    data,
    headers: {
      'Content-Type': 'application/json',
    },
  });
};

export const deleteRule = (id: number | string) => {
  return defHttp.delete({
    url: `${Api.OverProvisioningRule}/${id}`,
  });
};

// 执行超配检测
export const executeDetection = (execType) => {
  return defHttp.post({
    url: Api.ExecuteDetection,
    data: { execType },
    headers: {
      'Content-Type': 'application/json',
    },
  });
};

// --- Detection Records ---

export const getRecordList = (data: Recordable) => {
  return defHttp.post({ url: Api.GetRecordList, data });
};

export const exportRecordList = (data: Recordable) => {
  return defHttp.post({
    url: Api.RecordListExport,
    data,
    responseType: 'blob',
  });
};

// 无规则业务
export const getNoRuleBusiness = (params: Recordable) => {
  return defHttp.get({ url: Api.NoRuleBusiness, params });
};

// 更新无规则业务备注
export const updateNoRuleBusinessRemark = (data: Recordable) => {
  return defHttp.post({ url: Api.UpdateNoRuleBusinessRemark, data });
};
