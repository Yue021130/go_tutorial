// Package service 是业务逻辑层。
//
// 与 Java 对比：
// - Java: @Service
// - Go: 普通结构体 + 接口
package service

import (
	"context"
	"errors"
	"fmt"

	"go-tutorial/phase3/internal/domain"
	"go-tutorial/phase3/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// UserService 定义用户业务逻辑接口
type UserService interface {
	Register(ctx context.Context, req domain.CreateUserRequest) (*domain.UserResponse, error)
	Login(ctx context.Context, req domain.LoginRequest) (*domain.User, error)
	GetByID(ctx context.Context, id uint) (*domain.UserResponse, error)
	List(ctx context.Context, page, size int) ([]domain.UserResponse, int64, error)
	Update(ctx context.Context, id uint, req domain.UpdateUserRequest) (*domain.UserResponse, error)
	Delete(ctx context.Context, id uint) error
}

// userService 是 UserService 的实现
type userService struct {
	repo repository.UserRepository
}

// NewUserService 创建 UserService
func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

// Register 用户注册
func (s *userService) Register(ctx context.Context, req domain.CreateUserRequest) (*domain.UserResponse, error) {
	// 检查用户名是否已存在
	if _, err := s.repo.GetByUsername(ctx, req.Username); err == nil {
		return nil, errors.New("username already exists")
	}

	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password failed: %w", err)
	}

	role := req.Role
	if role == "" {
		role = "user"
	}

	user := &domain.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Email:    req.Email,
		Role:     role,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("create user failed: %w", err)
	}

	return toUserResponse(user), nil
}

// Login 用户登录
func (s *userService) Login(ctx context.Context, req domain.LoginRequest) (*domain.User, error) {
	user, err := s.repo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid username or password")
	}

	return user, nil
}

func (s *userService) GetByID(ctx context.Context, id uint) (*domain.UserResponse, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toUserResponse(user), nil
}

func (s *userService) List(ctx context.Context, page, size int) ([]domain.UserResponse, int64, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 10
	}

	offset := (page - 1) * size
	users, total, err := s.repo.List(ctx, offset, size)
	if err != nil {
		return nil, 0, err
	}

	var responses []domain.UserResponse
	for i := range users {
		responses = append(responses, *toUserResponse(&users[i]))
	}
	return responses, total, nil
}

func (s *userService) Update(ctx context.Context, id uint, req domain.UpdateUserRequest) (*domain.UserResponse, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Role != "" {
		user.Role = req.Role
	}

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}
	return toUserResponse(user), nil
}

func (s *userService) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

func toUserResponse(user *domain.User) *domain.UserResponse {
	return &domain.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}
}
