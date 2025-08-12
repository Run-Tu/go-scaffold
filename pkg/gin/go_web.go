package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1、先定义默认的Gin引擎
	router := gin.Default()

	// 2、定义一个路由和处理器函数
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// 3、运行服务
	router.Run("0.0.0.0:8080")
}
