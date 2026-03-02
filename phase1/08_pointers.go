// 08_pointers.go
// 本节目标：彻底理解 Go 指针，这是 Java 开发者最容易困惑的地方。
//
// 核心认知：
// - Java 中：基本类型存值，对象类型存引用（地址副本）。
// - Go 中：所有类型默认都是值传递；想修改原变量必须显式传指针。
// - Go 指针没有指针算术（不能 p++），比 C/C++ 安全。

package main

import "fmt"

type Counter struct {
	value int
}

// 值接收者：操作的是副本
func (c Counter) incrementByValue() {
	c.value++
	fmt.Printf("  值接收者内部：value=%d\n", c.value)
}

// 指针接收者：操作的是原对象
func (c *Counter) incrementByPointer() {
	c.value++ // 等价于 (*c).value++，Go 自动解引用
	fmt.Printf("  指针接收者内部：value=%d\n", c.value)
}

func main() {
	// ==================== 1. 指针基础 ====================
	// & 取地址，* 解引用
	//
	// 内存示意：
	//   栈
	//   ┌─────────┐
	//   │ x = 10  │
	//   │ p = &x ─┼──→ 同一块内存
	//   └─────────┘

	x := 10
	p := &x // p 的类型是 *int

	fmt.Printf("x=%d, x 的地址=%p\n", x, &x)
	fmt.Printf("p=%p, p 指向的值=%d, p 的类型=%T\n", p, *p, p)

	*p = 20 // 通过指针修改 x
	fmt.Printf("修改 *p 后，x=%d\n", x)

	// ==================== 2. 值传递 vs 指针传递 ====================
	// Go 函数参数默认是值传递：传的是副本。
	// Java 中对象参数传的是引用副本，可以通过方法修改对象内部状态。
	// Go 中如果要修改原变量，必须传指针。

	fmt.Println("\n--- 值传递 vs 指针传递 ---")
	num := 10
	fmt.Printf("修改前 num=%d\n", num)

	changeByValue(num)
	fmt.Printf("值传递修改后 num=%d（不变）\n", num)

	changeByPointer(&num)
	fmt.Printf("指针传递修改后 num=%d（变了）\n", num)

	// ==================== 3. new vs &T{} ====================
	// new(T) 返回 *T，分配零值内存
	// &T{} 也返回 *T，但可以初始化字段

	fmt.Println("\n--- new vs &T{} ---")
	pi := new(int) // *int，值为 0
	*pi = 100
	fmt.Printf("new(int): *pi=%d\n", *pi)

	pc1 := new(Counter) // *Counter，value 为 0
	pc1.incrementByPointer()
	fmt.Printf("new(Counter): value=%d\n", pc1.value)

	pc2 := &Counter{value: 10}
	pc2.incrementByPointer()
	fmt.Printf("&Counter{}: value=%d\n", pc2.value)

	// ==================== 4. 值接收者 vs 指针接收者 ====================
	fmt.Println("\n--- 值接收者 vs 指针接收者 ---")
	c := Counter{value: 0}

	c.incrementByValue()
	fmt.Printf("值接收者调用后：value=%d（不变）\n", c.value)

	c.incrementByPointer()
	fmt.Printf("指针接收者调用后：value=%d（变了）\n", c.value)

	// 即使 c 是值类型，调用指针接收者方法时 Go 会自动取地址
	c.incrementByPointer()
	fmt.Printf("再次调用指针接收者：value=%d\n", c.value)

	// ==================== 5. 指针接收者 vs 值接收者决策树 ====================
	//
	// 方法是否需要修改接收者的状态？
	//   是  → 指针接收者
	//   否  → 继续问：接收者是否很大？
	//          是  → 指针接收者（避免拷贝）
	//          否  → 值接收者
	//
	// 一致性原则：同一个类型的方法，尽量保持一致的接收者类型。
	//            不要有的用值、有的用指针，会让调用者困惑。

	// ==================== 6. 返回局部变量地址 ====================
	// Go 编译器会进行逃逸分析：如果局部变量在函数返回后仍被引用，
	// 会将其分配到堆上，而不是栈上。
	// Java 中基本类型不能返回地址，对象可以返回引用。

	pLocal := createInt()
	fmt.Printf("\n返回的局部变量地址=%p, 值=%d\n", pLocal, *pLocal)

	// ==================== 7. 哪些类型本身就是引用类型？====================
	// Go 中 map、slice、channel、function、interface 本身就是引用类型，
	// 通常不需要再取地址。
	//
	// 例如：
	// func modifySlice(s []int) { s[0] = 100 } // 不需要 * []
	// func modifyMap(m map[string]int) { m["a"] = 1 } // 不需要 * map

	// ==================== Go 指针 vs Java 引用对比表 ====================
	// | 特性            | Java 引用              | Go 指针              |
	// |----------------|----------------------|---------------------|
	// | 基本类型        | 值传递                 | 值传递               |
	// | 对象/结构体     | 引用传递（地址副本）     | 值传递，需显式指针     |
	// | 取地址          | 不支持（对象天然是引用）  | & 显式取地址          |
	// | 解引用          | 不支持                 | * 显式解引用          |
	// | 指针算术        | 无                     | 无（安全）            |
	// | 空值            | null                   | nil                  |
	// | 返回局部变量地址 | 对象可以，基本类型不可以  | 都可以（逃逸分析）     |
}

func changeByValue(n int) {
	n = 100
	fmt.Printf("  changeByValue 内部 n=%d\n", n)
}

func changeByPointer(n *int) {
	*n = 100
	fmt.Printf("  changeByPointer 内部 *n=%d\n", *n)
}

func createInt() *int {
	v := 42
	return &v // 返回局部变量地址，Go 会将其分配到堆上
}
