# 文档中心（Docs）

如果你是第一次接触这个仓库，建议按下面顺序阅读：

1. [新手起步指南](getting-started.md)：先看什么、先跑什么、每天怎么学
2. [代码阅读指南](how-to-read-code.md)：仓库该怎么看，`phase2/cmd` 和 `phase2/internal` 怎么配合看
3. [Phase 2 学习手册](phase2-playbook.md)：针对核心进阶阶段的章节级学习和练习方式
4. [Phase 3 工程实战](../phase3/README.md)：进入 `phase3/` 目录，先读 `README.md`，再跑 `go test ./...` 和 `go run ./cmd/server`
5. [Phase 4 高级与源码](../phase4/README.md)：进入 `phase4/` 目录，按 `cmd/20` 到 `cmd/31` 顺序运行，结合 `docs/` 源码阅读指南

---

如果你赶时间，只看这一段：

- 第一步：看根目录 `README.md` 的“学习路径总览”
- 第二步：执行 `phase1` 的 `go run`，把 00-09 全部跑一遍
- 第三步：进入 `phase2` 先 `go test ./...`，再按 `cmd/10` 到 `cmd/17` 顺序运行
- 第四步：进入 `phase3` 跑 `go test ./...`，再启动 `go run ./cmd/server`，用 curl 测试接口
- 第五步：进入 `phase4` 按 `cmd/20` 到 `cmd/31` 顺序运行示例，阅读 `docs/` 源码指南
- 第六步：每跑完一章，回到对应 `internal` 包看实现细节
