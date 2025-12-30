package middleware

import (
    "mini_go/pkg/cache"
    "mini_go/pkg/config"
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v4"
)

func JWTAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        auth := c.GetHeader("Authorization")
        if auth == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 10002, "msg": "missing token"})
            return
        }
        parts := strings.SplitN(auth, " ", 2)
        if len(parts) != 2 || parts[0] != "Bearer" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 10002, "msg": "invalid token format"})
            return
        }
        tokenStr := parts[1]
        token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
            return []byte(config.Conf.JWT.Secret), nil
        })
        if err != nil || !token.Valid {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 10002, "msg": "invalid or expired token"})
            return
        }
        if claims, ok := token.Claims.(jwt.MapClaims); ok {
            if jti, ok2 := claims["jti"].(string); ok2 && jti != "" {
                if v, _ := cache.GetString("jwt:blacklist:" + jti); v != "" {
                    c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 10002, "msg": "token revoked"})
                    return
                }
            }
            c.Set("user_id", claims["uid"])
        }
        c.Next()
    }
}
