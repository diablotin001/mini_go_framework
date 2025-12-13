package product

type BuyRequest struct {
    ID  int `json:"id" binding:"required"`
    Qty int `json:"qty" binding:"required,min=1"`
}
