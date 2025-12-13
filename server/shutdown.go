package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// WaitForShutdown 负责监听操作系统的中断信号 (如 Ctrl+C)，并实现 HTTP 服务器的优雅关闭。
// 优雅关闭意味着在接收到关闭信号后，服务器会停止接收新的请求，但会等待正在处理的请求完成，
// 从而避免数据丢失或服务中断。
func WaitForShutdown(srv *http.Server) {
	// 创建一个 Go 语言的 channel (通道)，用于接收操作系统信号。
	// channel 的容量设置为 1，确保能接收到至少一个信号。
	quit := make(chan os.Signal, 1)
	// 注册要监听的信号。这里我们监听 os.Interrupt (通常是 Ctrl+C) 信号。
	// 当接收到这些信号时，它们会被发送到上面创建的 'quit' channel。
	signal.Notify(quit, os.Interrupt)

	// 阻塞在这里，直到从 'quit' channel 接收到信号。
	// 这意味着程序会一直运行，直到用户按下 Ctrl+C 或收到其他中断信号。
	<-quit
	log.Println("Shutting down server...")

	// 创建一个带有超时的上下文 (Context)。
	// 这个上下文会在 5 秒后自动取消，用于限制服务器优雅关闭的最长时间。
	// 如果服务器在 5 秒内未能完成所有请求并关闭，它将被强制关闭。
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel() 确保在函数退出时调用 cancel()，释放上下文相关的资源。
	defer cancel()
	// 调用 HTTP 服务器的 Shutdown 方法，传入带有超时的上下文。
	// Shutdown 方法会尝试优雅地关闭服务器：停止接收新请求，等待现有请求完成。
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Forced shutdown:", err)
	}

	log.Println("Server exited gracefully")
}
