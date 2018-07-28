package common

const (
	StatusOk     = 0 //正常
	StatusNoDate = 1 //时间未到 白色
	StatusNoData = 2 //没有数据 灰色
	StatusWarn   = 3 //县级预警 黄色
	StatusAlarm  = 4 //市级报警 红色
	StatusCorpAlarm  = 5 //企业级报警 橙色
)
const (
	QRX_STATION=1000
	KSW_STATION=1001
	GTL_STATION=1002
)
