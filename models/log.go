package models

import (
	"time"
	"github.com/jinzhu/now"
)

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
//获取某个用户今天执行某个操作的次数
func GetOperCountOfToday(uid int, action int) (count int,err error) {
	//start:=now.BeginningOfDay().Unix()

	result:= struct {
		Total int `json:"total"`
	}{}
	if err:=g.Model(Log{}).Select("count(id) as total").
		Where("user_id=?",uid).
		Where("action=?",action).
		Where("log_data > ? and log_data < ?",now.BeginningOfDay().Unix(), now.EndOfDay().Unix()).Scan(&result).Error;err!=nil{
			return 0,err
	}
	return result.Total,nil
}

func AddLogByName(user string , action int) error {

	u,err:=GetUserByName(user)
	if err!=nil{
		return err
	}
	log:=&Log{
		UserID:u.ID,
		Action:action,
		LogDate:time.Now().Unix(),
	}

	return g.Create(log).Error
}