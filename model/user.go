package model

import (
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(255);default:'';not null;unique;index"`
	Password string `gorm:"type:varchar(255);default:'';not null"`
}

// 加密密码
func (user *User) EncryptPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	return nil
}

// 检查密码
func (user *User) CheckPassword(hashPassword, password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password)); err == nil{
		return true
	}
	return false
}
