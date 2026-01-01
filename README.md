# 使用说明

** 本项目是一个基于 Gin 框架的极简 Go 项目模板，用于快速搭建高并发服务。**
** 它的优势在于：**

* **极简**：代码量少，结构清晰，符合工程化最佳实践。
* **高并发**：基于 Gin 的最小化中间件(fast path)，避免了不必要的开销。
* **可扩展**：模块化目录结构，方便添加新功能模块。

** 如何使用**

1. **克隆项目**：`git clone https://github.com/diablotin001/mini_go_framework.git`
2. **安装依赖**：`go mod tidy`
3. **运行项目**：`go run main.go`

** 如何学习 **

* 建议先从 **STEP 1** 开始，理解为什么它能高并发。
* 然后逐步学习 **STEP 2**，掌握项目的模块化目录结构和优雅关闭。

** 每一步对应的git版本 **

* **STEP 1**：`git checkout step1`
* **STEP 2**：`git checkout step2`
* **STEP 3**：`git checkout step3`
* **STEP 4**：`git checkout step4`
* **STEP 5**：`git checkout step5`
* **STEP 6.1**：`git checkout step6.1`
* **STEP 6.2**：`git checkout step6.2`
* **STEP 7**：`git checkout step7`
* **STEP 8**：`git checkout step8`

---

# STEP 1

---

# 为什么它能高并发

## ** 使用 Gin 的最小化中间件(fast path) **

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

## 1. 需要修改/新增的文件列表

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

## 2. 修改内容（按文件分类）

---

### 🔧 **config/config.go（增加 JWT 配置项）**

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

### 🔧 **config.yaml（新增 JWT 配置）**

```yaml
jwt:
  secret: "my_super_secret_key_123"
  expire: "72h"
```

---

### **middleware/jwt.go（JWT 鉴权）**

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

### **internal/repo/user_repo.go（增加用于登录的查询）**

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

### **internal/service/user_service.go（登录逻辑 + JWT 生成）**

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

### **internal/handler/user_handler.go（新增 /login 接口）**

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

### **internal/server/router.go（路由分组）**

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

## **完整 JWT 工作流**

### 1. 登录

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

### 2. 访问受保护接口

客户端带着 token → middleware/JWTAuth
→ Token 有效 → 放行
→ Token 无效 → 401 返回错误

测试命令

```bash
curl -H "Authorization: Bearer <JWT>" http://localhost:8080/product/list
```

注意：替换<JWT>为user/login返回的token

---

# STEP 6.2 access_token,refresh_token 刷新与登出

---

**改动摘要**

- 新增 `/api/*` 统一分组保护，使用 JWT 中间件统一鉴权
- 完成 JWT 刷新与黑名单（登出/吊销）：登录返回 access/refresh，刷新生成新 access，登出将 jti 加入黑名单

**配置**

- `pkg/config/config.go:8` 新增 `JWT.RefreshExpire`，全局持有 `config.Conf`
- `config.yaml:14` 新增：
  - `jwt.secret`（签名秘钥）
  - `jwt.expire`（access token 有效期）
  - `jwt.refresh_expire`（refresh token 有效期）

**JWT 中间件**

- `middleware/jwt.go:1` 增强校验：
  - 解析 `Authorization: Bearer <token>`
  - 读取 `jti`，查询黑名单 `jwt:blacklist:<jti>`（Redis），命中则拒绝
  - 将 `uid` 注入 `Context` 供业务使用

**登录/刷新/登出**

- `modules/user/service.go:1`
  - `LoginService(username, password) (*TokenPair, error)` 生成并返回 `access_token` 和 `refresh_token`
  - `RefreshService(refreshToken string) (string, error)` 校验 refresh token 的 `typ=refresh` 并生成新 access token
  - `LogoutService(accessToken, refreshToken string) error` 解析 `exp` 计算剩余 TTL，将两个 token 的 `jti` 写入黑名单（TTL=剩余有效期）
  - 令牌携带 claims：`uid`、`typ`（access/refresh）、`jti`、`exp`、`iat`
- `modules/user/handler.go:10`
  - `/user/login` 返回 `access_token` 与 `refresh_token`
  - 新增 `/user/logout`（从 `Authorization` 头取 access，可选携带 refresh），写入黑名单
  - 新增 `/user/refresh`（必需 `refresh_token`），返回新的 `access_token`

**路由**

- `server/router.go:20`
  - `/user/login`、`/user/register`、`/user/logout`、`/user/refresh` 作为公共接口（登出依赖头部）
  - 新增 `/api` 分组统一保护
  - 将产品接口迁移至 `/api/product/list`、`/api/product/buy`

**数据库迁移（保留你要的模式）**

- `internal/database/mysql.go:1`
  - `Init(dsn)` 成功后在 `APP_ENV=dev` 时执行 `migrate()`，使用 `pkg/model.User` 与 `pkg/model.Product`
