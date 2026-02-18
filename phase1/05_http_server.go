package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// User 结构体用于 JSON 响应
type User struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// helloHandler 处理 GET /hello
// 类似 Spring Boot 的 @GetMapping("/hello")
func helloHandler(w http.ResponseWriter, r *http.Request) {
	// w 用于写响应，r 包含请求信息
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "World"
	}
	fmt.Fprintf(w, "Hello, %s!\n", name)
}

// usersHandler 处理 GET /users
// 返回 JSON 列表，类似 Spring Boot 的 ResponseEntity<List<User>>
func usersHandler(w http.ResponseWriter, r *http.Request) {
	users := []User{
		{ID: 1, Name: "张三", Age: 25},
		{ID: 2, Name: "李四", Age: 30},
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// createUserHandler 处理 POST /users
// 解析 JSON 请求体，类似 @RequestBody User user
func createUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Bad Request: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close() // 记得关闭请求体

	// 模拟保存成功，返回 201 Created
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	user.ID = 3
	json.NewEncoder(w).Encode(user)
}

// loggingMiddleware 是一个简单的中间件，记录请求日志
// 类似 Spring Boot 的拦截器 / Filter
func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%s] %s %s", r.Method, r.URL.Path, r.RemoteAddr)
		next(w, r)
	}
}

func main() {
	// 注册路由
	// Java Spring Boot: @RestController + @GetMapping
	// Go 标准库: http.HandleFunc("/path", handler)
	http.HandleFunc("/hello", loggingMiddleware(helloHandler))
	http.HandleFunc("/users", loggingMiddleware(func(w http.ResponseWriter, r *http.Request) {
		// 根据 HTTP 方法分发，类似 Spring Boot 中按注解分发
		switch r.Method {
		case http.MethodGet:
			usersHandler(w, r)
		case http.MethodPost:
			createUserHandler(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}))

	port := ":8080"
	fmt.Printf("HTTP 服务已启动，访问 http://localhost%s\n", port)
	fmt.Println("测试命令：")
	fmt.Println("  curl http://localhost:8080/hello")
	fmt.Println("  curl http://localhost:8080/hello?name=Go")
	fmt.Println("  curl http://localhost:8080/users")
	fmt.Println("  curl -X POST http://localhost:8080/users -H 'Content-Type: application/json' -d '{\"name\":\"王五\",\"age\":28}'")

	// 启动服务，第二个参数为 nil 时使用默认的 ServeMux
	// Java 类似：SpringApplication.run(App.class, args)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
