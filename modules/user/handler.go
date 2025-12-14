package user

import (
	"mini_go/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBind(&req); err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}
	response.Success(c, gin.H{
		"user": req.Username,
	})
	zap.L().Info("user logged in", zap.String("user", req.Username))
}

func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBind(&req); err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}
	response.Success(c, gin.H{
		"msg": "registered",
	})
}
