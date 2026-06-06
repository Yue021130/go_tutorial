# Phase 4：高级与源码（第 15-21 天）

## 学习目标

通过 7 天学习，达到以下能力：

1. 理解 Go 内存模型与 GC 原理，能与 JVM 内存模型对比说明差异。
2. 理解 GMP 调度器模型，能说清 goroutine、M、P 的协作关系。
3. 理解内存逃逸分析、栈扩容、sync.Pool 的适用场景。
4. 能使用 Go 标准库实现 TCP/UDP/HTTP2 服务，能跑通 gRPC 示例。
5. 能在项目中灵活运用 Option 模式、Functional Options、Pipeline、Worker Pool。
6. 能阅读中等规模开源项目（gin/etcd/prometheus）的核心源码，提炼可复用的设计。

## 预计耗时

- **Day 15**：Go 内存模型与 GC 原理
- **Day 16**：GMP 调度器模型
- **Day 17**：内存逃逸分析、栈扩容、sync.Pool
- **Day 18-19**：网络编程（TCP/UDP/HTTP2/gRPC）
- **Day 20**：Go 中的常用设计模式
- **Day 21**：开源项目源码阅读
- 每天约 2-3 小时

## 目录结构

```text
phase4/
├── README.md                  # 本文件
├── go.mod
├── api/
│   └── helloworld.proto       # gRPC 示例 proto
├── cmd/
│   ├── 20_memory_model_gc/    # 内存模型与 GC
│   ├── 21_gmp_scheduler/      # GMP 调度器
│   ├── 22_escape_analysis/    # 逃逸分析、栈扩容、sync.Pool
│   ├── 23_tcp_server/         # TCP Echo Server
│   ├── 24_udp_server/         # UDP Echo Server
│   ├── 25_http2_server/       # HTTP/2 (h2c) Server
│   ├── 26_grpc_server/        # gRPC Server
│   ├── 27_grpc_client/        # gRPC Client
│   ├── 28_option_pattern/     # Option 模式
│   ├── 29_functional_options/ # Functional Options
│   ├── 30_pipeline/           # Pipeline 模式
│   └── 31_worker_pool/        # Worker Pool 模式
├── internal/
│   ├── gc/                    # GC 工具函数
│   ├── gmp/                   # GMP 工具函数
│   └── escape/                # 逃逸分析示例
├── docs/
│   ├── reading-gin.md         # Gin 源码阅读指南
│   ├── reading-etcd.md        # etcd 源码阅读指南
│   └── reading-prometheus.md  # Prometheus Client 源码阅读指南
└── homework/
    ├── main.go                  # 作业入口
    └── logprocessor/            # 并发日志处理器
        ├── log_processor.go
        └── log_processor_test.go
```

## 快速开始

```bash
cd phase4

# 安装依赖（已执行过可跳过）
go mod tidy

# 运行内存与 GC 示例
go run ./cmd/20_memory_model_gc

# 运行 GMP 调度器示例
go run ./cmd/21_gmp_scheduler

# 运行逃逸分析示例（加 -gcflags=-m 查看逃逸信息）
go run ./cmd/22_escape_analysis
go run -gcflags=-m ./cmd/22_escape_analysis

# 运行网络示例
go run ./cmd/23_tcp_server
go run ./cmd/24_udp_server
go run ./cmd/25_http2_server

# gRPC 需要同时启动 server 和 client
go run ./cmd/26_grpc_server &
go run ./cmd/27_grpc_client

# 运行设计模式示例
go run ./cmd/28_option_pattern
go run ./cmd/29_functional_options
go run ./cmd/30_pipeline
go run ./cmd/31_worker_pool

# 运行测试
go test ./...

# 运行作业
echo -e "INFO hello\nERROR fail\nWARN slow" | go run ./homework
```

---

## Day 15：Go 内存模型与 GC 原理

### 15.1 Go 内存模型

Go 程序运行时的内存从逻辑上可分为：

- **栈（Stack）**：每个 goroutine 独占，存储局部变量、函数参数、返回地址等。
- **堆（Heap）**：全局共享，存储生命周期超出函数作用域的对象，由 GC 管理。
- **代码段、数据段**：与 C/Java 类似。

