package cancellation

import (
	"context"
	"fmt"
	"time"
)

// RunJobs 在 ctx 取消时停止处理后续任务。
func RunJobs(ctx context.Context, jobs []int, perJob time.Duration) []int {
	out := make([]int, 0, len(jobs))
	for _, job := range jobs {
		select {
		case <-ctx.Done():
			return out
		default:
		}

		time.Sleep(perJob)
		out = append(out, job*job)
	}
	return out
}

// TickUntilCancel 周期发出计数，直到 ctx 取消。
func TickUntilCancel(ctx context.Context, interval time.Duration) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		count := 0
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				count++
				out <- count
			}
		}
	}()
	return out
}

// ==================== context.WithValue ====================
//
// WithValue 用于在 context 中传递请求范围的元数据。
// 最佳实践：
//   - 只放真正跨层传递的数据（如 request_id、user_id）
//   - key 使用私有类型，避免冲突
//   - 不要把可选参数或业务配置放进 context

type contextKey string

const RequestIDKey contextKey = "request_id"
const UserIDKey contextKey = "user_id"

func GetRequestID(ctx context.Context) string {
	if id, ok := ctx.Value(RequestIDKey).(string); ok {
		return id
	}
	return ""
}

func GetUserID(ctx context.Context) int {
	if id, ok := ctx.Value(UserIDKey).(int); ok {
		return id
	}
	return 0
}

// ProcessRequest 演示 context value 的传递
func ProcessRequest(ctx context.Context) string {
	return fmt.Sprintf("processing request_id=%s user_id=%d",
		GetRequestID(ctx), GetUserID(ctx))
}

// ==================== context.Deadline 与 Timeout ====================
//
// WithTimeout = WithDeadline(parent, time.Now().Add(timeout))
// Deadline() 返回具体的截止时间。

func CheckDeadline(ctx context.Context) {
	if deadline, ok := ctx.Deadline(); ok {
		fmt.Printf("deadline: %s\n", deadline.Format("15:04:05.000"))
	} else {
		fmt.Println("no deadline")
	}
}

// ==================== 级联取消 ====================
//
// 父 context 取消时，所有子 context 都会收到取消信号。
// 子 context 取消不会影响父 context。

func CascadedCancellation(parent context.Context) []string {
	child1, cancel1 := context.WithCancel(parent)
	defer cancel1()

	child2, cancel2 := context.WithTimeout(parent, 200*time.Millisecond)
	defer cancel2()

	var results []string

	go func() {
		<-child1.Done()
		results = append(results, "child1 canceled")
	}()

	go func() {
		<-child2.Done()
		results = append(results, "child2 canceled")
	}()

	// 取消父 context，child1 和 child2 都会收到信号
	// 注意：这里只是演示，实际调用方不应该在这里取消
	return results
}
