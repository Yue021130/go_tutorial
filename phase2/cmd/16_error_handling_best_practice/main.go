package main

import (
	"errors"
	"fmt"

	"go-tutorial/phase2/internal/errdemo"
)

func main() {
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
}
