package product

import (
	"mini_go/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func List(c *gin.Context) {
	response.Success(c, gin.H{
		"products": []string{"apple", "banana", "milk"},
	})
	zap.L().Info("product list fetched")
}

func Buy(c *gin.Context) {
	var req BuyRequest
	if err := c.ShouldBind(&req); err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}
	response.Success(c, gin.H{
		"id":  req.ID,
		"qty": req.Qty,
	})
}
