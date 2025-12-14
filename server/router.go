package server

import (
    "net/http"
    "time"
    "mini_go/middleware"
    "mini_go/modules/product"
    "mini_go/modules/user"

    "github.com/gin-gonic/gin"
)

func NewHTTPServer() *http.Server {
    r := gin.New()
    r.Use(gin.Recovery())
    r.Use(middleware.ErrorHandler())
    r.Use(middleware.Validator())

    userGroup := r.Group("/user")
    {
        userGroup.POST("/login", user.Login)
        userGroup.POST("/register", user.Register)
    }

    productGroup := r.Group("/product")
    {
        productGroup.GET("/list", product.List)
        productGroup.POST("/buy", product.Buy)
    }

    return &http.Server{
        Addr:              ":8080",
        Handler:           r,
        ReadTimeout:       5 * time.Second,
        WriteTimeout:      10 * time.Second,
        ReadHeaderTimeout: 2 * time.Second,
    }
}
