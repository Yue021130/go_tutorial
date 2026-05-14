# 源码阅读：Prometheus Client Go

## 阅读目标

理解 Prometheus Go Client 如何设计指标收集、注册、导出和 HTTP 暴露。

## 推荐版本

- GitHub: https://github.com/prometheus/client_golang
- 建议阅读 v1.19+ 稳定版

## 核心文件与路线

```text
prometheus/
├── prometheus.go           # 默认注册表、全局快捷函数
├── registry.go             # Registry 注册表核心
├── metric.go               # Metric 接口
├── counter.go              # Counter 指标
├── gauge.go                # Gauge 指标
├── histogram.go            # Histogram 指标
├── summary.go              # Summary 指标
├── vec.go                  # 带标签的指标向量
├── desc.go                 # 指标描述
└── process_collector.go    # 进程级指标收集

promhttp/
├── http.go                 # HTTP handler（/metrics）
└── instrument_server.go    # HTTP 服务端指标装饰
```

## 重点解析

### 1. Metric 接口

所有指标都实现 `Metric` 和 `Collector` 接口：

```go
type Metric interface {
    Desc() *Desc
    Write(*dto.Metric) error
}

type Collector interface {
    Describe(chan<- *Desc)
    Collect(chan<- Metric)
}
```

### 2. Registry 注册表

- `Registry` 是 Collector 的集合，负责收集并导出指标。
- 默认全局注册表 `DefaultRegisterer`，支持 `prometheus.MustRegister()`。
- 支持自定义 Registry，用于隔离不同子系统的指标。

```go
reg := prometheus.NewRegistry()
reg.MustRegister(myCounter)
```

### 3. 四种核心指标类型

| 类型 | 特点 | 适用场景 |
|------|------|----------|
| Counter | 单调递增，只增不减 | 请求总数、错误总数 |
| Gauge | 可增可减 | 当前连接数、内存使用量 |
| Histogram | 采样分桶 | 请求延迟分布 |
| Summary | 滑动时间窗口分位数 | 请求延迟 P99（客户端计算） |

### 4. 标签向量（Vec）

- `CounterVec`、`GaugeVec`、`HistogramVec`、`SummaryVec` 支持动态标签。
- 通过 `WithLabelValues` 或 `With` 获取带标签的指标实例。
- 注意：标签基数爆炸会导致时间序列数量剧增（cardinality issue）。

### 5. HTTP 暴露 /metrics

- `promhttp.Handler()` 使用默认注册表。
- `promhttp.HandlerFor(reg, ...)` 使用自定义注册表。
- 输出格式为 Prometheus 文本格式（Content-Type: text/plain）。

## 阅读建议

1. 从 `prometheus/counter.go` 看 Counter 结构和 Inc/Add 实现。
2. 阅读 `prometheus/registry.go` 的 `Register`、`Gather`。
3. 阅读 `promhttp/http.go` 的 `Handler` 和 `HandlerFor`。
4. 阅读 `prometheus/histogram.go`，理解桶（bucket）和累积直方图。

## 可迁移到自己的项目中的设计

- Metric + Collector 接口设计，便于扩展自定义指标。
- Registry 隔离不同子系统指标。
- Vec 模式：延迟创建带标签指标，避免预创建所有组合。
- HTTP handler 统一暴露指标，与业务解耦。
