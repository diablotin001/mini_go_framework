package product

import (
    "mini_go/response"

    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

func List(c *gin.Context) {
    list, err := ListService()
    if err != nil {
        c.Error(err)
        return
    }
    response.Success(c, gin.H{"products": list})
    zap.L().Info("product list fetched")
}

func Buy(c *gin.Context) {
    var req BuyRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.Error(err).SetType(gin.ErrorTypeBind)
        return
    }
    if err := BuyService(uint(req.ID), req.Qty); err != nil {
        c.Error(err)
        return
    }
    response.Success(c, gin.H{"id": req.ID, "qty": req.Qty})
}
