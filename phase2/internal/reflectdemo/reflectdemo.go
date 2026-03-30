// Package reflectdemo 演示反射的基本用法和边界。
//
// 反射三定律（The Laws of Reflection）：
// 1. Reflection goes from interface value to reflection object.
// 2. Reflection goes from reflection object to interface value.
// 3. To modify a reflection object, the value must be settable.
package reflectdemo

import (
	"fmt"
	"reflect"
)

type User struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func (u User) Hello() string {
	return "hello, " + u.Name
}

// Inspect 使用 reflect.Type 和 reflect.Value 检查结构体
func Inspect(u User) {
	t := reflect.TypeOf(u)
	v := reflect.ValueOf(u)

	fmt.Println("type name:", t.Name())
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		fmt.Printf("field=%s kind=%s tag(json)=%q value=%v\n",
			f.Name, f.Type.Kind(), f.Tag.Get("json"), v.Field(i).Interface())
	}
}

// ModifyByReflect 通过反射修改结构体字段
// 注意：必须传入指针，并且 reflect.ValueOf 后调用 Elem() 获取可设置值
func ModifyByReflect(u *User) {
	v := reflect.ValueOf(u)
	if v.Kind() != reflect.Pointer || v.IsNil() {
		return
	}
	v = v.Elem() // 解引用指针

	if !v.CanSet() {
		fmt.Println("value is not settable")
		return
	}

	nameField := v.FieldByName("Name")
	if nameField.IsValid() && nameField.CanSet() && nameField.Kind() == reflect.String {
		nameField.SetString("ModifiedByReflect")
	}

	idField := v.FieldByName("ID")
	if idField.IsValid() && idField.CanSet() && idField.Kind() == reflect.Int64 {
		idField.SetInt(999)
	}
}

// CallMethod 通过反射调用方法
func CallMethod(u User) (string, error) {
	v := reflect.ValueOf(u)
	m := v.MethodByName("Hello")
	if !m.IsValid() {
		return "", fmt.Errorf("method Hello not found")
	}

	results := m.Call(nil)
	if len(results) > 0 {
		return results[0].String(), nil
	}
	return "", nil
}

// Compare 使用 reflect.DeepEqual 比较复杂结构
func Compare() {
	a := map[string][]int{"x": {1, 2}}
	b := map[string][]int{"x": {1, 2}}
	c := map[string][]int{"x": {1, 3}}

	fmt.Println("a == b ?", reflect.DeepEqual(a, b))
	fmt.Println("a == c ?", reflect.DeepEqual(a, c))
}
