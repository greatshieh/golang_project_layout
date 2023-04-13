package system

import (
	"golang_project_layout/pkg/errcode"
	"golang_project_layout/pkg/global"
	"golang_project_layout/pkg/model/common/request"
	"golang_project_layout/pkg/model/common/response"
	"golang_project_layout/pkg/utils"
	"strconv"

	"golang_project_layout/pkg/model/system"

	systemReq "golang_project_layout/pkg/model/system/request"
	systemRes "golang_project_layout/pkg/model/system/response"

	"github.com/gin-gonic/gin"
	"github.com/marmotedu/errors"
	"go.uber.org/zap"
)

type UserController struct{}

// Register
// @Tags     SysUser
// @Summary  用户注册账号
// @Produce   application/json
// @Param    data  body      systemReq.Register                                            true  "用户名, 昵称, 密码, 角色ID"
// @Success  200   {object}  response.Response{data=systemRes.SysUserResponse,msg=string}  "用户注册账号,返回包括用户信息"
// @Router   /user/admin_register [post]
func (b *UserController) Register(c *gin.Context) {
	var r systemReq.Register
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.WriteResponse(c, errors.WithCode(errcode.ErrBind, err.Error()), nil)
		return
	}
	err = utils.Verify(r, utils.RegisterVerify)
	if err != nil {
		// response.FailWithMessage(err.Error(), c)
		response.WriteResponse(c, errors.WithCode(errcode.ErrValidation, err.Error()), nil)
		return
	}
	var authorities []system.SysAuthority
	for _, v := range r.AuthorityIds {
		authorities = append(authorities, system.SysAuthority{
			AuthorityId: v,
		})
	}
	user := &system.SysUser{Username: r.Username, NickName: r.NickName, Password: r.Password, HeaderImg: r.HeaderImg, AuthorityId: r.AuthorityId, Authorities: authorities, Enable: r.Enable, Phone: r.Phone, Email: r.Email}
	userReturn, err := userService.Register(*user)
	if err != nil {
		global.GVA_LOG.Error("注册失败!", zap.Error(err))
		// response.FailWithDetailed(systemRes.SysUserResponse{User: userReturn}, "注册失败", c)
		response.WriteResponse(c, err, nil, "注册失败")
		return
	}
	// response.OkWithDetailed(systemRes.SysUserResponse{User: userReturn}, "注册成功", c)
	response.WriteResponse(c, nil, systemRes.SysUserResponse{User: userReturn}, "注册成功")
}

// ChangePassword
// @Tags      SysUser
// @Summary   用户修改密码
// @Security  ApiKeyAuth
// @Produce  application/json
// @Param     data  body      systemReq.ChangePasswordReq    true  "用户名, 原密码, 新密码"
// @Success   200   {object}  response.Response{msg=string}  "用户修改密码"
// @Router    /user/changePassword [post]
func (b *UserController) ChangePassword(c *gin.Context) {
	var req systemReq.ChangePasswordReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		// response.FailWithMessage(err.Error(), c)
		response.WriteResponse(c, errors.WithCode(errcode.ErrBind, err.Error()), nil)
		return
	}
	err = utils.Verify(req, utils.ChangePasswordVerify)
	if err != nil {
		// response.FailWithMessage(err.Error(), c)
		response.WriteResponse(c, errors.WithCode(errcode.ErrValidation, err.Error()), nil)
		return
	}
	uid := utils.GetUserID(c)
	u := &system.SysUser{GVA_MODEL: global.GVA_MODEL{ID: uid}, Password: req.Password}
	_, err = userService.ChangePassword(u, req.NewPassword)
	if err != nil {
		global.GVA_LOG.Error("修改失败!", zap.Error(err))
		// response.FailWithMessage("修改失败，原密码与当前账户不符", c)
		response.WriteResponse(c, errors.WithCode(errcode.ErrPasswordIncorrect, "修改失败，原密码与当前账户不符"), nil)
		return
	}
	// response.OkWithMessage("修改成功", c)
	response.WriteResponse(c, nil, nil, "修改成功")
}

// GetUserList
// @Tags      SysUser
// @Summary   分页获取用户列表
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.PageInfo                                        true  "页码, 每页大小"
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}  "分页获取用户列表,返回包括列表,总数,页码,每页数量"
// @Router    /user/getUserList [post]
func (b *UserController) GetUserList(c *gin.Context) {
	var pageInfo request.PageInfo
	err := c.ShouldBindJSON(&pageInfo)
	if err != nil {
		// response.FailWithMessage(err.Error(), c)
		response.WriteResponse(c, errors.WithCode(errcode.ErrBind, err.Error()), nil)
		return
	}
	err = utils.Verify(pageInfo, utils.PageInfoVerify)
	if err != nil {
		// response.FailWithMessage(err.Error(), c)
		response.WriteResponse(c, errors.WithCode(errcode.ErrValidation, err.Error()), nil)
		return
	}
	list, total, err := userService.GetUserInfoList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		// response.FailWithMessage("获取失败", c)
		response.WriteResponse(c, err, nil, "获取失败")
		return
	}

	response.WriteResponse(c, nil, response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功")
	// response.OkWithDetailed(response.PageResult{
	// 	List:     list,
	// 	Total:    total,
	// 	Page:     pageInfo.Page,
	// 	PageSize: pageInfo.PageSize,
	// }, "获取成功", c)
}

