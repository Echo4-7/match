package serializer

import (
	"Fire/config"
	"Fire/model"
)

// User 返回给前端的结构体
type User struct {
	UserID      int64  `json:"userId"`
	Email       string `json:"email"`
	NickName    string `json:"nick_name"`
	Gender      string `json:"gender"`
	TelNum      string `json:"telNum"`
	Status      string `json:"status"`
	Location    string `json:"location"`
	Description string `json:"description"`
	Avatar      string `json:"avatar"`
}

func BuildUser(user *model.User) *User {
	return &User{
		UserID:      user.UserId,
		Email:       user.Email,
		NickName:    user.NickName,
		Status:      user.Status,
		Avatar:      config.Config.Path.PhotoHost + config.Config.System.HttpPort + config.Config.Path.AvatarPath + user.Avatar,
		Gender:      user.Gender,
		TelNum:      user.TelNum,
		Location:    user.Location,
		Description: user.Description,
	}
}
