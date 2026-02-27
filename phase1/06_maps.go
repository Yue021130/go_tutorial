// 06_maps.go
// 本节目标：深入理解 map，包括 nil map、遍历、删除、线程安全等。
//
// 与 Java 对比：
// - Java HashMap 是引用类型，new HashMap<>() 后可直接使用。
// - Go map 也是引用类型，但 nil map 不能直接写入，必须用 make 初始化。
// - Go map 的遍历顺序是不确定的（即使同样的代码多次运行顺序也可能不同）。

package main

import (
	"fmt"
	"sort"
)

func main() {
	// ==================== 1. map 的创建 ====================
	// 方式 1：字面量（推荐，会自动初始化）
	scores := map[string]int{
		"Go":     95,
		"Java":   90,
		"Python": 88,
	}
	fmt.Printf("scores=%v\n", scores)

	// 方式 2：make
	ages := make(map[string]int)
	ages["Alice"] = 25
	ages["Bob"] = 30
	fmt.Printf("ages=%v\n", ages)

	// 方式 3：声明但不初始化（nil map）
	var nilMap map[string]int
	fmt.Printf("nilMap=%v, isNil=%v\n", nilMap, nilMap == nil)

	// nilMap["key"] = 1 // 运行时 panic：assignment to entry in nil map
	// 读取 nil map 中的 key 不会 panic，返回零值
	fmt.Printf("nilMap[\"key\"]=%d\n", nilMap["key"])

	// ==================== 2. map 的读写 ====================
	// 读取 map 时，可以返回两个值：value 和 ok（key 是否存在）
	if score, ok := scores["Go"]; ok {
		fmt.Printf("Go 的分数是 %d\n", score)
	} else {
		fmt.Println("Go 不存在")
	}

	// 不存在的 key 返回零值，这可能和 key 存在但值为零值混淆
	scores["Rust"] = 0
	fmt.Printf("Rust 的分数是 %d\n", scores["Rust"])

	// 判断 key 是否存在必须用两值模式
	if _, ok := scores["C++"]; !ok {
		fmt.Println("C++ 不存在（虽然 scores[\"C++\"] 返回 0）")
	}

	// ==================== 3. 删除 key ====================
	delete(scores, "Python")
	fmt.Printf("删除 Python 后: %v\n", scores)

	// 删除不存在的 key 不会报错
	delete(scores, "NotExist")

	// ==================== 4. map 的遍历顺序是不确定的 ====================
	fmt.Println("\n--- map 遍历顺序 ---")
	for i := 0; i < 3; i++ {
		fmt.Printf("第 %d 次遍历: ", i+1)
		for k, v := range scores {
			fmt.Printf("%s=%d ", k, v)
		}
		fmt.Println()
	}

	// 如果需要固定顺序，需要先取出 keys，排序后再遍历
	fmt.Println("\n--- 按 key 排序遍历 ---")
	keys := make([]string, 0, len(scores))
	for k := range scores {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Printf("%s=%d ", k, scores[k])
	}
	fmt.Println()

	// ==================== 5. map 的长度 ====================
	fmt.Printf("\nscores 长度=%d\n", len(scores))

	// ==================== 6. map 的 key 类型要求 ====================
	// map 的 key 必须是可比较类型（comparable）：
	// 允许：bool、数字、string、指针、channel、interface、包含这些类型的数组/结构体
	// 不允许：slice、map、function

	// 合法 key 示例
	type Point struct {
		X, Y int
	}
	pointMap := make(map[Point]string)
	pointMap[Point{1, 2}] = "A"
	fmt.Printf("pointMap=%v\n", pointMap)

	// 非法 key 示例（编译错误）：
	// badMap := make(map[[]int]string)

	// ==================== 7. map 不是线程安全的 ====================
	// 多个 goroutine 同时读写 map 会导致运行时 panic。
	// 并发安全方案（预告 Phase 2）：
	// - sync.RWMutex + map
	// - sync.Map

	// ==================== 常见坑 ====================
	// 坑 1：nil map 写入会 panic
	// 坑 2：map 取值单值模式无法区分"不存在"和"值为零值"
	// 坑 3：map 遍历顺序随机，不要依赖顺序
	// 坑 4：map 作为函数参数是引用传递，函数内修改会影响外部
	// 坑 5：map 不能并发读写
}
