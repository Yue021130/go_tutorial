package main

import (
	"fmt"
	"reflect"
	"time"
	"unsafe"

	"go-tutorial/phase2/internal/reflectdemo"
)

func main() {
	fmt.Println("===== 反射基础：TypeOf/ValueOf =====")
	u := reflectdemo.User{ID: 7, Name: "Go"}
	reflectdemo.Inspect(u)

	fmt.Println("\n===== 反射修改字段 =====")
	pu := &reflectdemo.User{ID: 1, Name: "Original"}
	fmt.Println("before:", pu)
	reflectdemo.ModifyByReflect(pu)
	fmt.Println("after:", pu)

	fmt.Println("\n===== 反射调用方法 =====")
	result, err := reflectdemo.CallMethod(u)
	if err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Println("method result:", result)
	}

	fmt.Println("\n===== reflect.DeepEqual =====")
	reflectdemo.Compare()

	fmt.Println("\n===== 反射性能对比 =====")
	start := time.Now()
	for i := 0; i < 100000; i++ {
		_ = u.Name
	}
	direct := time.Since(start)

	start = time.Now()
	v := reflect.ValueOf(u)
	for i := 0; i < 100000; i++ {
		_ = v.FieldByName("Name").String()
	}
	reflectTime := time.Since(start)

	fmt.Printf("direct access: %v\n", direct)
	if direct > 0 {
		fmt.Printf("reflect access: %v (roughly %d times slower)\n", reflectTime, reflectTime/direct)
	} else {
		fmt.Printf("reflect access: %v (direct access too fast to measure)\n", reflectTime)
	}

	fmt.Println("\n===== unsafe.Sizeof =====")
	fmt.Printf("unsafe sizeof(User)=%d bytes\n", unsafe.Sizeof(u))
	fmt.Println("提示：unsafe 仅用于理解底层机制或与 C 交互，业务代码应优先使用类型安全方案。")
}
