package main

import (
	"github.com/gin-gonic/gin"
	"week5/pkg/middleware"

	"net/http"
	"time"
)

func main() {
	//获取gin框架对象
	g := gin.Default()
	//注册中间件
	g.Use(middleware.Wrapper(100, 10, 0.5, time.Second))
	//注册方法
	g.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"massage": "hello",
		})
	})
	g.GET("/err", func(c *gin.Context) {
		c.Writer.WriteHeader(http.StatusNotFound)
	})
	//启动服务
	g.Run(":8080")
}
