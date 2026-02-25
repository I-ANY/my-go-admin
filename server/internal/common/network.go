package common

import (
	"biz-auto-api/internal/models"
	"biz-auto-api/pkg/tools"
	"fmt"
	"github.com/pkg/errors"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type PeakPlanInfo struct {
	StartDate          string  // 开始日期 (月.日)
	EndDate            string  // 结束日期 (月.日)
	PortCount          int64   // 端口数
	StartTime          string  // 开始时间
	EndTime            string  // 结束时间
	PeakValueGbps      float64 // 峰值
	ReductionValueGbps float64 // 削峰值
	MinSpeedLimitGbps  float64 // 保底限速
}

func ExtractPeakPlanInfo(input string) (*PeakPlanInfo, error) {
	peakPlanRegex := regexp.MustCompile(
		`^打峰方案-(\d{1,2}\.\d{1,2})-(\d{1,2}\.\d{1,2})日，(\d+)口打峰（(\d{1,2}:\d{2})-(\d{1,2}:\d{2})），峰值(\d+(?:\.\d+)?)G削峰(\d+(?:\.\d+)?)G，保底限速(\d+(?:\.\d+)?)G.*`)

	match := peakPlanRegex.FindStringSubmatch(input)
	if match == nil {
		return nil, errors.Errorf("input string does not match the expected format")
	}

	// 解析端口数
	ports, err := strconv.ParseInt(match[3], 10, 64)
	if err != nil {
		return nil, errors.Errorf("failed to parse ports: %w", err)
	}

	// 解析峰值
	peakValue, err := strconv.ParseFloat(match[6], 64)
	if err != nil {
		return nil, errors.Errorf("failed to parse peak value: %w", err)
	}

	// 解析削峰值
	reductionValue, err := strconv.ParseFloat(match[7], 64)
	if err != nil {
		return nil, errors.Errorf("failed to parse reduction value: %w", err)
	}

	// 解析保底限速
	minSpeedLimit, err := strconv.ParseFloat(match[8], 64)
	if err != nil {
		return nil, errors.Errorf("failed to parse min speed limit: %w", err)
	}

	return &PeakPlanInfo{
		StartDate:          FillZero(match[1], ".", "-"),
		EndDate:            FillZero(match[2], ".", "-"),
		PortCount:          ports,
		StartTime:          FillZero(match[4], ":", ":"),
		EndTime:            FillZero(match[5], ":", ":"),
		PeakValueGbps:      peakValue,
		ReductionValueGbps: reductionValue,
		MinSpeedLimitGbps:  minSpeedLimit,
	}, nil
}
func FillZero(str string, separator, newSeparator string) string {
	ss := strings.Split(str, separator)
	for i, s := range ss {
		if len(s) < 2 {
			ss[i] = "0" + s
		}
	}
	return strings.Join(ss, newSeparator)
}
func ValidateNetworkSpeedLimitJob(job *models.NetworkSpeedLimitJob) error {
	var (
		jobType  = job.JobType
		strategy = job.Strategies
		l, _     = time.LoadLocation("Asia/Shanghai")
	)
	if job.EcdnRoomId == nil || *job.EcdnRoomId == 0 {
		return errors.New("ecdn机房id不能为空")
	}
	// 校验限速类型是否合法
	validateLimit := func(limitType *models.LimitType) error {
		if limitType == nil {
			return errors.New("限速类型不能为空")
		}
		switch tools.ToValue(limitType.Type) {
		case models.NetworkSpeedLimitJobLimitTypeLimit:
			if limitType.LimitValue == nil {
				return errors.New("限速值不能为空")
			}
			if *limitType.LimitValue < 0 {
				return errors.New("限速值不能小于0")
			}
		case models.NetworkSpeedLimitJobLimitTypeUnLimit:
		default:
			return errors.New("限速类型错误")
		}
		return nil
	}

	// 校验时间范围内是否有重复的执行时间
	dateValidate := func(dates []*models.DateRange) error {
		allExecuteTimesSeconds := make([]int64, 0)
		allExecuteTimes, err := GetAllExecuteTimes(dates, l)
		if err != nil {
			return errors.New("获取所有执行时间失败")
		}
		// 转换成秒级的时间
		allExecuteTimesSeconds = tools.GetSlice(allExecuteTimes, func(t time.Time) int64 { return t.Unix() })
		if len(allExecuteTimesSeconds) != len(tools.RemoveDuplication(allExecuteTimesSeconds)) {
			return errors.New("存在重复的执行时间")
		}
		return nil
	}

	if tools.ToValue(job.RetryCount) < 0 || tools.ToValue(job.RetryCount) > 10 {
		return errors.New("任务重试次数不在合理范围内(0-10)")
	}
	if tools.ToValue(job.MaxExecuteDelayMinutes) < 1 || tools.ToValue(job.MaxExecuteDelayMinutes) > 45 {
		return errors.New("任务最大执行延迟时间不在合理范围内(1-45)")
	}
	if jobType == nil {
		return errors.New("任务类型不能为空")
	}
	if strategy == nil {
		return errors.New("执行配置不能为空")
	}
	if strategy.LimitTarget == nil {
		return errors.New("操作对象不能为空")
	}
	if len(strategy.LimitTarget.SwitchIds) != 1 {
		return errors.New("交换机数量不为1无法限速")
	}
	switch tools.ToValue(strategy.LimitTarget.OperateType) {
	case models.NetworkSpeedLimitJobTargetTypeSwitchPort: //端口
		if len(tools.ToValue(strategy.LimitTarget.SwitchPortRange)) == 0 && len(tools.ToValue(strategy.LimitTarget.SwitchPort)) == 0 {
			return errors.New("交换机端口和交换机端口范围不能同时为空")
		}
	case models.NetworkSpeedLimitJobTargetTypeBusinessPort:
		if len(strategy.LimitTarget.BusinessTag) == 0 {
			return errors.New("业务标签不能为空")
		}
	default:
		return errors.New("未知的操作对象")
	}

	switch tools.ToValue(jobType) {
	case models.NetworkJobTypeImmediateTask:
		if strategy.ImmediateTaskStrategy == nil || strategy.ImmediateTaskStrategy.LimitType == nil {
			return errors.New("执行配置不能为空")
		}
		err := validateLimit(strategy.ImmediateTaskStrategy.LimitType) // Validate limit value
		if err != nil {
			return err
		}
	case models.NetworkJobTypeOnceTask:
		if strategy.OnceTaskStrategy == nil || strategy.OnceTaskStrategy.LimitType == nil {
			return errors.New("执行配置不能为空")
		}
		err := validateLimit(strategy.OnceTaskStrategy.LimitType)
		if err != nil {
			return err
		}
		if strategy.OnceTaskStrategy.ExecuteTime == nil {
			return errors.New("执行时间不能为空")
		}
		_, err = time.ParseInLocation(time.DateTime, *strategy.OnceTaskStrategy.ExecuteTime, l)
		if err != nil {
			return errors.New("执行时间格式错误")
		}
	case models.NetworkJobTypeCycleTask:
		if strategy.CycleTaskStrategy == nil {
			return errors.New("执行配置不能为空")
		}
		if len(strategy.CycleTaskStrategy.Dates) == 0 {
			return errors.New("执行配置不能为空")
		}
		for _, dateRange := range strategy.CycleTaskStrategy.Dates {
			_, err := time.ParseInLocation(time.DateOnly, *dateRange.StartDate, l)
			if err != nil {
				return errors.New("开始日期格式错误")
			}
			_, err = time.ParseInLocation(time.DateOnly, *dateRange.EndDate, l)
			if err != nil {
				return errors.New("结束日期格式错误")
			}
			if len(dateRange.Times) == 0 {
				return errors.New("操作不能为空")
			}
			for _, timeRange := range dateRange.Times {
				_, err = time.ParseInLocation(time.TimeOnly, *timeRange.Time, l)
				if err != nil {
					return errors.New("时间格式错误")
				}
				err = validateLimit(timeRange.LimitType)
				if err != nil {
					return err
				}
			}
		}
		err := dateValidate(strategy.CycleTaskStrategy.Dates)
		if err != nil {
			return err
		}
	default:
		return errors.New("任务类型错误")
	}
	return nil
}

const (
	SpeedLimitJobExecStatusRunning = "running"
	SpeedLimitJobExecStatusDone    = "done"
	SpeedLimitJobExecStatusFailed  = "failed"
)

// GenerateNetworkSpeedLimitJobExecuteStatusKey 生成任务执行状态的key
func GenerateNetworkSpeedLimitJobExecuteStatusKey(jobId int64, executeTime time.Time) string {
	return fmt.Sprintf("network:job:execute:status:%d:%s", jobId, executeTime.Format("20060102150405"))
}

// GenerateNetworkSpeedLimitJobExecuteHeartbeatKey 生成任务执行心跳的key
func GenerateNetworkSpeedLimitJobExecuteHeartbeatKey(jobId int64, executeTime time.Time) string {
	return fmt.Sprintf("network:job:execute:heartbeat:%d:%s", jobId, executeTime.Format("20060102150405"))
}

func GetAllExecuteTimes(dates []*models.DateRange, l *time.Location) ([]time.Time, error) {
	var allExecuteTimes []time.Time
	for _, dateRange := range dates {
		startDate, err := time.ParseInLocation(time.DateOnly, *dateRange.StartDate, l)
		if err != nil {
			return nil, errors.New("开始日期格式错误")
		}
		endDate, err := time.ParseInLocation(time.DateOnly, *dateRange.EndDate, l)
		if err != nil {
			return nil, errors.New("结束日期格式错误")
		}
		for date := startDate; !date.After(endDate); date = date.AddDate(0, 0, 1) {
			for _, t := range dateRange.Times {
				timeStr := fmt.Sprintf("%s %s", date.Format(time.DateOnly), tools.ToValue(t.Time))
				executeTime, err := time.ParseInLocation(time.DateTime, timeStr, l)
				if err != nil {
					return nil, errors.New("执行时间格式错误")
				}
				executeTime = time.Unix(executeTime.Unix()/60*60, 0).In(l)
				allExecuteTimes = append(allExecuteTimes, executeTime)
			}
		}
	}
	return allExecuteTimes, nil
}
