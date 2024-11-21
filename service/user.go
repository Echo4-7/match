package service

import (
	"Fire/cache"
	"Fire/dao"
	"Fire/model"
	"Fire/pkg/e"
	"Fire/pkg/snowflake"
	"Fire/pkg/util"
	"Fire/serializer"
	"context"
	"fmt"
	"math/rand"
	"mime/multipart"
	"strconv"
	"time"
)

const (
	Register = "注册"
	Find     = "找回密码"
)

type UserService struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

type FindPwdService struct {
	Email  string `json:"email" form:"email"`
	NewPwd string `json:"new_pwd" form:"new_pwd"`
}

// Register 注册
func (service *UserService) Register(ctx context.Context) serializer.Response {
	var user *model.User
	code := e.SUCCESS

	userDao := dao.NewUserDao(ctx)
	_, exist, err := userDao.ExistOrNotExist(service.Email)
	if err != nil {
		code = e.ERROR
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	if exist {
		code = e.ErrorExistUser
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	uid := snowflake.GenID()
	// 激活用户
	user = &model.User{
		UserId: uid,
		Email:  service.Email,
		Avatar: "avatar.jpg",
		Status: model.Active,
	}

	// rsa解密
	decrypt, err := util.Decrypt(service.Password, util.GetPrivateKey())
	if err != nil {
		code = e.ERROR
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Data:   "rsa解析密码错误！",
		}
	}

	// 密码加密
	if err = user.SetPassword(decrypt); err != nil {
		code = e.ErrorFailEncryption
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}

	// 创建用户
	err = userDao.CreateUser(user)
	if err != nil {
		code = e.ERROR
	}
	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}
}

// Login 登陆
func (service *UserService) Login(ctx context.Context) serializer.Response {
	var user *model.User
	code := e.SUCCESS

	userDao := dao.NewUserDao(ctx)

	// 判断用户是否存在
	user, exist, err := userDao.ExistOrNotExist(service.Email)
	if err != nil {
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	if !exist {
		code = e.ErrorNotExistUser
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Data:   "用户不存在，请先注册！",
		}
	}

	// rsa解析密码
	decrypt, err := util.Decrypt(service.Password, util.GetPrivateKey())
	if err != nil {
		code = e.ERROR
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Data:   "rsa解析密码错误！",
		}
	}

	// 校验密码
	if user.CheckPassword(decrypt) == false {
		code = e.ErrorNotCompare
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Data:   "密码错误，请重新输入！",
		}
	}
	// http 无状态（认证，带上token)
	token, err := util.GenerateToken(uint(user.UserId), user.Status)
	if err != nil {
		code = e.ErrorAuthToken
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
		Data: serializer.TokenData{
			User:  serializer.BuildUser(user),
			Token: token,
		},
	}
}

// Update 更新用户信息
func (service *UserService) Update(ctx context.Context, uid uint) serializer.Response {
	var user *model.User
	// 找到用户
	userDao := dao.NewUserDao(ctx)
	code := e.SUCCESS

	user, err := userDao.GetUserByID(uid)

	err = userDao.UpdateUserByID(user, uid)
	if err != nil {
		code = e.ERROR
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
		Data:   serializer.BuildUser(user),
	}
}

// Post 头像更新
func (service *UserService) Post(ctx context.Context, uid uint, file multipart.File) serializer.Response {
	code := e.SUCCESS
	var user *model.User
	userDao := dao.NewUserDao(ctx)
	user, err := userDao.GetUserByID(uid)
	if err != nil {
		code = e.ERROR
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	// 保存图片到本地
	userID := strconv.Itoa(int(user.UserId))
	path, err := UploadAvatarToLocalStatic(file, uid, userID)
	if err != nil {
		code = e.ErrorUploadFile
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	user.Avatar = path
	err = userDao.UpdateUserByID(user, uid)
	if err != nil {
		code = e.ERROR
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
		Data:   serializer.BuildUser(user),
	}
}

// SendCheckCode 发送验证码
func (service *UserService) SendCheckCode(email string, status string) serializer.Response {
	code := e.SUCCESS
	var subject string

	if status == Register {
		subject = Register
	} else if status == Find {
		subject = Find
	}

	checkCode := fmt.Sprintf("%06d", rand.Intn(1000000)) // 生成 6 位数验证码

	// 做缓存
	if err := cache.RedisClient.Set("CHECK_CODE_MAIL:"+email, checkCode, 5*time.Minute).Err(); err != nil {
		code = e.ERROR
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Data:   "存储验证码失败",
		}
	}
	// 发送邮件
	err := util.SendEmail(email, checkCode, subject)
	if err != nil {
		code = e.ErrorSendEmail
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}
}

// Check 检查验证码
func (service *UserService) Check(email string, checkCode string) serializer.Response {
	code := e.SUCCESS
	// 检查验证码
	check, _ := cache.RedisClient.Get("CHECK_CODE_MAIL:" + email).Result()
	if check != checkCode {
		code = e.ERROR
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Data:   "验证码错误",
		}
	}
	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}
}

// FindPwd 找回密码
func (service *FindPwdService) FindPwd(ctx context.Context) serializer.Response {
	var user *model.User
	code := e.SUCCESS
	if service.Email == "" || service.NewPwd == "" {
		code = e.InvalidParams
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}

	// 查询用户是否存在
	userDao := dao.NewUserDao(ctx)
	user, exist, err := userDao.ExistOrNotExist(service.Email)
	if err != nil || !exist {
		code = e.ErrorNotExistUser
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Data:   "用户不存在",
		}
	}

	// rsa解析密码
	decrypt, err := util.Decrypt(service.NewPwd, util.GetPrivateKey())
	if err != nil {
		code = e.ERROR
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Data:   "rsa解析密码错误！",
		}
	}

	// 更新密码
	if err = user.SetPassword(decrypt); err != nil {
		code = e.ErrorFailEncryption
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}

	err = userDao.UpdateUserByID(user, uint(user.UserId))
	if err != nil {
		code = e.ERROR
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}
}
