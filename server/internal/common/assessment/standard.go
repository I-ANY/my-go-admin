package assessment

import (
	"math"

	"biz-auto-api/internal/models"
)

// GetDefaultRuleByRoomType 根据机房类型获取默认考核规则
func GetDefaultRuleByRoomType(roomType uint8) *models.AssessmentRuleStandard {
	switch roomType {
	case 1, 2: // IDC, ACDN
		return &models.AssessmentRuleStandard{
			UtilizationRateThreshold: 0.99, // 99%
			NightPeakPointsThreshold: 36,
		}
	case 3: // MCDN
		return &models.AssessmentRuleStandard{
			UtilizationRateThreshold: 0.95, // 95%
			NightPeakPointsThreshold: 36,
		}
	default:
		// 默认规则
		return &models.AssessmentRuleStandard{
			UtilizationRateThreshold: 0.90, // 90%
			NightPeakPointsThreshold: 36,
		}
	}
}

// IsSpecialApproval 判断是否为特批节点
// 比较利用率阈值和晚高峰点数，任一不等即为特批
func IsSpecialApproval(rule *models.AssessmentRuleStandard, roomType uint8) bool {
	if rule == nil {
		return false
	}
	defaultRule := GetDefaultRuleByRoomType(roomType)

	// 使用容差比较浮点数
	const epsilon = 0.0001
	if math.Abs(rule.UtilizationRateThreshold-defaultRule.UtilizationRateThreshold) > epsilon {
		return true
	}
	if rule.NightPeakPointsThreshold != defaultRule.NightPeakPointsThreshold {
		return true
	}
	return false
}
