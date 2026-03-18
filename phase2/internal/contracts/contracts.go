package contracts

import (
	"context"
	"errors"
	"fmt"
	"strconv"
)

var ErrNilNotifier = errors.New("nil notifier")

// Notifier 演示接口的隐式实现。
type Notifier interface {
	Notify(ctx context.Context, msg string) error
}

// AuditLogger 演示接口组合。
type AuditLogger interface {
	Log(ctx context.Context, msg string)
}

// Alerting = Notifier + AuditLogger
type Alerting interface {
	Notifier
	AuditLogger
}

type MessageSender struct {
	Notifier Notifier
}

func (s MessageSender) Send(ctx context.Context, msg string) error {
	if s.Notifier == nil {
		return ErrNilNotifier
	}
	return s.Notifier.Notify(ctx, msg)
}

type EmailNotifier struct {
	Prefix string
}

func (e EmailNotifier) Notify(ctx context.Context, msg string) error {
	if msg == "" {
		return errors.New("empty message")
	}

	select {
	case <-ctx.Done():
		return fmt.Errorf("email notify canceled: %w", ctx.Err())
	default:
	}

	fmt.Printf("%s[Email] %s\n", e.Prefix, msg)
	return nil
}

// MemoryNotifier 用指针接收者演示“状态更新需要指针接收者”。
type MemoryNotifier struct {
	Messages []string
}

func (m *MemoryNotifier) Notify(ctx context.Context, msg string) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("memory notify canceled: %w", ctx.Err())
	default:
	}
	m.Messages = append(m.Messages, msg)
	return nil
}

func (m *MemoryNotifier) Log(ctx context.Context, msg string) {
	_ = ctx
	m.Messages = append(m.Messages, "[LOG] "+msg)
}

// AnyToInt 演示 any + 类型断言/转换。
func AnyToInt(v any) (int, bool) {
	switch x := v.(type) {
	case int:
		return x, true
	case int8:
		return int(x), true
	case int16:
		return int(x), true
	case int32:
		return int(x), true
	case int64:
		return int(x), true
	case float32:
		return int(x), true
	case float64:
		return int(x), true
	case string:
		n, err := strconv.Atoi(x)
		if err != nil {
			return 0, false
		}
		return n, true
	default:
		return 0, false
	}
}

func DescribeValue(v any) string {
	switch x := v.(type) {
	case nil:
		return "nil"
	case string:
		return "string:" + x
	case int:
		return fmt.Sprintf("int:%d", x)
	case bool:
		return fmt.Sprintf("bool:%t", x)
	default:
		return fmt.Sprintf("unknown:%T", x)
	}
}

// Counter 用于方法集与接收者演示。
type Counter interface {
	Inc()
	Value() int
}

// ValueCounter 的 Inc 是值接收者，修改的是副本。
type ValueCounter struct {
	v int
}

func (c ValueCounter) Inc() {
	c.v++
}

func (c ValueCounter) Value() int {
	return c.v
}

type PointerCounter struct {
	v int
}

func (c *PointerCounter) Inc() {
	c.v++
}

func (c *PointerCounter) Value() int {
	return c.v
}

func BumpN(c Counter, n int) {
	for i := 0; i < n; i++ {
		c.Inc()
	}
}
