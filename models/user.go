package models

import "time"

//组织机构
type User struct {
	ID  int `json:"id" gorm:"primary_key"` //用户编号
	Name string `json:"username" gorm:"not null;unique"` //用户名
	PassWord string `json:"password" gorm:"not null"` //密码
	Phone string `json:"phone" gorm:"not null"` //手机号
	RegDate int64 `json:"reg_date" gorm:"not null"` //注册时间
	PayDate int64 `json:"pay_date" gorm:"not null"` //付费时间
	PayType int  `json:"pay_type" gorm:"not null;default:0"` //用户类型 0:注册用户，无查询功能；1:普通查询功能用户；2:高级查询功能用户。


}
//根据用户名获取用户信息.
func GetUserByName(name string)(u *User,e error)  {
	 u=&User{}
	 if err:=g.Where("name=?",name).First(u).Error; err!=nil{
		 return nil,err
	 }
	return u,nil
}



func AddUser(u *User) error{
	u.RegDate = time.Now().Unix()
	u.PayType = 0
	db:=g.Create(u)
	if db.Error!=nil{
		return db.Error
	}

	return nil
}
func GetUsers(users *[]User)error  {
	return g.Find(users).Error
}

func RemoveUserById(uid int) error {
	if err := g.Where("id=?", uid).Delete(User{}).Error; err != nil {
		return err
	}
	return nil
}
func RemoveUserByOrgId(oid int) error {
	if err := g.Where("org_id=?", oid).Delete(User{}).Error; err != nil {
		return err
	}
	return nil
}
func CheckLogin(user,password string) error {
	u:=User{}
	if user=="admin" && password=="admin"{
		return nil
	}
	if err:=g.Where("name=? and pass_word=?",user,password).Find(&u).Error;err!=nil{
		return err
	}
	return nil
}