与 JVM 内存模型对比：

| 区域 | JVM | Go |
|------|-----|-----|
| 线程栈 | JVM 栈（固定大小 -Xss） | goroutine 栈（动态扩容，初始 2KB） |
| 堆 | 新生代/老年代/元空间 | 单一堆，由 GC 管理 |
| 方法区/元空间 | 类信息、常量池、静态变量 | 无独立方法区，代码段+全局变量 |
| 字符串常量池 | 有 | 无（字符串不可变，但无专门池） |
| 直接内存 | Direct Buffer / MappedByteBuffer | 可通过 syscall 或 unsafe 分配 |

### 15.2 Go GC 原理

Go 当前使用**并发三色标记清除（Concurrent Tri-Color Mark-Sweep）**算法：

```text
三色标记过程：

初始：所有对象都是白色（候选垃圾）
      
1. STW 极短时间，扫描根对象（栈、全局变量），标记为灰色
   
2. 并发标记：灰色对象引用的对象标记为灰色，自身标记为黑色
   
3. 重复步骤 2，直到没有灰色对象
   
4. STW 极短时间，处理写屏障（write barrier）期间的引用变化
   
5. 清扫：白色对象回收
```

写屏障（Write Barrier）：在并发标记期间，如果黑色对象新引用了白色对象，会把白色对象标灰，避免漏标。

GC 触发条件（默认 GOGC=100）：

> 当堆内存相比上次 GC 后的存活对象增长 100% 时触发下一次 GC。

可通过 `GOGC` 环境变量或 `debug.SetGCPercent` 调整。

### 15.3 GC 调优核心指标

- **STW 时间**：Go 1.8+ 后通常小于 1ms。
- **CPU 占用**：GC 过程会占用 25% 的 CPU（GOGC 越小越频繁）。
- **堆内存**：观察 `runtime.MemStats.HeapAlloc`。

### 15.4 最佳实践

- 减少堆分配：小对象尽量留在栈上，复用对象（sync.Pool）。
- 避免内存泄漏：goroutine 未退出、channel 未关闭、全局 map 无限增长。
- 不要频繁调用 `runtime.GC()`，干扰 GC 自适应策略。

