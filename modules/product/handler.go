package product

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

func List(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "products": []string{"apple", "banana", "milk"},
    })
}

func Buy(c *gin.Context) {
    var req BuyRequest
    if err := c.ShouldBind(&req); err != nil {
        c.Error(err)
        return
    }
    c.JSON(http.StatusOK, gin.H{
        "msg":     "bought",
        "product": req.ID,
        "qty":     req.Qty,
    })
}
