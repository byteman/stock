package httpserver

import (
	"stock/controllers"
	"github.com/gin-gonic/gin"
	"time"
	"github.com/itsjamie/gin-cors"
	"stock/common"
	"fmt"
	"github.com/cihub/seelog"
	"runtime"
)

var engine *gin.Engine=nil

func Start(ctx *common.AppContext){
	var port = 8082
	var mode = "release"
	if ctx.Config!=nil{
		port = ctx.Config.Section("server").Key("port").MustInt(8082)
		mode = ctx.Config.Section("server").Key("mode").MustString("release")
	}

	gin.SetMode(mode)
	engine = gin.Default()
	Wrapper(engine)
	engine.Use(cors.Middleware(cors.Config{
		Origins:        "*",
		Methods:        "GET, PUT, POST, DELETE",
		RequestHeaders: "Origin, Authorization, Content-Type",
		ExposedHeaders: "",
		MaxAge: 50 * time.Second,
		Credentials: true,
		ValidateHeaders: false,
	}))

	controllers.RouteRegister(engine)


	seelog.Debugf("goroutine num =%d",runtime.NumGoroutine())
	//http.ListenAndServe(":8088",nil)

	engine.Run(fmt.Sprintf(":%d",port))
}
func Stop()  {

}