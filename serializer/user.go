package serializer

import (
	"Fire/config"
	"Fire/model"
)

// User 返回给前端的结构体
type User struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	NickName string `json:"nick_name"`
	Status   string `json:"status"`
	Avatar   string `json:"avatar"`
}

func BuildUser(user *model.User) *User {
	return &User{
		ID:       user.UserId,
		Email:    user.Email,
		NickName: user.NickName,
		Status:   user.Status,
		Avatar:   config.Config.Path.PhotoHost + config.Config.System.HttpPort + config.Config.Path.AvatarPath + user.Avatar,
	}
}
