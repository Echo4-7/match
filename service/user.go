package service

import (
	"Fire/cache"
	"Fire/dao"
	"Fire/model"
	"Fire/pkg/e"
	"Fire/pkg/snowflake"
	"Fire/pkg/util"
	"Fire/pkg/util/log"
	"Fire/serializer"
	"context"
	"fmt"
	"math/rand"
	"mime/multipart"
	"time"
)

const (
	Register = "注册"
	Find     = "找回密码"
)

type UserService struct {
	Password string `json:"password" form:"password"`
	Account  string `json:"account" form:"account"`
}

type FindPwdService struct {
	Account string `json:"account" form:"account"`
	NewPwd  string `json:"new_pwd" form:"new_pwd"`
}

type UserInfoService struct {
	NickName    string `json:"nick_name" form:"nick_name"`
	Gender      string `json:"gender" form:"gender"`
	Location    string `json:"location" form:"location"`
	Description string `json:"description" form:"description"`
}

// Register 注册
func (service *UserService) Register(ctx context.Context) serializer.Response {
	var user *model.User

	userDao := dao.NewUserDao(ctx)
	_, exist, err := userDao.IsExistByEmail(service.Account)
	if err != nil {
		log.LogrusObj.Infoln("IsExistByEmail failed:", err)
		return serializer.HandleError(e.ServerBusy)
	}
	if exist {
		return serializer.HandleError(e.ErrorExistUser)
	}
	userId := snowflake.GenID()
	// 激活用户
	user = &model.User{
		UserId: userId,
		Email:  service.Account,
		Avatar: "avatar.jpg",
		Status: model.Active,
	}

	// rsa解密
	decrypt, err := util.Decrypt(service.Password, util.GetPrivateKey())
	if err != nil {
		log.LogrusObj.Infoln("Decrypt failed:", err)
		return serializer.HandleError(e.ServerBusy)
	}

	// 密码加密
	if err = user.SetPassword(decrypt); err != nil {
		log.LogrusObj.Infoln("SetPassword failed:", err)
		return serializer.HandleError(e.ServerBusy)
	}

	// 创建用户
	err = userDao.CreateUser(user)
	if err != nil {
		log.LogrusObj.Infoln("CreateUser failed:", err)
		return serializer.HandleError(e.ServerBusy)
	}
	return serializer.HandleError(e.SUCCESS)
}

