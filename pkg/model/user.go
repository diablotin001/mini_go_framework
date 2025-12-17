package model

import "time"

type User struct {
    ID        uint `gorm:"primaryKey"`
    Username  string
    Password  string
    Email     string
    CreatedAt time.Time
}
