package app

import (
	"golang_project_layout/pkg/global"
	"golang_project_layout/pkg/router"
	"net/http"

	"golang_project_layout/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// 初始化路由
func InitRouter() *gin.Engine {
	Router := gin.Default()

	// 注册swagger路由
	docs.SwaggerInfo.BasePath = global.GVA_CONFIG.System.RouterPrefix
	Router.GET(global.GVA_CONFIG.System.RouterPrefix+"/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	global.GVA_LOG.Info("register swagger handler")

	// 添加统一的路由前缀
	PublicGroup := Router.Group(global.GVA_CONFIG.System.RouterPrefix)

	// 健康监测
	PublicGroup.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, "ok")
	})

	return Router
}

// 初始化系统路由
func InitSystemRouter(g gin.IRouter, prefix string, systemRouterGroup ...router.Router) {
	routerGroup := g.Group(prefix)
	for _, systemRouter := range systemRouterGroup {
		systemRouter.InitRouter(routerGroup)
	}
	global.GVA_LOG.Info("system router register success")
}

// 注册中间件
func InitMiddleware(R gin.IRoutes, middleware ...gin.HandlerFunc) {
	for _, m := range middleware {
		R.Use(m)
	}
	global.GVA_LOG.Info("middleware register success")
}
