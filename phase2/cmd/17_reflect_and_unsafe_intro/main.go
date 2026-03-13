package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

type User struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func main() {
	u := User{ID: 7, Name: "Go"}
	t := reflect.TypeOf(u)
	v := reflect.ValueOf(u)

	fmt.Println("type name:", t.Name())
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		fmt.Printf("field=%s kind=%s tag(json)=%q value=%v\n",
			f.Name, f.Type.Kind(), f.Tag.Get("json"), v.Field(i).Interface())
	}

	fmt.Printf("unsafe sizeof(User)=%d bytes\n", unsafe.Sizeof(u))
	fmt.Println("提示：unsafe 仅用于理解底层机制，业务代码应优先使用类型安全方案。")
}
