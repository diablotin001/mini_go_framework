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
# STEP5 — repo / service / handler（三层）

前置条件：
* 请自行准备可以访问的MySQL和Redis

主要添加：
* GORM(MySQL) + Redis 访问封装
* 模块化：modules/user、modules/product（每个含 handler/service/repo/dto）

---

## 项目树（新增内容）

```
...
├── pkg/
│   ├── db/
│   │   └── mysql.go
│   └── cache/
│       └── redis.go
├── modules/
│   ├── user/
│   │   ├── handler.go
│   │   ├── service.go
│   │   ├── repo.go
│   │   └── dto.go
│   └── product/
│       ├── handler.go
│       ├── service.go
│       ├── repo.go
│       └── dto.go
└── config.yaml
```

---

## config.yaml

```yaml
server:
  addr: ":8080"

db:
  dsn: "user:pass@tcp(127.0.0.1:3306)/yourdb?charset=utf8mb4&parseTime=True&loc=Local"

redis:
  addr: "127.0.0.1:6379"
  password: ""
  db: 0

logs:
  path: "logs/app.log"
```
请设置user:pass为正确的用户名和密码
---

## pkg/db/mysql.go

```go
package db

import (
    "log"

    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitMySQL(dsn string) error {
    var err error
    DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Silent),
    })
    if err != nil {
        return err
    }

    sqlDB, err := DB.DB()
    if err != nil {
        return err
    }
    sqlDB.SetMaxOpenConns(50)
    sqlDB.SetMaxIdleConns(10)
    return nil
}
```

---

## pkg/cache/redis.go

```go
package cache

import (
    "context"
    "time"

    "github.com/go-redis/redis/v8"
)

var RDB *redis.Client
var Ctx = context.Background()

func InitRedis(addr, password string, db int) error {
    RDB = redis.NewClient(&redis.Options{
        Addr:     addr,
        Password: password,
        DB:       db,
    })

    _, err := RDB.Ping(Ctx).Result()
    if err != nil {
        return err
    }
    return nil
}

func GetString(key string) (string, error) {
    return RDB.Get(Ctx, key).Result()
}

func SetString(key string, value interface{}, ttl time.Duration) error {
    return RDB.Set(Ctx, key, value, ttl).Err()
}
```

---

## main.go

```go
package main

import (
    "log"

    "yourapp/logger"
    "yourapp/pkg/cache"
    "yourapp/pkg/db"
    "yourapp/server"
)

func main() {
    // 1. init logger
    logger.InitLogger("logs/app.log")

    // 2. init mysql
    if err := db.InitMySQL("user:pass@tcp(127.0.0.1:3306)/yourdb?charset=utf8mb4&parseTime=True&loc=Local"); err != nil {
        log.Fatal("db init failed:", err)
    }

    // 3. init redis
    if err := cache.InitRedis("127.0.0.1:6379", "", 0); err != nil {
        log.Fatal("redis init failed:", err)
    }

    srv := server.NewHTTPServer()
    go func() {
        if err := srv.ListenAndServe(); err != nil && err.Error() != "http: Server closed" {
            log.Fatal(err)
        }
    }()

    server.WaitForShutdown(srv)
}
```

---

## modules/user/repo.go

```go
package user

import (
    "time"

    "yourapp/pkg/cache"
    "yourapp/pkg/db"
)

type UserModel struct {
    ID        uint `gorm:"primaryKey"`
    Username  string
    Password  string
    Email     string
    CreatedAt time.Time
}

func GetUserByUsername(username string) (*UserModel, error) {
    // try cache first
    key := "user:username:" + username
    if s, err := cache.GetString(key); err == nil && s != "" {
        // 简化：真实场景要 json.Unmarshal
        return &UserModel{Username: s}, nil
    }

    var u UserModel
    if err := db.DB.Where("username = ?", username).First(&u).Error; err != nil {
        return nil, err
    }
    // set cache
    _ = cache.SetString(key, u.Username, 60*60)
    return &u, nil
}

func CreateUser(u *UserModel) error {
    return db.DB.Create(u).Error
}
```

---

## modules/user/service.go

```go
package user

import (
    "errors"
)

func LoginService(username, password string) (*UserModel, error) {
    u, err := GetUserByUsername(username)
    if err != nil {
        return nil, err
    }
    // 简化密码比较
    if u.Password != password {
        return nil, errors.New("invalid credentials")
    }
    return u, nil
}

func RegisterService(username, password, email string) error {
    u := &UserModel{Username: username, Password: password, Email: email}
    return CreateUser(u)
}
```

---

## modules/user/handler.go

```go
package user

import (
    "github.com/gin-gonic/gin"
    "yourapp/response"
)

func Login(c *gin.Context) {
    var req LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.Error(err)
        return
    }
    u, err := LoginService(req.Username, req.Password)
    if err != nil {
        c.Error(err)
        return
    }
    response.Success(c, gin.H{"user": u.Username})
}

func Register(c *gin.Context) {
    var req RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.Error(err)
        return
    }
    if err := RegisterService(req.Username, req.Password, req.Email); err != nil {
        c.Error(err)
        return
    }
    response.Success(c, nil)
}
```

