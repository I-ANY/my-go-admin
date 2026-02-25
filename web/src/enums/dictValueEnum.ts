// network
export enum Network_SpeedLimitJob_Status {
  ENABLE = 1,
  DISABLE = 0,
}
export enum Network_SpeedLimitJob_JobType {
  // 立即执行
  IMMEDIATE_TASK = 1,
  // 定时执行(一次性)
  ONCE_TASK = 2,
  // 周期任务
  CYCLE_TASK = 3,
}

export enum Network_SpeedLimitJob_LimitType {
  LIMIT = 1,
  UNLIMIT = 2,
}

// 机房限速的操作操作对象
export enum Network_SpeedLimitJob_OperateTargetType {
  SWITCH_PORT = 1, // 交换机端口
  BUSINESS_TAG = 2, //业务标签
}
export enum Network_SpeedLimitJob_ExecuteNow {
  NOT_EXECUTE = 2, //不立即执行
  EXECUTE = 1, // 立即执行
}
