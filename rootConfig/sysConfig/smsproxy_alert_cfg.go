package sysConfig

type SmsproxyAlertConfig struct {
	CheckIntervalMinute int64 `json:"check_interval_minute"` // 检测分钟间隔
	MinAlertCount       int64 `json:"min_alert_count"`       // 条数大于指定值，再检测告警
	AlertRate           int64 `json:"alert_rate"`            // 回填率阈值 百分比
}
