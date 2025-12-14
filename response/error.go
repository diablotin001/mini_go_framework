package response

import "github.com/gin-gonic/gin"

type ErrorResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func Error(c *gin.Context, code int, msg string) {
	c.JSON(200, ErrorResponse{
		Code: code,
		Msg:  msg,
	})
}
