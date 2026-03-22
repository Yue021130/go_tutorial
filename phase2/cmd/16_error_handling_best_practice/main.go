package main

import (
	"errors"
	"fmt"

	"go-tutorial/phase2/internal/errdemo"
)

func main() {
	fmt.Println("===== 错误包装与分类 =====")
	for _, id := range []int{-1, 3, 1} {
		name, err := errdemo.LoadUserDisplayName(id)
		if err != nil {
			fmt.Printf("id=%d err=%v\n", id, err)

			switch {
			case errdemo.IsInvalidID(err):
				fmt.Println("  -> classified as invalid id")
			case errdemo.IsNotFound(err):
				fmt.Println("  -> classified as not found")
			default:
				var unwrapped error = errors.Unwrap(err)
				fmt.Println("  -> unwrapped one level:", unwrapped)
			}
			continue
		}
		fmt.Printf("id=%d name=%s\n", id, name)
	}

	fmt.Println("\n===== errors.Join 多错误聚合 =====")
	if err := errdemo.ValidateUser(-1, "", 200); err != nil {
		fmt.Println("validation errors:")
		// errors.Join 返回的 error 实现了 Unwrap() []error
		type unwrapErr interface {
			Unwrap() []error
		}
		if ue, ok := err.(unwrapErr); ok {
			for _, e := range ue.Unwrap() {
				fmt.Println("  -", e)
			}
		}
	}

	fmt.Println("\n===== panic vs error 决策 =====")
	result, err := errdemo.SafeDivide(10, 0)
	if err != nil {
		fmt.Println("safe divide error:", err)
	} else {
		fmt.Println("safe divide result:", result)
	}

	fmt.Println("\n===== 业务错误分类 =====")
	be := errdemo.BusinessError{Code: "USER_NOT_FOUND", Message: "user not found"}
	fmt.Println("business error:", be.Error())
	fmt.Println("is business error:", errdemo.IsBusinessError(be))
}
