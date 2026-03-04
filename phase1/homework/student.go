package main

import (
	"errors"
	"fmt"
)

// Student 表示一个学生
type Student struct {
	ID     int64
	Name   string
	Scores map[string]float64 // 科目 -> 分数
}

// AddScore 添加或更新一门成绩
func (s *Student) AddScore(subject string, score float64) error {
	if subject == "" {
		return errors.New("科目不能为空")
	}
	if score < 0 || score > 100 {
		return errors.New("分数必须在 0-100 之间")
	}
	if s.Scores == nil {
		s.Scores = make(map[string]float64)
	}
	s.Scores[subject] = score
	return nil
}

// Average 计算平均分
func (s *Student) Average() float64 {
	if len(s.Scores) == 0 {
		return 0
	}
	var sum float64
	for _, score := range s.Scores {
		sum += score
	}
	return sum / float64(len(s.Scores))
}

// String 返回学生信息的字符串表示
func (s Student) String() string {
	return fmt.Sprintf("Student{ID=%d, Name=%s, Scores=%v, Average=%.2f}",
		s.ID, s.Name, s.Scores, s.Average())
}

// StudentManager 管理所有学生
type StudentManager struct {
	students map[int64]*Student
	nextID   int64
}

// NewStudentManager 创建一个新的 StudentManager
func NewStudentManager() *StudentManager {
	return &StudentManager{
		students: make(map[int64]*Student),
		nextID:   1,
	}
}

// AddStudent 添加学生，返回学生 ID
func (sm *StudentManager) AddStudent(name string) (int64, error) {
	if name == "" {
		return 0, errors.New("学生姓名不能为空")
	}
	id := sm.nextID
	sm.nextID++
	sm.students[id] = &Student{
		ID:     id,
		Name:   name,
		Scores: make(map[string]float64),
	}
	return id, nil
}

// GetStudent 根据 ID 查询学生
func (sm *StudentManager) GetStudent(id int64) (*Student, error) {
	s, ok := sm.students[id]
	if !ok {
		return nil, fmt.Errorf("学生 ID=%d 不存在", id)
	}
	return s, nil
}

// DeleteStudent 删除学生
func (sm *StudentManager) DeleteStudent(id int64) error {
	if _, ok := sm.students[id]; !ok {
		return fmt.Errorf("学生 ID=%d 不存在", id)
	}
	delete(sm.students, id)
	return nil
}

// AddScore 给学生添加成绩
func (sm *StudentManager) AddScore(id int64, subject string, score float64) error {
	s, err := sm.GetStudent(id)
	if err != nil {
		return err
	}
	return s.AddScore(subject, score)
}

// GetAverage 获取学生平均分
func (sm *StudentManager) GetAverage(id int64) (float64, error) {
	s, err := sm.GetStudent(id)
	if err != nil {
		return 0, err
	}
	return s.Average(), nil
}

// GetTopStudent 返回平均分最高的学生
// 如果有多个学生并列最高，返回 ID 最小的那个
func (sm *StudentManager) GetTopStudent() (*Student, error) {
	if len(sm.students) == 0 {
		return nil, errors.New("没有学生")
	}

	var top *Student
	var topID int64
	for id, s := range sm.students {
		if top == nil || s.Average() > top.Average() ||
			(s.Average() == top.Average() && id < topID) {
			top = s
			topID = id
		}
	}
	return top, nil
}

// ListStudents 返回所有学生列表
func (sm *StudentManager) ListStudents() []*Student {
	list := make([]*Student, 0, len(sm.students))
	for _, s := range sm.students {
		list = append(list, s)
	}
	return list
}
