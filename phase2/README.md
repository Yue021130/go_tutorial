# Phase 2：核心进阶（第 4-7 天）

## 学习目标

通过 4 天学习，达到以下能力：

1. 深入理解 Go 接口系统：隐式实现、接口组合、类型断言、nil interface 陷阱
2. 掌握方法集规则，能正确选择值接收者/指针接收者
3. 掌握 goroutine + channel + select 的并发协作模式
4. 能用 `context` 实现超时、取消、级联控制
5. 熟练使用 `sync` 包构建并发安全逻辑
6. 建立 Go 错误处理工程规范
7. 了解反射与 `unsafe` 的适用边界

## 预计耗时

- 第 4 天：接口与类型系统、方法集
- 第 5 天：goroutine/channel、select 模式
- 第 6 天：context、sync 包
- 第 7 天：错误处理、反射/unsafe、作业
- 每天约 2-3 小时

## 目录结构

```text
phase2/
├── README.md              # 本文件
├── go.mod
├── cmd/                   # 各章节可运行入口
│   ├── 10_interfaces_and_typesystem/
│   ├── 11_method_sets_receivers/
│   ├── 12_goroutines_channels/
│   ├── 13_select_patterns/
│   ├── 14_context_and_cancellation/
│   ├── 15_sync_and_concurrency_safety/
│   ├── 16_error_handling_best_practice/
│   └── 17_reflect_and_unsafe_intro/
├── internal/              # 共享内部包
│   ├── cancellation/      # context 取消、WithValue、Deadline
│   ├── concurrency/       # channel、pipeline、select、worker pool
│   ├── contracts/         # 接口、方法集
│   ├── errdemo/           # 错误处理
│   ├── reflectdemo/       # 反射
│   └── safety/            # sync、atomic、并发安全
└── homework/              # 实战作业：并发任务调度器
```

## 环境检查

```bash
cd phase2

# 运行所有测试
go test ./...

# 使用 race detector 运行测试
go test -race ./...

# 格式化检查
go fmt ./...
```

---

## Day 4：接口与类型系统 + 方法集

### 4.1 接口与类型系统

**核心问题**：Go 的接口与 Java 接口有什么本质区别？

**与 Java 的关键差异**：

| 特性 | Java | Go |
|------|------|-----|
| 实现方式 | 显式 `implements` | 隐式实现（Duck Typing） |
| 接口位置 | 通常由实现方/框架定义 | 推荐由消费者定义 |
| 接口大小 | 倾向于大接口 | 推荐小接口（如 `io.Reader`） |
| 空接口 | 无 | `interface{}` / `any` |
| 多实现 | 一个类实现多个接口 | 一个类型隐式满足多个接口 |

**隐式实现的好处**：
- 解耦：实现方不需要依赖接口定义
- 灵活：任何类型都可以自动满足已有接口

**interface 内部结构（简化）**：

```text
interface 变量
┌─────────┬─────────┐
│  type   │  data   │
│ 指针    │ 指针    │
└─────────┴─────────┘
```

当 `type` 和 `data` 都为 nil 时，interface 才等于 nil。这是 typed nil 陷阱的根源。

**nil interface 陷阱**：

```go
var r Repository        // nil interface
var p *UserRepository   // typed nil
r = p                   // r != nil！
```

**运行命令**：

```bash
go run ./cmd/10_interfaces_and_typesystem
```

### 4.2 方法集与接收者

**核心问题**：一个类型什么时候满足某个接口？

**方法集规则**：

```text
对于类型 T：
  T  的方法集 = 所有值接收者方法
  *T 的方法集 = 所有值接收者方法 + 所有指针接收者方法
```

**关键推论**：
- `*T` 可以赋值给任何需要 `T` 方法集的接口
- `T` 只能赋值给只需要值接收者方法的接口
- 如果接口包含指针接收者方法，`T` 不能满足

**addressable（可寻址）规则**：
- 可寻址：变量、切片元素、指针解引用、可寻址字段
- 不可寻址：字面量、map 值、函数返回值
- 调用指针接收者方法时，Go 会自动对可寻址值取地址

**运行命令**：

```bash
go run ./cmd/11_method_sets_receivers
```

---

## Day 5：goroutine/channel + select

### 5.1 goroutine 与 channel

**核心问题**：Go 的并发模型是什么？

