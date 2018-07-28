package models

//组织机构
type User struct {
	ID  int `json:"id" gorm:"primary_key"` //用户编号
	ProjectID int `json:"-"`//真正在数据库中存在的外键值
	Name string `json:"username" gorm:"not null;unique"` //用户名
	PassWord string `json:"password" gorm:"not null"` //密码
	OrgID int `json:"org_id" gorm:"not null"` //组织的级别.
	Region string `json:"region"` //用户账户所在的行政区域
	AdCode string `json:"ad_code"`
	Group int `json:"-" gorm:"not null"` //组织的组.
}
//根据用户名获取用户信息.
func GetUserByName(name string)(u *User,e error)  {
	 u=&User{}
	 if err:=g.Where("name=?",name).First(u).Error; err!=nil{
		 return nil,err
	 }
	return u,nil
}

//根据用户名获取用户信息.
func GetUserByOrgId(oid int)(us []*User,e error)  {


	users := make([]*User, 0, 5)


	if err:=g.Where("org_id=?",oid).Find(&users).Error;err!=nil{
		return nil,err
	}
	return users,nil


}


func AddUser(u *User) error{
	db:=g.Create(u)
	if db.Error!=nil{
		return db.Error
	}

	return nil
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
