package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()

	// 生产环境应使用 Recovery() 以避免 panic 导致服务崩溃
	r.Use(gin.Recovery())

	// 示例路由：返回 JSON
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"msg": "pong"})
	})

	// 自定义高并发友好的 http.Server
	srv := &http.Server{
		Addr:              ":8080",
		Handler:           r,
		ReadTimeout:       5 * time.Second,  // 防止慢连接拖垮服务
		WriteTimeout:      10 * time.Second, // 超时避免阻塞
		ReadHeaderTimeout: 2 * time.Second,
		MaxHeaderBytes:    1 << 20, // 1 MB
	}

	srv.ListenAndServe()
}

