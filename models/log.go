package models

import "time"

/**
操作日志.
 */
type Log struct{
	ID  int `json:"id" gorm:"primary_key"` //操作日志编号
	//User User
	UserID int `json:"user_id"` //操作用户的id
	Action int `json:"action"` //行为类型 1普通查询 2高级查询
	LogDate int64 `json:"log_date"` //时间.

}
//添加一条操作日志.
func AddLog(uid int , action int) error {
	log:=&Log{
		UserID:uid,
		Action:action,
		LogDate:time.Now().Unix(),
	}

	return g.Create(log).Error
}
