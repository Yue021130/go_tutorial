# Phase 3：工程实战（第 8-14 天）

## 学习目标

通过 7 天学习，达到以下能力：

1. 掌握 Go Modules 依赖管理和版本控制
2. 理解 Standard Go Project Layout 项目结构
3. 能使用 Gin 框架搭建完整 RESTful API 服务
4. 掌握分层架构：handler → service → repository
5. 能使用 GORM 进行数据库操作
6. 能使用 Viper 管理配置、zap 记录日志
7. 能实现 JWT 认证、请求日志、错误恢复、限流等中间件
8. 会写单元测试、集成测试、基准测试，能看测试覆盖率
9. 能使用 pprof 进行性能分析

## 预计耗时

- 第 8 天：Go Modules + 项目布局 + 项目初始化
- 第 9 天：domain + repository + service 分层
- 第 10 天：Gin handler + router + 中间件
- 第 11 天：JWT 认证 + 限流 + 错误恢复
- 第 12 天：Viper 配置 + zap 日志
- 第 13 天：单元测试 + 集成测试 + 基准测试
- 第 14 天：pprof + 作业
- 每天约 2-3 小时

## 目录结构

```text
phase3/
├── README.md                  # 本文件
├── go.mod
├── go.sum
├── cmd/
│   └── server/
│       └── main.go            # 程序入口
├── internal/
│   ├── config/                # Viper 配置
│   ├── domain/                # 领域模型
│   ├── handler/               # HTTP handler
│   ├── logger/                # zap 日志
│   ├── middleware/            # 中间件
│   ├── repository/            # 数据访问层
│   ├── router/                # 路由
│   └── service/               # 业务逻辑层
├── configs/
│   └── config.yaml            # 配置文件
├── tests/
│   └── user_api_test.go       # 集成测试
├── homework/
│   ├── README.md              # 订单管理作业
│   ├── order.go
│   ├── main.go
│   └── order_test.go
└── docs/
    └── project-layout.md      # 项目布局说明
```

## 快速开始

```bash
cd phase3

# 安装依赖（已执行过可跳过）
go mod tidy

# 运行测试
go test ./...

# 启动服务
go run ./cmd/server

# 在另一个终端测试 API
curl http://localhost:8080/health
```

---

## Day 8：Go Modules + 项目布局

### 8.1 Go Modules 依赖管理

**核心问题**：Go 如何管理依赖？与 Maven/Gradle 有什么区别？

**与 Java 对比**：

| 特性 | Maven/Gradle | Go Modules |
|------|-------------|-----------|
| 描述文件 | `pom.xml` / `build.gradle` | `go.mod` |
| 依赖锁定 | 无（依赖本地仓库） | `go.sum` |
| 语义化版本 | 支持 | 支持 `v1.2.3` |
| 依赖仓库 | Maven Central / JCenter | 远程模块 + 本地缓存 |
| 私有仓库 | `settings.xml` | `GOPRIVATE` 环境变量 |
| 命令 | `mvn install` | `go mod tidy` |

**常用命令**：

```bash
go mod init <module-name>  # 初始化模块
go get <package>           # 添加依赖
go mod tidy                # 整理依赖
go mod download            # 下载依赖
go mod vendor              # 生成 vendor 目录
go list -m all             # 列出所有依赖
```

**私有仓库配置**：

```bash
# 配置 GOPRIVATE，不走 public sumdb
export GOPRIVATE="github.com/yourcompany/*"
export GONOSUMDB="github.com/yourcompany/*"
```

**语义化版本**：
- `v1.2.3`：主版本.次版本.补丁版本
- `v1.2.3-alpha.1`：预发布版本
- `v0.x.x`：不稳定版本

### 8.2 Standard Go Project Layout

详见 `docs/project-layout.md`。

核心目录：
- `/cmd`：程序入口
- `/internal`：私有代码
- `/configs`：配置文件
- `/docs`：文档
- `/tests`：集成测试

**分层架构**：

```text
HTTP 请求 → Router → Middleware → Handler → Service → Repository → Database
```

与 Java 对应：
- Handler ≈ Controller
- Service ≈ Service
- Repository ≈ DAO/Mapper

---

## Day 9：domain + repository + service

### 9.1 domain 领域模型

```go
type User struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    Username  string    `json:"username" gorm:"uniqueIndex"`
    Password  string    `json:"-" gorm:"not null"` // json:"-" 不序列化
    Email     string    `json:"email"`
    Role      string    `json:"role"`
    CreatedAt time.Time `json:"created_at"`
}
```

### 9.2 repository 数据访问层

使用接口定义 + GORM 实现：

```go
type UserRepository interface {
    Create(ctx context.Context, user *domain.User) error
    GetByID(ctx context.Context, id uint) (*domain.User, error)
    // ...
}
```

**与 Java 对比**：
- Java：MyBatis Mapper / JPA Repository
- Go：自定义接口 + GORM / sqlx / database/sql

### 9.3 service 业务逻辑层

负责：
- 业务规则校验
- 密码加密
- 调用 repository
- 事务协调

```go
func (s *userService) Register(ctx context.Context, req domain.CreateUserRequest) (*domain.UserResponse, error) {
    // 检查用户名是否存在
    // 密码加密
    // 创建用户
}
```

---

## Day 10：Gin handler + router + 中间件

### 10.1 Gin 框架

Gin 是 Go 最流行的 Web 框架之一，特点：
- 高性能（基于 httprouter）
- 中间件机制
- 参数绑定和校验
- 丰富的生态

### 10.2 Handler 处理 HTTP 请求

```go
func (h *UserHandler) Register(c *gin.Context) {
    var req domain.CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    // ...
}
```

