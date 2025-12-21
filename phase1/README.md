# Phase 1：Go 快速上手（第 1-3 天）

## 学习目标

- 完成 Go 开发环境配置，理解 GOPATH 与 Go Modules 的区别
- 掌握 Go 基础语法，能识别与 Java 的关键差异
- 理解函数、多返回值、defer、panic/recover 机制
- 掌握数组、切片、map、结构体的用法
- 理解指针基础，破除 Java 开发者对指针的恐惧
- 能用标准库 `net/http` 写出可运行的 HTTP 服务

## 预计耗时

- 理论 + 编码：约 6-8 小时
- 建议每天 2-3 小时，分 3 天完成

## 环境要求

- Go 1.18+（推荐使用 Go 1.20 或更高版本）
- 本示例使用 Go 1.26 验证通过

## 文件说明

| 文件 | 内容 |
|------|------|
| `01_basics.go` | 变量、常量、数据类型、控制流 |
| `02_functions.go` | 函数、多返回值、defer、panic/recover |
| `03_collections.go` | 数组、切片、map、结构体 |
| `04_pointers.go` | 指针基础、值接收者 vs 指针接收者 |
| `05_http_server.go` | 标准库 net/http 实现 RESTful 风格服务 |

## 快速开始

```bash
# 进入目录
cd go-tutorial/phase1

# 运行各个示例
go run 01_basics.go
go run 02_functions.go
go run 03_collections.go
go run 04_pointers.go

# 启动 HTTP 服务（会阻塞终端）
go run 05_http_server.go
```

## 各文件运行说明

### 01_basics.go

演示 Go 基础语法，包含：

- `var` / `:=` 变量声明
- `const` 常量与 `iota` 枚举
- 基本数据类型与零值
- `if`、`for`、`switch` 控制流
- 字符串的 `byte` 与 `rune` 遍历

运行示例：

```bash
go run 01_basics.go
```

### 02_functions.go

演示函数相关特性，包含：

- 函数声明与调用
- 多返回值与错误处理
- 命名返回值
- 可变参数
- `defer` 延迟执行
- `panic` / `recover` 异常机制
- 函数作为一等公民

运行示例：

```bash
go run 02_functions.go
```

### 03_collections.go

演示复合数据类型，包含：

- 数组（固定长度）
- 切片（动态数组，最常用）
- `map` 键值对
- 结构体（struct）与结构体方法
- 结构体嵌套与嵌入字段

运行示例：

```bash
go run 03_collections.go
```

### 04_pointers.go

专为 Java 开发者讲解指针，包含：

- 取地址 `&` 与解引用 `*`
- 值传递 vs 指针传递
- 结构体的值接收者 vs 指针接收者
- `new` 关键字
- 返回局部变量地址（逃逸分析）

运行示例：

```bash
go run 04_pointers.go
```

### 05_http_server.go

使用标准库 `net/http` 实现 HTTP 服务，包含：

- 路由注册
- 查询参数解析
- JSON 请求体解析
- JSON 响应
- 简单的中间件（请求日志）

启动服务：

```bash
go run 05_http_server.go
```

服务启动后，在另一个终端中测试：

```bash
# Hello World
curl http://localhost:8080/hello

# 带参数
curl "http://localhost:8080/hello?name=Go"

# 获取用户列表
curl http://localhost:8080/users

# 创建用户（Windows Git Bash 下中文输入可能有编码问题，建议先用 ASCII 测试）
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name":"WangWu","age":28}'
```

## 与 Java/Spring Boot 的关键对比

| 对比项 | Java/Spring Boot | Go |
|--------|------------------|-----|
| 依赖管理 | Maven/Gradle | Go Modules（`go.mod`/`go.sum`） |
| 入口函数 | `public static void main` | `func main()` |
| 访问控制 | `public`/`private` 关键字 | 首字母大小写 |
| 变量声明 | `int a = 10;` | `a := 10` 或 `var a int = 10` |
| 错误处理 | `try-catch` | `if err != nil` + `panic/recover` |
| 多返回值 | 需封装对象 | 原生支持 `(result, err)` |
| 集合 | `ArrayList`、`HashMap` | `slice`、`map` |
| POJO | `class` + getter/setter | `struct` |
| 指针 | 隐式引用 | 显式指针 `*T` 和 `&` |
| Web 路由 | `@GetMapping` 注解 | `http.HandleFunc` |
| 中间件 | Filter/Interceptor | 函数包装（高阶函数） |

## 最佳实践引用

- **Effective Go**：优先使用短变量声明 `:=`；`switch` 默认不穿透，避免遗漏 `break`
- **Go Code Review Comments**：普通错误用 `error` 返回值；方法是否需要指针接收者遵循一致性原则
- **Go 官方错误处理建议**：库代码不要吞掉 panic，除非能真正恢复

## 实战作业

### 作业一：学生成绩管理命令行程序

**要求**：

1. 定义 `Student` 结构体，包含：
   - `ID`（int64）
   - `Name`（string）
   - `Scores`（map[string]float64，科目→分数）

2. 实现以下函数/方法：
   - `AddStudent`：添加学生
   - `AddScore(subject string, score float64)`：添加/更新成绩
   - `Average()`：计算平均分
   - `GetTopStudent`：找出平均分最高的学生

3. 在 `main` 中演示添加、查询、计算功能

4. 添加单元测试（使用 `testing` 包）

**核心思路**：

```go
type Student struct {
    ID     int64
    Name   string
    Scores map[string]float64
}

func (s *Student) AddScore(subject string, score float64) {
    if s.Scores == nil {
        s.Scores = make(map[string]float64)
    }
    s.Scores[subject] = score
}

func (s *Student) Average() float64 {
    if len(s.Scores) == 0 {
        return 0
    }
    var sum float64
    for _, score := range s.Scores {
        sum += score
    }
    return sum / float64(len(s.Scores))
}
```

### 作业二：改造成 HTTP 服务（扩展挑战）

把命令行程序扩展为 HTTP 服务：

- `POST /students`：添加学生
- `GET /students/{id}`：查询学生
- `GET /students/{id}/average`：查询平均分
- `GET /students/top`：返回最高分学生

## 下一步

完成 Phase 1 后，继续学习 **Phase 2：核心进阶**，内容包括：

- 接口与类型系统
- goroutine 与 channel
- select、context、sync 包
- 错误处理最佳实践
- 反射与 unsafe
