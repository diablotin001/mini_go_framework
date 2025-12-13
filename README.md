---
# STEP 1
---
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
# STEP 2

* ✅ **模块化目录结构（user / product 分模块）**
* ✅ **优雅关闭（graceful shutdown）**
* ✅ **中间件统一参数验证（bind + validate）**
* ✅ **极简、高并发、可扩展**

代码很短，但结构干净、符合工程化最佳实践。
---

# ✅ 项目结构（推荐）

```
mini_go_framework/
│── main.go
│── server/
│     ├── router.go
│     └── shutdown.go
│
├── middleware/
│     └── validate.go
│
├── modules/
│     ├── user/
│     │     ├── handler.go
│     │     └── dto.go
│     │
│     └── product/
│           ├── handler.go
│           └── dto.go
```

---

# 🧩 **main.go（包含优雅关闭）**

# 🧩 server/router.go（模块化路由）

# 🧩 server/shutdown.go（优雅关闭）

# 🧩 middleware/validate.go（统一参数验证中间件）

这个中间件会：

* 自动解析 JSON / Query / Form
* 自动验证 struct 标签（binding:"required"）
* 参数错误时返回 400

# 🧩 modules/user/dto.go（参数定义）

# 🧩 modules/user/handler.go

# 🧩 modules/product/dto.go

# 🧩 modules/product/handler.go

# ⭐ 为什么这个结构适合高并发和可扩展？

### 1. **模块拆分清晰（DDD 轻量级风格）**

每个模块只处理自己的业务，不互相耦合。

### 2. **中间件统一负责错误/验证**

业务逻辑更干净，统一风格。

### 3. **原生 goroutine 高并发模型**

Gin + net/http 自带 goroutine 池，适合高负载。

### 4. **优雅关闭避免丢请求**

保障生产环境稳定性。

---

