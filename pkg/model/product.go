package model

import "time"

type Product struct {
    ID        uint `gorm:"primaryKey"`
    Name      string
    Stock     int
    Price     int64
    CreatedAt time.Time
}
