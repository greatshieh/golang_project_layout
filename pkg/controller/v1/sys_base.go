package system

import (
	"golang_project_layout/pkg/errcode"
	"golang_project_layout/pkg/global"
	"golang_project_layout/pkg/model/common/response"
	"golang_project_layout/pkg/model/system"
	systemReq "golang_project_layout/pkg/model/system/request"
	systemRes "golang_project_layout/pkg/model/system/response"
	"golang_project_layout/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/marmotedu/errors"
)

type BaseController struct{}

// Login
// @Tags     Base
// @Summary  用户登录
// @Produce  application/json
// @Param    data  body  systemReq.Login true  "用户名, 密码, 验证码"
// @Success  200   {object}  response.Response{data=systemRes.LoginResponse,msg=string}  "返回包括用户信息,token,过期时间"
// @Router   /base/login [post]
func (b *BaseController) Login(c *gin.Context) {
	var l systemReq.Login
	err := c.ShouldBindJSON(&l)
	// key := c.ClientIP()

	if err != nil {
		response.WriteResponse(c, errors.WithCode(errcode.ErrBind, err.Error()), nil)
		return
	}
	err = utils.Verify(l, utils.LoginVerify)
	if err != nil {
		response.WriteResponse(c, errors.WithCode(errcode.ErrValidation, err.Error()), nil)
		return
	}

	// 判断验证码是否开启
	// openCaptcha := global.GVA_CONFIG.Captcha.OpenCaptcha               // 是否开启防爆次数
	// openCaptchaTimeOut := global.GVA_CONFIG.Captcha.OpenCaptchaTimeOut // 缓存超时时间
	// v, ok := global.BlackCache.Get(key)
	// if !ok {
	// 	global.BlackCache.Set(key, 1, time.Second*time.Duration(openCaptchaTimeOut))
	// }

	// var oc bool = openCaptcha == 0 || openCaptcha < interfaceToInt(v)

	// if !oc || store.Verify(l.CaptchaId, l.Captcha, true) {
	u := &system.SysUser{Username: l.Username, Password: l.Password}
	user, err := userService.Login(u)
	if err != nil {
		// 登陆失败! 用户名不存在或者密码错误!
		// 验证码次数+1
		// global.BlackCache.Increment(key, 1)
		response.WriteResponse(c, errors.WithCode(errcode.ErrUserNotFound, "用户名不存在或者密码错误. "+err.Error()), nil)
		return
	}
	if user.Enable != 1 {
		// 登陆失败! 用户被禁止登录!
		// 验证码次数+1
		// global.BlackCache.Increment(key, 1)
		response.WriteResponse(c, errors.WithCode(errcode.ErrUserForbidden, "用户被禁止登录"), nil)
		return
	}
	b.TokenNext(c, *user)
	// return
	// }
	// 验证码次数+1
	// global.BlackCache.Increment(key, 1)
	// response.FailWithMessage("验证码错误", c)
}

// TokenNext 登录以后签发jwt
func (b *BaseController) TokenNext(c *gin.Context, user system.SysUser) {
	j := &utils.JWT{SigningKey: []byte(global.GVA_CONFIG.JWT.SigningKey)} // 唯一签名
	claims := j.CreateClaims(systemReq.BaseClaims{
		UUID:        user.UUID,
		ID:          user.ID,
		NickName:    user.NickName,
		Username:    user.Username,
		AuthorityId: user.AuthorityId,
	})
	token, err := j.CreateToken(claims)
	if err != nil {
		// global.GVA_LOG.Error("获取token失败!", zap.Error(err))
		response.WriteResponse(c, err, nil)
		return
	}
	if !global.GVA_CONFIG.System.UseMultipoint {
		response.WriteResponse(c, nil, systemRes.LoginResponse{
			User:      user,
			Token:     token,
			ExpiresAt: claims.StandardClaims.ExpiresAt * 1000,
		}, "登录成功")

		// response.OkWithDetailed(systemRes.LoginResponse{
		// 	User:      user,
		// 	Token:     token,
		// 	ExpiresAt: claims.StandardClaims.ExpiresAt * 1000,
		// }, "登录成功", c)
		return
	}

	if jwtStr, err := jwtService.GetRedisJWT(user.Username); err == redis.Nil {
		if err := jwtService.SetRedisJWT(token, user.Username); err != nil {
			// global.GVA_LOG.Error("设置登录状态失败!", zap.Error(err))
			// response.FailWithMessage("设置登录状态失败", c)
			response.WriteResponse(c, err, nil, "设置登录状态失败")
			return
		}
		response.WriteResponse(c, nil, systemRes.LoginResponse{
			User:      user,
			Token:     token,
			ExpiresAt: claims.StandardClaims.ExpiresAt * 1000,
		}, "登录成功")
		// response.OkWithDetailed(systemRes.LoginResponse{
		// 	User:      user,
		// 	Token:     token,
		// 	ExpiresAt: claims.StandardClaims.ExpiresAt * 1000,
		// }, "登录成功", c)
	} else if err != nil {
		// global.GVA_LOG.Error("设置登录状态失败!", zap.Error(err))
		// response.FailWithMessage("设置登录状态失败", c)
		response.WriteResponse(c, err, nil, "设置登录状态失败")
	} else {
		var blackJWT system.JwtBlacklist
		blackJWT.Jwt = jwtStr
		if err := jwtService.JsonInBlacklist(blackJWT); err != nil {
			// response.FailWithMessage("jwt作废失败", c)
			response.WriteResponse(c, errors.WithCode(errcode.ErrExpired, err.Error()), nil)
			return
		}
		if err := jwtService.SetRedisJWT(token, user.Username); err != nil {
			// response.FailWithMessage("设置登录状态失败", c)
			response.WriteResponse(c, err, nil, "设置登录状态失败")
			return
		}
		response.WriteResponse(c, nil, systemRes.LoginResponse{
			User:      user,
			Token:     token,
			ExpiresAt: claims.StandardClaims.ExpiresAt * 1000,
		}, "登录成功")
		// response.OkWithDetailed(systemRes.LoginResponse{
		// 	User:      user,
		// 	Token:     token,
		// 	ExpiresAt: claims.StandardClaims.ExpiresAt * 1000,
		// }, "登录成功", c)
	}
}