**CSP（Communicating Sequential Processes）**：
- 不要通过共享内存来通信，而要通过通信来共享内存
- goroutine 之间通过 channel 传递数据

**与 Java 的对比**：

| 特性 | Java | Go |
|------|------|-----|
| 并发单元 | Thread / Virtual Thread | goroutine（轻量级） |
| 调度 | OS 线程调度 | Go runtime 调度（GMP） |
| 通信方式 | 共享内存 + 锁 | channel + 锁 |
| 创建成本 | 高（~1MB 栈） | 低（~2KB 初始栈） |
| 异步编程 | CompletableFuture / Reactive | channel + goroutine |

**buffered vs unbuffered channel**：

| 类型 | 创建 | 特点 |
|------|------|------|
| unbuffered | `make(chan int)` | 同步，发送方阻塞直到接收 |
| buffered | `make(chan int, n)` | 异步，缓冲区满才阻塞 |

**channel 关闭语义**：
- 关闭后发送 panic
- 接收返回零值 + false（comma-ok）
- `range` 在 channel 关闭后自动结束

**goroutine 泄漏**：goroutine 永远阻塞在某个 channel 上，无法退出。解决方案：使用 context 或关闭 channel。

**运行命令**：

```bash
go run ./cmd/12_goroutines_channels
```

### 5.2 select 模式

**核心问题**：如何同时处理多个 channel？

**select 关键特性**：
- 多个 case 同时就绪时，随机选择一个
- `default` 分支实现非阻塞操作
- 可用于超时控制、 ticker、优雅关闭
- nil channel 在 select 中会被忽略

**常见模式**：

```go
select {
case v := <-ch:
    // 处理数据
case <-time.After(timeout):
    // 超时
case <-done:
    // 取消
}
```

**运行命令**：

```bash
go run ./cmd/13_select_patterns
```

---

## Day 6：context + sync

### 6.1 context

**核心问题**：如何控制 goroutine 的生命周期？

**context 树结构**：

```text
Background()
   └── parent
         ├── child1 (WithCancel)
         │      └── grandchild (WithTimeout)
         └── child2 (WithValue)

父 context 取消 → 所有子 context 都取消
子 context 取消 → 不影响父 context
```

**context 类型**：
- `context.Background()`：根 context，用于 main、初始化
- `context.TODO()`：临时占位，后续应替换
- `WithCancel`：手动取消
- `WithTimeout`：超时自动取消
- `WithDeadline`：截止时间自动取消
- `WithValue`：传递请求元数据

**WithValue 最佳实践**：
- 只放请求范围元数据（request_id、user_id）
- key 使用私有类型
- 不要把业务配置放进 context
- 不要把 context 存到结构体中

**运行命令**：

```bash
go run ./cmd/14_context_and_cancellation
```

### 6.2 sync 包与并发安全

**核心问题**：如何保护共享状态？

**sync 包工具箱**：

| 工具 | 用途 | 场景 |
|------|------|------|
| sync.Mutex | 互斥锁 | 保护临界区 |
| sync.RWMutex | 读写锁 | 读多写少 |
| sync.WaitGroup | 等待一组 goroutine | 任务汇聚 |
| sync.Once | 只执行一次 | 单例初始化 |
| sync.Pool | 对象缓存 | 减少 GC 压力 |
| sync.Map | 并发安全 map | 读多写少、key 单一 |
| atomic | 原子操作 | 计数器、标志位 |

**并发安全的 map 方案对比**：

1. `map + sync.RWMutex`：通用方案，适合大多数场景
2. `sync.Map`：读多写少、key 类型单一
3. `map + channel`：单 goroutine owner，完全串行化

**race detector**：

```bash
go test -race ./...
go run -race ./cmd/...
```

**运行命令**：

```bash
go run ./cmd/15_sync_and_concurrency_safety
go test -race ./internal/safety/...
```

---

## Day 7：错误处理 + 反射/unsafe

### 7.1 错误处理最佳实践

**核心问题**：Go 没有异常，如何处理错误？

**与 Java 的对比**：

| Java | Go |
|------|-----|
| `try-catch-finally` | `if err != nil` |
| 异常类继承 Throwable | error 接口 + 自定义类型 |
| 堆栈跟踪内置 | 需额外包或手动包装 |
| `throw` / `throws` | `return fmt.Errorf("...: %w", err)` |

