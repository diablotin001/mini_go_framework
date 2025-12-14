package middleware

import "github.com/gin-gonic/gin"

func Validator() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()
        if len(c.Errors) > 0 {
            c.Abort()
            return
        }
    }
}
