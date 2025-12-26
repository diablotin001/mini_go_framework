package user

import (
	"mini_go/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}
	token, err := LoginService(req.Username, req.Password)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, gin.H{"token": token})
	zap.L().Info("user logged in", zap.String("user", req.Username))
}

func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}
	if err := RegisterService(req.Username, req.Password, req.Email); err != nil {
		c.Error(err)
		return
	}
	response.Success(c, nil)
}
