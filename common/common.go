package common

import (
	"gopkg.in/ini.v1"
	"github.com/robfig/cron"
)

type AppContext struct {
		Config *ini.File
		Cron   *cron.Cron

}

func NewContext()(*AppContext)  {
	 ctx:=AppContext{

	 }
	var err error
	ctx.Config,err=ini.Load("wkmonitor.ini")
	if err!=nil {

	}
	ctx.Cron = cron.New()
	return &ctx
}

