package safety

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type SafeCounter struct {
	mu     sync.RWMutex
	counts map[string]int
}

func NewSafeCounter() *SafeCounter {
	return &SafeCounter{
		counts: make(map[string]int),
	}
}

func (c *SafeCounter) Inc(key string) {
	c.mu.Lock()
	c.counts[key]++
	c.mu.Unlock()
}

func (c *SafeCounter) Get(key string) int {
	c.mu.RLock()
	v := c.counts[key]
	c.mu.RUnlock()
	return v
}

func (c *SafeCounter) Snapshot() map[string]int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	out := make(map[string]int, len(c.counts))
	for k, v := range c.counts {
		out[k] = v
	}
	return out
}

func ParallelIncrements(c *SafeCounter, key string, goroutines int, each int) {
	var wg sync.WaitGroup
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < each; j++ {
				c.Inc(key)
			}
		}()
	}
	wg.Wait()
}

type AtomicCounter struct {
	v atomic.Int64
}

func (c *AtomicCounter) Inc() {
	c.v.Add(1)
}

func (c *AtomicCounter) Value() int64 {
	return c.v.Load()
}

// ==================== sync.Mutex 独立示例 ====================
//
// Mutex 用于保护临界区，确保同一时间只有一个 goroutine 访问。
// 注意：不要复制 Mutex。

type MutexCounter struct {
	mu    sync.Mutex
	value int
}

func (c *MutexCounter) Inc() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value++
}

func (c *MutexCounter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

func MutexParallelIncrements(c *MutexCounter, goroutines int, each int) {
	var wg sync.WaitGroup
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < each; j++ {
				c.Inc()
			}
		}()
	}
	wg.Wait()
}

// ==================== sync.Once ====================
//
// Once 保证某个函数只执行一次，常用于单例初始化。
// 与 Java 的 synchronized + double-checked locking 相比，Once 更简洁。

type Singleton struct {
	Value string
}

var (
	instance *Singleton
	once     sync.Once
)

func GetSingleton() *Singleton {
	once.Do(func() {
		instance = &Singleton{Value: "initialized"}
	})
	return instance
}

func (s *Singleton) GetValue() string {
	return s.Value
}

// ==================== sync.Pool ====================
//
// Pool 用于缓存临时对象，减少 GC 压力。
// 注意：Pool 中的对象可能被 GC 随时回收，不能假设对象一定存在。

var bufferPool = sync.Pool{
	New: func() any {
		return make([]byte, 1024)
	},
}

func UseBufferPool() {
	buf := bufferPool.Get().([]byte)
	copy(buf, "hello pool")
	fmt.Println("from pool:", string(buf[:11]))
	bufferPool.Put(buf)
}

// ==================== sync.Map ====================
//
// sync.Map 是并发安全的 map，适合以下场景：
//   - 读多写少
//   - key 类型单一
//   - 需要并发安全的简单缓存
//
// 普通 map + Mutex 通常性能更好，sync.Map 是特定场景的优化。

func SyncMapDemo() {
	var m sync.Map
	m.Store("key1", 100)
	m.Store("key2", 200)

	if v, ok := m.Load("key1"); ok {
		fmt.Println("key1 =", v)
	}

	m.Range(func(key, value any) bool {
		fmt.Printf("sync.Map: %s -> %v\n", key, value)
		return true // 继续遍历
	})
}

// ==================== 并发安全 map 方案对比 ====================
//
// 1. map + sync.RWMutex：通用方案，适合大多数场景
// 2. sync.Map：读多写少、key 类型单一
// 3. map + channel：单 goroutine owner，完全串行化

// ==================== race detector 说明 ====================
//
// 使用 `go test -race` 或 `go run -race` 检测数据竞争。
// 数据竞争：多个 goroutine 同时访问同一内存位置，且至少一个是写操作。

func init() {
	// 避免 fmt 未使用警告
	_ = fmt.Sprintf("")
}
