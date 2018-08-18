package controllers

import (
	"github.com/gin-gonic/gin"
	"stock/models"
	"io/ioutil"
)



/**
 * @apiDescription 获取指定组织下的账号信息（账号可以分配给组织机构和尾矿库).
 	这个接口调用前需要先登录获取jwt
 * @api {get} /api/v1/users 获取指定组织下的账号信息
 * @apiName GetUsers
 * @apiGroup User
  * @apiHeader {String} Authorization json web token.
  * @apiHeaderExample {json} Header-Example:
 *     {
 *       "Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0OTYzMDY4NjEsImlkIjoiYWRtaW4iLCJvcmlnX
 					2lhdCI6MTQ5NjMwMzI2MX0.z3CnD5zfZpNHHv38DT-ygjQG5M5OkclNV_ZXhwihC38"
 *     }
*@apiParam {String} [org_id]  组织的编号.URL中.(不传表示当前权限下的最顶级组织.)
 * @apiSuccess {object[]} users 账户信息数组
 * @apiSuccess {string} users.id 用户编号.
 * @apiSuccess {string} users.username 用户名.
 * @apiSuccess {String} users.password 密码.
 * @apiSuccess {String} users.org_id 账户所属于的组织编号.
 * @apiError 401  未授权
 */
func GetUsers(c *gin.Context)  {
	users:=make([]models.User,0)

	if models.GetUsers(&users) != nil{
		ErrorResponse(c, 401, "can not get users")
		return
	}
	for index,user:=range users{
		if user.PayType == 0 && user.NormalIsExpire(){
			//试用用户，并且过期了，就变成普通账号.
			users[index].PayType = 1
		}
	}
	SuccessQueryResponse(c,users)
}

func GetHelp(c *gin.Context)  {
	msg,err:=ioutil.ReadFile("help.txt")
	if err!=nil{
		ErrorResponse(c,404,"read file failed %v",err)
		return
	}
	SuccessQueryResponse(c,gin.H{
		"data":msg,
	})
}

/**
 * @apiDescription 添加用户.
 	这个接口调用前需要先登录获取jwt
 * @api {POST} /api/v1/users 添加用户
 * @apiName AddUsers
 * @apiGroup User
  * @apiHeader {String} Authorization json web token.
  * @apiHeaderExample {json} Header-Example:
 *     {
 *       "Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0OTYzMDY4NjEsImlkIjoiYWRtaW4iLCJvcmlnX
 					2lhdCI6MTQ5NjMwMzI2MX0.z3CnD5zfZpNHHv38DT-ygjQG5M5OkclNV_ZXhwihC38"
 *     }
*@apiParam {Number} org_id  		组织编号.(不传表示当前权限下的最顶级组织.)
*@apiParam {String} username  		用户名.(不传表示当前权限下的最顶级组织.)
*@apiParam {String} password  		密码 .
 * @apiSuccess 201 已经创建.
 * @apiError 401  未授权
 */
func AddUsers(c *gin.Context)  {

	user:=models.User{}

	var err error
	//绑定获取User信息
	err=c.BindJSON(&user)
	if err!=nil{
		ErrorInvalidParam(c,err)
	}

	err= models.AddUser(&user)
	if err!=nil{
		ErrorDB(c,err)
		return
	}
	SuccessCreated(c)
}

/**
 * @apiDescription 删除用户账户.
 	这个接口调用前需要先登录获取jwt
 * @api {DELETE} /api/v1/users/:id 删除用户账户
 * @apiName RemoveUsers
 * @apiGroup User
  * @apiHeader {String} Authorization json web token.
  * @apiHeaderExample {json} Header-Example:
 *     {
 *       "Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0OTYzMDY4NjEsImlkIjoiYWRtaW4iLCJvcmlnX
 					2lhdCI6MTQ5NjMwMzI2MX0.z3CnD5zfZpNHHv38DT-ygjQG5M5OkclNV_ZXhwihC38"
 *     }
 * @apiParam {string} id 需要删除的用户ID.
 * @apiSuccess 204 成功.
 * @apiError 401  未授权
 */

func RemoveUsers(c *gin.Context)  {
	user:=models.User{}
	DefaultRemoveById(c,&user)
}
func LoginOut(c *gin.Context)  {
	c.JSON(200,nil)
}
/**
 * @apiDescription 修改用户账户.
 	这个接口调用前需要先登录获取jwt
 * @api {PUT} /api/v1/users/:id 修改用户账户
 * @apiName UpdateUser
 * @apiGroup User
  * @apiHeader {String} Authorization json web token.
  * @apiHeaderExample {json} Header-Example:
 *     {
 *       "Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0OTYzMDY4NjEsImlkIjoiYWRtaW4iLCJvcmlnX
 					2lhdCI6MTQ5NjMwMzI2MX0.z3CnD5zfZpNHHv38DT-ygjQG5M5OkclNV_ZXhwihC38"
 *     }
 * @apiParam {string} id 需要更新的用户ID.
 * @apiParam {string} username 用户名.
 * @apiParam {String} password 密码.
 * @apiSuccess 204 成功.
 * @apiError 401  未授权
 */
func UpdateUser (c *gin.Context)  {
	user:=&models.User{}
	id,err:=ParseIntParamBindJSON(c,"id",user)
	if err!=nil{
		return
	}

	user.ID = id
	models.UpdateUserLevel(id,user.PayType)
	//DefaultUpdateById(c,id,user)
}
/**
 * @apiDescription 获取指定用户账户.
 	这个接口调用前需要先登录获取jwt
 * @api {GET} /api/v1/users/:id 获取指定用户账户
 * @apiName GetUserById
 * @apiGroup User
  * @apiHeader {String} Authorization json web token.
  * @apiHeaderExample {json} Header-Example:
 *     {
 *       "Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0OTYzMDY4NjEsImlkIjoiYWRtaW4iLCJvcmlnX
 					2lhdCI6MTQ5NjMwMzI2MX0.z3CnD5zfZpNHHv38DT-ygjQG5M5OkclNV_ZXhwihC38"
 *     }
 * @apiSuccess {string} id 用户编号.
 * @apiSuccess {string} username 用户名.
 * @apiSuccess {String} password 密码.
 * @apiSuccess {String} org_id 账户所属于的组织编号.
 * @apiError 401  未授权
 */
func GetUserById(c *gin.Context)  {
	user:=models.User{}
	DefaultGetById(c,&user)
}