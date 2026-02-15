package main

import (
	"fmt"
)

func main() {
	// ==================== 1. 变量声明 ====================
	// Java: int count = 10; final String name = "Tom";
	// Go 的变量声明有几种方式，最常用的是短变量声明 :=
	var a int = 10          // 完整声明，类似 Java 的显式类型
	var b = 20              // 类型推断，类似 Java 的 var b = 20;
	c := 30                 // 短变量声明，Go 特有，只能在函数内部使用
	name := "Go学员"         // 字符串类型 string
	isReady := true         // bool 类型

	fmt.Printf("a=%d, b=%d, c=%d, name=%s, isReady=%v\n", a, b, c, name, isReady)

	// 同时声明多个变量
	var x, y int = 1, 2
	m, n := "hello", 3.14
	fmt.Printf("x=%d, y=%d, m=%s, n=%f\n", x, y, m, n)

	// ==================== 2. 常量 ====================
	// Java: public static final int MAX_SIZE = 100;
	const MAX_SIZE = 100
	const (
		Monday    = iota // 0
		Tuesday          // 1
		Wednesday        // 2
	)
	fmt.Printf("MAX_SIZE=%d, Monday=%d, Tuesday=%d, Wednesday=%d\n", MAX_SIZE, Monday, Tuesday, Wednesday)

	// ==================== 3. 基本数据类型 ====================
	// 与 Java 对比：
	// Java 有 byte/short/int/long/float/double/char/boolean
	// Go 有 int/int8/int16/int32/int64, uint/uint8/uint16/uint32/uint64, uintptr
	//     float32/float64, complex64/complex128, byte(=uint8), rune(=int32, 表示 Unicode 码点), bool, string
	var age int32 = 25
	var score float64 = 89.5
	var ch rune = '中' // rune 是 Unicode 码点，类似 Java 的 char（但 Java char 是 UTF-16）
	var flag bool = false
	fmt.Printf("age=%d, score=%.2f, ch=%c, flag=%v\n", age, score, ch, flag)

	// ==================== 4. 零值（重要差异）====================
	// Java 中局部变量必须初始化后才能使用；Go 中声明未赋值的变量会有零值
	var zeroInt int
	var zeroString string
	var zeroBool bool
	var zeroSlice []int // nil
	fmt.Printf("零值：int=%d, string=%q, bool=%v, slice=%v\n", zeroInt, zeroString, zeroBool, zeroSlice)

	// ==================== 5. 控制流 ====================
	// 5.1 if：条件不需要括号（与 Java 不同）
	if score >= 60 {
		fmt.Println("及格")
	} else {
		fmt.Println("不及格")
	}

	// if 支持短变量声明，作用域只在 if 块内（Go 特有）
	if v := score + 10; v >= 90 {
		fmt.Printf("加分后优秀：%.2f\n", v)
	}

	// 5.2 for：Go 只有 for，没有 while/do-while
	// Java: for(int i=0; i<5; i++) { ... }
	for i := 0; i < 5; i++ {
		fmt.Printf("for i=%d ", i)
	}
	fmt.Println()

	// 类似 while 的写法
	j := 0
	for j < 3 {
		fmt.Printf("while-like j=%d ", j)
		j++
	}
	fmt.Println()

	// 无限循环：for { ... }
	// for { ... }

	// 5.3 switch：不需要 break，默认只匹配一个 case（Go 特有，避免 Java 的穿透问题）
	// Java 中每个 case 通常需要 break，否则 fallthrough
	level := "B"
	switch level {
	case "A":
		fmt.Println("优秀")
	case "B":
		fmt.Println("良好")
	case "C":
		fmt.Println("及格")
	default:
		fmt.Println("不及格")
	}

	// switch 也可以不带表达式，类似 if-else if
	switch {
	case score >= 90:
		fmt.Println("A")
	case score >= 80:
		fmt.Println("B")
	default:
		fmt.Println("C")
	}

	// ==================== 6. 字符串与 rune/byte ====================
	s := "Hello, 世界"
	fmt.Printf("字符串长度（字节）=%d\n", len(s)) // 字节长度，UTF-8 编码
	fmt.Printf("字符串长度（rune）=%d\n", len([]rune(s))) // 实际字符数

	// 遍历字符串：按 rune 遍历（推荐）
	for index, r := range s {
		fmt.Printf("index=%d, rune=%c\n", index, r)
	}
}