### 10.3 Router 路由注册

```go
api := r.Group("/api/v1")
{
    api.POST("/auth/register", userHandler.Register)
    api.POST("/auth/login", userHandler.Login)
    
    authorized := api.Group("")
    authorized.Use(middleware.JWTAuth())
    {
        authorized.GET("/users", userHandler.List)
        // ...
    }
}
```

---

## Day 11：JWT 认证 + 限流 + 错误恢复

### 11.1 JWT 认证中间件

流程：
1. 登录成功后生成 JWT token
2. 受保护路由校验 `Authorization: Bearer <token>`
3. 从 token 中解析用户信息，存入 gin.Context

### 11.2 限流中间件

使用 `golang.org/x/time/rate` 实现令牌桶限流：

```go
limiter := rate.NewLimiter(10, 20) // 每秒 10 个，突发 20 个
```

### 11.3 错误恢复中间件

捕获 handler 中的 panic，防止单个请求拖垮整个服务：

```go
defer func() {
    if err := recover(); err != nil {
        // 记录日志，返回 500
    }
}()
```

---

## Day 12：Viper 配置 + zap 日志

### 12.1 Viper 配置管理

Viper 支持：
- YAML/JSON/TOML 配置文件
- 环境变量
- 命令行参数
- 默认值

```go
viper.SetConfigFile("configs/config.yaml")
viper.AutomaticEnv()
viper.ReadInConfig()
viper.Unmarshal(&cfg)
```

### 12.2 zap 日志

zap 是 Uber 出品的高性能日志库：

```go
logger, _ := zap.NewProduction()
logger.Info("msg", zap.String("key", "value"))
```

---

## Day 13：测试

### 13.1 单元测试

使用 `testing` 包 + `testify`：

```go
func TestUserService_Register(t *testing.T) {
    repo := new(mockUserRepository)
    svc := NewUserService(repo)
    // ...
}
```

### 13.2 Mock

使用 `testify/mock` 创建 repository 的 mock：

```go
type mockUserRepository struct {
    mock.Mock
}
```

### 13.3 集成测试

使用 `httptest.Server` 启动完整服务，发送真实 HTTP 请求：

```go
server := httptest.NewServer(r)
resp, err := server.Client().Post(server.URL+"/api/v1/auth/login", ...)
```

### 13.4 基准测试

```go
func BenchmarkUserHandler_Register(b *testing.B) {
    for i := 0; i < b.N; i++ {
        // 执行被测代码
    }
}
```

### 13.5 覆盖率

```bash
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

---

## Day 14：pprof 性能分析

### 14.1 pprof 路由

```go
import _ "net/http/pprof"
```

本示例手动注册了 `/debug/pprof/*` 路由。

### 14.2 常用命令

```bash
# CPU 分析
curl http://localhost:8080/debug/pprof/profile?seconds=30 > cpu.prof
go tool pprof cpu.prof

# 内存分析
curl http://localhost:8080/debug/pprof/heap > heap.prof
go tool pprof heap.prof

# goroutine
curl http://localhost:8080/debug/pprof/goroutine > goroutine.prof
```

---

## API 接口列表

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | `/health` | 健康检查 | 否 |
| POST | `/api/v1/auth/register` | 用户注册 | 否 |
| POST | `/api/v1/auth/login` | 用户登录 | 否 |
| GET | `/api/v1/users` | 用户列表 | JWT |
| GET | `/api/v1/users/:id` | 单个用户 | JWT |
| POST | `/api/v1/users` | 创建用户 | JWT |
| PUT | `/api/v1/users/:id` | 更新用户 | JWT |
| DELETE | `/api/v1/users/:id` | 删除用户 | JWT |
| GET | `/debug/pprof/*` | 性能分析 | 否（生产应限制） |

## 测试命令

```bash
# 全部测试
go test ./...

# 覆盖率
go test -cover ./...

# race detector
go test -race ./...

# 基准测试
go test -bench=. ./internal/handler/
```

## 常见坑汇总

### 1. GORM 模型字段不迁移

确保 `AutoMigrate` 传入的是指针：`db.AutoMigrate(&domain.User{})`。

### 2. JWT secret 泄露

生产环境 secret 必须来自环境变量，不能硬编码。

### 3. gin.Context 在 goroutine 中使用

`gin.Context` 不能在请求处理完成后的 goroutine 中使用，需要复制。

### 4. 数据库连接未关闭

使用 `sqlDB, _ := db.DB()` 获取底层连接，程序退出时 `Close()`。

### 5. 测试依赖真实数据库

使用 SQLite 内存数据库或 testcontainers，避免污染开发数据库。

---

## 实战作业

详见 `homework/README.md`。

实现一个**订单管理模块**：
- `Order` 领域模型
- `OrderRepository`、`OrderService`、`OrderHandler`
- RESTful API：创建、查询、列表、更新状态
- 状态转换校验
- 单元测试覆盖

---

## 最佳实践引用

- **Effective Go**：保持包小、接口由使用者定义
- **Go Code Review Comments**：错误处理、context 使用、并发安全
- **Standard Go Project Layout**：https://github.com/golang-standards/project-layout

---

## 延伸阅读

- [Gin 官方文档](https://gin-gonic.com/docs/)
- [GORM 官方文档](https://gorm.io/docs/)
- [Viper 官方文档](https://github.com/spf13/viper)
- [zap 官方文档](https://pkg.go.dev/go.uber.org/zap)

---

## 下一步

完成 Phase 3 后，继续学习 **Phase 4：高级与源码**，内容包括：

- Go 内存模型与 GC 原理
- GMP 调度器
- 内存逃逸分析
- 网络编程与 gRPC
- 设计模式
- 开源项目源码阅读
