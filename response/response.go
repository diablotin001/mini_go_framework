package response

import "github.com/gin-gonic/gin"

type SuccessResponse struct {
    Code int         `json:"code"`
    Msg  string      `json:"msg"`
    Data interface{} `json:"data,omitempty"`
}

func Success(c *gin.Context, data interface{}) {
    c.JSON(200, SuccessResponse{
        Code: 0,
        Msg:  "success",
        Data: data,
    })
}
