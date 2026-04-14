package main

import (
	"context"
	"fmt"
)

func main() {
	ctx := context.Background()
	repo := NewInMemoryOrderRepository()
	svc := NewOrderService(repo)

	fmt.Println("===== 创建订单 =====")
	order, err := svc.CreateOrder(ctx, 1001, 299.99, []string{"iPhone 15", "充电器"})
	if err != nil {
		fmt.Println("创建订单失败:", err)
		return
	}
	fmt.Printf("创建订单成功: %+v\n", order)

	fmt.Println("\n===== 支付订单 =====")
	if err := svc.PayOrder(ctx, order.ID); err != nil {
		fmt.Println("支付失败:", err)
		return
	}
	fmt.Println("支付成功")

	fmt.Println("\n===== 查询订单 =====")
	found, err := svc.GetOrder(ctx, order.ID)
	if err != nil {
		fmt.Println("查询失败:", err)
		return
	}
	fmt.Printf("订单状态: %s\n", found.Status)

	fmt.Println("\n===== 非法状态转换 =====")
	if err := svc.CancelOrder(ctx, order.ID); err != nil {
		fmt.Println("取消失败（预期）:", err)
	}

	fmt.Println("\n===== 查询用户订单列表 =====")
	orders, err := svc.ListUserOrders(ctx, 1001)
	if err != nil {
		fmt.Println("查询失败:", err)
		return
	}
	fmt.Printf("用户 1001 共有 %d 个订单\n", len(orders))
}
