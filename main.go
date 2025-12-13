package main

import (
    "log"
    "mini_go/server"
)

func main() {
    srv := server.NewHTTPServer()

    go func() {
        log.Println("Server running at :8080")
        if err := srv.ListenAndServe(); err != nil && err.Error() != "http: Server closed" {
            log.Fatalf("Listen error: %s\n", err)
        }
    }()

    server.WaitForShutdown(srv)
}
