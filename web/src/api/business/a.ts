import { defHttp } from '@/utils/http/axios';

export enum Api {
  Composite = '/v1/price/composite',
  Owners = '/v1/price/composite/owners',
  Bizs = '/v1/price/composite/bizs',
  UtilizationOverview = '/v1/price/composite/bizReports',
  Export = '/v1/price/composite/reports/export',
  // eslint-disable-next-line @typescript-eslint/no-duplicate-enum-values
  BizReports = '/v1/price/composite/bizReports',
  Reports = '/v1/price/composite/reports',
  // 节点配置
  NODE_CONFIGS = '/v1/price/node-configs',
  NODE_CONFIG_BY_ROOM = '/v1/price/node-configs/room',
  NODE_CONFIGS_IMPORT = '/v1/price/node-configs/import',
  NODE_CONFIGS_BILLING_TYPES = '/v1/price/node-configs/billing-types',
  NODE_CONFIGS_ROOM_TYPES = '/v1/price/node-configs/room-types',
  NODE_CONFIGS_ORIGINS = '/v1/price/node-configs/origins',
  // 考核规则
  AssessmentRules = '/v1/price/assessment-rules',
  ImportAssessmentRules = '/v1/price/assessment-rules/import',
  ExportBizReports = '/v1/price/composite/bizReports/export',
  GetBizUtilRateList = '/v1/price/composite/biz-util/list',
  EditBizUtilRate = '/v1/price/composite/biz-util',
  ExportBizUtilRate = '/v1/price/composite/biz-util/export',
  // 历史利用率考核
  CompositeHistory = '/v1/price/composite/reports/history',
  // 评分接口
  Assessment = '/v1/price/composite/assessment',
  // 导出节点月利用率
  ExportNodeMonthRate = '/v1/price/composite/month-util/export',
  // 考核自动打分
  AssessmentAutoScore = '/v1/price/composite/auto-score',
  // 奖惩明细
  RewardPunishmentDetails = '/v1/price/composite/assessment/reward-punishment-details',
  // 奖惩历史方案查询
  RewardPunishmentHistoryPlans = '/v1/price/composite/assessment/plan-history',
  // 奖惩明显导出、
  RewardPunishmentDetailsExport = '/v1/price/composite/assessment/reward-punishment-details/export',
}

// category
export const getComposite = (params: Recordable) => {
  return defHttp.get({ url: Api.Composite + '/reports', params });
};

export const getOwners = () => {
  return defHttp.get({ url: Api.Owners });
};
export const getBizs = (params: Recordable) => {
  return defHttp.get({ url: Api.Bizs, params });
};
export const getUtilizationOverview = (params: Recordable) => {
  return defHttp.get({ url: Api.UtilizationOverview, params });
};

export const getExportData = (params: Recordable) => {
  return defHttp.post({ url: Api.Export, params });
};

export const EditBizReports = (params: Recordable) => {
  return defHttp.put({ url: Api.BizReports, params });
};
export const EditBizUtilRate = (params: Recordable) => {
  return defHttp.put({ url: Api.EditBizUtilRate, params });
};
export const editReports = (params: Recordable) => {
  return defHttp.put({ url: Api.Reports, params });
};

// 获取节点配置列表
export function getNodeConfigs(params?: any) {
  return defHttp.get({ url: Api.NODE_CONFIGS, params });
}

// 创建节点配置
export function createNodeConfig(params: any) {
  return defHttp.post({ url: Api.NODE_CONFIGS, params });
}

// 更新节点配置
export function updateNodeConfig(id: number, params: any) {
  return defHttp.put({ url: `${Api.NODE_CONFIGS}/${id}`, params });
}

// 批量导入节点配置
export function importNodeConfigs(params: any) {
  return defHttp.post({ url: Api.NODE_CONFIGS_IMPORT, params });
}

// 根据节点编号获取配置
export function getNodeConfigByRoomNo(roomNo: string) {
  return defHttp.get({ url: `${Api.NODE_CONFIG_BY_ROOM}/${roomNo}` });
}

// 获取计费方式列表
export function getBillingTypes() {
  return defHttp.get({ url: Api.NODE_CONFIGS_BILLING_TYPES });
}

// 获取业务利用率列表
export function getBizUtilRateList(params?: any) {
  return defHttp.get({ url: Api.GetBizUtilRateList, params });
}

// 获取机房类型列表
export function getRoomTypes() {
  return defHttp.get({ url: Api.NODE_CONFIGS_ROOM_TYPES });
}

// 获取机房归属列表
export function getOrigins() {
  return defHttp.get({ url: Api.NODE_CONFIGS_ORIGINS });
}

//获取考核规则列表
export function getAssessmentRules(params?: any) {
  return defHttp.get({ url: Api.AssessmentRules, params });
}

//更新考核规则
export function updateAssessmentRule(id: number, params: any) {
  return defHttp.put({
    url: `${Api.AssessmentRules}/${id}`,
    data: params,
  });
}

// 批量导入考核规则
export function importAssessmentRules(params: any) {
  return defHttp.post({ url: Api.ImportAssessmentRules, params });
}

// 获取历史利用率考核数据
export function getCompositeHistory(params: any) {
  return defHttp.get({ url: Api.CompositeHistory, params });
}

// 获取评分数据
export function getAssessment(params: any) {
  return defHttp.get({ url: Api.Assessment, params });
}

// 考核自动打分
export function assessmentAutoScore(params: any) {
  return defHttp.post({ url: Api.AssessmentAutoScore, params });
}

// 获取奖惩明细数据
export function getRewardPunishmentDetails(params: any) {
  return defHttp.get({ url: Api.RewardPunishmentDetails, params });
}

// 获取奖惩历史方案
export function getRewardPunishmentHistoryPlans(params: any) {
  return defHttp.get({ url: Api.RewardPunishmentHistoryPlans, params });
}

// 奖惩明细导出
export function exportRewardPunishmentDetails(data: any) {
  return defHttp.post({
    url: Api.RewardPunishmentDetailsExport,
    data,
  });
}
