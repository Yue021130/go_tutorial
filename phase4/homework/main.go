// Program homework 是 Phase 4 实战作业的可运行入口。
//
// 运行方式：
//   echo -e "INFO hello\nERROR fail" | go run ./homework
package main

import (
	"fmt"
	"os"

	"go-tutorial/phase4/homework/logprocessor"
)

func main() {
	stats := logprocessor.ProcessLogs(os.Stdin, 4)
	fmt.Print(logprocessor.FormatStats(stats))
}
