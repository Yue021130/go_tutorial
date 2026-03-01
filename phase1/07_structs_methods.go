// 07_structs_methods.go
// 本节目标：掌握结构体（struct）、方法（method）、嵌入字段、结构体标签。
//
// 与 Java 对比：
// - Java 用 class 定义对象，有封装、继承、多态。
// - Go 没有 class，用 struct 定义数据，用 method 绑定行为。
// - Go 没有继承，只有组合（embedding）。
// - Go 的可见性由首字母大小写控制，不需要 public/private 关键字。

package main

import (
	"encoding/json"
	"fmt"
)

// User 结构体：Go 的 POJO
// 大写开头 = exported（包外可见）
// 小写开头 = unexported（包内私有）
type User struct {
	ID    int64  `json:"id"` // 结构体标签（struct tag），用于反射/JSON 等
	Name  string `json:"name"`
	Email string `json:"email,omitempty"` // omitempty：空值时不序列化
	age   int    // 小写开头，包外不可见
}

// String 方法：值接收者
// 类似 Java 的 toString()
func (u User) String() string {
	return fmt.Sprintf("User{ID=%d, Name=%s, Email=%s}", u.ID, u.Name, u.Email)
}

// SetAge 方法：指针接收者
// 因为需要修改接收者本身
func (u *User) SetAge(age int) {
	u.age = age
}

// GetAge 方法：值接收者
func (u User) GetAge() int {
	return u.age
}

// ==================== 嵌入字段（Embedding）====================
// Go 没有继承，但可以通过嵌入字段实现类似"组合 + 方法提升"的效果。
//
// 嵌入字段的规则：
// - 匿名字段（只有类型名，没有字段名）
// - 嵌入字段的方法会被"提升"到外层结构体
// - 可以像访问自己的字段一样访问嵌入字段的字段

type Address struct {
	City   string
	Street string
}

func (a Address) FullAddress() string {
	return a.City + " " + a.Street
}

type Person struct {
	Name    string
	Age     int
	Address // 嵌入字段，匿名字段
}

func main() {
	// ==================== 1. 结构体创建方式 ====================
	// 方式 1：零值创建
	u1 := User{}
	u1.ID = 1
	u1.Name = "张三"
	u1.Email = "zhangsan@example.com"
	u1.SetAge(25)

	// 方式 2：按字段名赋值（推荐，可读性好）
	u2 := User{
		ID:    2,
		Name:  "李四",
		Email: "lisi@example.com",
	}

	// 方式 3：按字段顺序赋值（不推荐）
	u3 := User{3, "王五", "wangwu@example.com", 30}

	fmt.Println(u1.String())
	fmt.Println(u2.String())
	fmt.Println(u3.String())
	fmt.Printf("u1 的年龄=%d\n", u1.GetAge())

	// ==================== 2. 结构体指针 ====================
	// Go 会自动解引用，访问字段和方法都不需要显式写 (*pu).Name
	pu := &User{ID: 4, Name: "赵六", Email: "zhaoliu@example.com"}
	fmt.Println(pu.String())
	pu.SetAge(35)
	fmt.Printf("pu 的年龄=%d\n", pu.GetAge())

	// ==================== 3. JSON 序列化与反序列化 ====================
	fmt.Println("\n--- JSON 序列化 ---")
	user := User{
		ID:    100,
		Name:  "孙悟空",
		Email: "",
	}
	jsonBytes, err := json.Marshal(user)
	if err != nil {
		fmt.Println("序列化失败:", err)
		return
	}
	fmt.Printf("JSON: %s\n", string(jsonBytes))
	// 输出：{"id":100,"name":"孙悟空"}，email 因为 omitempty 且为空被省略

	fmt.Println("\n--- JSON 反序列化 ---")
	jsonStr := `{"id":200,"name":"猪八戒","email":"bajie@example.com"}`
	var parsed User
	if err := json.Unmarshal([]byte(jsonStr), &parsed); err != nil {
		fmt.Println("反序列化失败:", err)
		return
	}
	fmt.Printf("解析结果: %s\n", parsed.String())

	// ==================== 4. 嵌入字段 ====================
	fmt.Println("\n--- 嵌入字段 ---")
	p := Person{
		Name: "孙七",
		Age:  30,
		Address: Address{
			City:   "北京",
			Street: "长安街",
		},
	}

	// 可以直接访问嵌入字段的字段
	fmt.Printf("城市=%s, 街道=%s\n", p.City, p.Street)
	// 也可以显式通过嵌入字段类型访问
	fmt.Printf("城市=%s, 街道=%s\n", p.Address.City, p.Address.Street)
	// 嵌入字段的方法被提升
	fmt.Printf("完整地址=%s\n", p.FullAddress())

	// ==================== 5. 结构体比较 ====================
	// 结构体只有在所有字段都可比较时才能用 == 比较
	// slice、map、function 字段会导致结构体不可比较
	type Point struct {
		X, Y int
	}
	p1 := Point{1, 2}
	p2 := Point{1, 2}
	p3 := Point{2, 3}
	fmt.Printf("\np1 == p2 ? %v\n", p1 == p2)
	fmt.Printf("p1 == p3 ? %v\n", p1 == p3)

	// ==================== 常见坑 ====================
	// 坑 1：结构体是值类型，函数传参会复制整个结构体。
	// 坑 2：方法接收者用值还是指针，需要一致考虑（详细见 08_pointers.go）。
	// 坑 3：结构体标签只对反射可见，普通访问不影响。
	// 坑 4：嵌入字段的方法提升可能与外层方法冲突，需要注意。
}
