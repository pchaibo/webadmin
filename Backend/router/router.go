package router

import (
	"fmt"
	"webadmin/config"
	"webadmin/controller"

	"github.com/gin-gonic/gin"
)

//var R *gin.Engine

func Start() {
	var R *gin.Engine
	R = gin.Default()

	// === 公开路由（无需认证） ===
	R.POST("/login", controller.AdminLogin)
	R.POST("/api/login", controller.AdminLogin)

	// 前端静态文件（/ul/ 前缀）
	R.GET("/", func(c *gin.Context) {
		c.Redirect(301, "/ul")
	})
	// /ul/* 所有路径都返回 index.html（SPA 前端路由）
	R.GET("/ul", func(c *gin.Context) {
		c.File("./frontend/index.html")
	})
	R.Static("/ul/assets", "./frontend/assets")
	R.StaticFile("/ul/favicon.ico", "./frontend/favicon.ico")

	// SPA 回退 ── 所有未匹配 API 的路径都返回 index.html
	R.NoRoute(func(c *gin.Context) {
		c.File("./frontend/index.html")
	})
	ping := R.Group("/ping")
	ping.GET("", controller.Ping)

	// === 以下需要认证 ===

	// ── 无前缀路由（兼容旧版 Vite 代理 strip /api 的场景）──
	// admins := R.Group("/admin")
	// admins.Use(controller.AuthMiddleware())
	// {
	// 	admins.GET("", controller.AdminList)
	// 	admins.POST("", controller.AdminCreate)
	// 	admins.PUT("/:id", controller.AdminUpdate)
	// 	admins.DELETE("/:id", controller.AdminDelete)
	// }

	// authRules := R.Group("/auth_rule")
	// authRules.Use(controller.AuthMiddleware())
	// {
	// 	authRules.GET("", controller.AuthRuleList)
	// 	authRules.POST("", controller.AuthRuleCreate)
	// 	authRules.PUT("/:id", controller.AuthRuleUpdate)
	// 	authRules.DELETE("/:id", controller.AuthRuleDelete)
	// }

	// shellGroup := R.Group("/shell_group")
	// shellGroup.Use(controller.AuthMiddleware())
	// {
	// 	shellGroup.GET("", controller.ShellGroupList)
	// 	shellGroup.POST("", controller.ShellGroupCreate)
	// 	shellGroup.PUT("/:id", controller.ShellGroupUpdate)
	// 	shellGroup.DELETE("/:id", controller.ShellGroupDelete)
	// }

	// shell := R.Group("/shell")
	// shell.Use(controller.AuthMiddleware())
	// {
	// 	shell.GET("", controller.ShellList)
	// 	shell.POST("", controller.ShellCreate)
	// 	shell.PUT("/:id", controller.ShellUpdate)
	// 	shell.DELETE("/:id", controller.ShellDelete)
	// }

	// ── /api 前缀路由（兼容新版前端 api.ts 的 /api/xxx 调用）──
	apiAdmins := R.Group("/api/admin")
	apiAdmins.Use(controller.AuthMiddleware())
	{
		apiAdmins.GET("", controller.AdminList)
		apiAdmins.POST("", controller.AdminCreate)
		apiAdmins.PUT("/:id", controller.AdminUpdate)
		apiAdmins.DELETE("/:id", controller.AdminDelete)
	}

	apiAuthRules := R.Group("/api/auth_rule")
	apiAuthRules.Use(controller.AuthMiddleware())
	{
		apiAuthRules.GET("", controller.AuthRuleList)
		apiAuthRules.POST("", controller.AuthRuleCreate)
		apiAuthRules.PUT("/:id", controller.AuthRuleUpdate)
		apiAuthRules.DELETE("/:id", controller.AuthRuleDelete)
	}

	apiShellGroup := R.Group("/api/shell_group")
	apiShellGroup.Use(controller.AuthMiddleware())
	{
		apiShellGroup.GET("", controller.ShellGroupList)
		apiShellGroup.POST("", controller.ShellGroupCreate)
		apiShellGroup.PUT("/:id", controller.ShellGroupUpdate)
		apiShellGroup.DELETE("/:id", controller.ShellGroupDelete)
	}

	apiShell := R.Group("/api/shell")
	apiShell.Use(controller.AuthMiddleware())
	{
		apiShell.GET("", controller.ShellList)
		apiShell.POST("", controller.ShellCreate)
		apiShell.PUT("/:id", controller.ShellUpdate)
		apiShell.DELETE("/:id", controller.ShellDelete)
	}

	apiShellMax := R.Group("/api/shell_max")
	apiShellMax.Use(controller.AuthMiddleware())
	{
		apiShellMax.GET("", controller.ShellMaxList)
	}

	apiShellMin := R.Group("/api/shell_min")
	apiShellMin.Use(controller.AuthMiddleware())
	{
		apiShellMin.GET("", controller.ShellMinList)
	}

	apiCoin := R.Group("/api/coin")
	apiCoin.Use(controller.AuthMiddleware())
	{
		apiCoin.GET("", controller.CoinList)
		apiCoin.POST("", controller.CoinCreate)
		apiCoin.PUT("/:id", controller.CoinUpdate)
		apiCoin.DELETE("/:id", controller.CoinDelete)
	}

	apiHeyue := R.Group("/api/heyue")
	apiHeyue.Use(controller.AuthMiddleware())
	{
		apiHeyue.GET("", controller.HeyueList)
		apiHeyue.POST("", controller.HeyueCreate)
		apiHeyue.PUT("/:id", controller.HeyueUpdate)
		apiHeyue.DELETE("/:id", controller.HeyueDelete)
	}

	apiHeyuesorder := R.Group("/api/heyuesorder")
	//apiHeyuesorder.Use(controller.AuthMiddleware())
	{
		apiHeyuesorder.GET("", controller.HeyuesorderList)
		apiHeyuesorder.POST("", controller.HeyuesorderCreate)
		apiHeyuesorder.PUT("/:id", controller.HeyuesorderUpdate)
		apiHeyuesorder.DELETE("/:id", controller.HeyuesorderDelete)
	}

	apiUser := R.Group("/api/user")
	apiUser.Use(controller.AuthMiddleware())
	{
		apiUser.GET("", controller.UserList)
		apiUser.POST("", controller.UserCreate)
		apiUser.PUT("/:id", controller.UserUpdate)
		apiUser.DELETE("/:id", controller.UserDelete)
	}
	//ws
	apiWs := R.Group("/api/ws")
	//apiWs.Use(controller.AuthMiddleware())
	{
		apiWs.GET("", controller.WsHandler)
	}

	port := config.Get("port")
	if port == "" {
		port = "5000"
	}

	R.Run(fmt.Sprintf("0.0.0.0:%s", port))
}
