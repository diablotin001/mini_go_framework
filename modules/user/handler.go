package user

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
    var req LoginRequest
    if err := c.ShouldBind(&req); err != nil {
        c.Error(err)
        return
    }
    c.JSON(http.StatusOK, gin.H{
        "msg":  "login success",
        "user": req.Username,
    })
}

func Register(c *gin.Context) {
    var req RegisterRequest
    if err := c.ShouldBind(&req); err != nil {
        c.Error(err)
        return
    }
    c.JSON(http.StatusOK, gin.H{
        "msg": "registered",
    })
}
