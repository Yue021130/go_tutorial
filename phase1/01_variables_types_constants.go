// 01_variables_types_constants.go
// 本节目标：彻底掌握 Go 的变量、常量、数据类型、类型转换和零值。

package main

import (
	"fmt"
	"reflect"
)

func main() {
	// ==================== 1. 变量声明的四种方式 ====================
	//
	// Java 写法：
	//   int a = 10;
	//   var b = 20;          // Java 10+
	//   final int MAX = 100;
	//
	// Go 写法：

	// 方式 1：完整声明（显式类型）
	var a int = 10

	// 方式 2：类型推断
	var b = 20 // 编译器推断为 int

	// 方式 3：短变量声明（最常用，只能在函数内部使用）
	c := 30

	// 方式 4：声明多个变量
	var x, y int = 1, 2
	m, n := "hello", 3.14

	fmt.Printf("a=%d, b=%d, c=%d, x=%d, y=%d, m=%s, n=%f\n", a, b, c, x, y, m, n)

	// 重要：:= 只能在函数内部使用，包级变量必须用 var
	// 原因：包级变量需要明确的声明边界，:= 是语句不是声明。

	// ==================== 2. 常量与 iota ====================
	// Java: public static final int MAX_SIZE = 100;
	// Go: const MAX_SIZE = 100
	const MAX_SIZE = 100

	// iota 是常量声明中的计数器，从 0 开始递增
	// 对比 Java enum：更底层、更灵活，但也更容易写出难懂的代码
	const (
		Monday    = iota // 0
		Tuesday          // 1
		Wednesday        // 2
		Thursday         // 3
		Friday           // 4
		Saturday         // 5
		Sunday           // 6
	)
	fmt.Printf("Monday=%d, Sunday=%d\n", Monday, Sunday)

	// iota 经典用法：位掩码
	const (
		Read  = 1 << iota // 1 << 0 = 1
		Write             // 1 << 1 = 2
		Exec              // 1 << 2 = 4
	)
	fmt.Printf("Read=%d, Write=%d, Exec=%d\n", Read, Write, Exec)

	// ==================== 3. 数据类型对比表 ====================
	//
	// | Java 类型       | Go 类型                    | 说明                       |
	// |----------------|---------------------------|---------------------------|
	// | byte           | byte (= uint8)            | 0-255                     |
	// | short          | int16                     |                           |
	// | int            | int32 / int               | int 在 64 位系统为 64 位    |
	// | long           | int64                     |                           |
	// | float          | float32                   |                           |
	// | double         | float64                   |                           |
	// | char           | rune (= int32)            | Unicode 码点               |
	// | boolean        | bool                      |                           |
	// | String         | string                    | 不可变，UTF-8 编码          |
	// | Object         | interface{}               | 空接口，任意类型            |

	var age int32 = 25
	var score float64 = 89.5
	var ch rune = '中' // rune 是 Unicode 码点，类似 Java char 但 Java char 是 UTF-16
	var flag bool = false

	fmt.Printf("age=%d, score=%.2f, ch=%c, flag=%v\n", age, score, ch, flag)
	fmt.Printf("ch 的 Unicode 码点 = %d\n", ch)

	// ==================== 4. 类型转换：Go 是强类型，必须显式转换 ====================
	// Java 存在自动类型提升：int + long = long
	// Go 不允许隐式类型转换，必须显式写类型转换
	var i int = 10
	var f float64 = float64(i) // 必须显式转换
	var u uint = uint(i)
	fmt.Printf("i=%d, f=%f, u=%d\n", i, f, u)

	// 注意：浮点转整数会截断小数部分，不是四舍五入
	var pi float64 = 3.99
	fmt.Printf("int(pi) = %d\n", int(pi)) // 输出 3

	// 字符串与数字的转换需要 strconv 包，不是强制类型转换
	// s := string(65) // 得到 "A"，不是 "65"
	// n, _ := strconv.Atoi("123") // 得到 123

	// ==================== 5. 零值（Zero Value）====================
	// Java 中局部变量必须显式初始化后才能使用，否则编译报错。
	// Go 中声明变量但未赋值时，会自动赋予"零值"。
	//
	// 类型          零值
	// int/float    0 / 0.0
	// string       ""（空字符串）
	// bool         false
	// 指针         nil
	// slice/map/channel/func/interface  nil

	var zeroInt int
	var zeroFloat float64
	var zeroString string
	var zeroBool bool
	var zeroPointer *int
	var zeroSlice []int
	var zeroMap map[string]int
	var zeroInterface interface{}

	fmt.Printf("int=%d, float=%f, string=%q, bool=%v\n", zeroInt, zeroFloat, zeroString, zeroBool)
	fmt.Printf("pointer=%v, slice=%v, map=%v, interface=%v\n", zeroPointer, zeroSlice, zeroMap, zeroInterface)

	// 判断引用类型是否为 nil
	fmt.Printf("zeroSlice == nil ? %v\n", zeroSlice == nil)
	fmt.Printf("zeroMap == nil ? %v\n", zeroMap == nil)

	// ==================== 6. 字符串不可变性 ====================
	// 与 Java String 类似，Go 的 string 也是不可变的。
	// 修改字符串需要转换为 []byte 或 []rune。
	str := "Hello"
	// str[0] = 'h' // 编译错误：cannot assign to str[0]
	bytes := []byte(str)
	bytes[0] = 'h'
	fmt.Printf("原字符串=%s, 修改后=%s\n", str, string(bytes))

	// ==================== 7. reflect 查看类型 ====================
	// 类似 Java 的 obj.getClass().getName()
	fmt.Printf("a 的类型是 %s\n", reflect.TypeOf(a))
	fmt.Printf("n 的类型是 %s\n", reflect.TypeOf(n))

	// ==================== 常见坑 ====================
	// 坑 1：短变量声明 := 在相同作用域内不能对已有变量重新声明
	// a := 10 // 错误：no new variables on left side of :=
	// 但如果左边至少有一个新变量，就可以：
	a, z := 100, 200
	fmt.Printf("a=%d, z=%d\n", a, z)

	// 坑 2：:= 会遮蔽（shadow）外层同名变量
	if true {
		a := 999 // 这里的 a 是新的局部变量，和外层的 a 不同
		fmt.Printf("内层 a=%d\n", a)
	}
	fmt.Printf("外层 a=%d\n", a)

	// 坑 3：Go 的 int 大小依赖平台，跨平台代码建议用 int64/int32
	fmt.Printf("本机 int 占 %d 字节\n", reflect.TypeOf(0).Size())
}
