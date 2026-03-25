package errdemo

import (
	"errors"
	"fmt"
	"strings"
)

var ErrInvalidUserID = errors.New("invalid user id")

type NotFoundError struct {
	Resource string
	ID       int
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("%s id=%d not found", e.Resource, e.ID)
}

func ValidateUserID(id int) error {
	if id <= 0 {
		return fmt.Errorf("validate user id %d: %w", id, ErrInvalidUserID)
	}
	return nil
}

func FindUserName(id int) (string, error) {
	if err := ValidateUserID(id); err != nil {
		return "", err
	}

	switch id {
	case 1:
		return "alice", nil
	case 2:
		return "bob", nil
	default:
		return "", fmt.Errorf("find user: %w", NotFoundError{
			Resource: "user",
			ID:       id,
		})
	}
}

func LoadUserDisplayName(id int) (string, error) {
	name, err := FindUserName(id)
	if err != nil {
		return "", fmt.Errorf("load user display name: %w", err)
	}
	return strings.ToUpper(name), nil
}

func IsInvalidID(err error) bool {
	return errors.Is(err, ErrInvalidUserID)
}

func IsNotFound(err error) bool {
	var nf NotFoundError
	return errors.As(err, &nf)
}

// ==================== errors.Join（Go 1.20+）====================
//
// 用于聚合多个独立错误，常用于批量操作或校验多个字段。
// 与 fmt.Errorf("%w; %w", err1, err2) 不同，Join 保留每个错误的独立性。

func ValidateUser(userID int, name string, age int) error {
	var errs []error
	if userID <= 0 {
		errs = append(errs, ErrInvalidUserID)
	}
	if name == "" {
		errs = append(errs, errors.New("name cannot be empty"))
	}
	if age < 0 || age > 150 {
		errs = append(errs, errors.New("age out of range"))
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

// ==================== 错误分类与决策框架 ====================
//
// panic vs error 决策：
//   - 可预期、可恢复的问题 → error
//   - 编程错误/不可恢复的严重问题 → panic（通常只在 main 或顶层 recover）
//
// 错误分类：
//   - 用户输入错误：400 Bad Request
//   - 资源不存在：404 Not Found
//   - 权限错误：403 Forbidden
//   - 系统/第三方错误：500 Internal Server Error，可能需要重试

// BusinessError 用于区分业务错误和系统错误
type BusinessError struct {
	Code    string
	Message string
}

func (e BusinessError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func IsBusinessError(err error) bool {
	var be BusinessError
	return errors.As(err, &be)
}

// ==================== 常见反模式 ====================
//
// 1. 错误吞掉：if err != nil { return nil } // 丢失了错误信息
// 2. 重复包装：fmt.Errorf("layer1: %w", fmt.Errorf("layer2: %w", err)) // 链太长
// 3. 该返回 error 的地方用 panic
// 4. 在底层打印日志而不是包装后返回
// 5. 用字符串比较判断错误类型，而不是 errors.Is/errors.As

// SafeDivide 演示 panic 转 error 的边界
func SafeDivide(a, b int) (int, error) {
	defer func() {
		if r := recover(); r != nil {
			// 这里只是演示，实际不应该依赖 panic 来处理可预期错误
		}
	}()
	if b == 0 {
		return 0, errors.New("divide by zero")
	}
	return a / b, nil
}
