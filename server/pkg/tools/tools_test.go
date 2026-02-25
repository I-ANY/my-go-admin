package tools

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_GetNoRepeatSlice(t *testing.T) {
	s := []string{"1", "0", "1", "2", "3", "0", "4", "5", "6", "7", "8", "9", "0", "10"}
	fmt.Println(GetNoRepeatSlice(s, func(e string) string {
		return e
	}))

}

func TestIsTimeInDateRangeAtSpecificTime(t *testing.T) {
	// 设置测试时区
	loc, err := time.LoadLocation("Asia/Shanghai")
	assert.NoError(t, err)

	// 定义测试用例结构
	type testCase struct {
		name        string
		startDate   string
		endDate     string
		timePoint   string
		minTime     time.Time
		maxTime     time.Time
		expected    bool
		expectError bool
	}

	// 测试用例
	tests := []testCase{
		{
			name:        "valid range with match",
			startDate:   "2023-01-01",
			endDate:     "2023-01-03",
			timePoint:   "12:00:00",
			minTime:     time.Date(2023, 1, 1, 10, 0, 0, 0, loc),
			maxTime:     time.Date(2023, 1, 1, 14, 0, 0, 0, loc),
			expected:    true,
			expectError: false,
		},
		{
			name:        "valid range without match",
			startDate:   "2023-01-01",
			endDate:     "2023-01-03",
			timePoint:   "12:00:00",
			minTime:     time.Date(2023, 1, 1, 15, 0, 0, 0, loc),
			maxTime:     time.Date(2023, 1, 1, 16, 0, 0, 0, loc),
			expected:    false,
			expectError: false,
		},
		{
			name:        "single day match",
			startDate:   "2023-01-01",
			endDate:     "2023-01-01",
			timePoint:   "12:00:00",
			minTime:     time.Date(2023, 1, 1, 12, 0, 0, 0, loc),
			maxTime:     time.Date(2023, 1, 1, 12, 0, 0, 0, loc),
			expected:    true,
			expectError: false,
		},
		{
			name:        "start date after end date",
			startDate:   "2023-01-03",
			endDate:     "2023-01-01",
			timePoint:   "12:00:00",
			minTime:     time.Date(2023, 1, 1, 10, 0, 0, 0, loc),
			maxTime:     time.Date(2023, 1, 1, 14, 0, 0, 0, loc),
			expected:    false,
			expectError: false,
		},
		{
			name:        "invalid start date format",
			startDate:   "invalid-date",
			endDate:     "2023-01-03",
			timePoint:   "12:00:00",
			minTime:     time.Date(2023, 1, 1, 10, 0, 0, 0, loc),
			maxTime:     time.Date(2023, 1, 1, 14, 0, 0, 0, loc),
			expected:    false,
			expectError: true,
		},
		{
			name:        "invalid end date format",
			startDate:   "2023-01-01",
			endDate:     "invalid-date",
			timePoint:   "12:00:00",
			minTime:     time.Date(2023, 1, 1, 10, 0, 0, 0, loc),
			maxTime:     time.Date(2023, 1, 1, 14, 0, 0, 0, loc),
			expected:    false,
			expectError: true,
		},
		{
			name:        "time point equals min time",
			startDate:   "2023-01-01",
			endDate:     "2023-01-01",
			timePoint:   "12:00:00",
			minTime:     time.Date(2023, 1, 1, 12, 0, 0, 0, loc),
			maxTime:     time.Date(2023, 1, 1, 14, 0, 0, 0, loc),
			expected:    true,
			expectError: false,
		},
		{
			name:        "time point equals max time",
			startDate:   "2023-01-01",
			endDate:     "2023-01-01",
			timePoint:   "12:00:00",
			minTime:     time.Date(2023, 1, 1, 10, 0, 0, 0, loc),
			maxTime:     time.Date(2023, 1, 1, 12, 0, 0, 0, loc),
			expected:    true,
			expectError: false,
		},
		{
			name:        "time point before min time",
			startDate:   "2023-01-01",
			endDate:     "2023-01-01",
			timePoint:   "12:00:00",
			minTime:     time.Date(2023, 1, 1, 13, 0, 0, 0, loc),
			maxTime:     time.Date(2023, 1, 1, 14, 0, 0, 0, loc),
			expected:    false,
			expectError: false,
		},
		{
			name:        "time point after max time",
			startDate:   "2023-01-01",
			endDate:     "2023-01-01",
			timePoint:   "23:59:00",
			minTime:     time.Date(2023, 1, 1, 23, 59, 0, 0, loc),
			maxTime:     time.Date(2023, 1, 2, 0, 5, 0, 0, loc),
			expected:    true,
			expectError: false,
		},
	}

	// 执行测试用例
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, _, err := IsTimeInDateRangeAtSpecificTime(
				tc.startDate,
				tc.endDate,
				tc.timePoint,
				tc.minTime,
				tc.maxTime,
				loc,
			)

			if tc.expectError {
				assert.Error(t, err)
				assert.False(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}

func TestIsTimeInDateRangeAtSpecificTime_Timezone(t *testing.T) {
	//// 测试不同时区的情况
	//utcLoc, err := time.LoadLocation("UTC")
	//assert.NoError(t, err)

	chinaLoc, err := time.LoadLocation("Asia/Shanghai")
	assert.NoError(t, err)

	// UTC时间: 2023-01-01 16:00:00 (对应北京时间 2023-01-02 00:00:00)
	startDate := "2026-01-01"
	endDate := "2026-01-01"
	timePoint := "16:00:00"

	minTime := time.Date(2026, 1, 1, 16, 0, 0, 0, chinaLoc)
	maxTime := time.Date(2026, 1, 1, 16, 5, 0, 0, chinaLoc)

	result, _, err := IsTimeInDateRangeAtSpecificTime(
		startDate,
		endDate,
		timePoint,
		minTime,
		maxTime,
		chinaLoc,
	)

	assert.NoError(t, err)
	assert.True(t, result)
}

func TestIsTimeInDateRangeAtSpecificTime_EdgeCases(t *testing.T) {
	loc, err := time.LoadLocation("Asia/Shanghai")
	assert.NoError(t, err)

	// 测试跨越月份的日期范围
	result, _, err := IsTimeInDateRangeAtSpecificTime(
		"2023-01-30",
		"2023-02-02",
		"15:30:00",
		time.Date(2023, 2, 1, 14, 0, 0, 0, loc),
		time.Date(2023, 2, 1, 16, 0, 0, 0, loc),
		loc,
	)
	assert.NoError(t, err)
	assert.True(t, result)

	// 测试跨越年份的日期范围
	result, _, err = IsTimeInDateRangeAtSpecificTime(
		"2022-12-30",
		"2023-01-02",
		"10:00:00",
		time.Date(2023, 1, 1, 11, 0, 0, 0, loc),
		time.Date(2023, 1, 1, 12, 0, 0, 0, loc),
		loc,
	)
	assert.NoError(t, err)
	assert.False(t, result)
}
