package v1

import (
	"Fire/pkg/util"
	"Fire/pkg/util/log"
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
	var uploadAvatar service.UserService
	_, header, err := c.Request.FormFile("file")
	if err != nil {
		log.LogrusObj.Infoln("Failed to retrieve file: ", err)
		return
	}
	claim, _ := util.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&uploadAvatar); err == nil {
		res := uploadAvatar.UploadAvatar(c.Request.Context(), claim.UserID, header)
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
	switch status {
	case service.Register:
		res := userSendCheckCode.SendCheckCode(email, service.Register)
		c.JSON(http.StatusOK, res)
	case service.Find:
		res := userSendCheckCode.SendCheckCode(email, service.Find)
		c.JSON(http.StatusOK, res)
	case service.Modify:
		res := userSendCheckCode.SendCheckCode(email, service.Modify)
		c.JSON(http.StatusOK, res)
	}
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

// UserUpdateEmail 更改邮箱
func UserUpdateEmail(c *gin.Context) {
	var userUpdate service.UserModify
	claims, _ := util.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&userUpdate); err == nil {
		res := userUpdate.UpdateEmail(c.Request.Context(), claims.UserID)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, err)
	}
}

// UserUpdateTelNum 更改手机号
func UserUpdateTelNum(c *gin.Context) {
	var userUpdate service.UserModify
	claims, _ := util.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&userUpdate); err == nil {
		res := userUpdate.UpdateTelNum(c.Request.Context(), claims.UserID)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, err)
	}
}

// UserUpdatePwd 更改密码
func UserUpdatePwd(c *gin.Context) {
	var userUpdatePwd service.UserModify
	claim, _ := util.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&userUpdatePwd); err == nil {
		res := userUpdatePwd.UpdatePwd(c.Request.Context(), claim.UserID)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, err)
	}
}

// UserInfo 获取用户信息
func UserInfo(c *gin.Context) {
	var userService service.UserService
	claim, err := util.ParseToken(c.GetHeader("Authorization"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	} else {
		res := userService.Info(c.Request.Context(), claim.UserID)
		c.JSON(http.StatusOK, res)
	}
}
