package cache

import (
	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
)

var (
	RedisClient *redis.Client
	RedisDb     string
	RedisAddr   string
	RedisPw     string
	RedisDbName int
)

func Init() {
	// 从本地读取环境
	file, err := ini.Load("./conf/config.ini")
	if err != nil {
		log.Error("config.ini file not found", err.Error())
		panic(err)
	}
	// 导入redis配置
	LoadRedis(file)
	// redis连接
	Redis()
}

func LoadRedis(file *ini.File) {
	RedisDb = file.Section("redis").Key("RedisDb").String()
	RedisAddr = file.Section("redis").Key("RedisAddr").String()
	RedisPw = file.Section("redis").Key("RedisPw").String()
	RedisDbName, _ = file.Section("redis").Key("RedisDbName").Int()
}

func Redis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr: RedisAddr,
		//Password: RedisPw,
		DB: RedisDbName,
	})
	_, err := RedisClient.Ping().Result()
	if err != nil {
		log.Error(err.Error())
		panic(err)
	}
	log.Info("Redis 连接成功")
}
