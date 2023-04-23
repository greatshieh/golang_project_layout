package app

import (
	"golang_project_layout/pkg/global"
	"golang_project_layout/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func InitstallRateLimiter(engin *gin.Engine) {
	for _, rl := range global.GVA_CONFIG.RateLimiter {
		engin.Use(middleware.RateLimiter(rl.Type, rl.Rules...))
	}
}
