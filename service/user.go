package service

import (
	"chat/model"
	"chat/serializer"
	"net/http"
)

type LoginOrRegisterService struct {
	Username string `json:"username" form:"username" binding:"required"`
	Password string `json:"password" form:"password"  binding:"required"`
}

// 注册
func (service *LoginOrRegisterService) Register() serializer.Response {
	var user model.User
	count := 0
	// 判断是否注册
	model.Db.Model(&model.User{}).
		Where("username = ?", service.Username).
		First(&user).
		Count(&count)
	// 已存在用户，不可注册
	if count != 0 {
		return serializer.Response{
			Status: http.StatusOK,
			Msg:    "已存在用户，请直接登录",
		}
	}
	user.Username = service.Username
	// 加密密码
	if err := user.EncryptPassword(service.Password); err != nil {
		return serializer.Response{
			Status: http.StatusInternalServerError,
			Msg:    "密码加密错误",
		}
	}

	// 存库
	if err := model.Db.Create(&user).Error; err != nil {
		return serializer.Response{
			Status: http.StatusOK,
			Msg:    "新增失败",
		}
	}

	return serializer.Response{
		Status: http.StatusOK,
		Msg:    "新增成功",
	}
}
