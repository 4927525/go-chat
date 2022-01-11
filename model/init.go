package model

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"time"
)

var Db *gorm.DB

func DataBase(path string) {
	fmt.Println(path)
	db, err := gorm.Open("mysql", path)
	if err != nil {
		log.Error("数据库连接失败")
		panic(err)
	}
	log.Println("数据库连接成功")
	db.LogMode(true)
	if gin.Mode() == "release" {
		db.LogMode(false)
	}
	// 数据库配置
	db.SingularTable(true)                       // 表名不加s
	db.DB().SetMaxOpenConns(100)                 // 连接池最大连接数
	db.DB().SetMaxIdleConns(10)                  // 连接池最大空闲连接数
	db.DB().SetConnMaxIdleTime(time.Second * 30) // 超时时间
	Db = db
	// 迁移文件
	migration()
}
