package model

import "golang.org/x/crypto/bcrypt"

type User struct {
	UserId         int64
	NickName       string
	Gender         string
	TelNum         string
	Email          string
	PasswordDigest string
	Location       string
	Avatar         string
	Description    string
	Status         string
}

const (
	PasswordCost        = 12 // 密码加密难度
	Pending      string = "Pending"
	Active       string = "Active"
)

// 密码加密

func (user *User) SetPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), PasswordCost)
	if err != nil {
		return err
	}
	user.PasswordDigest = string(bytes)
	return nil
}

// 校验密码

func (user *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordDigest), []byte(password))
	return err == nil
}
