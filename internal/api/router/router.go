package router

import (
	"github.com/ResertCursorService/internal/api/handlers"
	"github.com/ResertCursorService/internal/api/middleware"
	"github.com/gin-gonic/gin"
)

type Router struct {
	adminHandler      *handlers.AdminHandler
	activationHandler *handlers.ActivationHandler
	appHandler        *handlers.AppHandler
}

func NewRouter(adminHandler *handlers.AdminHandler, activationHandler *handlers.ActivationHandler, appHandler *handlers.AppHandler) *Router {
	return &Router{
		adminHandler:      adminHandler,
		activationHandler: activationHandler,
		appHandler:        appHandler,
	}
}

func (r *Router) Setup(engine *gin.Engine) {
	// 静态文件服务
	engine.Static("/assets", "./web/dist/assets")
	engine.StaticFile("/", "./web/dist/index.html")

	// API路由
	api := engine.Group("/api")
	{
		// 公开路由
		api.POST("/login", r.adminHandler.Login)

		// App 相关路由
		app := api.Group("/app")
		{
			// 不需要认证的路由
			app.POST("/activate", r.appHandler.Activate)

			// 需要app token认证的路由
			appProtected := app.Group("")
			appProtected.Use(middleware.AppAuthMiddleware())
			{
				appProtected.GET("/account", r.appHandler.GetAccount)
				appProtected.POST("/account", r.appHandler.CreateAccount)
				appProtected.GET("/code-info", r.appHandler.GetCodeInfo)
			}
		}

		// 需要管理员认证的路由
		adminProtected := api.Group("")
		adminProtected.Use(middleware.AdminAuthMiddleware())
		{
			// 激活码相关路由
			adminProtected.POST("/activation-codes", r.activationHandler.CreateActivationCode)
			adminProtected.GET("/activation-codes", r.activationHandler.ListActivationCodes)
			adminProtected.GET("/activation-codes/:id", r.activationHandler.GetActivationCode)
			adminProtected.PUT("/activation-codes/:id/status", r.activationHandler.UpdateActivationCodeStatus)
		}
	}

	// 处理前端路由
	engine.NoRoute(func(c *gin.Context) {
		c.File("./web/dist/index.html")
	})
}
