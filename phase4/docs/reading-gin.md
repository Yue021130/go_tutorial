# 源码阅读：Gin

## 阅读目标

理解 Gin 如何实现一个高性能、中间件链式的 HTTP Web 框架。

## 推荐版本

- GitHub: https://github.com/gin-gonic/gin
- 建议阅读 v1.9.x 稳定版

## 核心文件与路线

```text
gin/
├── gin.go              # Engine 定义、核心入口
├── routergroup.go      # RouterGroup、路由组、中间件注册
├── tree.go             # 压缩前缀树路由匹配（radix tree）
├── context.go          # Context 上下文，请求/响应/参数/绑定
├── binding/            # 参数绑定（JSON、XML、Form、Query、Header）
├── render/             # 响应渲染（JSON、XML、HTML、String）
├── middleware/         # 内置中间件（Logger、Recovery）
└── response_writer.go  # 包装 http.ResponseWriter
```

## 重点解析

### 1. Engine 与 RouterGroup

- `Engine` 是 `RouterGroup` 的子类（组合关系）。
- `RouterGroup` 持有 `HandlersChain`（中间件链）和 `basePath`。
- `Engine` 实现了 `http.Handler` 接口，可直接传给 `http.ListenAndServe`。

```go
// 关键结构
type Engine struct {
    RouterGroup
    pool             sync.Pool       // Context 对象池
    trees            methodTrees     // 路由树
    ...
}
```

### 2. 中间件链（HandlersChain）

- 每个路由对应一个 `[]HandlerFunc`。
- 通过 `c.Next()` 在链中推进，`c.Abort()` 中断链。
- 中间件执行顺序：`use1 -> use2 -> handler -> use2 后置 -> use1 后置`。

```go
func (c *Context) Next() {
    c.index++
    for c.index < int8(len(c.handlers)) {
        c.handlers[c.index](c)
        c.index++
    }
}
```

### 3. 路由匹配：压缩前缀树（Radix Tree）

- Gin 不是用 map 做路由查找，而是用 radix tree。
- 支持 `:param` 和 `*wildcard`。
- 每个 HTTP 方法对应一棵树。

### 4. Context 对象池

- 每个请求复用 `Context` 对象，减少 GC 压力。
- `engine.pool.Get()` 获取，`c.reset()` 重置，`c.writermem.reset()` 重置 writer。

### 5. 参数绑定与验证

- `binding.Default` 根据 Content-Type 选择绑定器。
- 支持 tag：`json:`、`uri:`、`form:`、`query:`、`header:`。
- 验证使用 `go-playground/validator`。

## 阅读建议

1. 从 `gin.go` 的 `New()` 开始，看 Engine 初始化。
2. 跟一次完整请求：`ServeHTTP` -> `pool.Get()` -> `handleHTTPRequest` -> `c.Next()`。
3. 阅读 `tree.go` 的 `addRoute` 和 `getValue`，理解路由树。
4. 阅读 `logger.go` 和 `recovery.go`，学习如何写 Gin 中间件。

## 可迁移到自己的项目中的设计

- 中间件链模型（index + Next/Abort）。
- Context 对象池复用。
- Radix Tree 路由匹配思想（如需自研路由）。
- 绑定器 + 验证器的组合。
