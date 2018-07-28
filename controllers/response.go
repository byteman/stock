package controllers

import (
	"net/http"
	"fmt"
	"github.com/gin-gonic/gin"
)

type Response struct {
	//Code int	`json:"code"`
	Message  string `json:"message"`

}


func ErrorResponse(c *gin.Context,code int , format string,a ...interface{})  {
	//msg:=fmt.Printf("%d",1)
	msg:=format
	if len(a) > 0 {
		msg=fmt.Sprintf(format,a)
	}
	rsp:=Response{Message:msg}
	c.JSON(code,rsp)
	c.Abort()
}

func SuccessQueryResponse(c *gin.Context, data interface{})  {
		c.JSON(http.StatusOK,data)
}
func SuccessExecResponse(c *gin.Context,code  int)  {
	c.Status(code)
}