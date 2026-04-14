package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOrderService_CreateOrder(t *testing.T) {
	repo := NewInMemoryOrderRepository()
	svc := NewOrderService(repo)
	ctx := context.Background()

	order, err := svc.CreateOrder(ctx, 1001, 99.99, []string{"book"})
	require.NoError(t, err)
	assert.NotZero(t, order.ID)
	assert.Equal(t, OrderStatusPending, order.Status)
}

func TestOrderService_CreateOrder_InvalidInput(t *testing.T) {
	repo := NewInMemoryOrderRepository()
	svc := NewOrderService(repo)
	ctx := context.Background()

	_, err := svc.CreateOrder(ctx, 0, 99.99, []string{"book"})
	assert.Error(t, err)

	_, err = svc.CreateOrder(ctx, 1001, -10, []string{"book"})
	assert.Error(t, err)

	_, err = svc.CreateOrder(ctx, 1001, 99.99, nil)
	assert.Error(t, err)
}

func TestOrderService_PayAndCancel(t *testing.T) {
	repo := NewInMemoryOrderRepository()
	svc := NewOrderService(repo)
	ctx := context.Background()

	order, err := svc.CreateOrder(ctx, 1001, 99.99, []string{"book"})
	require.NoError(t, err)

	err = svc.PayOrder(ctx, order.ID)
	require.NoError(t, err)

	paid, err := svc.GetOrder(ctx, order.ID)
	require.NoError(t, err)
	assert.Equal(t, OrderStatusPaid, paid.Status)

	err = svc.CancelOrder(ctx, order.ID)
	require.NoError(t, err)

	cancelled, err := svc.GetOrder(ctx, order.ID)
	require.NoError(t, err)
	assert.Equal(t, OrderStatusCancelled, cancelled.Status)
}

func TestOrderService_CancelShippedOrder(t *testing.T) {
	repo := NewInMemoryOrderRepository()
	svc := NewOrderService(repo)
	ctx := context.Background()

	order, err := svc.CreateOrder(ctx, 1001, 99.99, []string{"book"})
	require.NoError(t, err)

	err = repo.UpdateStatus(ctx, order.ID, OrderStatusPaid)
	require.NoError(t, err)
	err = repo.UpdateStatus(ctx, order.ID, OrderStatusShipped)
	require.NoError(t, err)

	err = svc.CancelOrder(ctx, order.ID)
	assert.Error(t, err)
}

func TestOrderService_ListUserOrders(t *testing.T) {
	repo := NewInMemoryOrderRepository()
	svc := NewOrderService(repo)
	ctx := context.Background()

	_, err := svc.CreateOrder(ctx, 1001, 10, []string{"a"})
	require.NoError(t, err)
	_, err = svc.CreateOrder(ctx, 1001, 20, []string{"b"})
	require.NoError(t, err)
	_, err = svc.CreateOrder(ctx, 1002, 30, []string{"c"})
	require.NoError(t, err)

	orders, err := svc.ListUserOrders(ctx, 1001)
	require.NoError(t, err)
	assert.Len(t, orders, 2)
}

func TestCanTransition(t *testing.T) {
	assert.True(t, canTransition(OrderStatusPending, OrderStatusPaid))
	assert.True(t, canTransition(OrderStatusPending, OrderStatusCancelled))
	assert.False(t, canTransition(OrderStatusPaid, OrderStatusCompleted))
	assert.False(t, canTransition(OrderStatusCompleted, OrderStatusCancelled))
}
