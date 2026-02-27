// 04_defer_panic_recover.go
// 本节目标：深入理解 defer、panic、recover。
//
// 与 Java 对比：
// - Java: try-catch-finally，异常类继承 Throwable
// - Go:   普通错误用 error 返回值；严重错误用 panic + recover
// - defer 类似 finally，但更简洁，且 LIFO 执行

package main

import (
	"fmt"
)

// ==================== 1. defer 基础 ====================
// defer 语句在函数返回前执行，常用于关闭资源、解锁、记录日志等。
// 执行顺序：LIFO（后进先出），类似栈。
func deferOrderDemo() {
	fmt.Println("函数开始")
	defer fmt.Println("defer 1")
	defer fmt.Println("defer 2")
	defer fmt.Println("defer 3")
	fmt.Println("函数结束")
	// 输出顺序：
	// 函数开始
	// 函数结束
	// defer 3
	// defer 2
	// defer 1
}

// ==================== 2. defer 的参数在 defer 语句处求值 ====================
// 不是返回时才求值！这是常见坑。
func deferArgEvalDemo() {
	i := 0
	defer fmt.Printf("defer 中 i = %d\n", i) // 此时 i=0，后续修改不影响
	i = 100
	fmt.Printf("函数中 i = %d\n", i)
	// 输出：
	// 函数中 i = 100
	// defer 中 i = 0
}

// 如果希望 defer 看到最终值，用指针或闭包
func deferArgEvalFix() {
	i := 0
	defer func() {
		fmt.Printf("defer 闭包中 i = %d\n", i) // 闭包捕获变量，返回时 i=100
	}()
	i = 100
	fmt.Printf("函数中 i = %d\n", i)
}

// ==================== 3. defer 与 return 的交互（经典面试题）====================
// 如果函数使用命名返回值，defer 可以修改返回值。
func deferWithNamedReturn() (result int) {
	defer func() {
		result++ // 修改命名返回值
	}()
	return 10 // 实际返回 11
}

// 非命名返回值，defer 无法改变
func deferWithoutNamedReturn() int {
	result := 10
	defer func() {
		result++ // 修改的是局部变量，不影响返回值
	}()
	return result // 返回 10
}

// ==================== 4. panic 与 recover ====================
// panic 会中断当前函数执行，向上传播，直到被 recover 捕获或程序崩溃。
//
// panic 传播路径示意：
//   main -> funcA -> funcB -> panic
//   如果 funcB 没有 recover，则继续向上到 funcA
//   如果 funcA 有 recover，则捕获 panic，funcA 之后的代码继续执行
//   main 也会继续执行

func innerFunc() {
	fmt.Println("innerFunc: 即将 panic")
	panic("出现了严重错误") // 类似 Java throw new RuntimeException("...")
	fmt.Println("innerFunc: 这行不会执行")
}

func middleFunc() {
	fmt.Println("middleFunc: 调用 innerFunc")
	innerFunc()
	fmt.Println("middleFunc: 这行不会执行（因为 innerFunc panic 了）")
}

func outerFunc() {
	defer func() {
		// recover 必须在 defer 函数中调用
		if r := recover(); r != nil {
			fmt.Printf("outerFunc: 捕获到 panic: %v\n", r)
		}
	}()
	fmt.Println("outerFunc: 调用 middleFunc")
	middleFunc()
	fmt.Println("outerFunc: 继续执行")
}

// ==================== 5. panic/recover 使用原则 ====================
// - 普通可预期错误：用 error 返回值处理
// - 不可恢复的严重问题：用 panic
// - 库代码不要吞掉 panic，除非你确实能恢复
// - HTTP 服务等顶层代码通常会用 recover 防止单个请求拖垮整个进程

func safeDivide(a, b float64) (result float64, err error) {
	defer func() {
		if r := recover(); r != nil {
			// 捕获 panic 并转换为 error
			err = fmt.Errorf("除法 panic: %v", r)
		}
	}()
	if b == 0 {
		panic("除数为零") // 模拟某个深层调用 panic
	}
	return a / b, nil
}

func main() {
	fmt.Println("===== defer 执行顺序 =====")
	deferOrderDemo()

	fmt.Println("\n===== defer 参数求值时机 =====")
	deferArgEvalDemo()
	deferArgEvalFix()

	fmt.Println("\n===== defer 修改返回值 =====")
	fmt.Printf("命名返回值: %d\n", deferWithNamedReturn())
	fmt.Printf("非命名返回值: %d\n", deferWithoutNamedReturn())

	fmt.Println("\n===== panic/recover =====")
	outerFunc()
	fmt.Println("main: 程序继续执行")

	fmt.Println("\n===== panic 转 error =====")
	result, err := safeDivide(10, 0)
	if err != nil {
		fmt.Printf("错误: %v\n", err)
	} else {
		fmt.Printf("结果: %f\n", result)
	}
}
