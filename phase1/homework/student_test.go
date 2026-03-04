package main

import (
	"math"
	"testing"
)

// TestStudent_AddScore 测试添加成绩
func TestStudent_AddScore(t *testing.T) {
	s := &Student{ID: 1, Name: "Test"}

	if err := s.AddScore("语文", 85); err != nil {
		t.Errorf("添加成绩失败: %v", err)
	}

	if len(s.Scores) != 1 {
		t.Errorf("期望 1 门成绩，实际 %d", len(s.Scores))
	}

	// 测试无效分数
	if err := s.AddScore("数学", 101); err == nil {
		t.Error("期望分数无效错误，但没有返回错误")
	}

	// 测试空科目
	if err := s.AddScore("", 80); err == nil {
		t.Error("期望科目为空错误，但没有返回错误")
	}
}

// TestStudent_Average 测试平均分计算
func TestStudent_Average(t *testing.T) {
	tests := []struct {
		name   string
		scores map[string]float64
		want   float64
	}{
		{
			name:   "无成绩",
			scores: map[string]float64{},
			want:   0,
		},
		{
			name:   "单科",
			scores: map[string]float64{"语文": 90},
			want:   90,
		},
		{
			name:   "多科",
			scores: map[string]float64{"语文": 80, "数学": 90, "英语": 100},
			want:   90,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Student{Scores: tt.scores}
			got := s.Average()
			if math.Abs(got-tt.want) > 1e-9 {
				t.Errorf("Average() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestStudentManager 测试学生管理器
func TestStudentManager(t *testing.T) {
	manager := NewStudentManager()

	// 测试添加学生
	id1, err := manager.AddStudent("张三")
	if err != nil {
		t.Fatalf("添加学生失败: %v", err)
	}
	if id1 != 1 {
		t.Errorf("期望第一个学生 ID=1，实际 %d", id1)
	}

	// 测试空姓名
	if _, err := manager.AddStudent(""); err == nil {
		t.Error("期望空姓名错误")
	}

	// 测试添加成绩
	if err := manager.AddScore(id1, "语文", 80); err != nil {
		t.Errorf("添加成绩失败: %v", err)
	}

	// 测试查询
	s, err := manager.GetStudent(id1)
	if err != nil {
		t.Errorf("查询学生失败: %v", err)
	}
	if s.Name != "张三" {
		t.Errorf("期望姓名=张三，实际 %s", s.Name)
	}

	// 测试查询不存在的学生
	if _, err := manager.GetStudent(999); err == nil {
		t.Error("期望查询不存在学生返回错误")
	}

	// 测试平均分
	manager.AddScore(id1, "数学", 90)
	avg, err := manager.GetAverage(id1)
	if err != nil {
		t.Errorf("获取平均分失败: %v", err)
	}
	if math.Abs(avg-85) > 1e-9 {
		t.Errorf("期望平均分=85，实际 %.2f", avg)
	}

	// 测试最高分学生
	id2, _ := manager.AddStudent("李四")
	manager.AddScore(id2, "语文", 95)
	manager.AddScore(id2, "数学", 95)

	top, err := manager.GetTopStudent()
	if err != nil {
		t.Errorf("获取最高分学生失败: %v", err)
	}
	if top.ID != id2 {
		t.Errorf("期望最高分学生 ID=%d，实际 %d", id2, top.ID)
	}

	// 测试删除
	if err := manager.DeleteStudent(id1); err != nil {
		t.Errorf("删除学生失败: %v", err)
	}
	if _, err := manager.GetStudent(id1); err == nil {
		t.Error("期望删除后查询返回错误")
	}

	// 测试空管理器获取最高分
	emptyManager := NewStudentManager()
	if _, err := emptyManager.GetTopStudent(); err == nil {
		t.Error("期望空管理器返回错误")
	}
}
