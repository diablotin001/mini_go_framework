package product

import (
	"errors"
	"mini_go/pkg/model"
)

func ListService() ([]model.Product, error) {
	return ListProducts()
}

func BuyService(id uint, qty int) error {
	p, err := GetProductByID(id)
	if err != nil {
		return err
	}
	if p.Stock < qty {
		return errors.New("stock insufficient")
	}
	return DecreaseStock(id, qty)
}
