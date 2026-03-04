# Phase 1 实战作业：学生成绩管理系统

## 作业要求

实现一个学生成绩管理程序，包含以下功能：

1. **学生信息**
   - `ID`（int64）
   - `Name`（string）
   - `Scores`（map[string]float64，科目→分数）

2. **核心功能**
   - 添加学生
   - 添加/更新某科成绩
   - 计算学生平均分
   - 找出平均分最高的学生
   - 删除学生
   - 查询学生

3. **要求**
   - 使用结构体和方法实现
   - 命令行演示增删改查
   - 编写单元测试（覆盖率越高越好）

## 运行方式

```bash
cd phase1/homework

# 运行命令行演示
go run .

# 运行单元测试
go test -v

# 查看测试覆盖率
go test -cover
```

## 扩展挑战

把这个命令行程序改造成 HTTP 服务：

- `POST /students`：添加学生
- `GET /students/{id}`：查询学生
- `GET /students/{id}/average`：查询平均分
- `GET /students/top`：返回最高分学生
- `DELETE /students/{id}`：删除学生

## 参考答案

本目录下的 `student.go`、`main.go`、`student_test.go` 提供了一份完整参考答案。建议先自己实现，再对照查看。
