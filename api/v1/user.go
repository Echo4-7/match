package v1

import (
	"Fire/pkg/util"
	"Fire/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

// UserRegister 用户注册接口
func UserRegister(c *gin.Context) {
	var userRegister service.UserService
	if err := c.ShouldBind(&userRegister); err == nil {
		res := userRegister.Register(c.Request.Context())
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, err)
	}
}

// UserLogin 用户登陆接口
func UserLogin(c *gin.Context) {
	var userLogin service.UserService
	if err := c.ShouldBind(&userLogin); err == nil {
		if IsEmail(userLogin.Account) {
			res := userLogin.LoginWithEmail(c.Request.Context())
			c.JSON(http.StatusOK, res)
			return
		} else if IsTelNum(userLogin.Account) {
			res := userLogin.LoginWithTelNum(c.Request.Context())
			c.JSON(http.StatusOK, res)
			return
		} else if IsUserId(userLogin.Account) {
			res := userLogin.LoginWithUserId(c.Request.Context())
			c.JSON(http.StatusOK, res)
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, err)
	}
}

// UserUpdate 用户更新接口
func UserUpdate(c *gin.Context) {
	var userUpdate service.UserInfoService
	claims, _ := util.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&userUpdate); err == nil {
		res := userUpdate.Update(c.Request.Context(), claims.UserID)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, err)
	}
}

// UploadAvatar 上传头像
func UploadAvatar(c *gin.Context) {
	file, _, _ := c.Request.FormFile("file")
	var uploadAvatar service.UserService
	claim, _ := util.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&uploadAvatar); err == nil {
		res := uploadAvatar.UploadAvatar(c.Request.Context(), claim.UserID, file)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, err)
	}
}

// SendCheckCode 发送验证码
func SendCheckCode(c *gin.Context) {
	var userSendCheckCode service.UserService
	email := c.Query("email")
	status := c.Query("status")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "邮箱不能为空！"})
		return
	}
	res := userSendCheckCode.SendCheckCode(email, status)
	c.JSON(http.StatusOK, res)
}

// CheckCode 检验验证码
func CheckCode(c *gin.Context) {
	var userCheckCode service.UserService
	email := c.Query("email")
	checkCode := c.Query("code")
	if checkCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "验证码不能为空！"})
		return
	}
	res := userCheckCode.Check(email, checkCode)
	c.JSON(http.StatusOK, res)
}

// FindPwd 找回密码
func FindPwd(c *gin.Context) {
	var userFindPwd service.FindPwdService
	if err := c.ShouldBind(&userFindPwd); err == nil {
		res := userFindPwd.FindPwd(c.Request.Context())
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, err)
	}
}
