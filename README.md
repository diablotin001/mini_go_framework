# 为什么它能高并发

## ** 使用 Gin 的最小化中间件（fast path）**

```go
r := gin.New()
r.Use(gin.Recovery())
```

* 如果使用 `gin.Default()` 会自动加载 Logger，中间件越多开销越高。
* 高并发服务一般只用 `Recovery()`，避免 panic 后服务退出。

---

## ** 使用 http.Server 调优**

Gin 默认的 `r.Run()` 会使用默认配置，不够强壮。
这里自定义了：

### **ReadTimeout**

限制客户端读取时间，防止 **Slowloris** 攻击（极慢连接拖死服务）。

### **WriteTimeout**

避免响应写太慢，让 goroutine 占着不释放。

### **ReadHeaderTimeout**

限制读取请求头时间，让恶意请求更难拖延。

### **MaxHeaderBytes**

限制 header 尺寸，防止攻击者丢异常大的 header 造成内存消耗。

---

## ** Go 原生 goroutine = 高并发的核心**

Go 的 HTTP Server 每个请求自动分配 goroutine，不需要你手写线程池。
理论上轻松支撑 **几十万级 QPS 的并发阻塞场景（只要逻辑轻、IO 小）**。

---
