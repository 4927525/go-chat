package service

import (
	"chat/conf"
	"chat/model/ws"
	"context"
	"time"
)

// 插入
func InsertMsg(database, id, content string, read uint, expire int64) (err error) {
	// 插入到MongoDB中
	collection := conf.MongoDBClient.Database(database).Collection(id) // 没有集合创建集合
	comment := ws.Trainer{
		Content:   content,
		StartTime: time.Now().Unix(),
		EndTime:   time.Now().Unix() + expire,
		Read:      read,
	}
	_, err = collection.InsertOne(context.TODO(), comment)
	return
}
