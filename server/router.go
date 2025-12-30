package server

import (
	"mini_go/middleware"
	"mini_go/modules/product"
	"mini_go/modules/user"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func NewHTTPServer() *http.Server {
	r := gin.New()
	r.Use(middleware.ZapLogger())
	r.Use(middleware.ZapRecovery())
	r.Use(middleware.ErrorHandler())
	r.Use(middleware.Validator())

	userGroup := r.Group("/user")
	{
		userGroup.POST("/login", user.Login)
		userGroup.POST("/register", user.Register)
		userGroup.POST("/logout", user.Logout)
		userGroup.POST("/refresh", user.Refresh)
	}

	api := r.Group("/api")
	api.Use(middleware.JWTAuth())
	{
		api.GET("/product/list", product.List)
		api.POST("/product/buy", product.Buy)
	}

	return &http.Server{
		Addr:              ":8080",
		Handler:           r,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
	}
}
