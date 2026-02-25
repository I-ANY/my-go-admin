export enum PermissionCodeEnum {
  BUSINESS_INSPECT_TASK_RESULT = 'business:inspect:task:result',
  BUSINESS_INSPECT_SERVER = 'business:server:inspect',
  BUSINESS_SERVER_INSPECT_RECORD = 'business:server:inspect:record',
  BUSINESS_INSPECT_RESULT_EXPORT = 'business:inspect:result:export',
  BUSINESS_SERVER_INSPECT_RECORD_EXPORT = 'business:server:inspect:record:export',
}
export enum LAPermissionCodeEnum {
  TRAFFIC_DATA_EXPORT = 'business:la:traffic:data:export',
  SLA_DETAIL_EXPORT = 'business:la:sla:detail:export',
  TRAFFIC_DEVICE_INFO_EXPORT = 'business:la:traffic:device:info:export',
}
export enum KPermissionCodeEnum {
  BUSINESS_K_HDD_DAILYPEAK_EXPORT = 'business:k:hdd:dailyPeak:export',
  BUSINESS_K_HDD_DAILYPEAK_DETAIL_EXPORT = 'business:k:hdd:dailyPeak:detail:export',
  BUSINESS_K_HDD_5MIN_SUMMARY_EXPORT = 'business:k:hdd:5min:summary:export',
  BUSINESS_K_HDD_DIFF_BATCH_CONFIRM = 'business:k:hdd:diff:batch:confirm',
  BUSINESS_K_HDD_DIFF_CONFIRM = 'business:k:hdd:diff:confirm',
  BUSINESS_K_HDD_DIFF_EXPORT = 'business:k:hdd:diff:export',
}
export enum SystemPermissionCodeEnum {
  RESOURCE_ADD = 'system:resource:add',
  RESOURCE_EDIT = 'system:resource:edit',
  RESOURCE_VIEW = 'system:resource:view',
  RESOURCE_DELETE = 'system:resource:delete',
  ROLE_RESOURCE_PERMISSION_CONFIG = 'system:role:resource:permission:config',
  ROLE_BUSINESS_PERMISSION_CONFIG = 'system:role:business:permission:config',
}
export enum NetworkPermissionCodeEnum {
  NETWORK_ROOM_PEAK_EDIT = 'network:room:peak:edit',
  NETWORK_ROOM_PEAK_IMPROVE_95_PREDICTION = 'network:room:peak:improve:95:prediction',
  NETWORK_ROOM_PEAK_SPEED_LIMIT_CONFIGURATION = 'network:room:peak:speed:limit:configuration',
  NETWORK_ROOM_PEAK_SPEED_LIMIT_RECORD = 'network:room:peak:speed:limit:record',
  SINGLE_PORT_95_PEAK_SHAVING_PORT = 'network:room:peak:single:port:95:port',
  BIZ_DSCP_MODIFY = 'network:biz:dscp:modify',
  BIZ_DSCP_EXPORT = 'network:biz:dscp:export',
  SERVER_DSCP_MODIFY = 'network:server:dscp:modify',
  SERVER_DSCP_EXPORT = 'network:server:dscp:export',
  SERVER_DSCP_DETAIL = 'network:server:dscp:detail',
  SPEED_LIMIT_JOB_ADD = 'network:speed:limit:job:add',
  SPEED_LIMIT_JOB_EDIT = 'network:speed:limit:job:edit',
  SPEED_LIMIT_JOB_VIEW_EXEC_RECORD = 'network:speed:limit:job:view:exec:record',
}

export enum APermissionCodeEnum {
  BUSINESS_A_UTILIZATION_RATE_OVERVIEW_EXPORT = 'business:a:utilization:rate:overview:export',
  BUSINESS_A_UTILIZATION_RATE_OVERVIEW_NODE_DESCRIBE_EDIT = 'business:a:utilization:rate:overview:node:describe:edit',
  BUSINESS_A_UTILIZATION_RATE_OVERVIEW_DAY_DESCRIBE_EDIT = 'business:a:utilization:rate:overview:day:describe:edit',
  BUSINESS_A_UTILIZATION_RATE_BIZ_EXPORT = 'business:a:utilization:rate:biz:export',
  BUSINESS_A_UTILIZATION_RATE_BIZ_NODE_DESCRIBE_EDIT = 'business:a:utilization:rate:biz:node:describe:edit',
  BUSINESS_A_UTILIZATION_RATE_BIZ_DAY_DESCRIBE_EDIT = 'business:a:utilization:rate:biz:day:describe:edit',
  BUSINESS_A_NODE_MONTH_EXPORT = 'business:a:node:month:export',
  BUSINESS_A_UTILIZATION_RATE_ASSESSMENT_HISTORY_EDIT = 'business:a:utilization:rate:assessment:history:edit',
  BUSINESS_A_UTILIZATION_RATE_ASSESSMENT_SCORE = 'business:a:utilization:rate:assessment:score',
}