- `pkg/model/user.go:1`、`pkg/model/product.go:1` 存放数据模型，避免循环依赖

**关键代码定位**

- `pkg/config/config.go:8` 配置与全局 `Conf`
- `middleware/jwt.go:1` 黑名单校验与鉴权
- `modules/user/service.go:12` 登录/刷新/登出逻辑
- `modules/user/handler.go:10` 登录返回 token 对；`handler.go:33` 登出；`handler.go:45` 刷新
- `server/router.go:20` `/user/*`；`server/router.go:28` `/api/*`
- `internal/database/mysql.go:1` `Init + migrate` 风格（`APP_ENV=dev`）

**使用与验证**

- 启动（开发自动迁移）：
  - `APP_ENV=dev go run main.go`
- 登录获取令牌：
  - `curl -sS -X POST http://localhost:8080/user/login -H 'Content-Type: application/json' -d '{"username":"alice","password":"secret"}'`
  - 响应包含 `access_token` 与 `refresh_token`
- 使用 access 访问受保护接口：
  - `curl -sS http://localhost:8080/api/product/list -H "Authorization: Bearer <access_token>"`
- 刷新生成新 access：
  - `curl -sS -X POST http://localhost:8080/user/refresh -H 'Content-Type: application/json' -d '{"refresh_token":"<refresh_token>"}'`
- 登出（吊销当前 access，可选吊销 refresh）：
  - `curl -sS -X POST http://localhost:8080/user/logout -H "Authorization: Bearer <access_token>" -H 'Content-Type: application/json' -d '{"refresh_token":"<refresh_token>"}'`
  - 之后旧 access 再访问将返回 `{"code":10002,"msg":"token revoked"}`

**注意事项**

- Redis 未启动时，黑名单读写会被安全忽略（开发容错）；生产环境需开启 Redis
- refresh token 目前不做轮换，仅校验 `typ=refresh`；如需严格一次性刷新，可在刷新后将旧 refresh 的 `jti` 加入黑名单并返回新的 refresh
- 如需对 `/user/logout` 强制鉴权，可为该路由添加 `JWTAuth()` 中间件

**TODO**

- 目前refresh_token,access_token使用的是JWT格式，考虑是否需要切换为UUID（短token）格式
- 考虑blacklist基于uid，这样可以实现用户级别的登出

---

# STEP 7 单元测试

**最小可运行的单元测试示例**（涵盖 repo/service/handler 三层），全部是“基础示例”，可以复制后按功能继续扩展。

---

## 目录结构示例

使用：

* **github.com/DATA-DOG/go-sqlmock** → Mock MySQL（repo 测试用）
* **net/http/httptest + gin** → handler 测试
* **repo 使用 mock** → service 测试

全部是业界常用做法。

**说明**

- `pkg/db/mysql_mock.go` 提供创建 GORM+sqlmock 的辅助函数
- `modules/user/repo_test.go` 通过 sqlmock 验证查询逻辑
- `modules/user/service_test.go` 通过 mock repo 验证登录生成 token
- `modules/user/handler_test.go` 通过 httptest 验证 `/user/login` 响应结构

**关键实现**

- `pkg/db/mysql_mock.go:1` 新增 `NewMock(t)` 返回 `*gorm.DB` 与 `sqlmock`
- `modules/user/repo_test.go:1`
  - 设置 `internal/database.DB = NewMock(t)`
  - 期望 SQL：`SELECT * FROM \`users\` WHERE username = ? ORDER BY \`users\`.\`id\` LIMIT ?`，参数 `("amy", 1)`
  - 断言返回用户名正确
- `modules/user/service.go:20`
  - 新增 `IUserRepo` 接口与包级 `userRepo`（默认调用现有 `GetUserByUsername`）
  - `LoginService` 使用 `userRepo`，便于在测试中注入 mock
  - 登录返回 `TokenPair{access_token, refresh_token}`
  - 提供 `RefreshService` 和 `LogoutService`（刷新、吊销）
- `modules/user/service_test.go:1`
  - 注入 `mockRepo` 至 `userRepo`
  - 设置 `config.Conf.JWT.Secret`、`Expire`、`RefreshExpire`
  - 校验 `LoginService` 成功与失败场景
- `modules/user/handler_test.go:1`
  - 设置 `gin.TestMode`
  - 注入 `mockSvcRepo` 至 `userRepo`
  - 配置 `config.Conf.JWT`
  - 注册 `POST /login`，校验 200 与返回 `access_token` 非空

---

## repo 测试：Mock GORM + Mock MySQL

❗ 不需要真实数据库。

---

### internal/repo/user_repo_test.go

