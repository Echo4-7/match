package serializer

import (
	"Fire/model"
)

// User 返回给前端的结构体
type User struct {
	ID          int    `json:"id"`
	UserID      string `json:"userId"`
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
		ID:          user.ID,
		UserID:      user.UserId,
		Email:       user.Email,
		NickName:    user.NickName,
		Status:      user.Status,
		Avatar:      user.Avatar,
		Gender:      user.Gender,
		TelNum:      user.TelNum,
		Location:    user.Location,
		Description: user.Description,
	}
}
