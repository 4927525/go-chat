package main

import (
	"chat/conf"
	"chat/routers"
	"chat/service"
	_ "github.com/go-sql-driver/mysql"
)

//server
func main() {
	conf.Init()
	go service.Manager.Start()
	r := routers.NewRouters()
	r.Run(conf.HttpPort)
}
