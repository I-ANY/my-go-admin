package tools

import (
	"github.com/pkg/errors"
	"time"
)

func StrToTime(timeStr, format string) (time.Time, error) {
	if len(timeStr) == 0 || len(format) == 0 {
		return time.Time{}, errors.New("timeStr or format can't be empty")
	}
	l, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return time.Time{}, errors.Wrap(err, "load location failed")
	}
	t, err := time.ParseInLocation(format, timeStr, l)
	if err != nil {
		return time.Time{}, errors.Wrap(err, "parse time failed")
	}

	return t, nil
}

// CalMonthFiveMinuteAnd95PointsExact 计算本月的5分钟点数和95点数，优先使用billingStartAt作为开始时间，如果billingStartAt非本月时间，则使用本月开始时间
func CalMonthFiveMinuteAnd95PointsExact(billingStartAt *time.Time) (totalCount, maxPeakShavingPoint int64) {
	now := time.Now()
	startAt, endAt := GetMonthBoundaries(now.Year(), now.Month(), now.Location())
	// 如果账单开始时间在当月之前，则从当月开始计算
	if billingStartAt == nil || billingStartAt.Before(startAt) {
		return CalFiveMinuteAnd95PointsExact(startAt, endAt)
	} else {
		return CalFiveMinuteAnd95PointsExact(*billingStartAt, endAt)
	}
}

// GetMonthBoundaries 获取指定月份的开始时间和结束时间
func GetMonthBoundaries(year int, month time.Month, loc *time.Location) (start, end time.Time) {
	if loc == nil {
		loc = time.Local
	}
	// 月份开始时间
	start = time.Date(year, month, 1, 0, 0, 0, 0, loc)
	// 月份结束时间（最后一天的最后一秒），月+1 再减1秒
	end = time.Date(year, month+1, 1, 0, 0, 0, 0, loc).Add(-time.Second)
	return start, end
}

// CalFiveMinuteAnd95PointsExact 获取严格在时间范围内的5分钟点数和95点数
func CalFiveMinuteAnd95PointsExact(start, end time.Time) (totalCount, maxPeakShavingPoint int64) {
	totalCount = CalFiveMinutePointsExact(start, end)
	maxPeakShavingPoint = int64(float64(totalCount) * 0.05)
	return
}

// CalFiveMinutePointsExact 计算严格在时间范围内的整5分钟点数
func CalFiveMinutePointsExact(start, end time.Time) (totalCount int64) {
	// 找到范围内第一个整5分钟点
	firstPoint := findNextFiveMinutePoint(start)

	// 找到范围内最后一个整5分钟点
	lastPoint := findPrevFiveMinutePoint(end)

	// 如果没有点在范围内，返回0
	if firstPoint.After(lastPoint) {
		return
	}

	// 计算点数
	secondsDiff := lastPoint.Unix() - firstPoint.Unix()
	totalCount = secondsDiff/300 + 1
	return totalCount
}

// findNextFiveMinutePoint 找到大于等于给定时间的下一个整5分钟点
func findNextFiveMinutePoint(t time.Time) time.Time {
	minute := t.Minute()
	roundedMinute := ((minute + 4) / 5) * 5 // 向上取整到最近的5分钟

	// 如果超过了55分钟，需要进位到下一小时
	if roundedMinute >= 60 {
		t = t.Add(time.Hour)
		roundedMinute = 0
	}

	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), roundedMinute, 0, 0, t.Location())
}

// findPrevFiveMinutePoint 找到小于等于给定时间的上一个整5分钟点
func findPrevFiveMinutePoint(t time.Time) time.Time {
	minute := t.Minute()
	roundedMinute := (minute / 5) * 5 // 向下取整到最近的5分钟

	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), roundedMinute, 0, 0, t.Location())
}

// FloorTimeBySeconds 将时间向下取整到指定时间间隔
func FloorTimeBySeconds(t time.Time, stepSeconds int) time.Time {
	// 获取时间戳（秒）
	seconds := t.Unix()

	// 向下取整到最近的步长时间点
	flooredSeconds := (seconds / int64(stepSeconds)) * int64(stepSeconds)

	// 返回取整后的时间
	return time.Unix(flooredSeconds, 0).In(t.Location())
}

func IsTimeInDateRangeAtSpecificTime(startDate, endDate, timePoint string, minTime, maxTime time.Time, l *time.Location) (bool, time.Time, error) {
	// 解析开始日期
	startDateTime, err := time.ParseInLocation(time.DateTime, startDate+" "+timePoint, l)
	if err != nil {
		return false, time.Time{}, errors.Wrap(err, "parse start date failed")
	}
	startDateTime = time.Unix((startDateTime.Unix()/60)*60, 0).In(l)

	// 解析结束日期
	endDateTime, err := time.ParseInLocation(time.DateTime, endDate+" "+timePoint, l)
	if err != nil {
		return false, time.Time{}, errors.Wrap(err, "parse end date failed")
	}
	endDateTime = time.Unix((endDateTime.Unix()/60)*60, 0).In(l)

	for t := startDateTime; !t.After(endDateTime); t = t.AddDate(0, 0, 1) {
		// 创建一个取整后的时间用于比较，不要修改循环变量t
		if !t.Before(minTime) && !t.After(maxTime) {
			return true, t, nil
		}
	}
	return false, time.Time{}, nil
}
