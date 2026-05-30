// Program 31_worker_pool 演示 Worker Pool 并发模式。
//
// Worker Pool 模式：
//   - 固定数量的 worker goroutine 从任务队列中取出任务执行。
//   - 控制并发度，避免无限制创建 goroutine 导致资源耗尽。
//   - 适合批量处理任务（如图片处理、批量 HTTP 请求、消息消费）。
//
// 结构示意：
//
//        +--------+     +--------+     +--------+     +--------+
//   任务 |  Job 1 | --> |  Job 2 | --> |  Job 3 | --> |  Job 4 |  ...
//        +--------+     +--------+     +--------+     +--------+
//             |               |               |               |
//             v               v               v               v
//        +--------------------------------------------------------+
//        |                    Job Channel（任务队列）               |
//        +--------------------------------------------------------+
//             |               |               |
//             v               v               v
//        +---------+    +---------+    +---------+
//        | Worker  |    | Worker  |    | Worker  |
//        |    1    |    |    2    |    |    3    |
//        +---------+    +---------+    +---------+
//
// 与 Java ThreadPoolExecutor 对比：
//   - Java 的线程池核心参数：corePoolSize、maximumPoolSize、keepAliveTime、workQueue；
//   - Go 的 Worker Pool 通常自己实现：固定 worker 数 + channel 任务队列 + sync.WaitGroup。
//   - Go 没有内置 ThreadPoolExecutor，但社区有 panjf2000/ants 等成熟库。
package main

import (
	"fmt"
	"sync"
	"time"
)

// Job 表示一个任务。
type Job struct {
	ID int
}

// Result 表示任务执行结果。
type Result struct {
	JobID    int
	WorkerID int
	Output   int
}

// worker 从 jobs channel 取任务执行，结果写入 results channel。
func worker(id int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		// 模拟耗时工作
		time.Sleep(100 * time.Millisecond)
		results <- Result{JobID: job.ID, WorkerID: id, Output: job.ID * 10}
	}
}

func main() {
	fmt.Println("=== Worker Pool 模式演示 ===")

	const numJobs = 20
	const numWorkers = 4

	jobs := make(chan Job, numJobs)
	results := make(chan Result, numJobs)

	var wg sync.WaitGroup
	wg.Add(numWorkers)
	for w := 1; w <= numWorkers; w++ {
		go worker(w, jobs, results, &wg)
	}

	// 发送任务
	for j := 1; j <= numJobs; j++ {
		jobs <- Job{ID: j}
	}
	close(jobs)

	// 等待所有 worker 完成，然后关闭 results channel
	go func() {
		wg.Wait()
		close(results)
	}()

	// 收集结果
	for r := range results {
		fmt.Printf("Worker %d 处理 Job %d，结果: %d\n", r.WorkerID, r.JobID, r.Output)
	}

	fmt.Println("=== 演示结束 ===")
}
