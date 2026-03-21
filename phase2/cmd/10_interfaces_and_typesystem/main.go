package main

import (
	"context"
	"fmt"
	"time"

	"go-tutorial/phase2/internal/contracts"
)

func notifyAndAudit(ctx context.Context, a contracts.Alerting, msg string) {
	a.Log(ctx, "about to notify: "+msg)
	if err := a.Notify(ctx, msg); err != nil {
		fmt.Println("notify error:", err)
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	fmt.Println("===== 隐式实现 =====")
	// 隐式实现：EmailNotifier 没有 implements 关键字，也能赋值给接口。
	sender := contracts.MessageSender{
		Notifier: contracts.EmailNotifier{Prefix: "[Phase2] "},
	}
	_ = sender.Send(ctx, "hello interface")

	fmt.Println("\n===== 接口组合 =====")
	// 接口组合：MemoryNotifier 同时实现 Notify 和 Log。
	mem := &contracts.MemoryNotifier{}
	notifyAndAudit(ctx, mem, "alert message")
	fmt.Println("memory messages:", mem.Messages)

	fmt.Println("\n===== any + 类型断言 + 类型切换 =====")
	// any + 类型断言 + 类型切换。
	values := []any{42, "100", true, 3.14, nil}
	for _, v := range values {
		fmt.Printf("describe=%s", contracts.DescribeValue(v))
		if n, ok := contracts.AnyToInt(v); ok {
			fmt.Printf(", as int=%d", n)
		}
		fmt.Println()
	}

	fmt.Println("\n===== nil interface 陷阱 =====")
	// typed nil 与 nil interface 的区别
	contracts.DescribeNilInterface()

	fmt.Println("\n===== 消费者端定义接口 =====")
	// 小接口哲学：Process 只依赖 Reader 接口
	r := contracts.NewStringReader("hello go interface")
	contracts.Process(r)
}
