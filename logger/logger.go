package logger

import (
	"github.com/cihub/seelog"
	"runtime/debug"
)

var Logger seelog.LoggerInterface
func InitFromFile(file string) error {
	var err error
	Logger = seelog.Disabled
	Logger,err=seelog.LoggerFromConfigAsFile(file)
	if err!=nil{
		return err
	}
	seelog.ReplaceLogger(Logger)
	return nil
}
func Flush()  {
	seelog.Flush()
}
func CheckPanic()  {
	if err:=recover();err!=nil{
		stack:=debug.Stack()
		seelog.Error(string(stack))
		seelog.Flush()
		panic("error")
	}
}