**错误处理关键工具**：
- `errors.New`：创建 sentinel error
- `fmt.Errorf("...: %w", err)`：包装错误
- `errors.Is`：判断错误链中是否包含某个错误
- `errors.As`：将错误转换为特定类型
- `errors.Join`（Go 1.20+）：聚合多个错误

**panic vs error 决策树**：

```text
问题是否可预期？
  ├── 否（编程错误/不可恢复）→ panic（通常只在顶层 recover）
  └── 是 → 用 error 返回值
           ├── 调用方可以处理 → 返回 error
           └── 调用方无法处理 → 包装后向上传播
```

**常见反模式**：
- 吞掉错误：`if err != nil { return nil }`
- 过度包装：错误链过长
- 用字符串比较判断错误类型
- 在底层随意打印日志

**运行命令**：

```bash
go run ./cmd/16_error_handling_best_practice
```

### 7.2 反射与 unsafe

**核心问题**：什么时候用反射？什么时候用 unsafe？

**反射三定律**：
1. 从接口值到反射对象：`reflect.TypeOf` / `reflect.ValueOf`
2. 从反射对象到接口值：`v.Interface()`
3. 修改反射对象需要值是可设置的（settable）

**反射适用场景**：
- JSON/XML 序列化
- 配置文件映射到结构体
- 通用测试框架
- 动态代理（有限）

**反射成本**：比直接访问慢 1-2 个数量级，避免在热路径使用。

**unsafe 适用场景**：
- 与 C 代码交互（cgo）
- 极致性能优化（如字节切片转字符串零拷贝）
- 序列化库底层实现

**unsafe 风险**：
- 破坏类型安全
- GC 不友好
- 可移植性差
- 业务代码应避免

**运行命令**：

```bash
go run ./cmd/17_reflect_and_unsafe_intro
```

---

## 实战作业

### 作业要求

实现一个**并发任务调度器**，要求：

1. 定义 `Task` 接口：
   - `ID() string`
   - `Execute(ctx context.Context) error`

2. 实现 `Scheduler`：
   - 支持指定 worker 数量
   - 使用 worker pool 并发执行 Task
   - 支持 `context.Context` 取消
   - 所有任务返回结果，不因某个失败而中断

3. 实现 `RunWithAggregation`：
   - 使用 `errors.Join` 聚合失败任务的错误

4. 编写单元测试，并通过 `go test -race`

### 参考答案

本目录下 `homework/` 提供了完整参考答案：

```bash
cd homework
go run .
go test -v
go test -race
```

### 扩展挑战

1. 增加最大并发数限制（semaphore 模式）
2. 给每个任务增加重试机制
3. 实现任务优先级
4. 给每个任务独立超时控制

---

## 常见错误汇总

### 1. typed nil interface

```go
var p *MyType = nil
var i MyInterface = p
if i == nil { // false！
}
```

### 2. 值接收者 vs 指针接收者

```go
func (t T) Method() {} // T 和 *T 都有此方法
func (t *T) Method() {} // 只有 *T 有
```

### 3. 关闭已关闭的 channel

```go
close(ch)
close(ch) // panic
```

### 4. 未接收 goroutine 的 panic

单个 goroutine panic 会导致整个程序崩溃，顶层服务要用 recover。

### 5. context 存到结构体

```go
type Service struct {
    ctx context.Context // 反模式
}
```

context 应该作为函数第一个参数传递。

### 6. 错误吞掉

```go
if err != nil {
    return nil // 错误信息丢失了
}
```

---

## 最佳实践引用

- **Effective Go**：
  - 接口由使用者定义
  - 小接口优于大接口
  - 错误应该提供足够的上下文
- **Go Code Review Comments**：
  - 不要吞掉错误
  - context 应该作为第一个参数
  - 不要复制 sync 类型
- **Go 并发模式**：
  - 不要通过共享内存通信，要通过通信共享内存
  - 使用 context 控制 goroutine 生命周期

---

## 延伸阅读

- [Effective Go](https://go.dev/doc/effective_go)
- [Go Concurrency Patterns](https://go.dev/blog/pipelines)
- [Go Memory Model](https://go.dev/ref/mem)
- [Context Package](https://pkg.go.dev/context)

---

## 下一步

完成 Phase 2 后，继续学习 **Phase 3：工程实战**，内容包括：

- Go Modules 进阶
- Standard Go Project Layout
- RESTful API 服务（Gin）
- 分层架构
- 数据库操作
- 配置、日志、中间件
- 测试与性能分析
