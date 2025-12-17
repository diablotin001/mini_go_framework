package product

import (
	"errors"
	"fmt"
	"mini_go/internal/database"
	"mini_go/pkg/cache"
	"mini_go/pkg/model"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func ListProducts() ([]model.Product, error) {
	// 简化：显示不使用缓存
	var list []model.Product
	if err := database.DB.Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func GetProductByID(id uint) (*model.Product, error) {
	key := fmt.Sprintf("product:%d", id)
	if s, err := cache.GetString(key); err == nil && s != "" {
		return &model.Product{ID: id, Name: s}, nil
	}
	var p model.Product
	if err := database.DB.First(&p, id).Error; err != nil {
		return nil, err
	}
	_ = cache.SetString(key, p.Name, 300*time.Second)
	return &p, nil
}

func DecreaseStock(id uint, qty int) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		var p model.Product
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&p, id).Error; err != nil {
			return err
		}
		if p.Stock < qty {
			return errors.New("stock insufficient")
		}
		p.Stock -= qty
		if err := tx.Save(&p).Error; err != nil {
			return err
		}
		return nil
	})
}
