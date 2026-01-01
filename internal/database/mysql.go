package database

import (
    conf "mini_go/pkg/config"
    mydb "mini_go/pkg/db"
    "mini_go/pkg/model"
    "os"

    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    gormlog "gorm.io/gorm/logger"
)

var DB *gorm.DB

func Init(dsn string) {
    var err error
    gLogger := mydb.NewZapGormLogger(conf.Conf.DB.SlowThresholdMS, gormlog.LogLevel(conf.Conf.DB.LogLevel))
    DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: gLogger})
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