// SetUserAuthority
// @Tags      SysUser
// @Summary   更改用户权限
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      systemReq.SetUserAuth          true  "用户UUID, 角色ID"
// @Success   200   {object}  response.Response{msg=string}  "设置用户权限"
// @Router    /user/setUserAuthority [post]
func (b *UserController) SetUserAuthority(c *gin.Context) {
	var sua systemReq.SetUserAuth
	err := c.ShouldBindJSON(&sua)
	if err != nil {
		// response.FailWithMessage(err.Error(), c)
		response.WriteResponse(c, errors.WithCode(errcode.ErrBind, err.Error()), nil)
		return
	}
	if UserVerifyErr := utils.Verify(sua, utils.SetUserAuthorityVerify); UserVerifyErr != nil {
		// response.FailWithMessage(UserVerifyErr.Error(), c)
		response.WriteResponse(c, errors.WithCode(errcode.ErrValidation, err.Error()), nil)
		return
	}
	userID := utils.GetUserID(c)
	err = userService.SetUserAuthority(userID, sua.AuthorityId)
	if err != nil {
		global.GVA_LOG.Error("修改失败!", zap.Error(err))
		// response.FailWithMessage(err.Error(), c)
		response.WriteResponse(c, err, nil, "修改失败")
		return
	}
	claims := utils.GetUserInfo(c)
	j := &utils.JWT{SigningKey: []byte(global.GVA_CONFIG.JWT.SigningKey)} // 唯一签名
	claims.AuthorityId = sua.AuthorityId
	if token, err := j.CreateToken(*claims); err != nil {
		global.GVA_LOG.Error("修改失败!", zap.Error(err))
		// response.FailWithMessage(err.Error(), c)
		response.WriteResponse(c, err, nil, "修改失败")
	} else {
		c.Header("new-token", token)
		c.Header("new-expires-at", strconv.FormatInt(claims.ExpiresAt, 10))
		// response.OkWithMessage("修改成功", c)
		response.WriteResponse(c, nil, nil, "修改成功")
	}
}

// SetUserAuthorities
// @Tags      SysUser
// @Summary   设置用户权限
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      systemReq.SetUserAuthorities   true  "用户UUID, 角色ID"
// @Success   200   {object}  response.Response{msg=string}  "设置用户权限"
// @Router    /user/setUserAuthorities [post]
func (b *UserController) SetUserAuthorities(c *gin.Context) {
	var sua systemReq.SetUserAuthorities
	err := c.ShouldBindJSON(&sua)
	if err != nil {
		// response.FailWithMessage(err.Error(), c)
		response.WriteResponse(c, errors.WithCode(errcode.ErrBind, err.Error()), nil)
		return
	}
	err = userService.SetUserAuthorities(sua.ID, sua.AuthorityIds)
	if err != nil {
		global.GVA_LOG.Error("修改失败!", zap.Error(err))
		// response.FailWithMessage("修改失败", c)
		response.WriteResponse(c, err, nil, "修改失败")
		return
	}
	// response.OkWithMessage("修改成功", c)
	response.WriteResponse(c, nil, nil, "修改成功")
}

// DeleteUser
// @Tags      SysUser
// @Summary   删除用户
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.GetById                true  "用户ID"
// @Success   200   {object}  response.Response{msg=string}  "删除用户"
// @Router    /user/deleteUser [delete]
func (b *UserController) DeleteUser(c *gin.Context) {
	var reqId request.GetById
	err := c.ShouldBindJSON(&reqId)
	if err != nil {
		// response.FailWithMessage(err.Error(), c)
		response.WriteResponse(c, err, nil)
		return
	}
	err = utils.Verify(reqId, utils.IdVerify)
	if err != nil {
		// response.FailWithMessage(err.Error(), c)
		response.WriteResponse(c, err, nil)
		return
	}
	jwtId := utils.GetUserID(c)
	if jwtId == uint(reqId.ID) {
		// response.FailWithMessage("删除失败", c)
		response.WriteResponse(c, err, nil, "删除失败")
		return
	}
	err = userService.DeleteUser(reqId.ID)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		// response.FailWithMessage("删除失败", c)
		response.WriteResponse(c, err, nil, "删除失败")
		return
	}
	// response.OkWithMessage("删除成功", c)
	response.WriteResponse(c, nil, nil, "删除成功")
}

