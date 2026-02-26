// 03_functions.go
// 本节目标：掌握 Go 函数的所有核心特性。
//
// 与 Java 对比：
// - Java 方法属于类；Go 函数是独立的，也可以绑定到类型（方法）。
// - Java 方法只能返回一个值；Go 支持多返回值。
// - Go 函数是一等公民：可以赋值给变量、作为参数、作为返回值。

package main

import (
	"errors"
	"fmt"
)

// ==================== 1. 函数声明 ====================
// Java: public int add(int a, int b) { return a + b; }
// Go:  func add(a int, b int) int { return a + b }
// 连续同类型参数可简写：func add(a, b int) int
func add(a, b int) int {
	return a + b
}

// ==================== 2. 多返回值（Go 标志性特性）====================
// Java 要返回多个值，通常封装对象或数组；Go 原生支持。
// 最常见模式：result, err := doSomething()
// 这正是 Go 错误处理的核心。
func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("除数不能为零")
	}
	return a / b, nil
}

// ==================== 3. 命名返回值 ====================
// 在函数签名中声明返回值变量名，可以直接使用，也可以用裸 return。
// 优点：代码更短；缺点：大函数中可读性差。
// Effective Go 建议：命名返回值只在短函数中使用。
func rectangle(width, height float64) (area float64, perimeter float64) {
	area = width * height
	perimeter = 2 * (width + height)
	return // 裸 return，返回命名返回值
}

// ==================== 4. 可变参数 ====================
// Java: public int sum(int... nums) { ... }
// Go:  func sum(nums ...int) int
// nums 在函数内部实际上是 []int 切片。
func sum(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

// 可变参数 + 普通参数混用，可变参数必须放在最后
func printInfo(prefix string, nums ...int) {
	fmt.Printf("%s: ", prefix)
	for _, n := range nums {
		fmt.Printf("%d ", n)
	}
	fmt.Println()
}

// ==================== 5. 函数作为一等公民 ====================
// 类似 Java 的 Function<T, R>、Predicate<T>、Consumer<T>，但更直接。

// 函数作为参数（高阶函数）
func apply(a, b int, op func(int, int) int) int {
	return op(a, b)
}

// 函数作为返回值（工厂函数）
func makeMultiplier(factor int) func(int) int {
	return func(x int) int {
		return x * factor
	}
}

// ==================== 6. 闭包 ====================
// 闭包 = 函数 + 其引用的外部变量。
// Java 中类似 lambda 捕获外部变量，但 Go 的闭包更灵活。
func makeCounter() func() int {
	n := 0
	return func() int {
		n++
		return n
	}
}

func main() {
	fmt.Println("===== 函数基础 =====")
	fmt.Printf("add(3, 5) = %d\n", add(3, 5))

	fmt.Println("\n===== 多返回值与错误处理 =====")
	result, err := divide(10, 3)
	if err != nil {
		fmt.Printf("错误：%v\n", err)
	} else {
		fmt.Printf("10 / 3 = %.2f\n", result)
	}

	if result, err = divide(10, 0); err != nil {
		fmt.Printf("10 / 0 错误：%v\n", err)
	} else {
		fmt.Printf("10 / 0 = %.2f\n", result)
	}

	fmt.Println("\n===== 命名返回值 =====")
	area, perimeter := rectangle(4, 5)
	fmt.Printf("矩形面积=%.2f, 周长=%.2f\n", area, perimeter)

	fmt.Println("\n===== 可变参数 =====")
	fmt.Printf("sum(1,2,3)=%d\n", sum(1, 2, 3))
	fmt.Printf("sum()=%d\n", sum())
	printInfo("成绩", 85, 90, 78)

	// 把切片展开为可变参数
	nums := []int{1, 2, 3, 4, 5}
	fmt.Printf("sum(nums...)=%d\n", sum(nums...))

	fmt.Println("\n===== 函数作为参数/返回值 =====")
	mul := func(a, b int) int { return a * b }
	fmt.Printf("apply(4, 5, mul) = %d\n", apply(4, 5, mul))

	double := makeMultiplier(2)
	triple := makeMultiplier(3)
	fmt.Printf("double(5)=%d, triple(5)=%d\n", double(5), triple(5))

	fmt.Println("\n===== 闭包 =====")
	counter1 := makeCounter()
	counter2 := makeCounter()
	fmt.Printf("counter1: %d, %d, %d\n", counter1(), counter1(), counter1())
	fmt.Printf("counter2: %d, %d\n", counter2(), counter2())

	// ==================== 常见坑：for-range + 闭包 ====================
	// 坑：range 的 value 变量在每次迭代中复用，闭包捕获的是同一个变量。
	funcs := make([]func(), 3)
	for i := 0; i < 3; i++ {
		v := i // 必须创建局部副本
		funcs[i] = func() {
			fmt.Printf("%d ", v)
		}
	}
	fmt.Println("\n===== 闭包捕获副本 =====")
	for _, f := range funcs {
		f()
	}
	fmt.Println()
}
