package limiter

import (
	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
)

// 基于https://github.com/juju/ratelimit实现限流器

type LimiterIface interface {
	// 获取对应的限流器的键值名称
	Key(c *gin.Context) string
	// 获取令牌桶
	GetBucket(key string) (*ratelimit.Bucket, bool)
	// 新增多个令牌桶
	AddBuckets(rules ...LimiterBucketRules)
	// 限流器执行
	Register() gin.HandlerFunc
}

type Limiter struct {
	limiterBuckets map[string]*ratelimit.Bucket
}

type LimiterBucketRules struct {
	// 自定义键值对名称
	Key string
	// 间隔多久时间放N个令牌
	FillInterval int64
	// 令牌桶容量
	Capacity int64
	// 每次到达时间间隔后所放的具体令牌数量
	Quantum int64
}
