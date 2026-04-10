package service

import (
	"context"
	"errors"
	"testing"

	"go-tutorial/phase3/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserService_Register(t *testing.T) {
	repo := new(mockUserRepository)
	svc := NewUserService(repo)
	ctx := context.Background()

	repo.On("GetByUsername", ctx, "alice").Return(nil, errors.New("not found"))
	repo.On("Create", ctx, mock.AnythingOfType("*domain.User")).Return(nil)

	req := domain.CreateUserRequest{
		Username: "alice",
		Password: "password123",
		Email:    "alice@example.com",
	}

	resp, err := svc.Register(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, "alice", resp.Username)
	assert.Equal(t, "alice@example.com", resp.Email)

	repo.AssertExpectations(t)
}

func TestUserService_Register_UsernameExists(t *testing.T) {
	repo := new(mockUserRepository)
	svc := NewUserService(repo)
	ctx := context.Background()

	repo.On("GetByUsername", ctx, "alice").Return(&domain.User{Username: "alice"}, nil)

	req := domain.CreateUserRequest{
		Username: "alice",
		Password: "password123",
	}

	_, err := svc.Register(ctx, req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
}

func TestUserService_Login(t *testing.T) {
	repo := new(mockUserRepository)
	svc := NewUserService(repo)
	ctx := context.Background()

	// 先注册一个用户
	repo.On("GetByUsername", ctx, "alice").Return(&domain.User{
		ID:       1,
		Username: "alice",
		Password: "$2a$10$abcdefghijklmnopqrstuuxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", // invalid hash for test
	}, nil)

	req := domain.LoginRequest{
		Username: "alice",
		Password: "wrongpassword",
	}

	_, err := svc.Login(ctx, req)
	assert.Error(t, err)
}

func TestUserService_List(t *testing.T) {
	repo := new(mockUserRepository)
	svc := NewUserService(repo)
	ctx := context.Background()

	users := []domain.User{
		{ID: 1, Username: "alice"},
		{ID: 2, Username: "bob"},
	}

	repo.On("List", ctx, 0, 10).Return(users, int64(2), nil)

	resp, total, err := svc.List(ctx, 1, 10)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, resp, 2)
}
