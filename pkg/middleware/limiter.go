package middleware

import (
	"golang_project_layout/pkg/errcode"
	"golang_project_layout/pkg/model/common/response"
	"golang_project_layout/pkg/plugin/limiter"

	"github.com/gin-gonic/gin"
	"github.com/marmotedu/errors"
)

func RateLimiter(l limiter.LimiterIface, rule limiter.LimiterBucketRules) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := l.Key(c)
		bucket, ok := l.GetBucket(key)
		if !ok {
			rule.Key = key
			l.AddBuckets(rule)
			bucket, _ = l.GetBucket(key)
		}

		count := bucket.TakeAvailable(1)
		if count == 0 {
			response.WriteResponse(c, errors.WithCode(errcode.ErrTooManyRequests, "限流等待"), nil)
			c.Abort()
			return
		}

		c.Next()
	}
}