> 参考：[Go GC Guide](https://go.dev/doc/gc-guide)、Effective Go 中关于内存和 GC 的章节。

---

## Day 16：GMP 调度器模型

### 16.1 G、M、P 是什么

```text
        +-----+     +-----+     +-----+
        |  G  |     |  G  |     |  G  |   <- 全局可运行队列
        +-----+     +-----+     +-----+
           \          |          /
            v         v         v
        +-----------------------------+
        |              P              |   <- 逻辑处理器（GOMAXPROCS 个）
        |  +-----------------------+  |
        |  |   Local Run Queue     |  |   <- 每个 P 的本地队列
        |  | [G1] [G2] [G3] ...    |  |
        |  +-----------------------+  |
        +-----------------------------+
                     |
                     v
        +-----------------------------+
        |              M              |   <- OS 线程
        |  (running on CPU core)      |
        +-----------------------------+
```

- **G（Goroutine）**：用户态轻量线程，初始栈 2KB，由 runtime 调度。
- **M（Machine）**：操作系统线程，M 必须绑定 P 才能执行 G。
- **P（Processor）**：逻辑处理器，持有本地可运行 G 队列，数量默认等于 CPU 核心数。

### 16.2 调度流程

1. 程序启动时创建 `runtime.GOMAXPROCS(0)` 个 P。
2. 创建 goroutine 时，优先放入当前 P 的本地队列。
3. M 绑定 P 后，从 P 的本地队列取 G 执行。
4. 若本地队列为空，M 尝试从全局队列或其他 P **偷取（work stealing）** G。
5. 当 G 阻塞（系统调用、channel 等待）时，M 与 P 解绑，新的 M 接管 P 继续执行。

### 16.3 与 Java 线程调度对比

| 特性 | Java 线程 | Go goroutine |
|------|-----------|--------------|
| 模型 | 1:1 OS 线程 | M:N 用户态协程 |
| 创建成本 | 高（MB 级栈） | 低（KB 级栈） |
| 切换成本 | 需陷入内核 | 用户态切换 |
| 调度器 | 操作系统 | Go runtime |
| 并发数量 | 通常数百上千 | 可轻松数十万 |
| 控制并发 | ThreadPoolExecutor | channel + sync + context |

### 16.4 调度陷阱

- CPU 密集型任务数量不要超过 GOMAXPROCS，否则无收益。
- 大量 goroutine 阻塞在 IO 时，M 数量会增加，注意监控 `runtime.NumGoroutine()`。
- 不要滥用 `runtime.GOMAXPROCS(1)`，除非有特殊测试需求。

> 参考：[Go Scheduler](https://www.ardanlabs.com/blog/2018/08/scheduling-in-go-part1.html)（William Kennedy 系列文章）。

---

## Day 17：内存逃逸分析、栈扩容、sync.Pool

### 17.1 逃逸分析

编译器决定变量分配在栈上还是堆上的过程称为**逃逸分析**。

常见逃逸场景：

- 返回局部变量指针。
- 闭包捕获局部变量。
- 切片/数组/map 大小不确定或过大。
- 接口类型参数或返回值。
- 调用 `fmt.Println` 等接收 `interface{}` 的函数。

查看逃逸分析：

```bash
go build -gcflags=-m ./cmd/22_escape_analysis
```

### 17.2 栈扩容

goroutine 栈初始仅 2KB，按需要**连续栈（contiguous stack）**扩容：

```text
   初始栈 2KB              需要更多空间              扩容后栈
   +----------+            +----------+             +----------+
   |          |            |          |             |          |
   |  小帧    |    -->     |  溢出！  |     -->     |  更大帧  |
   |          |            |          |             |          |
   +----------+            +----------+             +----------+
```

- 栈扩容时，旧栈内容复制到新栈，调整指针。
- 栈最大可达 1GB（64 位）。

### 17.3 sync.Pool

`sync.Pool` 是临时对象池，用于复用对象、减少 GC 压力。

特点：

- 线程安全。
- 池中对象可能被 GC 回收，**不能当作持久缓存**。
- 适合临时缓冲区、解析器、大对象等。

与 Java 对象池对比：

| 特性 | Java Commons Pool | Go sync.Pool |
|------|-------------------|--------------|
| 持久性 | 可长期持有对象 | 对象可能被 GC 回收 |
| 配置 | 复杂（最大/最小数量、过期策略） | 简单，无配置 |
| 适用 | 数据库连接池等 | 临时对象复用 |

> 参考：Effective Go 中关于并发和 sync 包的说明。

---

## Day 18-19：网络编程

### 18.1 TCP/UDP

Go 的 `net` 包屏蔽了底层差异：

- 在 Linux 使用 epoll，macOS 使用 kqueue，Windows 使用 IOCP。
- 对开发者呈现阻塞式 API，但底层由 runtime netpoller 管理。
- 每个连接一个 goroutine，代码简洁，无需像 Java NIO 那样写 Reactor 模式。

### 18.2 HTTP/2

- Go 标准库 `net/http` 在 TLS 下自动协商 HTTP/2。
- 非 TLS 场景可用 `golang.org/x/net/http2/h2c` 启用 h2c。

### 18.3 gRPC

- 基于 HTTP/2 + Protocol Buffers。
- 强类型、高性能、支持流式调用。
- 服务定义使用 `.proto` 文件，通过 `protoc` 生成客户端/服务端代码。

与 REST 对比：

| 特性 | REST | gRPC |
|------|------|------|
| 协议 | HTTP/1.1 或 HTTP/2 | HTTP/2 |
| 序列化 | JSON / XML | Protobuf |
| 性能 | 一般 | 高 |
| 调试 | 方便（curl） | 需 grpcurl / Postman / BloomRPC |
| 适用 | 开放 API、浏览器 | 微服务内部通信 |

与 Java Dubbo 对比：

| 特性 | Dubbo | gRPC |
|------|-------|------|
| 传输 | TCP | HTTP/2 |
| 序列化 | Hessian / Protobuf | Protobuf |
| 服务发现 | ZooKeeper / Nacos / etcd | 需额外组件 |
| 生态 | 国内主流 | 云原生主流 |

---

## Day 20：常用设计模式

### 20.1 Option 模式 / Functional Options

Go 中替代 Java Builder 的常用方式：

```go
type Option func(*Config)

func WithHost(host string) Option {
    return func(c *Config) { c.host = host }
}

func NewConfig(opts ...Option) *Config {
    c := &Config{host: "127.0.0.1", port: 8080}
    for _, opt := range opts { opt(c) }
    return c
}
```

优点：

- API 简洁，扩展时不需要改构造函数签名。
- 被广泛应用于 grpc.NewServer、zap.New 等知名库。

> 参考：Dave Cheney《Functional options for friendly APIs》。

### 20.2 Pipeline

将数据处理拆分为多个阶段，通过 channel 连接：

```text
[Source] --chan--> [Stage 1] --chan--> [Stage 2] --chan--> [Sink]
```

适合流式、可并行、可组合的数据处理。

### 20.3 Worker Pool

固定数量 worker 从任务队列取任务执行：

```text
        +--------+     +--------+     +--------+
   任务 |  Job 1 | --> |  Job 2 | --> |  Job 3 |
        +--------+     +--------+     +--------+
             |               |               |
             v               v               v
        +-----------------------------------------+
        |              Job Channel                |
        +-----------------------------------------+
             |               |               |
             v               v               v
        +---------+    +---------+    +---------+
        | Worker 1|    | Worker 2|    | Worker 3|
        +---------+    +---------+    +---------+
```

用于控制并发度，避免资源耗尽。

---

## Day 21：开源项目源码阅读

阅读三个项目源码，详见 `docs/` 目录：

- [Gin 源码阅读指南](docs/reading-gin.md)
- [etcd 源码阅读指南](docs/reading-etcd.md)
- [Prometheus Client Go 源码阅读指南](docs/reading-prometheus.md)

阅读方法：

1. 先明确目标：我要学路由？中间件？还是 Raft？
2. 从入口开始，跟一次完整请求链路。
3. 画出核心结构图和数据流图。
4. 把可复用的设计点记录到自己的项目中。

---

## 实战作业

### 作业：并发日志处理器

**需求**：

1. 从标准输入读取日志，每行格式为 `LEVEL message`。
2. 使用 Worker Pool 并发解析日志，提取 `LEVEL` 和 `message`。
3. 使用 channel 流水线统计 `INFO` / `WARN` / `ERROR` 数量。
4. 输出统计结果。

**参考答案**：见 `homework/logprocessor/`。

**运行**：

```bash
echo -e "INFO request received\nERROR database timeout\nWARN slow query\nINFO response sent" | go run ./homework
```

**进阶挑战**：

1. 增加按小时/分钟聚合统计。
2. 把解析器和统计器拆成独立 Pipeline Stage。
3. 使用 `sync.Pool` 复用日志解析缓冲区，对比基准测试提升。

---

## 常见坑（需重点规避）

- 认为 `sync.Pool` 是缓存，存放需要长期保留的对象。
- 大量 goroutine 泄漏：只启动不管理退出条件。
- 在 CPU 密集型任务中创建远超 GOMAXPROCS 的 goroutine，导致无效切换。
- 写 gRPC 服务时不设置超时/取消，导致调用 hanging。
- 阅读源码时一开始就陷入细节，没有先建立整体数据流。

---

## 参考资源

- [Effective Go](https://go.dev/doc/effective_go)
- [Go Memory Model](https://go.dev/ref/mem)
- [Go GC Guide](https://go.dev/doc/gc-guide)
- [Scheduling In Go](https://www.ardanlabs.com/blog/2018/08/scheduling-in-go-part1.html)（William Kennedy）
- [Functional options for friendly APIs](https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis)（Dave Cheney）
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
