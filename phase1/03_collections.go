package main

import (
	"fmt"
)

// ==================== 结构体：Go 的 POJO ====================
// Java: public class User { private Long id; private String name; ... }
// Go: 用 struct 定义，字段公开（大写）或私有（小写）
type User struct {
	ID     int64  // 大写开头 = public
	Name   string // 大写开头 = public
	age    int    // 小写开头 = private（包内可见）
	Email  string
}

// 结构体方法（值接收者）
func (u User) String() string {
	return fmt.Sprintf("User{ID=%d, Name=%s, Email=%s}", u.ID, u.Name, u.Email)
}

func main() {
	// ==================== 1. 数组：固定长度 ====================
	// Java: int[] arr = new int[3]; 或 int[] arr = {1, 2, 3};
	// Go 数组长度是类型的一部分：[3]int 和 [5]int 是不同类型
	var arr1 [3]int
	arr1[0] = 10
	arr2 := [3]int{1, 2, 3}
	arr3 := [...]int{1, 2, 3, 4, 5} // 编译器推断长度

	fmt.Printf("arr1=%v, 长度=%d\n", arr1, len(arr1))
	fmt.Printf("arr2=%v, 长度=%d\n", arr2, len(arr2))
	fmt.Printf("arr3=%v, 长度=%d\n", arr3, len(arr3))

	// ==================== 2. 切片（Slice）：动态数组，最常用 ====================
	// 类似 Java 的 ArrayList<int> / List<Integer>
	// 切片是引用类型，底层引用一个数组，包含指针、长度、容量三部分
	s1 := []int{1, 2, 3}              // 字面量创建
	s2 := make([]int, 3, 5)           // 长度3，容量5，类似 Java new int[5] 但只暴露前3个
	s3 := arr2[0:2]                   // 从数组切片，左闭右开 [0,2)

	fmt.Printf("s1=%v, len=%d, cap=%d\n", s1, len(s1), cap(s1))
	fmt.Printf("s2=%v, len=%d, cap=%d\n", s2, len(s2), cap(s2))
	fmt.Printf("s3=%v, len=%d, cap=%d\n", s3, len(s3), cap(s3))

	// 追加元素：append 可能触发底层数组重新分配
	s1 = append(s1, 4, 5)
	fmt.Printf("append后 s1=%v, len=%d, cap=%d\n", s1, len(s1), cap(s1))

	// 遍历切片
	for i, v := range s1 {
		fmt.Printf("s1[%d]=%d\n", i, v)
	}

	// 切片共享底层数组的陷阱
	original := []int{1, 2, 3, 4, 5}
	share := original[1:3] // [2,3]
	share[0] = 99
	fmt.Printf("修改 share 后，original=%v\n", original) // 原始切片也被修改

	// 复制切片：避免共享底层数组
	copySlice := make([]int, len(original))
	copy(copySlice, original)
	copySlice[0] = 100
	fmt.Printf("copySlice=%v, original=%v\n", copySlice, original)

	// ==================== 3. Map：键值对 ====================
	// 类似 Java 的 HashMap<String, Integer>
	// Go map 是引用类型，未初始化时为 nil，不能直接写入
	var m1 map[string]int // nil map
	fmt.Printf("m1=%v, isNil=%v\n", m1, m1 == nil)

	m2 := make(map[string]int) // 初始化空 map
	m2["go"] = 100
	m2["java"] = 90
	fmt.Printf("m2=%v\n", m2)

	// 取值时返回第二个布尔值表示 key 是否存在
	if score, ok := m2["go"]; ok {
		fmt.Printf("go 的分数=%d\n", score)
	} else {
		fmt.Println("go 不存在")
	}

	// 删除 key
	delete(m2, "java")
	fmt.Printf("删除 java 后 m2=%v\n", m2)

	// 字面量创建 map
	m3 := map[string]string{
		"CN": "China",
		"US": "United States",
	}
	fmt.Printf("m3=%v\n", m3)

	// ==================== 4. 结构体 ====================
	// 创建结构体实例的多种方式
	u1 := User{} // 零值：ID=0, Name="", age=0, Email=""
	u1.ID = 1
	u1.Name = "张三"
	u1.Email = "zhangsan@example.com"

	u2 := User{
		ID:    2,
		Name:  "李四",
		Email: "lisi@example.com",
	}

	u3 := User{3, "王五", 25, "wangwu@example.com"} // 按字段顺序赋值，不推荐可读性差

	fmt.Println(u1.String())
	fmt.Println(u2.String())
	fmt.Println(u3.String())

	// 结构体指针
	pu := &User{ID: 4, Name: "赵六", Email: "zhaoliu@example.com"}
	fmt.Println(pu.String()) // 自动解引用，等价于 (*pu).String()

	// ==================== 5. 嵌套结构体 ====================
	type Address struct {
		City   string
		Street string
	}

	type Person struct {
		Name    string
		Age     int
		Address // 嵌入字段（匿名字段），类似 Java 的继承/组合
	}

	p := Person{
		Name: "孙七",
		Age:  30,
		Address: Address{
			City:   "北京",
			Street: "长安街",
		},
	}
	fmt.Printf("Person=%+v\n", p)
	fmt.Printf("城市=%s\n", p.City) // 可直接访问嵌入字段的属性
}
