package main

import "fmt"

func main() {
	// 创建管理器
	manager := NewStudentManager()

	// 添加学生
	id1, _ := manager.AddStudent("张三")
	id2, _ := manager.AddStudent("李四")
	id3, _ := manager.AddStudent("王五")

	// 添加成绩
	manager.AddScore(id1, "语文", 85)
	manager.AddScore(id1, "数学", 90)
	manager.AddScore(id1, "英语", 78)

	manager.AddScore(id2, "语文", 92)
	manager.AddScore(id2, "数学", 88)

	manager.AddScore(id3, "语文", 70)
	manager.AddScore(id3, "数学", 75)
	manager.AddScore(id3, "英语", 80)
	manager.AddScore(id3, "物理", 85)

	// 查询所有学生
	fmt.Println("===== 所有学生 =====")
	for _, s := range manager.ListStudents() {
		fmt.Println(s)
	}

	// 查询平均分
	fmt.Println("\n===== 平均分查询 =====")
	avg, _ := manager.GetAverage(id1)
	fmt.Printf("张三平均分: %.2f\n", avg)

	// 最高分学生
	fmt.Println("\n===== 最高分学生 =====")
	top, _ := manager.GetTopStudent()
	fmt.Println(top)

	// 更新成绩
	fmt.Println("\n===== 更新成绩后 =====")
	manager.AddScore(id3, "数学", 95)
	top, _ = manager.GetTopStudent()
	fmt.Println(top)

	// 删除学生
	fmt.Println("\n===== 删除学生后 =====")
	manager.DeleteStudent(id2)
	for _, s := range manager.ListStudents() {
		fmt.Println(s)
	}

	// 错误处理示例
	fmt.Println("\n===== 错误处理 =====")
	if _, err := manager.GetStudent(999); err != nil {
		fmt.Printf("查询失败: %v\n", err)
	}
}
