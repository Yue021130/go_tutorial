# Phase 2：核心进阶（第 4-7 天）

> 本阶段采用“**章节入口（cmd）+ 共享内部包（internal）+ 基础测试**”的工程化组织方式。

## 学习目标

1. 理解接口系统：隐式实现、接口组合、类型断言与类型切换
2. 理解方法集：值接收者/指针接收者与接口实现关系
3. 掌握 goroutine + channel + select 的并发协作模式
4. 能用 `context` 与 `sync` 构建可取消、可回收、可并发安全的逻辑
5. 建立错误处理工程规范：包装、分类、传播与边界日志
6. 对反射与 `unsafe` 的适用边界有明确认知

## 目录结构

```text
phase2/
├── go.mod
├── README.md
├── cmd/
│   ├── 10_interfaces_and_typesystem/
│   ├── 11_method_sets_receivers/
│   ├── 12_goroutines_channels/
│   ├── 13_select_patterns/
│   ├── 14_context_and_cancellation/
│   ├── 15_sync_and_concurrency_safety/
│   ├── 16_error_handling_best_practice/
│   └── 17_reflect_and_unsafe_intro/
└── internal/
    ├── cancellation/
    ├── concurrency/
    ├── contracts/
    ├── errdemo/
    └── safety/
```

## 快速开始

```bash
cd phase2
go test ./...
```

## 配套文档（建议先看）

- [Docs 文档中心](../docs/)
- [新手起步指南](../docs/getting-started.md)
- [代码阅读指南](../docs/how-to-read-code.md)
- [Phase 2 学习手册](../docs/phase2-playbook.md)

运行单个章节：

```bash
go run ./cmd/10_interfaces_and_typesystem
go run ./cmd/11_method_sets_receivers
go run ./cmd/12_goroutines_channels
go run ./cmd/13_select_patterns
go run ./cmd/14_context_and_cancellation
go run ./cmd/15_sync_and_concurrency_safety
go run ./cmd/16_error_handling_best_practice
go run ./cmd/17_reflect_and_unsafe_intro
```

## 学习路径（Day 4 - Day 7）

| Day | 章节入口 | 主题 | 重点 |
|------|------|------|------|
| Day 4 | `cmd/10_interfaces_and_typesystem` | 接口与类型系统 | 隐式实现、`any`、断言、类型切换、接口组合 |
| Day 4 | `cmd/11_method_sets_receivers` | 方法集与接收者 | 值/指针接收者、接口实现判定 |
| Day 5 | `cmd/12_goroutines_channels` | goroutine 与 channel | pipeline、fan-in、worker pool |
| Day 5 | `cmd/13_select_patterns` | select 模式 | 超时控制、分支选择 |
| Day 6 | `cmd/14_context_and_cancellation` | context 取消控制 | 超时取消、级联取消、泄漏防护 |
| Day 6 | `cmd/15_sync_and_concurrency_safety` | sync 与并发安全 | Mutex/RWMutex/atomic |
| Day 7 | `cmd/16_error_handling_best_practice` | 错误处理最佳实践 | `%w`、`errors.Is/As`、错误分类 |
| Day 7 | `cmd/17_reflect_and_unsafe_intro` | 反射与 unsafe 入门 | 结构体标签、反射成本、unsafe 边界 |

## 与 Java 的重点映射

- Java `implements`（显式） ↔ Go 接口实现（隐式）
- Java 线程池/锁 ↔ Go goroutine/channel/sync/atomic
- Java 异常体系 ↔ Go `error` 显式返回链 + `errors.Is/As`
- Java 反射 ↔ Go `reflect`（能力相似，成本敏感）

## 常见坑提醒

1. 值接收者更新的是副本，不会修改原对象状态。
2. channel 未关闭或关闭时机错误会造成阻塞/ panic。
3. 忽略 `context` 取消信号会导致 goroutine 泄漏。
4. 多 goroutine 共享 map 必须加锁或采用并发安全方案。
5. 错误只打印不包装/返回，会破坏调用链定位能力。
