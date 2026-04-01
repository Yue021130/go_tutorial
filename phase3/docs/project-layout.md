# Standard Go Project Layout 说明

Go 官方并没有像 Maven 那样强制约定目录结构，但社区形成了广泛认可的 **Standard Go Project Layout**。掌握它对于理解和参与真实 Go 项目非常重要。

## 与 Java Maven 的对比

| Java Maven | Go 标准布局 | 说明 |
|------------|------------|------|
| `src/main/java` | `internal/` / `pkg/` | 业务代码 |
| `src/main/resources` | `configs/` / `assets/` | 配置文件、静态资源 |
| `src/test/java` | `*_test.go` 分布在各包中 | 单元测试 |
| `pom.xml` | `go.mod` / `go.sum` | 依赖管理 |
| `Application.java` | `cmd/<app>/main.go` | 程序入口 |

## 核心目录

### `/cmd`

每个子目录对应一个可执行程序。例如：

```text
cmd/
├── server/
│   └── main.go       # API 服务入口
└── worker/
    └── main.go       # 后台 worker 入口（如果有）
```

**原则**：`cmd/` 下尽量少放业务逻辑，只做初始化和依赖注入。

### `/internal`

`internal` 是 Go 的特殊目录，里面的代码只能被本项目（父目录及子目录）导入，外部项目无法导入。

这是 Go 强制实现封装的方式，相当于 Java 的 package-private。

```text
internal/
├── domain/           # 领域模型（Entity/POJO）
├── repository/       # 数据访问层（DAO）
├── service/          # 业务逻辑层
├── handler/          # HTTP 处理层（Controller）
├── middleware/       # 中间件
├── config/           # 配置
├── logger/           # 日志
└── router/           # 路由
```

### `/pkg`

如果有一些代码希望被外部项目复用，可以放在 `pkg/` 下。

本教程项目暂不涉及对外暴露的库，所以没有 `pkg/`。

### `/configs`

配置文件存放目录：

```text
configs/
└── config.yaml
```

### `/docs`

项目文档，包括设计文档、API 文档、部署文档等。

### `/tests`

集成测试、端到端测试。单元测试通常和源码放在同一个包内（`*_test.go`）。

### `/api`

API 协议定义，如 protobuf、OpenAPI/Swagger 文件。

## 分层架构对应关系

```text
HTTP 请求
   ↓
Router（路由分发）
   ↓
Middleware（认证、日志、限流）
   ↓
Handler（参数校验、组装响应）  ≈ Java Controller
   ↓
Service（业务逻辑、事务）      ≈ Java Service
   ↓
Repository（数据库操作）        ≈ Java DAO/Mapper
   ↓
Database
```

## 依赖注入

Go 没有 Spring 的自动依赖注入，通常采用**手动构造函数注入**：

```go
userRepo := repository.NewUserRepository(db)
userService := service.NewUserService(userRepo)
userHandler := handler.NewUserHandler(userService)
```

这种方式简单、明确、易于测试。

## 为什么这样组织？

1. **清晰的责任边界**：每个目录只负责一层
2. **易于测试**：每层可以独立 mock
3. **可维护性**：修改数据访问不影响业务逻辑
4. **工程化**：与大型 Go 项目（Kubernetes、Docker、etcd）结构一致
