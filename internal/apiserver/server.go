package apiserver

import (
	"fmt"
	"golang_project_layout/internal/apiserver/router"
	"golang_project_layout/pkg/app"
	"golang_project_layout/pkg/global"
	"time"

	"golang_project_layout/pkg/plugin/email"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type server interface {
	ListenAndServe() error
}

type GenericServer struct {
	Name        string
	Description string
	*gin.Engine
}

func NewServer(name string) *GenericServer {
	genericServer := &GenericServer{Name: name, Engine: gin.New()}

	if global.GVA_CONFIG.System.UseMultipoint || global.GVA_CONFIG.System.UseRedis {
		// 初始化redis服务
		app.Redis()
	}

	// TODO
	// 从db加载jwt数据
	// if global.GVA_DB != nil {
	// 	system.LoadAll()
	// }

	// 注册默认系统基础路由
	app.InstalltSystemRouter(genericServer.Engine, global.GVA_CONFIG.System.SystemRouters)

	// 注册中间件
	app.InstalltMiddleware(genericServer.Engine, global.GVA_CONFIG.System.Middleware)

	// 注册限流器
	app.InitstallRateLimiter(genericServer.Engine)

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
	// 注册需要的插件
	app.InstallPlugin(genericServer.Engine, emailPlugin)

	// 注册应用路由
	router.AppRouter(genericServer.Engine)

	return genericServer
}

func RunServer(name string) {
	genericServer := &GenericServer{Name: name, Engine: gin.New()}

	if global.GVA_CONFIG.System.UseMultipoint || global.GVA_CONFIG.System.UseRedis {
		// 初始化redis服务
		app.Redis()
	}

	// TODO
	// 从db加载jwt数据
	// if global.GVA_DB != nil {
	// 	system.LoadAll()
	// }

	// 注册默认系统基础路由
	app.InstalltSystemRouter(genericServer.Engine, global.GVA_CONFIG.System.SystemRouters)

	// 注册中间件
	app.InstalltMiddleware(genericServer.Engine, global.GVA_CONFIG.System.Middleware)

	// TODO 注册插件

	// 注册应用路由
	router.AppRouter(genericServer.Engine)

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
	app.InstallPlugin(genericServer.Engine, emailPlugin) // 安装插件

	address := fmt.Sprintf(":%d", global.GVA_CONFIG.System.Addr)
	s := initServer(address, genericServer.Engine)
	// 保证文本顺序输出
	// In order to ensure that the text order output can be deleted
	time.Sleep(10 * time.Microsecond)
	global.GVA_LOG.Info("server run success on ", zap.String("address", address))

	global.GVA_LOG.Error(s.ListenAndServe().Error())
}

func (s *GenericServer) Run() server {
	address := fmt.Sprintf(":%d", global.GVA_CONFIG.System.Addr)
	srv := endless.NewServer(address, s.Engine)
	srv.ReadHeaderTimeout = 20 * time.Second
	srv.WriteTimeout = 20 * time.Second
	srv.MaxHeaderBytes = 1 << 20

	time.Sleep(10 * time.Microsecond)
	global.GVA_LOG.Info("server run success on ", zap.String("address", address))

	global.GVA_LOG.Error(srv.ListenAndServe().Error())
	return srv
}

func initServer(address string, router *gin.Engine) server {
	s := endless.NewServer(address, router)
	s.ReadHeaderTimeout = 20 * time.Second
	s.WriteTimeout = 20 * time.Second
	s.MaxHeaderBytes = 1 << 20
	return s
}
