package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-tutorial/phase3/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockUserService 是 UserService 的 mock
type mockUserService struct {
	mock.Mock
}

func (m *mockUserService) Register(ctx context.Context, req domain.CreateUserRequest) (*domain.UserResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.UserResponse), args.Error(1)
}

func (m *mockUserService) Login(ctx context.Context, req domain.LoginRequest) (*domain.User, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *mockUserService) GetByID(ctx context.Context, id uint) (*domain.UserResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.UserResponse), args.Error(1)
}

func (m *mockUserService) List(ctx context.Context, page, size int) ([]domain.UserResponse, int64, error) {
	args := m.Called(ctx, page, size)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]domain.UserResponse), args.Get(1).(int64), args.Error(2)
}

func (m *mockUserService) Update(ctx context.Context, id uint, req domain.UpdateUserRequest) (*domain.UserResponse, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.UserResponse), args.Error(1)
}

func (m *mockUserService) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func setupGin() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestUserHandler_Register(t *testing.T) {
	svc := new(mockUserService)
	h := NewUserHandler(svc)
	r := setupGin()
	r.POST("/register", h.Register)

	svc.On("Register", mock.Anything, domain.CreateUserRequest{
		Username: "alice",
		Password: "password123",
		Email:    "alice@example.com",
	}).Return(&domain.UserResponse{
		ID:       1,
		Username: "alice",
		Email:    "alice@example.com",
		Role:     "user",
	}, nil)

	body := `{"username":"alice","password":"password123","email":"alice@example.com"}`
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp domain.UserResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "alice", resp.Username)
}

func TestUserHandler_Register_BadRequest(t *testing.T) {
	svc := new(mockUserService)
	h := NewUserHandler(svc)
	r := setupGin()
	r.POST("/register", h.Register)

	body := `{"username":"al","password":"short"}`
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserHandler_GetByID(t *testing.T) {
	svc := new(mockUserService)
	h := NewUserHandler(svc)
	r := setupGin()
	r.GET("/users/:id", h.GetByID)

	svc.On("GetByID", mock.Anything, uint(1)).Return(&domain.UserResponse{
		ID:       1,
		Username: "alice",
	}, nil)

	req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUserHandler_GetByID_NotFound(t *testing.T) {
	svc := new(mockUserService)
	h := NewUserHandler(svc)
	r := setupGin()
	r.GET("/users/:id", h.GetByID)

	svc.On("GetByID", mock.Anything, uint(999)).Return(nil, errors.New("not found"))

	req := httptest.NewRequest(http.MethodGet, "/users/999", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
