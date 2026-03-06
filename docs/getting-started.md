# 新手起步指南（先看什么、怎么学）

## 1. 你应该先看什么

先看这个顺序，不要跳：

1. 根目录 `README.md`：知道全局学习路径（Phase 1 到 Phase 5）
2. `phase1/README.md`：先学基础语法和标准库
3. `phase2/README.md`：再进入并发、接口、错误处理等进阶内容

> 原因：Go 的进阶概念（接口方法集、并发取消、错误包装）都建立在基础语法之上。

## 2. 第一次运行建议（Windows）

在仓库根目录打开终端后，按这个顺序执行：

```bash
go version
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

然后进入 Phase 2：

```bash
cd ..\phase2
go test ./...
go run ./cmd/10_interfaces_and_typesystem
go run ./cmd/11_method_sets_receivers
go run ./cmd/12_goroutines_channels
go run ./cmd/13_select_patterns
go run ./cmd/14_context_and_cancellation
go run ./cmd/15_sync_and_concurrency_safety
go run ./cmd/16_error_handling_best_practice
go run ./cmd/17_reflect_and_unsafe_intro
```

## 3. 每天怎么学（建议模板）

每天 2-3 小时可以这样分配：

- 20 分钟：读当日章节 README 内容，明确概念和目标
- 40 分钟：直接运行示例，观察输出
- 40 分钟：改代码（改参数、改流程、故意造错误）验证理解
- 20 分钟：做笔记（今天学了什么、踩了什么坑）

## 4. 学习时一定要做的 5 件事

1. 每段代码必须亲手运行，不要只看不跑
2. 每个示例至少改一次再跑一次
3. 把 Go 概念映射到 Java（接口、错误、并发）理解
4. 并发相关代码要学会看“退出条件”
5. 不懂就回到最小示例，不要一上来堆复杂逻辑

## 5. 完成标准（你可以自测）

如果你能做到下面这些，说明真的入门了：

- 能解释“为什么 Go 接口不需要 implements”
- 能写出带 `context.WithTimeout` 的并发任务
- 能用 `errors.Is` / `errors.As` 分类错误
- 能区分值接收者和指针接收者对行为的影响
- 能独立运行并读懂 `phase2/cmd` 和 `phase2/internal` 的关系
