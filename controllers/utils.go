package controllers

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"net/http"
	"stock/models"
	"time"
	"fmt"
	"errors"
	"github.com/cihub/seelog"
)

//获取
func DefaultQueryInt(c *gin.Context, key string, defValue int)int  {
	sv,exist:=c.GetQuery(key)
	if exist {
		iv,err:=strconv.Atoi(sv)
		if err!=nil{
			return defValue
		}
		return iv
	}else{
		return defValue
	}

}
func DefaultQueryInt2(c *gin.Context, key string)(oid int,exist bool)  {
	sv,exist:=c.GetQuery(key)
	if exist {
		iv,err:=strconv.Atoi(sv)
		if err!=nil{
			return -1,false
		}
		return iv,true
	}else{
		return -1,false
	}

}
func DefaultQueryIntWithFunc(c *gin.Context, key string, GetIDFunc func (c* gin.Context) int)int  {
	sv,exist:=c.GetQuery(key)
	if exist {
		iv,err:=strconv.Atoi(sv)
		if err!=nil{
			return GetIDFunc(c)
		}
		return iv
	}else{
		return GetIDFunc(c)
	}

}

func ParseIntParam(c *gin.Context, key string)(id int, err error){
	value:=c.Param(key)
	id,err=strconv.Atoi(value)
	if err!=nil{
		ErrorResponse(c,http.StatusBadRequest,"invalid %s %v",key,err)
	}
	return
}
func ParseIntParamBindJSON(c *gin.Context, key string,o interface{})(id int, err error){

	id,err=ParseIntParam(c,key)
	if err!=nil{
		return
	}
	err = ParseBindJSON(c,o)
	return
}
func ParseBindJSON(c *gin.Context,o interface{})( err error) {
	err=c.BindJSON(o)
	if err!=nil{
		ErrorResponse(c,http.StatusBadRequest,"BindJSON failed %v",err)
	}
	return
}
func ParseDefaultPageQuery(c *gin.Context)  (page,page_size int, start,end int){
	day := int64(3600*24)

	start=DefaultQueryInt(c,"start",int(time.Now().Unix()-day))
	end=DefaultQueryInt(c,"end",int(time.Now().Unix()))
	page=DefaultQueryInt(c,"page",1)
	page_size=DefaultQueryInt(c,"page_size",20)

	if page < 1 {
		page=1
	}
	return
}
func ErrorInvalidParam(c* gin.Context,err error) {
	seelog.Errorf("ErrorInvalidParam %v",err)
	ErrorResponse(c,http.StatusBadRequest,"invalid param %v",err)
	c.Abort()
}
func ErrorDB(c* gin.Context,err error) {
	seelog.Errorf("ErrorDB %v",err)
	ErrorResponse(c,http.StatusInternalServerError,"db operation failed! %v",err)
	c.Abort()
}
func SuccessCreated(c* gin.Context){
	c.Status(http.StatusCreated)
}

func DefaultAddOperation(c *gin.Context, o interface{})error{

	err:=c.BindJSON(o)
	if err!=nil{

		ErrorInvalidParam(c,err)
		return err
	}
	err= models.DefaultAdd(o)
	if err!=nil{

		ErrorDB(c,err)
		return err
	}
	SuccessCreated(c)
	return nil
}

func DefaultGetAll(c *gin.Context, o interface{}){

	err:=models.DefaultGetAll(o)
	if err!=nil{
		ErrorResponse(c,http.StatusNotFound,"err=%v",err)
		return
	}
	SuccessQueryResponse(c,o)
}
func DefaultRemoveById(c *gin.Context,o interface{}){
	id:=c.Param("id")
	cid,err:=strconv.Atoi(id)
	if err!=nil{
		ErrorResponse(c,http.StatusBadRequest,"id [%s] %s",id,err)
		return
	}
	if err=models.DefaultRemoveById(cid,o);err!=nil{
		ErrorResponse(c,http.StatusBadRequest,"Remove SaveDB failed %s",err)
		return
	}
	c.Status(http.StatusNoContent)
}
func DefaultUpdateById(c *gin.Context,id int,o interface{}){

	if err:=models.DefaultUpdateById(id,o);err!=nil{
		ErrorResponse(c,http.StatusBadRequest,"Update SaveDB failed %s",err)
		return
	}
	c.Status(http.StatusNoContent)
}

func DefaultGetById(c *gin.Context,o interface{}){

	id:=c.Param("id")
	cid,err:=strconv.Atoi(id)
	if err!=nil{
		ErrorResponse(c,http.StatusBadRequest,"id [%s] %v",id,err)
		return
	}
	if err=models.DefaultGetById(cid,o);err!=nil{
		ErrorResponse(c,http.StatusBadRequest,"Remove SaveDB failed %s",err)
		return
	}
	SuccessQueryResponse(c,o)
}

func GeCurrentUser(c *gin.Context)(u *models.User, err error)  {
	uid,exist:=c.Get("userID")
	if !exist {
		return nil,errors.New("can not get userId")
	}
	uid2,ok:=uid.(string)
	if !ok{
		return nil,errors.New("can not get userId")
	}
	//根据用户名获取到用户的信息
	u,err=models.GetUserByName(uid2)
	if err!=nil{
		return nil,fmt.Errorf("can not get user %s",uid2)
	}
	return u,nil
}

type  api_page_data struct{
	//Corps []*models.Corp `json:"corps"`
	Data interface{} `json:"data"`
	TotalPage int `json:"total_pages"`
	TotalRecords int `json:"total_records"`
	CurrentPage int `json:"current_page"`
}

func SuccessResponsePageData(c *gin.Context, data interface{}, page,pageSz,total int)()  {
	apd:=api_page_data{
		Data:data,
		TotalPage:(total+pageSz-1)/pageSz,
		TotalRecords:total,
		CurrentPage:page,
	}
	SuccessQueryResponse(c,apd)
}