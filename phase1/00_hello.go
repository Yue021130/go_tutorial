// 00_hello.go
// 本节目标：理解 Go 的包（package）、导入（import）和可见性规则。
//
// 与 Java 对比：
// - Java 用 package 组织类，import 引入类；Go 用 package 组织文件，import 引入包。
// - Java 的 public/private 是关键字；Go 用首字母大小写控制可见性。
// - Java 一个文件通常一个 public 类；Go 一个文件属于一个 package，可以包含多个类型/函数。

package main

// import 的几种写法：
// 1. 单行导入
// import "fmt"
// 2. 分组导入（推荐，gofmt 会自动整理）
// import (
//     "fmt"
//     "os"
// )
// 3. 别名导入
// import f "fmt"
// 4. 匿名导入（只执行包的 init 函数，不使用包名）
// import _ "net/http/pprof"

import (
	"fmt"
	"os"
	"strings"
)

// 可见性规则（重要）：
// - 首字母大写的标识符：包外可访问（exported，类似 Java public）
// - 首字母小写的标识符：仅包内可访问（unexported，类似 Java private）
//
// 这个规则适用于：函数、变量、常量、类型、结构体字段、方法等。

// Greeting 是 exported 函数（大写开头），可以被其他包调用。
func Greeting(name string) string {
	return fmt.Sprintf("Hello, %s!", name)
}

// secretMessage 是 unexported 函数（小写开头），只能在 main 包内使用。
func secretMessage() string {
	return "这是包内私有函数"
}

func main() {
	// 获取命令行参数，os.Args[0] 是程序本身的路径，os.Args[1:] 是用户传入的参数。
	// 类似 Java 的 public static void main(String[] args)
	name := "World"
	if len(os.Args) > 1 {
		name = strings.Join(os.Args[1:], " ")
	}

	fmt.Println(Greeting(name))
	fmt.Println(secretMessage())

	// 小练习：运行下面的命令观察输出：
	// go run 00_hello.go
	// go run 00_hello.go Go
	// go run 00_hello.go Go 语言
}
