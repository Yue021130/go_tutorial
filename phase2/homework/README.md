# Phase 2 实战作业：并发任务调度器

## 作业要求

实现一个并发任务调度器，要求：

1. **Task 接口**
   - `ID() string`
   - `Execute(ctx context.Context) error`

2. **Scheduler 调度器**
   - 支持指定 worker 数量（并发度）
   - 使用 worker pool 并发执行多个 Task
   - 支持 `context.Context` 取消
   - 所有任务都要返回结果，不因某个任务失败而中断

3. **错误聚合**
   - 提供 `RunWithAggregation` 方法
   - 使用 `errors.Join` 聚合所有失败任务的错误

4. **单元测试**
   - 覆盖正常执行、context 取消、错误聚合、空任务等场景
   - 使用 `go test -race` 验证并发安全

## 运行方式

```bash
cd phase2/homework

# 运行命令行演示
go run .

# 运行单元测试
go test -v

# 使用 race detector 验证并发安全
go test -race
```

## 参考答案

本目录下的 `scheduler.go`、`main.go`、`scheduler_test.go` 提供了一份完整参考答案。建议先自己实现，再对照查看。

## 扩展挑战

1. 给 Scheduler 增加"最大执行任务数"限制（semaphore 模式）
2. 给每个任务增加重试机制（失败自动重试 N 次）
3. 给任务增加优先级，高优先级任务优先执行
4. 实现任务执行超时控制（每个任务独立超时）
