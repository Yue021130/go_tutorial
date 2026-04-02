// Package domain 定义领域模型。
//
// 与 Java 对比：
// - Java: Entity / POJO / DTO
// - Go: struct + tag，通常不区分那么细
package domain

import "time"

// User 用户领域模型
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"uniqueIndex;size:64;not null"`
	Password  string    `json:"-" gorm:"size:128;not null"` // json:"-" 表示不序列化
	Email     string    `json:"email" gorm:"size:128"`
	Role      string    `json:"role" gorm:"size:32;default:user"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserResponse 是返回给客户端的用户信息
type UserResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=64"`
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email" binding:"email"`
	Role     string `json:"role"`
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	Email string `json:"email" binding:"email"`
	Role  string `json:"role"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}
