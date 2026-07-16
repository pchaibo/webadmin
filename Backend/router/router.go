package router

import (
	"fmt"
	"webadmin/config"
	"webadmin/controller"
	"webadmin/mobilecontr"

	"github.com/gin-gonic/gin"
)

//var R *gin.Engine

// 前端用户
func mobile(R *gin.Engine) {
	R.POST("/user/login", mobilecontr.UserLogin)
 	R.POST("/user/register", mobilecontr.UserCreate)
	userAuth := R.Group("/user")
 	userAuth.Use(mobilecontr.AuthMiddleware())
 	{
 		userAuth.GET("/info", mobilecontr.UserInfo)
 		userAuth.PUT("/password", mobilecontr.UserUpdatePassword)
 		userAuth.PUT("/binance", mobilecontr.UserUpdateBinance)
 	}
	apiCoin := R.Group("/user/coin")
	//apiCoin.Use(mobilecontr.AuthMiddleware())
	{
		apiCoin.GET("", mobilecontr.CoinList)
	}

	apiHeyue := R.Group("/user/heyue")
	apiHeyue.Use(mobilecontr.AuthMiddleware())
	{
	apiHeyue.GET("", mobilecontr.HeyueList)
 		apiHeyue.GET("/:id", mobilecontr.HeyueGet)
		apiHeyue.POST("", mobilecontr.HeyueCreate)
		apiHeyue.PUT("/:id", mobilecontr.HeyueUpdate)
		apiHeyue.DELETE("/:id", mobilecontr.HeyueDelete)
	}

	apiHeyuesorder := R.Group("/user/heyuesorder")
	apiHeyuesorder.Use(mobilecontr.AuthMiddleware())
	{
		apiHeyuesorder.GET("", mobilecontr.HeyuesorderList)
		//apiHeyuesorder.POST("", mobilecontr.HeyuesorderCreate)
		//apiHeyuesorder.PUT("/:id", mobilecontr.HeyuesorderUpdate)
		//apiHeyuesorder.DELETE("/:id", mobilecontr.HeyuesorderDelete)
	}
	apiWs := R.Group("/user/ws")
	//apiWs.Use(controller.AuthMiddleware())
	{
		apiWs.GET("", controller.WsHandler)
	}
}

func Start() {
	var R *gin.Engine
	R = gin.Default()
	mobile(R)
	Admin(R)

	port := config.Get("port")
	if port == "" {
		port = "5000"
	}

	R.Run(fmt.Sprintf("0.0.0.0:%s", port))
}
