package model

func migration()  {
	// 自动迁移
	Db.AutoMigrate(&User{})
}