package main

import (
	"fmt"

	"go-tutorial/phase2/internal/contracts"
)

func main() {
	// 值接收者的“修改副本”问题：
	vc := contracts.ValueCounter{}
	contracts.BumpN(vc, 3)
	fmt.Printf("ValueCounter after BumpN(vc,3) = %d (期望仍是 0)\n", vc.Value())

	// 指针接收者能修改原对象状态：
	pc := &contracts.PointerCounter{}
	contracts.BumpN(pc, 3)
	fmt.Printf("PointerCounter after BumpN(pc,3) = %d (期望是 3)\n", pc.Value())

	// 注意：contracts.PointerCounter（非指针）不满足 Counter 接口，无法传给 BumpN。
}
