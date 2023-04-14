package middleware

import (
	"fmt"
	"golang_project_layout/pkg/errcode"
	"golang_project_layout/pkg/global"
	"golang_project_layout/pkg/model/common/response"
	"golang_project_layout/pkg/plugin/email"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/marmotedu/errors"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				global.GVA_LOG.Error(fmt.Sprintf("panic recover err: %-v", err))
				emailService := new(email.EmailService)
				err := emailService.SendEmail(global.GVA_CONFIG.Email.To, fmt.Sprintf("异常发生. 时间: %s", time.Now().Local().Format("2006-01-02 15:04:05")), fmt.Sprintf("错误信息: %-v", err))
				if err != nil {
					response.WriteResponse(c, errors.WithCode(errcode.ErrEmailSend, err.Error()), nil)
					c.Abort()
				} else {
					response.WriteResponse(c, errors.WithCode(errcode.ErrInternalServer, "内部异常"), nil)
				}
			}
		}()
		c.Next()
	}
}
