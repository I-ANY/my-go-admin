package tools

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

func TestParseTables(t *testing.T) {
	changeData := fmt.Sprintf("%s", "主题：SA需求缺口更新差异（单位：G）\n\n新增缺口: \n类型,省份,电信,联通,移动\n\n只跑UDP,,150,-,-\n\n\n本省本网,新疆,-,-,14\n\n\n本省本网,江苏,-,-,107\n\n\n本省本网,浙江,24,-,-\n\n\n本省本网,海南,-,-,3\n\n\n本省本网,湖南,-,-,19\n\n\n扩大缺口(增量/最新缺口): \n类型,省份,电信,联通,移动\n\n只跑UDP,,-,132/150,-\n\n\n本省本网,内蒙古,-,11/151,10/70\n\n\n本省本网,北京,2/23,-,-\n\n\n本省本网,吉林,2/61,-,-\n\n\n本省本网,四川,1/251,4/9,14/100\n\n\n本省本网,安徽,8/13,-,-\n\n\n本省本网,山东,8/156,19/76,-\n\n\n本省本网,山西,-,1/82,-\n\n\n本省本网,江苏,3/238,-,-\n\n\n本省本网,河南,-,-,53/96\n\n\n本省本网,湖北,3/70,-,-\n\n\n本省本网,甘肃,-,4/16,51/141\n\n\n本省本网,福建,-,-,8/22\n\n\n本省本网,辽宁,-,-,3/117\n\n\n本省本网,陕西,15/179,-,-\n\n\n本省本网,青海,2/66,-,-\n\n\n缩小缺口(缩量/最新缺口): \n类型,省份,电信,联通,移动\n\n只跑UDP,,-,-,41/150\n\n\n本省本网,上海,1/42,2/9,8/14\n\n\n本省本网,云南,7/157,1/16,11/264\n\n\n本省本网,内蒙古,3/46,-,-\n\n\n本省本网,北京,-,3/43,6/46\n\n\n本省本网,吉林,-,22/105,-\n\n\n本省本网,天津,-,5/58,-\n\n\n本省本网,宁夏,4/45,-,-\n\n\n本省本网,安徽,-,3/25,18/163\n\n\n本省本网,山东,-,-,72/645\n\n\n本省本网,山西,3/37,-,6/161\n\n\n本省本网,广东,36/321,8/48,-\n\n\n本省本网,广西,17/289,-,31/246\n\n\n本省本网,新疆,6/57,-,-\n\n\n本省本网,江苏,-,26/4,-\n\n\n本省本网,江西,-,-,11/148\n\n\n本省本网,河北,16/208,32/362,40/292\n\n\n本省本网,河南,34/63,34/59,-\n\n\n本省本网,浙江,-,5/36,39/79\n\n\n本省本网,湖北,-,4/7,-\n\n\n本省本网,湖南,4/262,4/29,-\n\n\n本省本网,甘肃,4/189,-,-\n\n\n本省本网,福建,5/78,1/7,-\n\n\n本省本网,贵州,2/144,5/38,32/257\n\n\n本省本网,重庆,3/76,-,9/56\n\n\n本省本网,青海,-,10/1,-\n\n\n本省本网,黑龙江,-,-,11/190\n\n\n跨省,,83/3189,152/1199,86/3262\n\n\n需求关闭: \n类型,省份,运营商\n\n本省本网,青海,移动\n本省本网,吉林,移动\n本省本网,广西,联通\n\n跨网：\n跨网类型,类型,变动/缺口\n\n联通转移动,缩小,59/1652\n电信转移动,扩大,45/644\n移动转联通,缩小,144/9\n电信转联通,缩小,168/75\n移动转电信,新增,55")
	imgBase64, md5sum, err := DrawTablesAsImage(changeData)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Printf("imgBase64: %s, md5sum: %s", imgBase64, md5sum)
	sender := Sender(&Tencent{
		MsgType:   "image",
		ImgBase64: imgBase64,
		ImgMd5:    md5sum,
	})
	webhookUrl := "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=b5a44fc8-6041-46e8-8584-214017825dd1"
	err = sender.SendMsg("", webhookUrl, "")
}

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
