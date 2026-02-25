# Phase 1：Go 快速上手（第 1-3 天）

## 学习目标

通过 3 天学习，达到以下能力：

1. 理解 Go 的工程组织方式：包、导入、可见性、Go Modules
2. 掌握 Go 基础语法：变量、常量、类型、控制流
3. 掌握函数、多返回值、defer、panic/recover
4. 深入理解数组、切片、map、结构体及其常见陷阱
5. 理解指针本质，能正确选择值接收者/指针接收者
6. 能用标准库 `net/http` 编写可运行的 HTTP 服务
7. 会写简单的单元测试

## 预计耗时

- 第 1 天：环境 + 基础语法（文件 00-02）
- 第 2 天：函数 + 复合类型（文件 03-06）
- 第 3 天：结构体 + 指针 + HTTP 服务 + 作业（文件 07-09 + homework）
- 每天约 2-3 小时

## 学习路径

```text
Day 1: 00_hello.go → 01_variables_types_constants.go → 02_control_flow.go
Day 2: 03_functions.go → 04_defer_panic_recover.go → 05_arrays_and_slices.go → 06_maps.go
Day 3: 07_structs_methods.go → 08_pointers.go → 09_http_server.go → homework/
```

## 文件说明

| 文件 | 主题 | 重点 |
|------|------|------|
| `00_hello.go` | 第一个程序 | package、import、可见性 |
| `01_variables_types_constants.go` | 变量、常量、类型 | 类型转换、零值、字符串不可变 |
| `02_control_flow.go` | 控制流 | if 短变量、for 四种形式、switch/range |
| `03_functions.go` | 函数 | 多返回值、命名返回值、可变参数、闭包 |
| `04_defer_panic_recover.go` | defer、panic、recover | defer 求值时机、命名返回值交互、panic 传播 |
| `05_arrays_and_slices.go` | 数组与切片 | 切片 header、扩容、共享底层数组陷阱 |
| `06_maps.go` | map | nil map、遍历无序、线程不安全 |
| `07_structs_methods.go` | 结构体与方法 | struct、method、tag、JSON、嵌入字段 |
| `08_pointers.go` | 指针 | &/*、new、逃逸分析、接收者决策树 |
| `09_http_server.go` | HTTP 服务 | net/http、中间件、RESTful 风格 |
| `homework/` | 实战作业 | 完整参考答案 + 单元测试 |

## 环境检查

```bash
# 检查 Go 版本
go version

# 进入 Phase 1 目录
cd phase1

# 格式化检查所有文件
go fmt ./...

# 运行示例
go run 00_hello.go
```

---

## Day 1：环境 + 基础语法

### 1.1 包（Package）与可见性

**核心问题**：Go 如何组织代码？哪些东西能被外部访问？

**与 Java 对比**：

| 特性 | Java | Go |
|------|------|-----|
| 代码组织 | package + class | package + file |
| 一个文件 | 通常一个 public class | 一个 package，可包含多个类型/函数 |
| 访问控制 | `public`/`private` 关键字 | 首字母大小写 |
| 入口 | `public static void main` | `func main()` |

**Go 可见性规则**：
- 首字母大写：`exported`，包外可访问（类似 `public`）
- 首字母小写：`unexported`，仅包内可访问（类似 `private`）
- 适用于：函数、变量、常量、类型、结构体字段、方法

**运行命令**：

```bash
go run 00_hello.go
go run 00_hello.go Go
go run 00_hello.go "Go 语言"
```

### 1.2 变量、常量与类型

**核心问题**：Go 如何声明变量？类型转换有什么限制？零值是什么？

**与 Java 的关键差异**：

1. **变量声明**：Go 的类型在变量名后面
   ```go
   var a int = 10   // 完整声明
   var b = 20       // 类型推断
   c := 30          // 短变量声明（函数内最常用）
   ```

2. **常量枚举**：Go 用 `const` + `iota`，Java 用 `enum`
   ```go
   const (
       Monday = iota  // 0
       Tuesday        // 1
       Wednesday      // 2
   )
   ```

3. **类型转换**：Go 是强类型，不允许隐式转换
   ```go
   var i int = 10
   var f float64 = float64(i)  // 必须显式转换
   ```

4. **零值**：Go 中声明未赋值的变量有零值；Java 局部变量必须显式初始化
   ```go
   var s string  // s = ""
   var p *int    // p = nil
   ```

**完整数据类型对比表**：

| Java | Go | 说明 |
|------|-----|------|
| byte | byte (= uint8) | 0-255 |
| short | int16 | |
| int | int32 / int | int 在 64 位系统为 64 位 |
| long | int64 | |
| float | float32 | |
| double | float64 | |
| char | rune (= int32) | Unicode 码点 |
| boolean | bool | |
| String | string | 不可变，UTF-8 编码 |
| Object | interface{} | 空接口 |

**常见坑**：
- `:=` 只能在函数内部使用，包级变量必须用 `var`
- `:=` 会遮蔽外层同名变量
- `string(65)` 得到 `"A"`，不是 `"65"`
- `int` 的大小依赖平台，跨平台代码建议用 `int64`

**运行命令**：

```bash
go run 01_variables_types_constants.go
```

### 1.3 控制流

**核心问题**：Go 如何控制程序流程？

**与 Java 的关键差异**：

1. **Go 没有 `while`、`do-while`，只有 `for`**：
   ```go
   // C-style
   for i := 0; i < 5; i++ { }

   // while-like
   for condition { }

   // 无限循环
   for { }
   ```

2. **`if` 支持短变量声明**：
   ```go
   if v := score + 10; v >= 90 {
       // v 的作用域只在 if-else 块内
   }
   ```

3. **`switch` 默认不穿透**：不需要 `break`，避免 Java 中常见的漏写 break bug
   - 需要穿透时显式使用 `fallthrough`
   - `switch` 可以没有表达式，变成 `if-else if` 的清晰写法

4. **`range` 遍历**：
   - 数组/切片：返回 `index, value`
   - map：返回 `key, value`，顺序随机
   - string：按 `rune` 遍历，正确处理中文

**常见坑**：
- `for-range` 中的 `value` 是副本，对值类型修改无效
- `for-range` 的 `index/value` 变量在每次迭代中复用，闭包中需要创建局部副本

**运行命令**：

```bash
go run 02_control_flow.go
```

---

## Day 2：函数 + 复合类型

### 2.1 函数

**核心问题**：Go 函数有什么特别之处？

**与 Java 的关键差异**：

| 特性 | Java | Go |
|------|------|-----|
| 归属 | 类的方法 | 独立的函数，也可绑定到类型 |
| 返回值 | 一个 | 多个（原生支持） |
| 可变参数 | `int... nums` | `nums ...int` |
| 一等公民 | 需要 Lambda/Function 接口 | 原生支持 |
| 闭包 | Lambda 捕获 | 更直接自然 |

**Go 错误处理模式**：

```go
result, err := divide(10, 3)
if err != nil {
    // 处理错误
}
```

这是 Go 最重要的惯用法之一，对应 Java 的异常处理。

**命名返回值**：
- 优点：代码更短
- 缺点：大函数可读性差
- Effective Go 建议：只在短函数中使用

**闭包陷阱**：
```go
for i := 0; i < 3; i++ {
    v := i  // 必须创建局部副本
    funcs[i] = func() { fmt.Println(v) }
}
```

**运行命令**：

```bash
go run 03_functions.go
```

### 2.2 defer、panic、recover

**核心问题**：Go 如何做资源清理和异常处理？

**defer 要点**：
- 在函数返回前执行
- 执行顺序：LIFO（后进先出）
- **参数在 defer 语句处求值**，不是返回时求值
- 可以修改命名返回值

**panic/recover 要点**：
- `panic` 类似 Java 的 `throw RuntimeException`
- `recover` 类似 `catch`，但必须在 `defer` 中调用
- 普通错误用 `error` 返回值；不可恢复的严重问题才用 `panic`

**panic 传播路径**：

```text
main
  └── funcA
        └── funcB
              └── panic

如果 funcB 没有 recover → 传播到 funcA
如果 funcA 有 recover → 捕获 panic，funcA 后续代码继续执行
main 也会继续执行
```

**运行命令**：

```bash
go run 04_defer_panic_recover.go
```

### 2.3 数组与切片

**核心问题**：Go 的"数组"和"切片"有什么区别？为什么切片那么容易踩坑？

**数组 vs 切片**：

| 特性 | 数组 `[3]int` | 切片 `[]int` |
|------|--------------|-------------|
| 长度 | 固定，是类型一部分 | 动态 |
| 类型 | 值类型 | 引用类型 |
| 传递 | 整体复制 | 只复制 header |
| 常用度 | 低 | 高 |

**切片 Header 结构**：

```text
切片变量 s
┌─────────────┐
│ ptr → 底层数组 │
│ len         │
│ cap         │
└─────────────┘
```

**扩容机制**：
- `cap` 不足时，`append` 会分配新的底层数组
- 新容量通常按当前容量的 1.25~2 倍增长（小于 1024 时翻倍）

**常见坑**：
1. 切片共享底层数组，修改一个可能影响另一个
2. `append` 可能触发重新分配，导致原切片"失联"
3. `nil` slice 和 empty slice 在 JSON 序列化时不同
4. `range` 中修改切片元素副本无效

**运行命令**：

```bash
go run 05_arrays_and_slices.go
```

### 2.4 map

**核心问题**：Go 的 map 有什么特点和限制？

**与 Java HashMap 的差异**：

| 特性 | Java HashMap | Go map |
|------|-------------|--------|
| 空容器 | `new HashMap<>()` 后可用 | nil map 不能直接写入 |
| 取值 | `map.get(key)` 可能返回 null | `value, ok := m[key]` 两值模式 |
| 遍历顺序 | 通常稳定 | 随机，不要依赖顺序 |
| 线程安全 | 非线程安全 | 非线程安全 |

**nil map 陷阱**：

```go
var m map[string]int
m["key"] = 1  // panic: assignment to entry in nil map
```

**map key 要求**：必须是可比较类型（comparable）
- 允许：`bool`、数字、`string`、指针、channel、interface、包含这些类型的数组/结构体
- 不允许：slice、map、function

**运行命令**：

```bash
go run 06_maps.go
```

---

## Day 3：结构体 + 指针 + HTTP 服务

### 3.1 结构体与方法

**核心问题**：Go 没有 class，如何定义"对象"？

**与 Java 的对比**：

| 特性 | Java | Go |
|------|------|-----|
| 类型定义 | class | struct |
| 方法 | 类内部 | 绑定到类型（receiver） |
| 构造 | 构造函数 | 普通函数或字面量 |
| 继承 | extends | 嵌入字段（组合） |
| 封装 | private/public 关键字 | 首字母大小写 |

**结构体标签（tag）**：

```go
type User struct {
    ID    int64  `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email,omitempty"`
}
```

标签用于反射，最常见的用途是 JSON 序列化/反序列化。

**嵌入字段**：
- 匿名字段
- 嵌入字段的方法会被"提升"到外层结构体
- 可以直接访问嵌入字段的字段

**运行命令**：

```bash
go run 07_structs_methods.go
```

### 3.2 指针

**核心问题**：Go 的指针和 Java 的引用有什么本质区别？

**Java 开发者常见误解**：
- "Java 没有指针" → 错误，Java 对象变量本质上就是引用（受限的指针）
- "Go 指针很危险" → 错误，Go 指针没有指针算术，比 C/C++ 安全得多

**Go 指针的本质**：
- `&`：取地址
- `*`：解引用
- 所有类型默认都是值传递
- 想修改原变量，必须传指针

**值接收者 vs 指针接收者决策树**：

```text
方法是否需要修改接收者的状态？
  ├── 是 → 指针接收者
  └── 否 → 接收者是否很大？
           ├── 是 → 指针接收者（避免拷贝）
           └── 否 → 值接收者
```

**一致性原则**：同一个类型的方法，尽量保持一致的接收者类型。

**逃逸分析**：
- Go 可以安全返回局部变量地址
- 编译器会判断变量是否"逃逸"到函数外部
- 如果会逃逸，分配到堆上；否则分配到栈上

**运行命令**：

```bash
go run 08_pointers.go
```

### 3.3 HTTP 服务

**核心问题**：不用框架，如何用 Go 标准库写一个 HTTP 服务？

**与 Spring Boot 的对比**：

| Spring Boot | Go 标准库 |
|-------------|----------|
| `@RestController` | `http.HandleFunc` |
| `@GetMapping` | 在 handler 中判断 `r.Method` |
| `@RequestParam` | `r.URL.Query().Get` |
| `@PathVariable` | 手动解析 `r.URL.Path` |
| `@RequestBody` | `json.NewDecoder(r.Body).Decode` |
| `Filter` | 函数包装中间件 |

**中间件模式**：

```go
func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        log.Printf("[%s] %s", r.Method, r.URL.Path)
        next(w, r)
    }
}
```

这是 Go 中非常常见的高阶函数应用。

**运行命令**：

```bash
# 启动服务（会阻塞终端）
go run 09_http_server.go

# 在另一个终端测试
curl http://localhost:8080/hello
curl "http://localhost:8080/hello?name=Go"
curl http://localhost:8080/users/
curl http://localhost:8080/users/1
curl -X POST http://localhost:8080/users/ \
  -H "Content-Type: application/json" \
  -d '{"name":"WangWu","age":28}'
```

---

## 实战作业

### 作业要求

实现一个学生成绩管理系统，要求：

1. 定义 `Student` 结构体：
   - `ID`（int64）
   - `Name`（string）
   - `Scores`（map[string]float64）

2. 实现 `Student` 方法：
   - `AddScore(subject string, score float64) error`
   - `Average() float64`

3. 实现 `StudentManager`：
   - `AddStudent(name string) (int64, error)`
   - `GetStudent(id int64) (*Student, error)`
   - `DeleteStudent(id int64) error`
   - `AddScore(id int64, subject string, score float64) error`
   - `GetAverage(id int64) (float64, error)`
   - `GetTopStudent() (*Student, error)`
   - `ListStudents() []*Student`

4. 编写单元测试，覆盖率越高越好

### 参考答案

本目录下 `homework/` 提供了完整参考答案：

```bash
cd homework

# 运行命令行演示
go run .

# 运行测试
go test -v

# 查看覆盖率
go test -cover
```

---

## 常见错误汇总

### 1. 变量声明错误

```go
// 错误：:= 不能在包级使用
package main
a := 10  // 编译错误

// 正确
var a = 10
```

### 2. nil map 写入

```go
var m map[string]int
m["key"] = 1  // panic

// 正确
m := make(map[string]int)
m["key"] = 1
```

### 3. 切片共享底层数组

```go
a := []int{1, 2, 3, 4, 5}
b := a[1:3]
b[0] = 99
// a 变成 [1, 99, 3, 4, 5]
```

### 4. 值接收者无法修改状态

```go
func (c Counter) increment() { c.value++ }
// 调用后原 Counter 的 value 不变

// 正确
func (c *Counter) increment() { c.value++ }
```

### 5. defer 参数提前求值

```go
i := 0
defer fmt.Println(i)  // 输出 0，不是 100
i = 100
```

---

## 最佳实践引用

- **Effective Go**：
  - 优先使用短变量声明 `:=`
  - `switch` 默认不穿透，避免遗漏 `break`
  - 命名返回值只在短函数中使用
- **Go Code Review Comments**：
  - 普通错误用 `error` 返回值
  - 同一个类型的方法尽量保持接收者类型一致
  - 库代码不要吞掉 panic
- **Go 官方建议**：
  - 不可恢复的严重问题才用 panic
  - HTTP 服务等顶层代码用 recover 防止单个请求拖垮进程

---

## 延伸阅读

- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Go by Example](https://gobyexample.com/)
- [Go 官方文档](https://go.dev/doc/)

---

## 下一步

完成 Phase 1 后，继续学习 **Phase 2：核心进阶**，内容包括：

- 接口与类型系统
- goroutine 与 channel
- select、context、sync 包
- 错误处理最佳实践
- 反射与 unsafe
