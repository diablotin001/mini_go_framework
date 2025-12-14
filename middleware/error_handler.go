package middleware

import (
	"mini_go/response"
	"strings"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors[0]
		if err.Type == gin.ErrorTypeBind || strings.Contains(err.Error(), "EOF") {
			response.Error(c, response.CodeParamError, "invalid parameters")
			return
		}

		response.Error(c, response.CodeServerError, err.Error())
	}
}
