package main

import (
    "mini_go/logger"
    "mini_go/server"
    "go.uber.org/zap"
)

func main() {
    logger.InitLogger()
    srv := server.NewHTTPServer()

    go func() {
        zap.L().Info("Server started", zap.String("addr", ":8080"))
        if err := srv.ListenAndServe(); err != nil && err.Error() != "http: Server closed" {
            zap.L().Fatal("Listen error", zap.Error(err))
        }
    }()

    server.WaitForShutdown(srv)
}
