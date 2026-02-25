// 02_control_flow.go
// 本节目标：掌握 Go 的控制流：if、for、switch、range。
//
// 关键认知：Go 没有 while、do-while，只有 for；
// Go 的 switch 默认不穿透（不需要 break），这是为了减少 Java 中常见的漏写 break 的 bug。

package main

import "fmt"

func main() {
	// ==================== 1. if 语句 ====================
	// 与 Java 的差异：
	// 1. 条件表达式不需要括号
	// 2. 可以在 if 中声明短作用域变量

	score := 85.0
	if score >= 60 {
		fmt.Println("及格")
	} else {
		fmt.Println("不及格")
	}

	// if 短变量声明：v 的作用域只在 if-else 块内
	// 类似 Java 中 if 块内用 {} 限制局部变量，但 Go 更简洁
	if v := score + 10; v >= 90 {
		fmt.Printf("加分后优秀：%.2f\n", v)
	} else {
		fmt.Printf("加分后未达优秀：%.2f\n", v)
	}
	// fmt.Println(v) // 编译错误：undefined: v

	// ==================== 2. for 语句的四种形式 ====================

	// 形式 1：C-style（最常用）
	// Java: for (int i = 0; i < 5; i++) { ... }
	fmt.Println("\n--- for 形式 1: C-style ---")
	for i := 0; i < 5; i++ {
		fmt.Printf("i=%d ", i)
	}
	fmt.Println()

	// 形式 2：while-like（只有条件表达式）
	fmt.Println("\n--- for 形式 2: while-like ---")
	j := 0
	for j < 3 {
		fmt.Printf("j=%d ", j)
		j++
	}
	fmt.Println()

	// 形式 3：无限循环
	fmt.Println("\n--- for 形式 3: 无限循环（带 break）---")
	k := 0
	for {
		if k >= 3 {
			break
		}
		fmt.Printf("k=%d ", k)
		k++
	}
	fmt.Println()

	// 形式 4：range 遍历（Go 特有，类似 Java 的 for-each + Iterator）
	fmt.Println("\n--- for 形式 4: range ---")
	nums := []int{10, 20, 30}
	for index, value := range nums {
		fmt.Printf("index=%d, value=%d\n", index, value)
	}

	// 如果不需要 index，用 _ 占位
	for _, value := range nums {
		fmt.Printf("value=%d ", value)
	}
	fmt.Println()

	// range 遍历 map（注意：map 遍历顺序是不确定的！）
	scores := map[string]int{"Go": 95, "Java": 90, "Python": 88}
	for subject, score := range scores {
		fmt.Printf("%s=%d ", subject, score)
	}
	fmt.Println()

	// range 遍历字符串：按 rune 遍历（正确处理中文）
	str := "Go语言"
	for index, r := range str {
		fmt.Printf("index=%d, rune=%c\n", index, r)
	}

	// ==================== 3. switch 语句 ====================
	// Go 的 switch 默认不穿透，不需要写 break。
	// 这是 Go 设计者基于经验做的选择：Java 中 90% 的 case 都需要 break，
	// 漏写 break 导致的 bug 很常见。

	fmt.Println("\n--- switch 基础 ---")
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

	// 如果需要穿透，显式使用 fallthrough
	fmt.Println("\n--- switch + fallthrough ---")
	n := 1
	switch n {
	case 1:
		fmt.Println("一")
		fallthrough
	case 2:
		fmt.Println("二")
		fallthrough
	case 3:
		fmt.Println("三")
	}

	// switch 可以没有表达式，变成 if-else if 的清晰写法
	fmt.Println("\n--- switch 无表达式 ---")
	switch {
	case score >= 90:
		fmt.Println("A")
	case score >= 80:
		fmt.Println("B")
	case score >= 60:
		fmt.Println("C")
	default:
		fmt.Println("D")
	}

	// case 可以是多个值
	fmt.Println("\n--- switch case 多值 ---")
	day := "Saturday"
	switch day {
	case "Saturday", "Sunday":
		fmt.Println("周末")
	case "Monday", "Tuesday", "Wednesday", "Thursday", "Friday":
		fmt.Println("工作日")
	default:
		fmt.Println("未知")
	}

	// ==================== 4. 跳转语句 ====================
	// break：跳出当前循环或 switch
	// continue：跳过当前迭代
	// goto：可以跳转到标签（不推荐，容易写出面条代码）
	fmt.Println("\n--- break 与 continue ---")
	for i := 0; i < 10; i++ {
		if i%2 == 0 {
			continue
		}
		if i > 7 {
			break
		}
		fmt.Printf("%d ", i)
	}
	fmt.Println()

	// ==================== 常见坑 ====================
	// 坑 1：在 for-range 中修改切片元素，修改的是副本（对值类型无效）
	arr := []int{1, 2, 3}
	for _, v := range arr {
		v = v * 10 // 修改的是副本，原切片不变
	}
	fmt.Printf("\nrange 中修改副本后 arr=%v\n", arr)

	// 正确做法：通过索引修改
	for i := range arr {
		arr[i] = arr[i] * 10
	}
	fmt.Printf("通过索引修改后 arr=%v\n", arr)

	// 坑 2：for-range 的 index/value 变量在每次迭代中被复用
	// 这个坑在闭包中尤其常见，后面函数章节会再讲。
}
