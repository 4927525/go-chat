package routers

import (
	"chat/api"
	"chat/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewRouters() *gin.Engine {
	r := gin.Default()
	// Recovery 中间件会恢复(recovers) 任何恐慌(panics) 如果存在恐慌，中间件将会写入500
	r.Use(gin.Recovery(), gin.Logger())
	v1 := r.Group("/")
	{
		v1.GET("ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, "ok")
		})
		v1.POST("user/register", api.UserRegister)
		v1.GET("ws", service.Handler)
	}

















	return r
}
