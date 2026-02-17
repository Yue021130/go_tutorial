package main

import "fmt"

// 值接收者：不会修改原对象
func (c Counter) incrementByValue() {
	c.value++
	fmt.Printf("  值接收者内部：c.value=%d\n", c.value)
}

// 指针接收者：会修改原对象
func (c *Counter) incrementByPointer() {
	c.value++ // 等价于 (*c).value++，Go 自动解引用
	fmt.Printf("  指针接收者内部：c.value=%d\n", c.value)
}

type Counter struct {
	value int
}

func main() {
	// ==================== 1. 什么是指针？====================
	// Java 中：对象变量存的是引用（地址），但没有显式指针操作
	// Go 中：指针是显式的，用 & 取地址，用 * 解引用
	x := 10
	p := &x // p 是 *int 类型，存储 x 的内存地址

	fmt.Printf("x=%d, x 的地址=%p\n", x, &x)
	fmt.Printf("p=%p, p 指向的值=%d\n", p, *p)

	*p = 20 // 通过指针修改 x
	fmt.Printf("修改 *p 后，x=%d\n", x)

	// ==================== 2. 值传递 vs 指针传递 ====================
	// Go 函数参数默认是值传递：传的是副本
	// Java 中对象参数传的是引用副本，可以通过方法修改对象内部状态
	// Go 中如果要修改原变量，必须传指针
	num := 10
	fmt.Println("\n===== 值传递 vs 指针传递 =====")
	fmt.Printf("修改前 num=%d\n", num)
	changeByValue(num)
	fmt.Printf("值传递修改后 num=%d（不变）\n", num)
	changeByPointer(&num)
	fmt.Printf("指针传递修改后 num=%d（变了）\n", num)

	// ==================== 3. 结构体指针 ====================
	fmt.Println("\n===== 结构体指针 =====")
	c1 := Counter{value: 0}
	c1.incrementByValue()
	fmt.Printf("值接收者调用后：c1.value=%d（不变）\n", c1.value)

	c1.incrementByPointer()
	fmt.Printf("指针接收者调用后：c1.value=%d（变了）\n", c1.value)

	// 注意：即使 c1 是值类型，调用指针接收者方法时 Go 会自动取地址
	c1.incrementByPointer()
	fmt.Printf("值类型自动取地址后：c1.value=%d\n", c1.value)

	// ==================== 4. new 关键字 ====================
	// Java: Object obj = new Object();
	// Go: new(T) 返回 *T，分配零值内存
	pi := new(int) // pi 是 *int
	*pi = 100
	fmt.Printf("\nnew(int) 分配的值：%d\n", *pi)

	pc := new(Counter)
	pc.incrementByPointer()
	fmt.Printf("new(Counter) 的值：%d\n", pc.value)

	// ==================== 5. 指针的常见使用场景 ====================
	// 1) 函数内修改调用者变量
	// 2) 避免大结构体拷贝，提高性能
	// 3) 实现链表、树等数据结构
	// 4) 方法需要修改接收者状态时用指针接收者

	// ==================== 6. Go 指针 vs Java 引用的差异 ====================
	// - Java 没有指针算术，Go 也没有（不允许 p++ 等操作）
	// - Go 可以返回局部变量地址（编译器会分配到堆上），Java 基本类型不能
	// - Go 指针可以是 nil，Java 对象引用可以是 null
	// - Go 的 map、slice、channel 本身就是引用类型，通常不需要再取地址

	// 示例：返回局部变量地址（Go 会自动逃逸到堆）
	pLocal := createInt()
	fmt.Printf("\n返回的局部变量地址=%p, 值=%d\n", pLocal, *pLocal)
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
	return &v // 返回局部变量地址，Go 会将其分配到堆上（逃逸分析）
}