```go
package repo_test

import (
    "regexp"
    "testing"

    "github.com/DATA-DOG/go-sqlmock"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"

    "yourapp/internal/model"
    "yourapp/internal/repo"
)

func newMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
    mockDB, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("mock db err: %v", err)
    }
    dial := mysql.New(mysql.Config{
        Conn:                      mockDB,
        SkipInitializeWithVersion: true,
    })
    db, err := gorm.Open(dial, &gorm.Config{})
    if err != nil {
        t.Fatalf("gorm open err: %v", err)
    }
    return db, mock
}

func TestUserRepo_GetByUsername(t *testing.T) {
    db, mock := newMockDB(t)

    r := repo.NewUserRepo(db)

    rows := sqlmock.NewRows([]string{"id", "username", "password"}).
        AddRow(1, "amy", "123456")

    mock.ExpectQuery(regexp.QuoteMeta(
        "SELECT * FROM `users` WHERE username = ? ORDER BY `users`.`id` LIMIT 1",
    )).
        WithArgs("amy").
        WillReturnRows(rows)

    user, err := r.GetByUsername("amy")
    if err != nil {
        t.Fatalf("unexpected err: %v", err)
    }

    if user.Username != "amy" {
        t.Fatalf("wrong username: %v", user.Username)
    }
}
```

---

## service 测试：Mock Repo，不需要数据库

使用最简单 mock（不引入 gomock）：

---

### internal/service/user_service_test.go

```go
package service_test

import (
    "errors"
    "testing"

    "yourapp/internal/model"
    "yourapp/internal/service"
)

// ---- mock repo ----
type mockUserRepo struct {
    mockUser *model.User
    mockErr  error
}

func (m *mockUserRepo) GetByUsername(username string) (*model.User, error) {
    return m.mockUser, m.mockErr
}

func TestUserService_Login_Success(t *testing.T) {
    r := &mockUserRepo{
        mockUser: &model.User{
            ID:       1,
            Username: "amy",
            Password: "123456",
        },
    }

    s := service.NewUserService(r)

    token, err := s.Login("amy", "123456")
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    if token == "" {
        t.Fatalf("empty token")
    }
}

func TestUserService_Login_Fail(t *testing.T) {
    r := &mockUserRepo{
        mockErr: errors.New("not found"),
    }
    s := service.NewUserService(r)

    _, err := s.Login("not_exist", "xx")
    if err == nil {
        t.Fatal("expected error, got none")
    }
}
```

---

## handler 测试：使用 httptest + gin

---

### internal/handler/user_handler_test.go

```go
package handler_test

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"

    "yourapp/internal/handler"
    "yourapp/internal/service"
    "yourapp/internal/model"
)

// mock service
type mockUserService struct {
    token string
    err   error
}

func (m *mockUserService) Login(username, password string) (string, error) {
    return m.token, m.err
}

func TestUserHandler_Login(t *testing.T) {
    gin.SetMode(gin.TestMode)

    svc := &mockUserService{token: "mock_token"}

    h := handler.NewUserHandler(svc)

    router := gin.Default()
    router.POST("/login", h.Login)

    payload := map[string]string{
        "username": "amy",
        "password": "123456",
    }
    body, _ := json.Marshal(payload)

    req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
    req.Header.Set("Content-Type", "application/json")

    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    if w.Code != 200 {
        t.Fatalf("expected 200, got %d", w.Code)
    }

    var resp struct {
        Code int `json:"code"`
        Data struct {
            Token string `json:"token"`
        }
    }
    json.Unmarshal(w.Body.Bytes(), &resp)

    if resp.Data.Token != "mock_token" {
        t.Fatalf("token mismatch: %v", resp.Data.Token)
    }
}
```

## 使用说明

- 运行全部测试：`go test ./... -v`

---

# STEP8 多环境日志（dev、prod）、 MySQL + Redis + zap 打印 SQL 执行时间

---

内容包括：
**1 多环境日志 dev/prod
2 GORM + Redis 的 SQL 执行时间拦截 + Zap 打印
3 支持自动根据环境启用不同的输出格式（console/json）**

---

## 1. 配置文件新增字段（config.yaml）

```yaml
app:
  env: dev   # or prod

log:
  level: info
  format: console  # dev = console，prod = json
  file: "./logs/app.log"

database:
  dsn: "user:pass@tcp(127.0.0.1:3306)/demo?charset=utf8mb4&parseTime=True&loc=Local"
  log_level: 4   # 1=Silent 2=Error 3=Warn 4=Info
  slow_threshold_ms: 200

redis:
  addr: "127.0.0.1:6379"
  db: 0
```

---

## 2. 新增日志初始化（pkg/logger/logger.go）

Zap Logger 增强版：支持 console/json 格式 + dev/prod 自动切换

