package database

import (
    "os"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
    "mini_go/pkg/model"
)

var DB *gorm.DB

func Init(dsn string) {
    var err error
    DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
    if err != nil {
        panic(err)
    }
    if os.Getenv("APP_ENV") == "dev" {
        migrate()
    }
}

func migrate() {
    err := DB.AutoMigrate(
        &model.User{},
        &model.Product{},
    )
    if err != nil {
        panic(err)
    }
}
