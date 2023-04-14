package apiserver

import (
	"fmt"
	"golang_project_layout/internal/apiserver/router"
	"golang_project_layout/pkg/app"
	"golang_project_layout/pkg/global"
	"golang_project_layout/pkg/middleware"
	sysRouter "golang_project_layout/pkg/router"
	"time"

	"golang_project_layout/pkg/plugin/email"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type server interface {
	ListenAndServe() error
}

func RunServer() {
	if global.GVA_CONFIG.System.UseMultipoint || global.GVA_CONFIG.System.UseRedis {
		// 初始化redis服务
		app.Redis()
	}

	// TODO
	// 从db加载jwt数据
	// if global.GVA_DB != nil {
	// 	system.LoadAll()
	// }

	// 初始化路由
	Router := app.InitRouter()
	// 声明不需要鉴权的基础路由
	baseRouter := new(sysRouter.BaseRouter)
	// 注册基础路由
	app.InitSystemRouter(Router, global.GVA_CONFIG.System.RouterPrefix, baseRouter)
	// 声明需要注册的系统路由
	jwtRouter := new(sysRouter.JwtRouter)
	userRouter := new(sysRouter.UserRouter)
	// 注册系统路由
	app.InitSystemRouter(Router, global.GVA_CONFIG.System.RouterPrefix, jwtRouter, userRouter)
	// 注册全局中间件
	privateGroup := Router.Group(global.GVA_CONFIG.System.RouterPrefix)
	app.InitMiddleware(privateGroup, middleware.JWTAuth(), middleware.CasbinHandler())
	// 注册应用路由
	Router.Use(middleware.AccessLog())
	router.AppRouter(Router)

	// 创建email插件
	emailPlugin := email.CreateEmailPlug(
		global.GVA_CONFIG.Email.To,
		global.GVA_CONFIG.Email.From,
		global.GVA_CONFIG.Email.Host,
		global.GVA_CONFIG.Email.Secret,
		global.GVA_CONFIG.Email.Nickname,
		global.GVA_CONFIG.Email.Port,
		global.GVA_CONFIG.Email.IsSSL,
	)
	// 安装需要的插件
	app.InstallPlugin(Router, emailPlugin) // 安装插件

	address := fmt.Sprintf(":%d", global.GVA_CONFIG.System.Addr)
	s := initServer(address, Router)
	// 保证文本顺序输出
	// In order to ensure that the text order output can be deleted
	time.Sleep(10 * time.Microsecond)
	global.GVA_LOG.Info("server run success on ", zap.String("address", address))

	global.GVA_LOG.Error(s.ListenAndServe().Error())
}

func initServer(address string, router *gin.Engine) server {
	s := endless.NewServer(address, router)
	s.ReadHeaderTimeout = 20 * time.Second
	s.WriteTimeout = 20 * time.Second
	s.MaxHeaderBytes = 1 << 20
	return s
}