// SetUserInfo
// @Tags      SysUser
// @Summary   设置用户信息
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      system.SysUser                                             true  "ID, 用户名, 昵称, 头像链接"
// @Success   200   {object}  response.Response{data=map[string]interface{},msg=string}  "设置用户信息"
// @Router    /user/setUserInfo [put]
func (b *UserController) SetUserInfo(c *gin.Context) {
	var user systemReq.ChangeUserInfo
	err := c.ShouldBindJSON(&user)
	if err != nil {
		// response.FailWithMessage(err.Error(), c)
		response.WriteResponse(c, errors.WithCode(errcode.ErrBind, err.Error()), nil)
		return
	}
	err = utils.Verify(user, utils.IdVerify)
	if err != nil {
		// response.FailWithMessage(err.Error(), c)
		response.WriteResponse(c, errors.WithCode(errcode.ErrValidation, err.Error()), nil)
		return
	}

	if len(user.AuthorityIds) != 0 {
		err = userService.SetUserAuthorities(user.ID, user.AuthorityIds)
		if err != nil {
			global.GVA_LOG.Error("设置失败!", zap.Error(err))
			// response.FailWithMessage("设置失败", c)
			response.WriteResponse(c, err, nil, "设置失败")
			return
		}
	}
	err = userService.SetUserInfo(system.SysUser{
		GVA_MODEL: global.GVA_MODEL{
			ID: user.ID,
		},
		NickName:  user.NickName,
		HeaderImg: user.HeaderImg,
		Phone:     user.Phone,
		Email:     user.Email,
		SideMode:  user.SideMode,
		Enable:    user.Enable,
	})
	if err != nil {
		global.GVA_LOG.Error("设置失败!", zap.Error(err))
		// response.FailWithMessage("设置失败", c)
		response.WriteResponse(c, err, nil, "设置失败")
		return
	}
	// response.OkWithMessage("设置成功", c)
	response.WriteResponse(c, nil, nil, "设置成功")
}

// SetSelfInfo
// @Tags      SysUser
// @Summary   设置用户信息
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      system.SysUser                                             true  "ID, 用户名, 昵称, 头像链接"
// @Success   200   {object}  response.Response{data=map[string]interface{},msg=string}  "设置用户信息"
// @Router    /user/SetSelfInfo [put]
func (b *UserController) SetSelfInfo(c *gin.Context) {
	var user systemReq.ChangeUserInfo
	err := c.ShouldBindJSON(&user)
	if err != nil {
		// response.FailWithMessage(err.Error(), c)
		response.WriteResponse(c, errors.WithCode(errcode.ErrBind, err.Error()), nil)
		return
	}
	user.ID = utils.GetUserID(c)
	err = userService.SetSelfInfo(system.SysUser{
		GVA_MODEL: global.GVA_MODEL{
			ID: user.ID,
		},
		NickName:  user.NickName,
		HeaderImg: user.HeaderImg,
		Phone:     user.Phone,
		Email:     user.Email,
		SideMode:  user.SideMode,
		Enable:    user.Enable,
	})
	if err != nil {
		global.GVA_LOG.Error("设置失败!", zap.Error(err))
		// response.FailWithMessage("设置失败", c)
		response.WriteResponse(c, err, nil, "设置失败")
		return
	}
	// response.OkWithMessage("设置成功", c)
	response.WriteResponse(c, nil, nil, "设置成功")
}

// GetUserInfo
// @Tags      SysUser
// @Summary   获取用户信息
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Success   200  {object}  response.Response{data=map[string]interface{},msg=string}  "获取用户信息"
// @Router    /user/getUserInfo [get]
func (b *UserController) GetUserInfo(c *gin.Context) {
	uuid := utils.GetUserUuid(c)
	ReqUser, err := userService.GetUserInfo(uuid)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		// response.FailWithMessage("获取失败", c)
		response.WriteResponse(c, err, nil, "获取失败")
		return
	}
	// response.OkWithDetailed(gin.H{"userInfo": ReqUser}, "获取成功", c)
	response.WriteResponse(c, nil, gin.H{"userInfo": ReqUser}, "获取成功")
}

// ResetPassword
// @Tags      SysUser
// @Summary   重置用户密码
// @Security  ApiKeyAuth
// @Produce  application/json
// @Param     data  body      system.SysUser                 true  "ID"
// @Success   200   {object}  response.Response{msg=string}  "重置用户密码"
// @Router    /user/resetPassword [post]
func (b *UserController) ResetPassword(c *gin.Context) {
	var user system.SysUser
	err := c.ShouldBindJSON(&user)
	if err != nil {
		// response.FailWithMessage(err.Error(), c)
		response.WriteResponse(c, errors.WithCode(errcode.ErrBind, err.Error()), nil)
		return
	}
	err = userService.ResetPassword(user.ID)
	if err != nil {
		global.GVA_LOG.Error("重置失败!", zap.Error(err))
		// response.FailWithMessage("重置失败"+err.Error(), c)
		response.WriteResponse(c, err, nil, "重置失败")
		return
	}
	// response.OkWithMessage("重置成功", c)
	response.WriteResponse(c, nil, nil, "重置成功")
}
