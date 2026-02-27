// 05_arrays_and_slices.go
// 本节目标：深入理解数组和切片，这是 Go 中最容易踩坑的地方之一。
//
// 核心差异：
// - Java 数组是引用类型；Go 数组是值类型。
// - Go 切片是引用类型，底层基于数组，包含 ptr/len/cap 三部分。
// - Go 切片类似 Java ArrayList，但更底层、更灵活。

package main

import "fmt"

func main() {
	// ==================== 1. 数组：固定长度，值类型 ====================
	// Java: int[] arr = new int[3]; 或 int[] arr = {1, 2, 3};
	// Go 数组长度是类型的一部分：[3]int 和 [5]int 是不同类型！

	var arr1 [3]int
	arr1[0] = 10

	arr2 := [3]int{1, 2, 3}
	arr3 := [...]int{1, 2, 3, 4, 5} // 编译器推断长度

	fmt.Printf("arr1=%v, len=%d, type=%T\n", arr1, len(arr1), arr1)
	fmt.Printf("arr2=%v, len=%d, type=%T\n", arr2, len(arr2), arr2)
	fmt.Printf("arr3=%v, len=%d, type=%T\n", arr3, len(arr3), arr3)

	// 数组是值类型，赋值/传参会复制整个数组
	arr4 := arr2
	arr4[0] = 999
	fmt.Printf("arr2=%v, arr4=%v（arr2 未被修改）\n", arr2, arr4)

	// ==================== 2. 切片（Slice）：动态数组，最常用 ====================
	//
	// 切片在内存中的结构（Slice Header）：
	//
	//   切片变量 s
	//   ┌─────────────┐
	//   │ ptr → 底层数组 │
	//   │ len = 3      │
	//   │ cap = 5      │
	//   └─────────────┘
	//
	// 切片是引用类型，赋值/传参只复制 header，不复制底层数据。

	s1 := []int{1, 2, 3}    // 字面量创建
	s2 := make([]int, 3, 5) // 长度3，容量5，元素初始化为零值
	fmt.Printf("s1=%v, len=%d, cap=%d\n", s1, len(s1), cap(s1))
	fmt.Printf("s2=%v, len=%d, cap=%d\n", s2, len(s2), cap(s2))

	// 从数组创建切片
	arr := [5]int{10, 20, 30, 40, 50}
	s3 := arr[1:3] // 左闭右开 [1,3)，包含 arr[1], arr[2]
	fmt.Printf("s3=%v, len=%d, cap=%d\n", s3, len(s3), cap(s3))
	// cap = 4 因为从索引1开始到数组末尾还有 4 个位置

	// ==================== 3. append：切片的动态扩容 ====================
	fmt.Println("\n--- append 与扩容 ---")
	s := make([]int, 0, 2)
	fmt.Printf("初始: %v, len=%d, cap=%d\n", s, len(s), cap(s))

	for i := 1; i <= 5; i++ {
		s = append(s, i)
		fmt.Printf("append %d: %v, len=%d, cap=%d\n", i, s, len(s), cap(s))
	}
	// 你会观察到 cap 不足时翻倍扩容：2 -> 4 -> 8

	// ==================== 4. 切片共享底层数组的陷阱 ====================
	fmt.Println("\n--- 共享底层数组陷阱 ---")
	original := []int{1, 2, 3, 4, 5}
	share := original[1:3] // [2, 3]
	fmt.Printf("修改前: original=%v, share=%v\n", original, share)

	share[0] = 99
	fmt.Printf("修改 share[0] 后: original=%v, share=%v\n", original, share)
	// original 也被修改了！因为 share 和 original 共享底层数组

	// 如何避免？使用 copy 创建独立副本
	copySlice := make([]int, len(original))
	copy(copySlice, original)
	copySlice[0] = 100
	fmt.Printf("copySlice=%v, original=%v（互不影响）\n", copySlice, original)

	// ==================== 5. append 导致重新分配 ====================
	fmt.Println("\n--- append 重新分配 ---")
	a := []int{1, 2, 3, 4, 5}
	b := a[1:3] // [2, 3], cap=4
	fmt.Printf("a=%v, b=%v, b cap=%d\n", a, b, cap(b))

	b = append(b, 999) // 此时 cap 足够，修改的是 a 的底层数组
	fmt.Printf("append 999 后: a=%v, b=%v\n", a, b)

	b = append(b, 888, 777) // cap 不够了，重新分配底层数组
	fmt.Printf("继续 append 后: a=%v, b=%v\n", a, b)
	// 此时修改 b 不再影响 a

	// ==================== 6. nil slice vs empty slice ====================
	fmt.Println("\n--- nil slice vs empty slice ---")
	var nilSlice []int    // nil slice，未初始化
	emptySlice := []int{} // empty slice，已初始化但长度为0
	madeSlice := make([]int, 0)

	fmt.Printf("nilSlice: %v, len=%d, cap=%d, isNil=%v\n", nilSlice, len(nilSlice), cap(nilSlice), nilSlice == nil)
	fmt.Printf("emptySlice: %v, len=%d, cap=%d, isNil=%v\n", emptySlice, len(emptySlice), cap(emptySlice), emptySlice == nil)
	fmt.Printf("madeSlice: %v, len=%d, cap=%d, isNil=%v\n", madeSlice, len(madeSlice), cap(madeSlice), madeSlice == nil)

	// JSON 序列化时，nil slice 会变成 null，empty slice 会变成 []
	// 这是实际开发中需要注意的地方。

	// ==================== 7. 切片操作技巧 ====================
	fmt.Println("\n--- 切片操作 ---")
	data := []int{1, 2, 3, 4, 5}

	// 删除元素（保持顺序）
	idx := 2
	data = append(data[:idx], data[idx+1:]...)
	fmt.Printf("删除索引 %d 后: %v\n", idx, data)

	// 截取前 n 个
	fmt.Printf("前 3 个: %v\n", data[:3])

	// 截取后 n 个
	fmt.Printf("后 2 个: %v\n", data[len(data)-2:])

	// ==================== 常见坑总结 ====================
	// 1. 数组是值类型，大数组传参会复制，性能差。
	// 2. 切片共享底层数组，修改一个可能影响另一个。
	// 3. append 可能触发重新分配，导致原切片"失联"。
	// 4. nil slice 和 empty slice 在 JSON 序列化时表现不同。
	// 5. range 遍历切片时，value 是副本，对值类型修改无效。
}
