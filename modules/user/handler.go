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
	pair, err := LoginService(req.Username, req.Password)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, gin.H{"access_token": pair.AccessToken, "refresh_token": pair.RefreshToken})
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

func Logout(c *gin.Context) {
	var req LogoutRequest
	_ = c.ShouldBindJSON(&req)
	auth := c.GetHeader("Authorization")
	token := ""
	if auth != "" && len(auth) > 7 {
		token = auth[7:]
	}
	if err := LogoutService(token, req.RefreshToken); err != nil {
		c.Error(err)
		return
	}
	response.Success(c, nil)
}

func Refresh(c *gin.Context) {
	// var req struct {
	// 	RefreshToken string `json:"refresh_token" binding:"required"`
	// }
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}
	token, err := RefreshService(req.RefreshToken)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, gin.H{"access_token": token})
}
