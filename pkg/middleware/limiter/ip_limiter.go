package limiter

import (
	"context"
	"golang_project_layout/pkg/errcode"
	"golang_project_layout/pkg/global"
	"golang_project_layout/pkg/model/common/response"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
	"github.com/marmotedu/errors"
	"go.uber.org/zap"
)

type IPLimiter struct {
	LimitCountIP    int64
	LimitIntervalIP int64
}

func NewIPLimiter() LimiterIface {
	return &IPLimiter{LimitCountIP: global.GVA_CONFIG.System.LimitCountIP, LimitIntervalIP: global.GVA_CONFIG.System.LimitIntervalIP}
}

func (l *IPLimiter) Key(c *gin.Context) string {
	return "GVA_Limit" + c.ClientIP()
}

func (l *IPLimiter) GetBucket(key string) (*ratelimit.Bucket, bool) {
	return nil, true
}

func (l *IPLimiter) AddBuckets(rules ...LimiterBucketRules) {
	for _, rule := range rules {
		if rule.FillInterval > 0 {
			l.LimitIntervalIP = rule.FillInterval
		}

		if rule.Capacity > 0 {
			l.LimitCountIP = rule.Capacity
		}
	}
}

func (l *IPLimiter) Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := l.Key(c)
		if err := checkOrMark(key, l.LimitIntervalIP, l.LimitCountIP); err != nil {
			response.WriteResponse(c, errors.WithCode(errcode.ErrLimitedIP, err.Error()), nil)
			c.Abort()
			return
		} else {
			c.Next()
		}
	}
}

func checkOrMark(key string, interval int64, count int64) (err error) {
	// ! 判断是否开启redis, 需要在config中配置开启redis
	if global.GVA_REDIS == nil {
		return errors.New("没有开启redis")
	}
	if err = setLimitWithTime(key, count, time.Duration(interval)*time.Second); err != nil {
		global.GVA_LOG.Error("limit", zap.Error(err))
	}
	return err
}

// setLimitWithTime 设置访问次数
func setLimitWithTime(key string, limit int64, interval time.Duration) error {
	count, err := global.GVA_REDIS.Exists(context.Background(), key).Result()
	if err != nil {
		return err
	}
	if count == 0 {
		pipe := global.GVA_REDIS.TxPipeline()
		pipe.Incr(context.Background(), key)
		pipe.Expire(context.Background(), key, interval)
		_, err = pipe.Exec(context.Background())
		return err
	} else {
		// 次数
		if times, err := global.GVA_REDIS.Get(context.Background(), key).Int(); err != nil {
			return err
		} else {
			if int64(times) >= limit {
				if t, err := global.GVA_REDIS.PTTL(context.Background(), key).Result(); err != nil {
					return errors.New("请求太过频繁，请稍后再试")
				} else {
					return errors.New("请求太过频繁, 请 " + t.String() + " 秒后尝试")
				}
			} else {
				return global.GVA_REDIS.Incr(context.Background(), key).Err()
			}
		}
	}
}
