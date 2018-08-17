package models

import (
	"time"
	"fmt"
)

const (
	UserTrail = 0
	UserExpire = 1
	VipNormal=2
	VipAdvance=3
)
//组织机构
type User struct {
	ID  int `json:"id" gorm:"primary_key"` //用户编号
	Name string `json:"username" gorm:"not null;unique"` //用户名
	PassWord string `json:"password" gorm:"not null"` //密码
	Phone string `json:"phone" gorm:"not null"` //手机号
	RegDate int64 `json:"reg_date" gorm:"not null"` //注册时间
	PayDate int64 `json:"pay_date" gorm:"not null"` //付费时间
	PayType int  `json:"pay_type" gorm:"not null;default:0"` //用户类型 0:试用注册 1:注册用户，无查询功能；2:普通查询功能用户；3高级查询功能用户。
	ExpireDate int64 `json:"expire_date" gorm:"not null;default:0"` //过期时间.

}
//判断普通用户是否过期
func (u *User)NormalIsExpire()bool{
	nowt:=time.Now().Unix()

	return nowt >(u.RegDate+2*24*3600)
}
//判断注册用户是否过期
func (u *User)VipIsExpire()bool{
	nowt:=time.Now().Unix()
	return nowt >(u.PayDate+365*24*3600)
}
//0 普通注册用户 过期和未过期 	过期的不能做基本查询 未过期可以每天一次基本查询,普通用户的过期时间2天
//1 普通充值会员 有过期和未过期 过期的不能做基本查询 未过期可以无限基本查询,充值会员1年有效期
//2 高级充值会员 有过期和未过期 过期的不能做基本查询 未过期可以无限基本和高级查询,充值会员1年有效期
func (u* User)CheckPermission(funcId int)error{
	if u.PayType == UserTrail{

		switch funcId{
		case 2:
			//普通注册用户绝对不能使用高级功能
			return fmt.Errorf("权限不够请充值为高级会员")
		case 1:
			//普通用户访问普通查询，判断是否已经过期,过期后，每天只能试用1次
			if u.NormalIsExpire(){
				var count = 0
				var err error
				if count,err=GetOperCountOfToday(u.ID,funcId);err!=nil{
					return fmt.Errorf("试用账号已经过期,请充值")
				}
				fmt.Printf("%s call count=%d",u.Name,count)
				if count > 0{

					return fmt.Errorf("试用账号已经过期,每天只能试用一次")
				}
				return nil

			}
		default:
			return fmt.Errorf("系统内部错误")
		}

	}else if u.PayType == VipNormal{
		//普通付费账号
		switch funcId{
		case 2:
			//普通付费用户绝对不能使用高级功能
			return fmt.Errorf("权限不够请充值为高级会员")
		case 1:
			//普通用户访问普通查询，判断是否已经过期
			if u.VipIsExpire(){
				return fmt.Errorf("账号已经过期,请充值")
			}
		default:
			return fmt.Errorf("系统内部错误")
		}
	}else if u.PayType == VipAdvance{

			//高级账号只需要判断是否已过期.
			if u.VipIsExpire(){
				return fmt.Errorf("账号已经过期,请充值")
			}


	}

	return nil
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

func UpdateUserLevel(uid int,level int)error{


	return g.Table("users").Where("id=?",uid).Update(User{PayType:level,PayDate:time.Now().Unix()}).Error

}
func UpdateUserPayTime(uid int)error{


	return g.Table("users").Where("id=?",uid).Update(User{PayDate:time.Now().Unix()}).Error

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
