package main

import (
	"fmt"
	"stock/cache"
	"stock/common"
	"stock/httpserver"
	"stock/logger"
	"stock/models"


	"github.com/cihub/seelog"
)

func main() {

	defer func() {
		logger.CheckPanic()
		logger.Flush()
	}()
	err := logger.InitFromFile("seelog.xml")
	if err != nil {
		fmt.Errorf("Init logger failed %v", err)
	}

	seelog.Info("App starting .....")
	ctx := common.NewContext()
	seelog.Infof("App start database .....")
	models.InitDao(ctx)

	cache.InitCaches(models.DBConn())

	seelog.Error("App start Scanners.....")

	seelog.Infof("App start http server .....")
	httpserver.Start(ctx)

}
