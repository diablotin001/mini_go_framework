package middleware

import (
    "net/http"
    "runtime/debug"

    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

func ZapRecovery() gin.HandlerFunc {
    return func(c *gin.Context) {
        defer func() {
            if r := recover(); r != nil {
                zap.L().Error("panic recovered",
                    zap.Any("error", r),
                    zap.ByteString("stack", debug.Stack()),
                )
                c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
                    "code": 10000,
                    "msg":  "internal server error",
                })
            }
        }()
        c.Next()
    }
}
