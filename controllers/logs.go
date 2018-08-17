package controllers

import (
	"github.com/gin-gonic/gin"
	"stock/models"
)

func GetLogs(c *gin.Context)  {
	logs:=make([]models.Log,0)

	if models.GetLogs(&logs) != nil{
		ErrorResponse(c, 401, "can not get logs")
		return
	}

	SuccessQueryResponse(c,logs)
}