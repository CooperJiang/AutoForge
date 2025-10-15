package user

import (
	"auto-forge/internal/dto/request"
	"auto-forge/internal/dto/response"
	"auto-forge/internal/services/user"
	"auto-forge/pkg/common"
	"auto-forge/pkg/errors"

	"github.com/gin-gonic/gin"
)


func Register(c *gin.Context) {

	req, err := common.ValidateRequest[request.RegisterRequest](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	if err := user.RegisterUser(req.Username, req.Email, req.Password, req.Code); err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, nil, "注册成功")
}


func Login(c *gin.Context) {

	req, err := common.ValidateRequest[request.LoginRequest](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	userInfo, token, expiresAt, err := user.Login(req.Account, req.Password)
	if err != nil {
		errors.HandleError(c, err)
		return
	}


	username, _ := userInfo["username"].(string)
	email, _ := userInfo["email"].(string)
	avatar, _ := userInfo["avatar"].(string)
	status, _ := userInfo["status"].(int)
	role, _ := userInfo["role"].(int)

	resp := response.LoginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User: response.UserInfo{
			Username: username,
			Email:    email,
			Avatar:   avatar,
			Status:   status,
			Role:     role,
		},
	}

	errors.ResponseSuccess(c, resp, "登录成功")
}


func GetUserInfo(c *gin.Context) {

	userID, exists := c.Get("user_id")
	if !exists {
		errors.HandleError(c, errors.New(errors.CodeUnauthorized, "未授权"))
		return
	}

	userInfo, err := user.GetUserInfo(userID.(string))
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, userInfo, "获取用户信息成功")
}


func UpdateProfile(c *gin.Context) {

	userID, exists := c.Get("user_id")
	if !exists {
		errors.HandleError(c, errors.New(errors.CodeUnauthorized, "未授权"))
		return
	}

	req, err := common.ValidateRequest[request.UpdateProfileRequest](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	userInfo, err := user.UpdateProfile(userID.(string), req.Username, req.Email, req.Avatar, req.Code)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, userInfo, "更新资料成功")
}


func ChangePassword(c *gin.Context) {

	userID, exists := c.Get("user_id")
	if !exists {
		errors.HandleError(c, errors.New(errors.CodeUnauthorized, "未授权"))
		return
	}

	req, err := common.ValidateRequest[request.ChangePasswordRequest](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	if err := user.ChangePassword(userID.(string), req.OldPassword, req.NewPassword); err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, nil, "密码修改成功")
}


func SendRegistrationCode(c *gin.Context) {
	req, err := common.ValidateRequest[request.SendCodeRequest](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	if err := user.SendRegistrationCode(req.Email); err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, nil, "验证码已发送")
}


func SendResetPasswordCode(c *gin.Context) {
	req, err := common.ValidateRequest[request.SendCodeRequest](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	if err := user.SendResetPasswordCode(req.Email); err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, nil, "验证码已发送")
}


func ResetPassword(c *gin.Context) {
	req, err := common.ValidateRequest[request.ResetPasswordRequest](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	if err := user.ResetPassword(req.Email, req.Code, req.NewPassword); err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, nil, "密码重置成功")
}


func SendChangeEmailCode(c *gin.Context) {

	userID, exists := c.Get("user_id")
	if !exists {
		errors.HandleError(c, errors.New(errors.CodeUnauthorized, "未授权"))
		return
	}

	req, err := common.ValidateRequest[request.SendCodeRequest](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	if err := user.SendChangeEmailCode(userID.(string), req.Email); err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, nil, "验证码已发送")
}
