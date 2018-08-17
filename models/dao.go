package models

import (
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/jinzhu/gorm"
	"fmt"
	"stock/common"
	"github.com/cihub/seelog"
)


var (
	g         *gorm.DB
	tables  []interface{}
)

func init()  {
	tables=append(tables,new(User),new(Log),
		)

}
func DBConn()*gorm.DB{
	return g.Debug()
}
func AutoMigrate(g    *gorm.DB)  {
	for  _,t:=range tables{
		g.AutoMigrate(t)
	}
}
//https://github.com/go-sql-driver/mysql#dsn-data-source-name
//[username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
type MyLogger struct{

}
var mylog MyLogger
// Print format & print log
func (logger MyLogger) Print(values ...interface{}) {
	seelog.Debug(values)
}
func ConnectDB(engine string,host string, port int,database string, user string, pass string,debug bool,sync bool) (err error) {

	var dbConnString = fmt.Sprintf("%s:%s@tcp(%v:%d)/%s?charset=utf8&parseTime=True", user,pass, host,port,database)
	if engine == "mysql"{
		dbConnString = fmt.Sprintf("%s:%s@tcp(%v:%d)/%s?charset=utf8&parseTime=True", user,pass, host,port,database)
	}else if engine=="postgres" {
		dbConnString = fmt.Sprintf("host=%v port=%d user=%s dbname=%s password=%s",host,port,user,database,pass)
	}else if engine=="sqlite3"{
		dbConnString = fmt.Sprintf("%s.db3",database)
	}

	fmt.Println(dbConnString)

	//db, err := gorm.Open("mysql", "user:password@/dbname?charset=utf8&parseTime=True&loc=Local")
	db, err := gorm.Open(engine, dbConnString)
	if err !=nil{
		panic("open database failed!!")
		return err
	}
	fmt.Println("database connect ok")
	if debug == true{
		g = db.Debug()
		g.SetLogger(mylog)
		g.LogMode(true)
	}else{
		g = db
	}
	//连接远程数据库的时候，为了连接速度，不进行automigrate操作.
	if host=="localhost" || sync{
		AutoMigrate(g)
	}

	return err
}
type DataBase struct {
	Engine string `ini:"engine"`
	Host string `ini:"host"`
	Port int `ini:"port"`
	Db string `ini:"db"`
	UserName string `ini:"username"`
	PassWord string `ini:"password"`
	Debug bool `ini:"debug"`
	Sync bool `ini:"sync"`
}

func InitDao(ctx *common.AppContext)error  {
	db:=DataBase{
		Engine:"mysql",
		Host:"localhost",
		Port:3306,
		Db:"mine",
		UserName:"root",
		PassWord:"wangcheng123",
		Debug:false,
	}
	var err error
	if ctx.Config!=nil{

		err = ctx.Config.Section("database").MapTo(&db)
		seelog.Debugf("%+v",db)
		if err != nil{
			seelog.Debugf("load database failed %v",err)
		}
	}
	//db.Debug = true
	err=ConnectDB(db.Engine,db.Host,db.Port,db.Db,db.UserName,db.PassWord,db.Debug,db.Sync)
	if err!=nil{
		seelog.Critical("Connect database failed %v",err)
	}
	return nil
}