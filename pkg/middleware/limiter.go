package middleware

import (
	"golang_project_layout/pkg/middleware/limiter"

	"github.com/gin-gonic/gin"
)

// 新建一个限流器
func RateLimiter(limiterType string, rules ...limiter.LimiterBucketRules) gin.HandlerFunc {
	var l limiter.LimiterIface
	switch limiterType {
	case "router":
		l = limiter.NewRouterLimiter()
	case "ip":
		l = limiter.NewIPLimiter()
	default:
		l = limiter.NewRouterLimiter()
	}

	l.AddBuckets(rules...)
	return l.Register()
}
