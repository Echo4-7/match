package dao

import (
	"Fire/model"
	"context"
	"errors"
	"gorm.io/gorm"
)

type UserDao struct {
	*gorm.DB
}

// 根据上下文 ctx 创建一个新的 UserDao 实例

func NewUserDao(ctx context.Context) *UserDao {
	return &UserDao{NewDBClient(ctx)}
}

func NewUserDaoByDB(db *gorm.DB) *UserDao {
	return &UserDao{db}
}

// IsExistByEmail 邮箱检查用户是否存在
func (dao *UserDao) IsExistByEmail(email string) (user *model.User, exist bool, err error) {
	var count int64
	err = dao.DB.Model(&model.User{}).Where("email = ?", email).Find(&user).Count(&count).Error
	if count == 0 {
		return nil, false, err
	}
	return user, true, nil
}

// IsExistByTelNum 电话检查用户是否存在
func (dao *UserDao) IsExistByTelNum(telNum string) (user *model.User, exist bool, err error) {
	var count int64
	err = dao.DB.Model(&model.User{}).Where("tel_num = ?", telNum).Find(&user).Count(&count).Error
	if count == 0 {
		return nil, false, err
	}
	return user, true, nil
}

// IsExistByUserId userId检查用户是否存在
func (dao *UserDao) IsExistByUserId(userId string) (user *model.User, exist bool, err error) {
	var count int64
	err = dao.DB.Model(&model.User{}).Where("user_id = ?", userId).Find(&user).Count(&count).Error
	if count == 0 {
		return nil, false, err
	}
	return user, true, nil
}

// CreateUser 创建新用户
func (dao *UserDao) CreateUser(user *model.User) error {
	return dao.DB.Model(&model.User{}).Create(&user).Error
}

// GetUserByID 根据ID获取用户
func (dao *UserDao) GetUserByID(userId string) (user *model.User, err error) {
	if err = dao.DB.Model(&model.User{}).Where("user_id = ?", userId).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户未查询到！")
		}
		return nil, err
	}
	return user, nil
}

// UpdateUserByID 根据ID更新用户信息
func (dao *UserDao) UpdateUserByID(user *model.User, userId string) error {
	err := dao.DB.Model(&model.User{}).Where("user_id = ?", userId).Updates(&user).Error
	return err
}

// NickNameIsExist 检查用户名是否存在
func (dao *UserDao) NickNameIsExist(nickname string) (exist bool, err error) {
	var user *model.User
	if err = dao.DB.Model(&model.User{}).Where("nick_name = ?", nickname).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil

}

// EmailIsExist 检查邮箱是否存在
func (dao *UserDao) EmailIsExist(email string) (exist bool, err error) {
	var user *model.User
	if err = dao.DB.Model(&model.User{}).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil

}

// TelNumIsExist 检查电话号码是否存在
func (dao *UserDao) TelNumIsExist(telNum string) (exist bool, err error) {
	var user *model.User
	if err = dao.DB.Model(&model.User{}).Where("tel_num = ?", telNum).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil

}
