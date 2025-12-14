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
# STEP3: 添加全局错误处理器（error handler）、统一返回结构（success/error struct）

---

最终效果：

* 业务 handler 不需要 `c.JSON` 来重复写结构
* 只需返回 `c.Error(err)` 或 `response.Success(c, data)` 即可
* 所有错误格式统一
* 错误码可扩展（业务错误码、系统错误码）

---

# ✅ **一、统一返回结构 response/**

目录

```
response/
│── response.go
│── error.go
│── codes.go
```

---

# 🧩 response/response.go（统一 Success 返回）

---

# 🧩 response/codes.go（自定义错误码：可无限扩展）

---

# 🧩 response/error.go（错误统一格式）

---

# ✅ **二、全局错误处理器 middleware/error_handler.go**

---

# ⭐ ** server/router.go 加这个中间件**

```go
...
r.Use(middleware.ErrorHandler())   // 全局错误处理
...
```

---

# 📌 三、修改 Validator 中间件

---

# 📌 四、业务 Handler 改为使用统一响应

## user/login：

```go
func Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBind(&req); err != nil {
		c.Error(err) // 自动交给全局错误处理
		return
	}

	response.Success(c, gin.H{
		"user": req.Username,
	})
}
```

---

# 🎉 最终效果展示

## ✔ 成功统一格式：

```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "user": "alice"
  }
}
```

## ❌ 系统错误统一格式：

```json
{
  "code": 10000,
  "msg": "something wrong"
}
```

---

# STEP4: **生产级 Zap 日志系统的完整接入方案**

* 全局统一 Zap Logger（JSON 输出、按天轮转、支持 Info/Error）
* Gin 日志改为 Zap（替代默认 Logger）
* 请求日志：包含 method/path/status/latency/ip
* Panic 日志：完整 stack trace
* 所有业务模块可直接使用：`zap.L().Info(...)`

---

# 一、目录结构（新增 logger/）

```
yourapp/
├── logger/
│     └── logger.go
├── middleware/
│     ├── zap_request.go
│     └── zap_recovery.go
└── server/router.go
```

---

# 二、生产级 Zap 日志初始化（JSON + 文件分割）

位置：`logger/logger.go`

使用 **lumberjack** 实现日志切割
Zap 使用 JSON 格式（适合 ELK / CloudWatch / Loki）

---

# 三、在 main.go 初始化 Zap 日志

```go
func main() {
	logger.InitLogger() // 添加这一行

	...
}
```

---

# 四、接入 Gin 请求日志（用 Zap 替换 Gin 的默认日志）

新增文件：

## `middleware/zap_request.go`

---

# 五、Zap Panic 恢复（带 stacktrace）

新增文件：

## `middleware/zap_recovery.go`

---

# 六、在 router.go 中启用 Zap 中间件

```go
r := gin.New()

r.Use(middleware.ZapLogger())    // 请求日志
r.Use(middleware.ZapRecovery())  // panic 日志

...
```

---

# 七、业务模块可以直接使用 Zap

任何地方都可以：

```go
zap.L().Info("user logged in", zap.String("user", req.Username))
zap.L().Error("db failed", zap.Error(err))
```

---

# 八、日志输出示例（JSON，可直接进入 ELK）

```json
{
  "level": "info",
  "time": "2025-12-11T16:34:22+08:00",
  "caller": "user/handler.go:21",
  "msg": "user logged in",
  "user": "alice"
}
```

---