// LoginWithEmail 使用email登陆
func (service *UserService) LoginWithEmail(ctx context.Context) serializer.Response {
	var user *model.User
	code := e.SUCCESS

	userDao := dao.NewUserDao(ctx)

	// 判断用户是否存在
	user, exist, err := userDao.IsExistByEmail(service.Account)
	if err != nil {
		log.LogrusObj.Infoln("IsExistByEmail failed:", err)
		return serializer.HandleError(e.ServerBusy)
	}
	if !exist {
		code = e.ErrorNotExistUser
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Data: serializer.TokenData{
				User: user,
			},
		}
	}

	// rsa解析密码
	decrypt, err := util.Decrypt(service.Password, util.GetPrivateKey())
	if err != nil {
		log.LogrusObj.Infoln("Decrypt failed:", err)
		return serializer.HandleError(e.ServerBusy)
	}

	// 校验密码
	if user.CheckPassword(decrypt) == false {
		code = e.ErrorNotCompare
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Data: serializer.TokenData{
				User: user,
			},
		}
	}
	// http 无状态（认证，带上token)
	token, err := util.GenerateToken(user.UserId, user.Status)
	if err != nil {
		log.LogrusObj.Infoln("GenerateToken failed:", err)
		return serializer.HandleError(e.ServerBusy)
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

// LoginWithTelNum 使用telNum登陆
func (service *UserService) LoginWithTelNum(ctx context.Context) serializer.Response {
	code := e.SUCCESS
	var user *model.User

	UserDao := dao.NewUserDao(ctx)

	user, exist, err := UserDao.IsExistByTelNum(service.Account)
	if err != nil {
		log.LogrusObj.Infoln("IsExistByTelNum failed:", err)
		return serializer.HandleError(e.ServerBusy)
	}
	if !exist {
		code = e.ErrorNotExistUser
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Data: serializer.TokenData{
				User: user,
			},
		}
	}
	// rsa解析密码
	decrypt, err := util.Decrypt(service.Password, util.GetPrivateKey())
	if err != nil {
		log.LogrusObj.Infoln("Decrypt failed:", err)
		return serializer.HandleError(e.ServerBusy)
	}

	// 校验密码
	if user.CheckPassword(decrypt) == false {
		code = e.ErrorNotCompare
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Data: serializer.TokenData{
				User: user,
			},
		}
	}
	// http 无状态（认证，带上token)
	token, err := util.GenerateToken(user.UserId, user.Status)
	if err != nil {
		log.LogrusObj.Infoln("GenerateToken failed:", err)
		return serializer.HandleError(e.ServerBusy)
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

// LoginWithUserId 使用userId登陆
func (service *UserService) LoginWithUserId(ctx context.Context) serializer.Response {
	code := e.SUCCESS
	var user *model.User

	UserDao := dao.NewUserDao(ctx)

	user, exist, err := UserDao.IsExistByUserId(service.Account)
	if err != nil {
		log.LogrusObj.Infoln("IsExistByUserId failed:", err)
		return serializer.HandleError(e.ServerBusy)
	}
	if !exist {
		code = e.ErrorNotExistUser
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Data: serializer.TokenData{
				User: user,
			},
		}
	}
	// rsa解析密码
	decrypt, err := util.Decrypt(service.Password, util.GetPrivateKey())
	if err != nil {
		log.LogrusObj.Infoln("Decrypt failed:", err)
		return serializer.HandleError(e.ServerBusy)
	}

	// 校验密码
	if user.CheckPassword(decrypt) == false {
		code = e.ErrorNotCompare
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Data: serializer.TokenData{
				User: user,
			},
		}
	}
	// http 无状态（认证，带上token)
	token, err := util.GenerateToken(user.UserId, user.Status)
	if err != nil {
		log.LogrusObj.Infoln("GenerateToken failed:", err)
		return serializer.HandleError(e.ServerBusy)
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

// Update 更新用户信息 //
func (service *UserInfoService) Update(ctx context.Context, userId string) serializer.Response {
	var user *model.User
	// 找到用户
	userDao := dao.NewUserDao(ctx)

	user, err := userDao.GetUserByID(userId)
	if err != nil {
		log.LogrusObj.Infoln("GetUserByID failed:", err)
		return serializer.HandleError(e.ServerBusy)
	}
	if user == nil {
		return serializer.HandleError(e.ErrorNotExistUser)
	}

	// 昵称更新
	if service.NickName != "" {
		// 检查昵称是否被使用
		exist, err := userDao.NickNameIsExist(service.NickName)
		if err != nil {
			log.LogrusObj.Infoln("NickNameIsExist failed:", err)
			return serializer.HandleError(e.ServerBusy)
		}
		if exist {
			return serializer.HandleError(e.ErrorExistNick)
		}
		user.NickName = service.NickName
	}

	// 性别更新
	if service.Gender != "" {
		user.Gender = service.Gender
	}
	// 个性签名更新
	if service.Description != "" {
		user.Description = service.Description
	}
	// 位置更新
	if service.Location != "" {
		user.Location = service.Location
	}

	err = userDao.UpdateUserByID(user, userId)
	if err != nil {
		log.LogrusObj.Infoln("UpdateUserByID failed:", err)
		return serializer.HandleError(e.ServerBusy)
	}
	return serializer.HandleError(e.SUCCESS)
}

//func (service *UserService) UpdateTelNum(ctx context.Context, userId string, telNum string) serializer.Response {
//	var user *model.User
//	userDao := dao.NewUserDao(ctx)
//
//	user, err := userDao.GetUserByID(userId)
//	if err != nil {
//		log.LogrusObj.Infoln("GetUserByID failed:", err)
//		return serializer.HandleError(e.ServerBusy)
//	}
//	if user == nil {
//		return serializer.HandleError(e.ErrorNotExistUser)
//	}
//	exist, err := userDao.TelNumIsExist(telNum)
//	if err != nil {
//		log.LogrusObj.Infoln("TelNumIsExist failed:", err)
//		return serializer.HandleError(e.ServerBusy)
//	}
//	if exist {
//		return serializer.HandleError(e.ErrorExistTelNum)
//	}
//	user.TelNum = telNum
//
//	err = userDao.UpdateUserByID(user, userId)
//	if err != nil {
//		log.LogrusObj.Infoln("UpdateUserByID failed:", err)
//		return serializer.HandleError(e.ServerBusy)
//	}
//	return serializer.HandleError(e.SUCCESS)
//
//}
//
//func (service *UserService) UpdateEmail(ctx context.Context, userId string, email string) serializer.Response {
//	var user *model.User
//	userDao := dao.NewUserDao(ctx)
//
//	user, err := userDao.GetUserByID(userId)
//	if err != nil {
//		log.LogrusObj.Infoln("GetUserByID failed:", err)
//		return serializer.HandleError(e.ServerBusy)
//	}
//	if user == nil {
//		return serializer.HandleError(e.ErrorNotExistUser)
//	}
//	exist, err := userDao.EmailIsExist(email)
//	if err != nil {
//		log.LogrusObj.Infoln("EmailIsExist failed:", err)
//		return serializer.HandleError(e.ServerBusy)
//	}
//	if exist {
//		return serializer.HandleError(e.ErrorExistTelNum)
//	}
//	user.Email = email
//
//	err = userDao.UpdateUserByID(user, userId)
//	if err != nil {
//		log.LogrusObj.Infoln("UpdateUserByID failed:", err)
//		return serializer.HandleError(e.ServerBusy)
//	}
//	return serializer.HandleError(e.SUCCESS)
//}

// UploadAvatar 头像更新
func (service *UserService) UploadAvatar(ctx context.Context, userId string, header *multipart.FileHeader) serializer.Response {
	var user *model.User
	userDao := dao.NewUserDao(ctx)

	user, err := userDao.GetUserByID(userId)
	if err != nil {
		log.LogrusObj.Infoln("GetUserByID failed:", err)
		return serializer.HandleError(e.ServerBusy)
	}

	// 保存图片到本地
	//path, err := UploadAvatarToLocalStatic(file, user.ID, userId)
	//if err != nil {
	//	return serializer.HandleError(e.ErrorUploadFile)
	//}
	//user.Avatar = path
	//err = userDao.UpdateUserByID(user, userId)
	//if err != nil {
	//	log.LogrusObj.Infoln("UpdateUserByID failed:", err)
	//	return serializer.HandleError(e.ServerBusy)
	//}

	// 保存到oss
	var minioService MinioService
	// 检查桶是否存在
	err = minioService.EnsureBucket()
	if err != nil {
		log.LogrusObj.Infoln("EnsureBucket failed: ", err)
	}
	// 上传头像
	objectName, err := minioService.Upload(header)
	if err != nil {
		log.LogrusObj.Infoln("minio Upload failed: ", err)
		return serializer.HandleError(e.ServerBusy)
	}
	fileUrl, err := minioService.Preview(objectName)
	if err != nil {
		log.LogrusObj.Infoln("minio Preview failed: ", err)
		return serializer.HandleError(e.ServerBusy)
	}
	// 返回文件访问 URL
	//URL := fmt.Sprintf("http://%s/%s/%s", config.Config.Minio.Endpoint, config.Config.Minio.BucketName, fileUrl)
	// 保存到数据库
	user.Avatar = fileUrl
	err = userDao.UpdateUserByID(user, userId)
	if err != nil {
		log.LogrusObj.Infoln("UpdateUserByID failed:", err)
		return serializer.HandleError(e.ServerBusy)
	}
	return serializer.Response{
		Status: e.SUCCESS,
		Data:   fileUrl,
		Msg:    e.GetMsg(e.SUCCESS),
	}
}

// SendCheckCode 发送验证码
func (service *UserService) SendCheckCode(email string, status string) serializer.Response {
	var subject string

	if status == Register {
		subject = Register
	} else if status == Find {
		subject = Find
	}

	checkCode := fmt.Sprintf("%06d", rand.Intn(1000000)) // 生成 6 位数验证码

	// 做缓存
	if err := cache.RedisClient.Set("CHECK_CODE_MAIL:"+email, checkCode, 5*time.Minute).Err(); err != nil {
		log.LogrusObj.Infoln("RedisClient.Set failed:", err)
		return serializer.HandleError(e.ServerBusy)
	}
	// 发送邮件
	err := util.SendEmail(email, checkCode, subject)
	if err != nil {
		return serializer.HandleError(e.ErrorSendEmail)
	}
	return serializer.HandleError(e.SUCCESS)
}

// Check 检查验证码
func (service *UserService) Check(email string, checkCode string) serializer.Response {
	// 检查验证码
	check, _ := cache.RedisClient.Get("CHECK_CODE_MAIL:" + email).Result()
	if check != checkCode {
		return serializer.HandleError(e.ErrorCheckCode)
	}
	return serializer.HandleError(e.SUCCESS)
}

// FindPwd 找回密码
func (service *FindPwdService) FindPwd(ctx context.Context) serializer.Response {
	var user *model.User
	if service.Account == "" || service.NewPwd == "" {
		return serializer.HandleError(e.InvalidParams)
	}

	// 查询用户是否存在
	userDao := dao.NewUserDao(ctx)
	user, exist, err := userDao.IsExistByEmail(service.Account)
	if err != nil || !exist {
		return serializer.HandleError(e.ErrorNotExistUser)
	}

	// rsa解析密码
	decrypt, err := util.Decrypt(service.NewPwd, util.GetPrivateKey())
	if err != nil {
		log.LogrusObj.Infoln("Decrypt failed:", err)
		return serializer.HandleError(e.ServerBusy)
	}

	// 新密码不能和原密码一致
	if user.CheckPassword(decrypt) == true {
		return serializer.HandleError(e.ErrorComparePassword)
	}

	// 更新密码
	if err = user.SetPassword(decrypt); err != nil {
		log.LogrusObj.Infoln("SetPassword failed:", err)
		return serializer.HandleError(e.ErrorFailEncryption)
	}
	err = userDao.UpdateUserByID(user, user.UserId)
	if err != nil {
		log.LogrusObj.Infoln("UpdateUserByID failed:", err)
		return serializer.HandleError(e.ServerBusy)
	}
	return serializer.HandleError(e.SUCCESS)
}

// Info 获取用户信息
func (service *UserService) Info(ctx context.Context, userId string) serializer.Response {
	var user *model.User
	userDao := dao.NewUserDao(ctx)
	user, err := userDao.GetUserByID(userId)
	if err != nil {
		log.LogrusObj.Infoln("GetUserByID failed: ", err)
		return serializer.HandleError(e.ServerBusy)
	}
	return serializer.Response{
		Status: e.SUCCESS,
		Data:   serializer.BuildUser(user),
		Msg:    e.GetMsg(e.SUCCESS),
	}
}
