# Go 语言从入门到精通教程

> 本教程面向有 Java/Spring Boot 后端开发经验的工程师，通过 30 天系统性学习快速掌握 Go 语言并投入实际项目开发。
>
> 特点：每个知识点均附带**可直接复制运行的代码示例**，关键概念与 Java 做对比说明，工程导向、务实落地。

## 适用人群

- 有 Java/Spring Boot 后端开发经验
- 熟悉 MySQL、Redis、Docker、Linux
- 希望快速掌握 Go 语言并投入项目开发

## 学习路径总览

| 阶段 | 天数 | 主题 | 目录 |
|------|------|------|------|
| Phase 1 | 第 1-3 天 | 快速上手 | [`phase1/`](phase1/) |
| Phase 2 | 第 4-7 天 | 核心进阶 | [`phase2/`](phase2/)（项目与代码已补全） |
| Phase 3 | 第 8-14 天 | 工程实战 | `phase3/`（待更新） |
| Phase 4 | 第 15-21 天 | 高级与源码 | `phase4/`（待更新） |
| Phase 5 | 第 22-30 天 | 云原生与部署 | `phase5/`（待更新） |

## 环境要求

- **Go 1.18+**（推荐 Go 1.20 或更高版本）
- 本教程所有代码示例使用 Go 1.26 验证通过
- 操作系统：Windows / macOS / Linux 均可

## 仓库结构

```text
go-tutorial/
├── README.md              # 本文件：总览与学习指南
├── docs/                  # 新手文档中心（先看这里）
│   ├── README.md
│   ├── getting-started.md
│   ├── how-to-read-code.md
│   └── phase2-playbook.md
├── phase1/                # 快速上手（已重构，深度版）
│   ├── README.md          # Phase 1 完整教程文档
│   ├── 00_hello.go        # 包、导入、可见性
│   ├── 01_variables_types_constants.go  # 变量、常量、类型、零值
│   ├── 02_control_flow.go # if/for/switch/range
│   ├── 03_functions.go    # 函数、多返回值、闭包
│   ├── 04_defer_panic_recover.go        # defer、panic、recover
│   ├── 05_arrays_and_slices.go          # 数组与切片深度
│   ├── 06_maps.go         # map 深度
│   ├── 07_structs_methods.go            # 结构体、方法、tag、JSON
│   ├── 08_pointers.go     # 指针深度
│   ├── 09_http_server.go  # 标准库 net/http 服务
│   └── homework/          # 实战作业完整参考答案
│       ├── README.md
│       ├── main.go
│       ├── student.go
│       └── student_test.go
├── phase2/                # 核心进阶（多包工程 + 章节入口 + 测试）
│   ├── README.md          # Phase 2 学习指南
│   ├── go.mod
│   ├── cmd/               # 10-17 章节可运行入口
│   └── internal/          # 共享内部包与测试
├── phase3/                # RESTful API、Gin、数据库、测试、pprof
├── phase4/                # 内存模型、GMP 调度、网络编程、设计模式
└── phase5/                # Docker、K8s、可观测性、CI/CD
```

## 快速开始

### 1. 克隆或进入仓库

```bash
cd go-tutorial
```

### 2. 验证 Go 环境

```bash
go version
```

### 3. 按阶段学习

每个 `phase/` 目录都是独立的可运行示例，进入对应目录即可运行：

```bash
cd phase1
go run 00_hello.go
go run 01_variables_types_constants.go
go run 02_control_flow.go
go run 03_functions.go
go run 04_defer_panic_recover.go
go run 05_arrays_and_slices.go
go run 06_maps.go
go run 07_structs_methods.go
go run 08_pointers.go
go run 09_http_server.go
```

### 4. 新手先看文档（强烈建议）

- [Docs 文档中心](docs/)
- [新手起步指南](docs/getting-started.md)
- [代码阅读指南](docs/how-to-read-code.md)
- [Phase 2 学习手册](docs/phase2-playbook.md)

## 与 Java/Spring Boot 的核心差异速查

| 对比项 | Java/Spring Boot | Go |
|--------|------------------|-----|
| 依赖管理 | Maven / Gradle | Go Modules（`go.mod` / `go.sum`） |
| 入口函数 | `public static void main(String[] args)` | `func main()` |
| 访问控制 | `public` / `private` 关键字 | 首字母大小写 |
| 变量声明 | `int a = 10;` | `a := 10` 或 `var a int = 10` |
| 错误处理 | `try-catch-finally` | `if err != nil` + `panic/recover` |
| 多返回值 | 需封装对象或数组 | 原生支持 `(result, err)` |
| 集合 | `ArrayList`、`HashMap` | `slice`、`map` |
| POJO | `class` + getter/setter | `struct` |
| 指针 | 隐式引用 | 显式指针 `*T` 和 `&` |
| 接口 | 显式实现（`implements`） | 隐式实现（Duck Typing） |
| 并发 | `Thread`、`ExecutorService` | `goroutine` + `channel` |
| Web 路由 | `@GetMapping` 等注解 | `http.HandleFunc` / Gin 路由 |
| 中间件 | Filter / Interceptor | 函数包装（高阶函数） |