```go
package logger

import (
    "os"

    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

var Log *zap.Logger

type LogConfig struct {
    Level  string
    Format string // console / json
    File   string
    Env    string // dev / prod
}

func Init(cfg LogConfig) error {
    var level zapcore.Level
    if err := level.Set(cfg.Level); err != nil {
        level = zap.InfoLevel
    }

    // encoder
    var encoder zapcore.Encoder
    if cfg.Format == "json" || cfg.Env == "prod" {
        encoder = zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
    } else {
        encoder = zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
    }

    // log file writer
    fileWriter := zapcore.AddSync(&lumberjack.Logger{
        Filename:   cfg.File,
        MaxSize:    50,
        MaxBackups: 5,
        MaxAge:     28,
    })

    core := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(fileWriter, zapcore.AddSync(os.Stdout)), level)

    Log = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
    return nil
}
```

> dev 默认 console，prod 默认 json
> 文件输出 + 控制台混合输出

---

## 3. GORM SQL 执行时间拦截（pkg/db/gorm_logger.go）

让 MySQL 的 SQL、耗时、慢查询全部写入 Zap

```go
package db

import (
    "context"
    "time"

    "github.com/your_project/pkg/logger"
    "go.uber.org/zap"
    "gorm.io/gorm/logger"
)

type ZapGormLogger struct {
    SlowThreshold time.Duration
    LogLevel      logger.LogLevel
}

func NewZapGormLogger(slowMS int, level logger.LogLevel) *ZapGormLogger {
    return &ZapGormLogger{
        SlowThreshold: time.Duration(slowMS) * time.Millisecond,
        LogLevel:      level,
    }
}

func (l *ZapGormLogger) LogMode(level logger.LogLevel) logger.Interface {
    l.LogLevel = level
    return l
}

func (l *ZapGormLogger) Info(ctx context.Context, msg string, args ...interface{}) {
    logger.Log.Sugar().Infof(msg, args...)
}

func (l *ZapGormLogger) Warn(ctx context.Context, msg string, args ...interface{}) {
    logger.Log.Sugar().Warnf(msg, args...)
}

func (l *ZapGormLogger) Error(ctx context.Context, msg string, args ...interface{}) {
    logger.Log.Sugar().Errorf(msg, args...)
}

func (l *ZapGormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
    elapsed := time.Since(begin)
    sql, rows := fc()

    fields := []zap.Field{
        zap.Duration("elapsed", elapsed),
        zap.String("sql", sql),
        zap.Int64("rows", rows),
    }

    switch {
    case err != nil:
        logger.Log.Error("SQL Error", append(fields, zap.Error(err))...)
    case elapsed > l.SlowThreshold:
        logger.Log.Warn("Slow SQL", fields...)
    default:
        logger.Log.Info("SQL", fields...)
    }
}
```

---

## 4. GORM 初始化（pkg/db/mysql.go）

```go
package db

import (
    "log"

    "github.com/spf13/viper"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    gormLogger "gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitMySQL() {
    cfg := viper.GetStringMapString("database")

    dsn := cfg["dsn"]

    logLevel := gormLogger.LogLevel(viper.GetInt("database.log_level"))
    slowMS := viper.GetInt("database.slow_threshold_ms")

    gLogger := NewZapGormLogger(slowMS, logLevel)

    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
        Logger: gLogger,
    })
    if err != nil {
        log.Fatalf("MySQL connection failed: %v", err)
    }

    DB = db
}
```

---

## 5. Redis 打印执行时间（pkg/redis/redis.go）

```go
package redis

import (
    "time"

    "github.com/redis/go-redis/v9"
    "github.com/spf13/viper"
    "github.com/your_project/pkg/logger"
)

var Rdb *redis.Client

func InitRedis() {
    addr := viper.GetString("redis.addr")
    db := viper.GetInt("redis.db")

    Rdb = redis.NewClient(&redis.Options{
        Addr: addr,
        DB:   db,
    })

    // Wrap Do for time logging
    origDo := Rdb.Process
    Rdb.AddHook(redis.Hook{
        BeforeProcess: func(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
            return context.WithValue(ctx, "start", time.Now()), nil
        },
        AfterProcess: func(ctx context.Context, cmd redis.Cmder) error {
            start := ctx.Value("start").(time.Time)
            cost := time.Since(start)
            logger.Log.Info("Redis",
                zap.String("cmd", cmd.FullName()),
                zap.String("args", fmt.Sprint(cmd.Args())),
                zap.Duration("cost", cost),
            )
            return nil
        },
    })
}
```

---

## 6. 最终效果

### Dev 环境输出（控制台）

```
2025-01-01 SQL elapsed=18.2ms sql="SELECT ..." rows=1
2025-01-01 Redis cmd=GET args=[user:1] cost=1.3ms
```

### Prod 环境输出（json）

```json
{
  "level": "info",
  "sql": "SELECT ...",
  "elapsed": "20ms",
  "rows": 1,
  "timestamp": "2025-01-01T12:00:00"
}
```

---