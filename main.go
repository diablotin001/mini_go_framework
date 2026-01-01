package main

import (
    "mini_go/internal/database"
    "mini_go/logger"
    "mini_go/pkg/cache"
    "mini_go/pkg/config"
    "mini_go/server"

    "go.uber.org/zap"
)

func main() {
    cfg, err := config.Load("config.yaml")
    if err != nil {
        panic(err)
    }

    _ = logger.Init(logger.LogConfig{Level: cfg.Logs.Level, Format: cfg.Logs.Format, File: cfg.Logs.Path, Env: cfg.App.Env})

    database.Init(cfg.DB.DSN)
    if err := cache.InitRedis(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB); err != nil {
        zap.L().Error("redis init failed", zap.Error(err))
    }

    srv := server.NewHTTPServer()
    go func() {
        zap.L().Info("Server started", zap.String("addr", cfg.Server.Addr))
        if err := srv.ListenAndServe(); err != nil && err.Error() != "http: Server closed" {
            zap.L().Fatal("Listen error", zap.Error(err))
        }
    }()
    server.WaitForShutdown(srv)
}