## 各阶段内容概览

### Phase 1：快速上手（第 1-3 天）

- Go 环境搭建：GOPATH vs Go Modules
- 基础语法：变量、常量、数据类型、控制流
- 函数、多返回值、defer、panic/recover
- 数组、切片、map、结构体
- 指针基础
- 标准库 `net/http` 实现 HTTP Hello World Server

👉 [进入 Phase 1 学习](phase1/)

### Phase 2：核心进阶（第 4-7 天）

#### 学习目标

- 理解 Go 接口系统与 Java 接口体系的设计差异，掌握隐式实现的工程价值。
- 掌握 goroutine、channel、select 的并发协作模型，能写出可终止、可控的并发代码。
- 学会使用 `context`、`sync`、`atomic` 构建可靠并发组件，避免竞态和 goroutine 泄漏。
- 建立面向工程的错误处理规范（错误包装、分类、传播、日志边界）。
- 对反射与 `unsafe` 建立边界认知：知道何时用、何时不用。

#### 建议学习节奏（4 天）

- **Day 4**：接口与类型系统、方法集与接收者
- **Day 5**：goroutine、channel、select、并发模式（pipeline/worker pool）
- **Day 6**：`context` 取消传播、`sync` 原语、`atomic` 与并发安全 map 方案
- **Day 7**：错误处理最佳实践、反射基础与 `unsafe` 风险认知

#### 核心知识点清单

- 接口：隐式实现、空接口 `any`、类型断言、类型切换、接口组合
- 接收者：值接收者/指针接收者与方法集关系，避免“看似实现却未实现”问题
- 并发：channel 缓冲与阻塞语义、关闭约定、select 超时/退出分支
- 取消：`context.WithCancel/WithTimeout` 的生命周期管理
- 同步：`sync.Mutex/RWMutex/WaitGroup/Once/Cond/Pool` 的适用边界
- 错误：`errors.Is/As/Join`、`fmt.Errorf("...: %w", err)`、哨兵错误与错误类型
- 反射：动态调用与 tag 读取的成本认知；`unsafe` 仅作机制理解

#### 常见坑（需重点规避）

- 在循环中错误捕获迭代变量导致并发逻辑错乱
- channel 发送/接收方向不匹配，或提前关闭导致 panic
- 只启动 goroutine 不管理退出条件，造成泄漏
- 在共享 map 上并发读写触发数据竞争
- 错误只打印不返回，导致调用链失真

👉 [进入 Phase 2 学习](phase2/)

### Phase 3：工程实战（第 8-14 天）

- Go Modules 依赖管理进阶
- Standard Go Project Layout 项目目录结构
- 完整 RESTful API 服务（推荐 Gin）
- 分层架构：handler → service → repository
- 数据库操作：GORM / sqlx
- 配置管理：Viper
- 日志：zap / logrus
- 中间件：JWT 认证、请求日志、错误恢复、限流
- 单元测试、基准测试与覆盖率
- 性能分析：pprof

### Phase 4：高级与源码（第 15-21 天）

- Go 内存模型与 GC 原理
- GMP 调度器模型
- 内存逃逸分析、栈扩容、sync.Pool 原理
- 网络编程：TCP/UDP、HTTP/2、gRPC
- 常用设计模式：Option 模式、Pipeline、Worker Pool
- 开源项目源码阅读（gin / etcd / prometheus 部分源码）

### Phase 5：云原生与部署（第 22-30 天）

- Docker 多阶段构建 Go 应用
- 优雅关闭与信号处理
- Kubernetes 部署：Deployment、Service、ConfigMap、健康探针
- 可观测性：Prometheus 指标、Jaeger 链路追踪（OpenTelemetry）
- CI/CD：GitHub Actions / GitLab CI 流水线

## 学习建议

1. **每段代码都要亲手跑一遍**：只看不动手等于没学。
2. **对比 Java 理解**：把 Go 的特性映射到你熟悉的 Java 概念上。
3. **完成实战作业**：作业是检验理解的关键。
4. **善用官方资源**：
   - [Effective Go](https://go.dev/doc/effective_go)
   - [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
   - [Go 官方文档](https://go.dev/doc/)
   - [Go by Example](https://gobyexample.com/)

## 代码规范

- 代码注释使用中文
- 每个示例文件均为 `package main`，可直接 `go run`
- 优先使用标准库，减少外部依赖
- 遵循 `gofmt` 格式化规范

## 贡献与反馈

如果你在学习过程中发现代码有误或希望补充内容，欢迎提出。

---

**当前进度**：Phase 1 已完成 ✅，Phase 2 项目与代码已补全 ✅  
**下一步**：[开始学习 Phase 2](phase2/)
