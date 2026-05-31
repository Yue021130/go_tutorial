// Package logprocessor 提供 Phase 4 实战作业参考答案：并发日志处理器。
//
// 需求：
//   1. 读取一个日志文件（每行一条日志）。
//   2. 使用 Worker Pool 并发解析每行日志，提取级别（INFO/WARN/ERROR）和消息。
//   3. 使用 Pipeline 统计各级别日志数量。
//   4. 输出统计结果。
package logprocessor

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"sync"
)

// LogEntry 表示一条解析后的日志。
type LogEntry struct {
	Level   string
	Message string
}

// LogStats 保存统计结果。
type LogStats struct {
	mu     sync.RWMutex
	counts map[string]int
}

// NewLogStats 创建统计器。
func NewLogStats() *LogStats {
	return &LogStats{counts: make(map[string]int)}
}

// Add 增加某级别的计数。
func (s *LogStats) Add(level string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.counts[level]++
}

// Counts 返回统计结果副本。
func (s *LogStats) Counts() map[string]int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	copy := make(map[string]int, len(s.counts))
	for k, v := range s.counts {
		copy[k] = v
	}
	return copy
}

// ParseLogLine 解析单行日志，格式假设为 "LEVEL message"。
func ParseLogLine(line string) (LogEntry, bool) {
	parts := strings.SplitN(strings.TrimSpace(line), " ", 2)
	if len(parts) < 2 {
		return LogEntry{}, false
	}
	return LogEntry{Level: strings.ToUpper(parts[0]), Message: parts[1]}, true
}

// ProcessLogs 使用 Worker Pool 并发解析日志并统计。
// workers 控制并发度。
func ProcessLogs(r io.Reader, workers int) *LogStats {
	lines := make(chan string, 100)
	entries := make(chan LogEntry, 100)
	stats := NewLogStats()

	var wg sync.WaitGroup

	// 启动 worker 解析日志
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			for line := range lines {
				if entry, ok := ParseLogLine(line); ok {
					entries <- entry
				}
			}
		}()
	}

	// 启动一个 goroutine 关闭 entries channel
	go func() {
		wg.Wait()
		close(entries)
	}()

	// 启动统计 goroutine
	var statWg sync.WaitGroup
	statWg.Add(1)
	go func() {
		defer statWg.Done()
		for entry := range entries {
			stats.Add(entry.Level)
		}
	}()

	// 读取输入并发送给 worker
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lines <- scanner.Text()
	}
	close(lines)

	statWg.Wait()
	return stats
}

// FormatStats 格式化输出统计结果。
func FormatStats(stats *LogStats) string {
	counts := stats.Counts()
	var sb strings.Builder
	fmt.Fprintln(&sb, "日志统计结果:")
	for level, count := range counts {
		fmt.Fprintf(&sb, "  %s: %d\n", level, count)
	}
	return sb.String()
}
