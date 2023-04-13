package router

import (
	system "golang_project_layout/pkg/controller/v1"

	"github.com/gin-gonic/gin"
)

type BaseRouter struct{}

func (s *BaseRouter) InitRouter(Router gin.IRouter) {
	baseRouter := Router.Group("base")
	baseController := system.ControllerGroupApp.BaseController
	baseRouter.POST("login", baseController.Login)
}
