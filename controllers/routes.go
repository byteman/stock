package controllers

import (
	"time"
	"github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
)

//api接口采用resultful风格，遵循这篇文章中的建议
//https://blog.igevin.info/posts/restful-api-get-started-to-write/
//http://www.ruanyifeng.com/blog/2014/05/restful_api.html
func RouteRegister(e *gin.Engine)  {

	e.Static("/doc","./apidoc")
	//e.Static("/","./cloudwalk")
	e.Static("/web","./static/dist")
	e.GET("/", func(c *gin.Context) {
		e.LoadHTMLFiles("./static/web/redirect.html")
		c.HTML(200,"redirect.html",nil)
	})
	authMiddleware := &jwt.GinJWTMiddleware{
		Realm:      "test zone",
		Key:        []byte("secret key"),
		Timeout:    time.Hour,
		MaxRefresh: time.Hour,
		Authenticator: MyAuthenticator,
		Authorizator: MyAuthorizator,
		Unauthorized: MyUnauthorized,
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		TokenLookup: "header:Authorization",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	}
	e.POST("/login", authMiddleware.LoginHandler)
	e.Static("/static","./static")
	e.Static("/admin","./static/admin")
	e.GET("/stock/basic/:id",basicQueryStock)
	e.GET("/stock/advance",advanceQueryStock)
	e.POST("/upload",uploadFile)
	e.POST("/stock/users",AddUsers)
	e.GET("/stock/users",GetUsers)
	auth := e.Group("/api/v1")
	auth.Use(authMiddleware.MiddlewareFunc())
	{
		auth.GET("/hello", helloHandler)
		auth.GET("/refresh_token", authMiddleware.RefreshHandler)
		auth.POST("/upload",uploadFile)
		auth.GET("/user/info", func(c *gin.Context) {
			c.JSON(200,gin.H{
				"roles":"admin",
				"name":"admin",
				"avatar":"avatar.png",
			})
		})

	}

}
