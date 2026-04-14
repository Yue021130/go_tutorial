// Package main 是 Phase 3 实战作业：订单管理模块。
//
// 本作业演示如何在一个简化场景中实现：
// - domain 模型
// - repository 数据访问层
// - service 业务逻辑层
// - 单元测试
package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

// OrderStatus 订单状态
type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusPaid      OrderStatus = "paid"
	OrderStatusShipped   OrderStatus = "shipped"
	OrderStatusCompleted OrderStatus = "completed"
	OrderStatusCancelled OrderStatus = "cancelled"
)

// Order 订单领域模型
type Order struct {
	ID        uint64
	UserID    uint64
	Amount    float64
	Status    OrderStatus
	Items     []string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// OrderRepository 定义订单数据访问接口
type OrderRepository interface {
	Create(ctx context.Context, order *Order) error
	GetByID(ctx context.Context, id uint64) (*Order, error)
	ListByUser(ctx context.Context, userID uint64) ([]Order, error)
	UpdateStatus(ctx context.Context, id uint64, status OrderStatus) error
}

// InMemoryOrderRepository 是内存实现，便于测试
type InMemoryOrderRepository struct {
	mu     sync.RWMutex
	orders map[uint64]*Order
	nextID uint64
}

func NewInMemoryOrderRepository() *InMemoryOrderRepository {
	return &InMemoryOrderRepository{
		orders: make(map[uint64]*Order),
		nextID: 1,
	}
}

func (r *InMemoryOrderRepository) Create(ctx context.Context, order *Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	order.ID = r.nextID
	r.nextID++
	order.Status = OrderStatusPending
	order.CreatedAt = time.Now()
	order.UpdatedAt = order.CreatedAt

	r.orders[order.ID] = order
	return nil
}

func (r *InMemoryOrderRepository) GetByID(ctx context.Context, id uint64) (*Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	order, ok := r.orders[id]
	if !ok {
		return nil, fmt.Errorf("order %d not found", id)
	}
	return order, nil
}

func (r *InMemoryOrderRepository) ListByUser(ctx context.Context, userID uint64) ([]Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []Order
	for _, order := range r.orders {
		if order.UserID == userID {
			result = append(result, *order)
		}
	}
	return result, nil
}

func (r *InMemoryOrderRepository) UpdateStatus(ctx context.Context, id uint64, status OrderStatus) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	order, ok := r.orders[id]
	if !ok {
		return fmt.Errorf("order %d not found", id)
	}

	if !canTransition(order.Status, status) {
		return fmt.Errorf("invalid status transition from %s to %s", order.Status, status)
	}

	order.Status = status
	order.UpdatedAt = time.Now()
	return nil
}

// canTransition 检查状态转换是否合法
func canTransition(from, to OrderStatus) bool {
	transitions := map[OrderStatus][]OrderStatus{
		OrderStatusPending:   {OrderStatusPaid, OrderStatusCancelled},
		OrderStatusPaid:      {OrderStatusShipped, OrderStatusCancelled},
		OrderStatusShipped:   {OrderStatusCompleted},
		OrderStatusCompleted: {},
		OrderStatusCancelled: {},
	}

	allowed, ok := transitions[from]
	if !ok {
		return false
	}
	for _, s := range allowed {
		if s == to {
			return true
		}
	}
	return false
}

// OrderService 订单业务逻辑
type OrderService struct {
	repo OrderRepository
}

func NewOrderService(repo OrderRepository) *OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) CreateOrder(ctx context.Context, userID uint64, amount float64, items []string) (*Order, error) {
	if userID == 0 {
		return nil, errors.New("user id is required")
	}
	if amount <= 0 {
		return nil, errors.New("amount must be positive")
	}
	if len(items) == 0 {
		return nil, errors.New("items cannot be empty")
	}

	order := &Order{
		UserID: userID,
		Amount: amount,
		Items:  items,
	}

	if err := s.repo.Create(ctx, order); err != nil {
		return nil, err
	}
	return order, nil
}

func (s *OrderService) GetOrder(ctx context.Context, id uint64) (*Order, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *OrderService) ListUserOrders(ctx context.Context, userID uint64) ([]Order, error) {
	return s.repo.ListByUser(ctx, userID)
}

func (s *OrderService) PayOrder(ctx context.Context, id uint64) error {
	return s.repo.UpdateStatus(ctx, id, OrderStatusPaid)
}

func (s *OrderService) CancelOrder(ctx context.Context, id uint64) error {
	order, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if order.Status == OrderStatusShipped || order.Status == OrderStatusCompleted {
		return errors.New("cannot cancel shipped or completed order")
	}
	return s.repo.UpdateStatus(ctx, id, OrderStatusCancelled)
}
