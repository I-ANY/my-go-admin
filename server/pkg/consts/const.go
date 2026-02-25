package consts

const (
	SecretKey = "WcUfsh8UtBtGM0dc53sVqnkjP5LyhmU9s"
	Salt      = "Dk4ynio0iAmxusZeCbB7uHTpBMUv5H41B"
)
const (
	AdminRoleIdentify = "ADMIN"
	BaseRoleIdentify  = "BASE_ROLE"
)

const (
	LeKeyPrefix = "business:le:sn:"
)

// ResourceIdentity 资源资源标识
type ResourceIdentity string

// Subcategory 业务小类资源标识
const (
	// Subcategory 业务小类资源
	Subcategory ResourceIdentity = "SUBCATEGORY"
	// Category 业务大类资源
	Category ResourceIdentity = "CATEGORY"
)

// ResourcePermissionCode 资源权限类型
type ResourcePermissionCode string

// SubcategoryGeneralPermission 资源权限类型
const (
	// SubcategoryGeneralPermission 业务小类通用权限
	SubcategoryGeneralPermission ResourcePermissionCode = "GENERAL_PERMISSION"
	// CategoryGeneralPermission 业务大类通用权限
	CategoryGeneralPermission ResourcePermissionCode = "GENERAL_PERMISSION"
)

// 队列名称
const (
	// QueneName_DSCPReport  队列名称 DSNP
	QueneName_DSCPReport = "queue:dscp:report"
	// QueneName_SpeedLimitTask  队列名称 限速任务执行
	QueneName_SpeedLimitTask = "queue:speedLimitTask:exec"
)

// 队列函数名称
const (
	// QueneFunctionName_DSCPReportExecutor 队列函数名称 DSNP
	QueneFunctionName_DSCPReportExecutor = "DSCPReportExecutor"
	// QueneFunctionName_SpeedLimitTaskExecutor 队列函数名称 限速任务执行
	QueneFunctionName_SpeedLimitTaskExecutor = "SpeedLimitTaskExecutor"
)

// 设备业务为空标识
const (
	Business_Empty = "无业务归属"
)
