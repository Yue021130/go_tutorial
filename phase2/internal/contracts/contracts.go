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

// ==================== 方法集规则演示 ====================
//
// 对于类型 T：
//   T  的方法集 = 所有值接收者方法
//   *T 的方法集 = 所有值接收者方法 + 所有指针接收者方法
//
// 这意味着：
//   - *T 可以赋值给任何需要 T 方法集的接口
//   - T 只能赋值给只需要值接收者方法的接口
//   - 如果接口包含指针接收者方法，T 不能满足

// Mover 接口包含一个值接收者方法
type Mover interface {
	Move()
}

// Jumper 接口包含一个指针接收者方法
type Jumper interface {
	Jump()
}

// Athlete 同时有值接收者和指针接收者方法
type Athlete struct {
	position int
}

func (a Athlete) Move() {
	a.position++ // 修改副本，不会生效
}

func (a *Athlete) Jump() {
	a.position += 10
}

// CheckMethodSets 演示方法集规则
func CheckMethodSets() {
	// 值类型 Athlete 满足 Mover（因为 Move 是值接收者）
	var m Mover = Athlete{}
	_ = m

	// 值类型 Athlete 不满足 Jumper（因为 Jump 是指针接收者）
	// var j Jumper = Athlete{} // 编译错误

	// 指针类型 *Athlete 同时满足 Mover 和 Jumper
	var j Jumper = &Athlete{}
	var m2 Mover = &Athlete{}
	_ = j
	_ = m2

	fmt.Println("方法集检查通过")
}

// ==================== addressable 与非 addressable ====================
//
// 可寻址（addressable）的值：变量、切片元素、指针解引用、可寻址字段
// 不可寻址的值：字面量、map 值、函数返回值、类型转换结果
//
// 当调用指针接收者方法时，Go 会自动对可寻址值取地址。
// 不可寻址值不能自动取地址，因此不能调用指针接收者方法。

func AddressableDemo() {
	a := Athlete{}
	a.Jump() // 变量可寻址，自动取地址：(&a).Jump()
	fmt.Printf("变量调用 Jump 后 position=%d\n", a.position)

	// Athlete{}.Jump() // 编译错误：字面量不可寻址

	// map 值不可寻址
	m := map[string]Athlete{"tom": {}}
	// m["tom"].Jump() // 编译错误
	// 正确做法：取出来修改再放回去
	v := m["tom"]
	v.Jump()
	m["tom"] = v

	// 切片元素可寻址
	s := []Athlete{{}}
	s[0].Jump()
	fmt.Printf("切片元素调用 Jump 后 position=%d\n", s[0].position)
}

// ==================== nil interface 陷阱 ====================
//
// interface 内部可以抽象为 (type, data) 两部分：
//   - type: 动态类型信息
//   - data: 指向具体值的指针
//
// 当 interface 的 type 和 data 都为 nil 时，interface 才等于 nil。
// 当一个具体类型的指针是 nil，但被赋值给 interface 时，interface 的 type 不为 nil，
// 此时 interface != nil，但调用其方法可能会 panic（typed nil）。

// Repository 用于演示 nil interface
type Repository interface {
	Find(id int) (string, error)
}

type UserRepository struct{}

func (r *UserRepository) Find(id int) (string, error) {
	if r == nil {
		return "", errors.New("repository is nil")
	}
	return fmt.Sprintf("user:%d", id), nil
}

func DescribeNilInterface() {
	var r1 Repository      // type=nil, data=nil → nil interface
	var r2 *UserRepository // typed nil：type=*UserRepository, data=nil
	var r3 Repository = r2 // interface 的 type 不为 nil！

	fmt.Printf("r1 == nil ? %v\n", r1 == nil)
	fmt.Printf("r2 == nil ? %v\n", r2 == nil)
	fmt.Printf("r3 == nil ? %v (typed nil 陷阱)\n", r3 == nil)

	// 安全做法：在方法内部检查 typed nil
	if _, err := r3.Find(1); err != nil {
		fmt.Println("r3.Find error:", err)
	}
}

// ==================== 接口设计哲学 ====================
//
// Go 推荐"消费者端定义接口"：
//   - Java：接口通常由实现方定义（如 JDBC 接口由 JDK 定义）
//   - Go：接口应该由使用方定义，越小越好
//
// 经典例子：io.Reader 只有一个方法 Read，任何能读的类型都自动满足。
// 这被称为"小接口"哲学，也是 Go 接口灵活性的来源。

// Reader 是一个消费者端定义的小接口示例
type Reader interface {
	Read(p []byte) (n int, err error)
}

// StringReader 是 Reader 的一个简单实现
type StringReader struct {
	Data string
	pos  int
}

func NewStringReader(data string) *StringReader {
	return &StringReader{Data: data}
}

func (s *StringReader) Read(p []byte) (n int, err error) {
	if s.pos >= len(s.Data) {
		return 0, errors.New("EOF")
	}
	n = copy(p, s.Data[s.pos:])
	s.pos += n
	return n, nil
}

// Process 只依赖 Reader 接口，不依赖具体实现
func Process(r Reader) {
	buf := make([]byte, 8)
	n, _ := r.Read(buf)
	fmt.Printf("Process read: %q\n", string(buf[:n]))
}
