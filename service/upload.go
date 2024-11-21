package service

import (
	"Fire/config"
	"io"
	"mime/multipart"
	"os"
	"strconv"
)

// UploadAvatarToLocalStatic 更新头像到本地
func UploadAvatarToLocalStatic(file multipart.File, id int, userId string) (filePath string, err error) {
	bId := strconv.Itoa(id) // 路径拼接
	basePath := "." + config.Config.Path.AvatarPath + "user" + bId + "/"
	if !DirExistOrNot(basePath) {
		CreateDir(basePath)
	}
	avatarPath := basePath + userId + ".jpg"
	content, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}
	err = os.WriteFile(avatarPath, content, 06666)
	if err != nil {
		return
	}
	return "user" + bId + "/" + userId + ".jpg", err
}

// DirExistOrNot 判断路径是否存在
func DirExistOrNot(fileAddr string) bool {
	s, err := os.Stat(fileAddr)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// CreateDir 创建文件夹
func CreateDir(dirName string) bool {
	err := os.MkdirAll(dirName, 0755)
	if err != nil {
		return false
	}
	return true
}
