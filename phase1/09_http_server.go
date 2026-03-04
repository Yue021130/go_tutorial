// 09_http_server.go
// 本节目标：使用 Go 标准库 net/http 实现一个可运行的 HTTP 服务。
//
// 与 Spring Boot 对比：
// - Spring Boot: @RestController + @GetMapping/@PostMapping + @RequestBody + Filter
// - Go 标准库: http.HandleFunc + handler function + 自定义中间件
//
// 本示例包含：
// - 路由注册与 HTTP 方法分发
// - 查询参数、路径参数（简单实现）
// - JSON 请求体解析与响应
// - 中间件：日志、panic 恢复

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// User 结构体用于 JSON 序列化/反序列化
type User struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// 内存中的用户数据（模拟数据库）
var users = map[int64]User{
	1: {ID: 1, Name: "张三", Age: 25},
	2: {ID: 2, Name: "李四", Age: 30},
}
var nextUserID int64 = 3

// ==================== Handler 函数 ====================
// Handler 函数签名：func(w http.ResponseWriter, r *http.Request)
// w 用于写响应，r 包含请求信息

func helloHandler(w http.ResponseWriter, r *http.Request) {
	// 获取查询参数，类似 Spring Boot 的 @RequestParam
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "World"
	}
	fmt.Fprintf(w, "Hello, %s!\n", name)
}

func listUsersHandler(w http.ResponseWriter, r *http.Request) {
	// 将 map 转为 slice
	userList := make([]User, 0, len(users))
	for _, u := range users {
		userList = append(userList, u)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(userList); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	// 简单路径参数解析：/users/1
	// 生产环境推荐使用第三方路由库如 gorilla/mux 或 chi
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 2 {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, ok := users[id]
	if !ok {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	// 解析 JSON 请求体，类似 Spring Boot 的 @RequestBody
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Bad Request: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close() // 记得关闭请求体

	user.ID = nextUserID
	nextUserID++
	users[user.ID] = user

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// usersRouter 根据 HTTP 方法分发到不同 handler
// 类似 Spring Boot 中按注解自动分发
func usersRouter(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// /users 或 /users/1
		if r.URL.Path == "/users" || r.URL.Path == "/users/" {
			listUsersHandler(w, r)
		} else {
			getUserHandler(w, r)
		}
	case http.MethodPost:
		createUserHandler(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// ==================== 中间件 ====================
// 中间件是高阶函数：接收一个 handler，返回一个新的 handler。
// 类似 Spring Boot 的 Filter / Interceptor。

type Middleware func(http.HandlerFunc) http.HandlerFunc

// loggingMiddleware 记录请求日志
func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%s] %s %s", r.Method, r.URL.Path, r.RemoteAddr)
		next(w, r)
	}
}

// recoveryMiddleware 捕获 panic，防止单个请求拖垮整个服务
func recoveryMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic recovered: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next(w, r)
	}
}

// chain 把多个中间件组合起来
func chain(h http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}

func main() {
	// 注册路由
	// http.HandleFunc 使用默认的 ServeMux
	// 第一个参数是路径前缀，第二个参数是 handler
	http.HandleFunc("/hello", chain(helloHandler, loggingMiddleware, recoveryMiddleware))
	// 使用 /users/ 作为前缀，可以同时匹配 /users/ 和 /users/1
	// 访问 /users 会被自动重定向到 /users/
	http.HandleFunc("/users/", chain(usersRouter, loggingMiddleware, recoveryMiddleware))

	port := ":8080"
	fmt.Printf("HTTP 服务已启动，访问 http://localhost%s\n", port)
	fmt.Println("测试命令：")
	fmt.Println("  curl http://localhost:8080/hello")
	fmt.Println("  curl \"http://localhost:8080/hello?name=Go\"")
	fmt.Println("  curl http://localhost:8080/users/")
	fmt.Println("  curl http://localhost:8080/users/1")
	fmt.Println("  curl -X POST http://localhost:8080/users/ -H \"Content-Type: application/json\" -d '{\"name\":\"WangWu\",\"age\":28}'")

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