---

## modules/product/repo.go

## modules/product/service.go

## modules/product/handler.go

---

# STEP6.1 — **完整的 JWT 鉴权方案**
* middleware：JWTAuth
* user/login：生成 JWT
* config：JWT 配置
* router：如何给路由使用 JWT
* service：登录逻辑
* repo：查询用户

---

# 1. 需要修改/新增的文件列表

```
config/config.go         ← 增加 JWT 配置读取
config.yaml              ← 增加 JWT 配置项

internal/middleware/jwt.go    ← 新增 JWT 鉴权中间件

internal/repo/user_repo.go    ← 增加 GetByUsername
internal/service/user_service.go  ← 登录逻辑（验证密码 + 生成 token）
internal/handler/user_handler.go  ← 新增 /login API

internal/server/router.go     ← /login 不需要鉴权，其他路由需要
```

---

# 2. 修改内容（按文件分类）

---

# 🔧 **config/config.go（增加 JWT 配置项）**

```go
type JWTConfig struct {
    Secret string
    Expire time.Duration
}

type Config struct {
    ...
    JWT  JWTConfig
}

func Init(path string) {
    ...
    Conf.JWT.Secret = v.GetString("jwt.secret")
    Conf.JWT.Expire = v.GetDuration("jwt.expire")
}
```

---

# 🔧 **config.yaml（新增 JWT 配置）**

```yaml
jwt:
  secret: "my_super_secret_key_123"
  expire: "72h"
```

---

# **middleware/jwt.go（JWT 鉴权）**

```go
package middleware

import (
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v4"
    "yourapp/config"
)

func JWTAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        auth := c.GetHeader("Authorization")
        if auth == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 10002, "msg": "missing token"})
            return
        }

        parts := strings.SplitN(auth, " ", 2)
        if len(parts) != 2 || parts[0] != "Bearer" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 10002, "msg": "invalid token format"})
            return
        }

        tokenStr := parts[1]
        token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
            return []byte(config.Conf.JWT.Secret), nil
        })

        if err != nil || !token.Valid {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 10002, "msg": "invalid or expired token"})
            return
        }

        // 解析 claims
        if claims, ok := token.Claims.(jwt.MapClaims); ok {
            c.Set("user_id", claims["uid"])
        }

        c.Next()
    }
}
```

---

# **internal/repo/user_repo.go（增加用于登录的查询）**

```go
func (r *UserRepo) GetByUsername(username string) (*model.User, error) {
    var user model.User
    if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
        return nil, err
    }
    return &user, nil
}
```

---

# **internal/service/user_service.go（登录逻辑 + JWT 生成）**

```go
package service

import (
    "errors"
    "time"

    "github.com/golang-jwt/jwt/v4"
    "yourapp/config"
    "yourapp/internal/repo"
)

type UserService struct {
    repo *repo.UserRepo
}

func NewUserService(r *repo.UserRepo) *UserService {
    return &UserService{repo: r}
}

func (s *UserService) Login(username, password string) (string, error) {
    user, err := s.repo.GetByUsername(username)
    if err != nil {
        return "", errors.New("user not found")
    }

    // 简化：生产环境用 bcrypt
    if user.Password != password {
        return "", errors.New("incorrect password")
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "uid":  user.ID,
        "exp":  time.Now().Add(config.Conf.JWT.Expire).Unix(),
        "iat":  time.Now().Unix(),
    })

    return token.SignedString([]byte(config.Conf.JWT.Secret))
}
```

---

# **internal/handler/user_handler.go（新增 /login 接口）**

```go
func (h *UserHandler) Login(c *gin.Context) {
    var req struct {
        Username string `json:"username" binding:"required"`
        Password string `json:"password" binding:"required"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"code": 10001, "msg": "invalid params"})
        return
    }

    token, err := h.service.Login(req.Username, req.Password)
    if err != nil {
        c.JSON(401, gin.H{"code": 10002, "msg": err.Error()})
        return
    }

    c.JSON(200, gin.H{"code": 0, "msg": "ok", "data": gin.H{"token": token}})
}
```

---

# **internal/server/router.go（路由分组）**

```go
r := gin.New()

// 公共接口
userGroup := r.Group("/user")
{
    userGroup.POST("/login", userHandler.Login)
}

// 私有接口（需要 JWT）
authGroup := r.Group("/api", middleware.JWTAuth())
{
    authGroup.GET("/profile", userHandler.Profile)
    // ...
}
```

---

# **完整 JWT 工作流**

## 1. 登录

POST `/user/login`
→ 校验用户
→ 生成 JWT 返回客户端

客户端保存 token（放 header）

```
Authorization: Bearer xxxxx
```

测试命令

```bash
curl -X POST http://localhost:8080/user/login -H 'Content-Type: application/json' -d '{"username":"alice","password":"secret"}'
```

## 2. 访问受保护接口

客户端带着 token → middleware/JWTAuth
→ Token 有效 → 放行
→ Token 无效 → 401 返回错误

测试命令
```bash
curl -H "Authorization: Bearer <JWT>" http://localhost:8080/product/list
```
注意：替换<JWT>为user/login返回的token

---
