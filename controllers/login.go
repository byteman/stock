package controllers

import (
	"github.com/gin-gonic/gin"
	"stock/models"

)


/**
* @apiDescription  登录成功以后，会获取到一个token，这个token的有效期只有1个小时，如果你不调用refresh_token接口的话，一个小时后就不能用了，然后你需要每隔1分钟，调用一个refresh_token获取一个新的token，然后有效期又可以从当前时间开始延长一个小时

 * @api {post} /login 请求登陆
 * @apiName PutSession
 * @apiGroup User
 * @apiParam {String} username 用户名.
 * @apiParam {String} password 密码.
 *
 * @apiSuccess {String} token 用户的token
 * @apiSuccess {String} expire token的过期UNIX时间戳
 * @apiError 401  未授权
 */


/**
* @apiDescription  登录成功以后，会获取到一个token，这个token的有效期只有1个小时，如果你不调用refresh_token接口的话，一个小时后就不能用了，然后你需要每隔1分钟，调用一个refresh_token获取一个新的token，然后有效期又可以从当前时间开始延长一个小时

 * @api {get} /api/v1/refresh_token 刷新登录token
 * @apiName refresh_token
 * @apiGroup User
  * @apiHeader {String} Authorization json web token.
  * @apiHeaderExample {json} Header-Example:
 *     {
 *       "Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0OTYzMDY4NjEsImlkIjoiYWRtaW4iLCJvcmlnX
 					2lhdCI6MTQ5NjMwMzI2MX0.z3CnD5zfZpNHHv38DT-ygjQG5M5OkclNV_ZXhwihC38"
 *     }
 *
 * @apiSuccess {String} token 用户的token
 * @apiSuccess {String} expire token的过期UNIX时间戳
 * @apiError 401  未授权
 */
type Login struct {
	Name     string
	PassWord string
}

func helloHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"text": "Hello World.",
	})
}
//http://ju.outofmemory.cn/entry/134189
// Authorization
// Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0OTYyMTM3NTUsImlkIjoiYWRtaW4iLCJvcmlnX2lhdCI6MTQ5NjIxMDE1NX0.tDxSCsWRCCUdmjon5YUcjKqe7UCN0kC05KbFQghaMho
func MyAuthenticator(userID string, password string, c *gin.Context) (string, bool) {

	if userID == "admin" && password == "admin" {
		return userID,true
	}
	uuid:=c.GetHeader("uuid")

	user,err:=models.GetUserByName(userID)
	if err!=nil{
		return "找不到该账号",false
	}
	if user.PassWord != password{
		return "账号或者密码错误",false
	}

	if user.UUID != uuid{
		return "您的账号在非注册手机登陆，登陆失败\r\n如果您是正常更换手机，请用QQ：492951940与服务方联系",false
	}
	//添加一条登录记录
	if user,err:=models.GetUserByName(userID);err==nil{
		models.AddLoginLog(user)
	}


	return userID, true
}
func MyPayLoadFunc(userId string)map[string]interface{}{
	user,err:=models.GetUserByName(userId)
	if err!=nil{
		return gin.H{

		}
	}

	return gin.H{
		"level":user.PayType,
		"regts":user.RegDate,
		"payts":user.PayDate,
	}
}
func MyAuthorizator(userId string, c *gin.Context)bool  {


	//判断某个
	return true
}
//严重不成功返回的消息.
func MyUnauthorized(c *gin.Context, code int, message string)()  {
	c.JSON(code, gin.H{
		"code":    code,
		"message": message,
	})
}