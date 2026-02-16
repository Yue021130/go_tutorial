package main

import (
	"errors"
	"fmt"
)

// ==================== 1. 函数声明 ====================
// Java 方法：public int add(int a, int b) { return a + b; }
// Go 函数：func add(a int, b int) int { return a + b }
// 连续同类型参数可简写：func add(a, b int) int
func add(a, b int) int {
	return a + b
}

// ==================== 2. 多返回值（Go 标志性特性）====================
// Java 一个方法只能返回一个值；Go 可以返回多个值
// 常见模式：result, err := doSomething()
func divide(a, b float64) (float64, error) {
	if b == 0 {
		// errors.New 创建简单错误，类似 Java 的 new RuntimeException("msg")
		return 0, errors.New("除数不能为零")
	}
	return a / b, nil
}

// 命名返回值：result 和 err 在函数签名中声明，可直接 return
func rectangle(width, height float64) (area float64, perimeter float64) {
	area = width * height
	perimeter = 2 * (width + height)
	return // 裸 return，返回命名返回值（小函数可用，大函数不建议）
}

// ==================== 3. 可变参数 ====================
// Java: public int sum(int... nums) { ... }
// Go: func sum(nums ...int) int
func sum(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

// ==================== 4. defer：延迟执行 ====================
// 类似 Java 的 try-finally，但更简洁，常用于资源释放
// defer 语句在函数返回前按 LIFO（后进先出）顺序执行
func deferDemo() {
	fmt.Println("deferDemo 开始")
	defer fmt.Println("defer 1：最后执行")
	defer fmt.Println("defer 2：中间执行")
	defer fmt.Println("defer 3：最先执行")
	fmt.Println("deferDemo 结束")
}

// defer 常用于关闭文件、释放锁、记录耗时等
func deferUseCase() {
	fmt.Println("开始处理")
	defer fmt.Println("清理资源（无论是否出错都会执行）")
	fmt.Println("业务处理中...")
}

// ==================== 5. panic 与 recover ====================
// Java：throw 异常；调用栈向上传播，直到被 try-catch 捕获
// Go：panic 会中断当前函数执行，逐层向上传播，直到被 recover 捕获或程序崩溃
// Go 推荐：普通错误用 error 返回值处理；不可恢复的严重问题才用 panic
func mayPanic() {
	defer func() {
		// recover 必须在 defer 函数中调用，用于捕获 panic
		if r := recover(); r != nil {
			fmt.Printf("捕获到 panic: %v\n", r)
		}
	}()
	fmt.Println("即将 panic")
	panic("出现了严重错误") // 类似 Java 的 throw new RuntimeException("...")
	fmt.Println("这行不会执行")
}

// ==================== 6. 函数作为一等公民 ====================
// Go 函数可以赋值给变量、作为参数、作为返回值（类似 Java 的 Function/Predicate）
func apply(a, b int, op func(int, int) int) int {
	return op(a, b)
}

func main() {
	fmt.Println("===== 函数基础 =====")
	fmt.Printf("add(3, 5) = %d\n", add(3, 5))

	fmt.Println("\n===== 多返回值 =====")
	// Java 要返回多个值，通常封装成对象或数组；Go 直接返回多个值
	result, err := divide(10, 3)
	if err != nil {
		fmt.Printf("错误：%v\n", err)
	} else {
		fmt.Printf("10 / 3 = %.2f\n", result)
	}

	if result, err := divide(10, 0); err != nil {
		fmt.Printf("10 / 0 错误：%v\n", err)
	} else {
		fmt.Printf("10 / 0 = %.2f\n", result)
	}

	area, perimeter := rectangle(4, 5)
	fmt.Printf("矩形面积=%.2f, 周长=%.2f\n", area, perimeter)

	fmt.Println("\n===== 可变参数 =====")
	fmt.Printf("sum(1,2,3)=%d\n", sum(1, 2, 3))
	fmt.Printf("sum()=%d\n", sum())

	fmt.Println("\n===== defer 示例 =====")
	deferDemo()
	deferUseCase()

	fmt.Println("\n===== panic/recover 示例 =====")
	mayPanic()
	fmt.Println("程序继续执行（panic 已被 recover 捕获）")

	fmt.Println("\n===== 函数作为参数 =====")
	mul := func(a, b int) int { return a * b }
	fmt.Printf("apply(4, 5, mul) = %d\n", apply(4, 5, mul))
}
