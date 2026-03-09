# 代码阅读指南（仓库到底怎么读）

## 1. 先理解目录，再读代码

这个仓库目前最核心的目录是：

- `phase1/`：单文件入门示例（语法和标准库基础）
- `phase2/`：工程化进阶示例（多包结构 + 测试）
- `docs/`：文档导航和学习手册

## 2. Phase 2 应该怎么读

`phase2/` 里有两类目录：

- `cmd/<chapter>/main.go`：章节入口，负责“演示怎么用”
- `internal/<pkg>`：共享实现，负责“具体怎么做”

建议每章都按这个固定动作：

1. 先 `go run ./cmd/<chapter>` 看行为
2. 再打开对应 `internal/<pkg>` 看实现
3. 最后看 `*_test.go` 理解边界和预期

## 3. “从输出反推代码”的阅读法（新手友好）

很多同学一上来就看实现，容易迷路。建议反过来：

1. 先运行程序，记录输出
2. 找输出对应的 `fmt.Println`
3. 顺着调用链跳到 `internal` 包
4. 看测试如何断言这个行为

这样更容易建立“功能 -> 代码 -> 测试”的完整认知。

## 4. 章节对照表（入口 -> 实现）

| 章节入口 | 主要实现包 | 应该重点看什么 |
|------|------|------|
| `cmd/10_interfaces_and_typesystem` | `internal/contracts` | 接口隐式实现、类型断言/类型切换 |
| `cmd/11_method_sets_receivers` | `internal/contracts` | 值/指针接收者对状态修改的影响 |
| `cmd/12_goroutines_channels` | `internal/concurrency` | pipeline、fan-in、worker pool |
| `cmd/13_select_patterns` | `internal/concurrency` | select 超时分支与错误处理 |
| `cmd/14_context_and_cancellation` | `internal/cancellation` | 取消传播与 goroutine 收敛 |
| `cmd/15_sync_and_concurrency_safety` | `internal/safety` | mutex/atomic 的线程安全实践 |
| `cmd/16_error_handling_best_practice` | `internal/errdemo` | 错误包装与分类判断 |
| `cmd/17_reflect_and_unsafe_intro` | 章节内代码 | reflect 与 unsafe 的边界认知 |

## 5. 初学者常见误区（读代码时）

1. 只看 `main.go` 不看 `internal`：会误以为逻辑很简单
2. 只看实现不看测试：会看不清“预期行为”
3. 并发代码只看 happy path：忽略取消和超时分支
4. 错误只看字符串：不看 `errors.Is/As` 的语义分类
