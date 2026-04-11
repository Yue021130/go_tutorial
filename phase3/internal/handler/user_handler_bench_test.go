package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-tutorial/phase3/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
)

func BenchmarkUserHandler_Register(b *testing.B) {
	gin.SetMode(gin.TestMode)

	svc := new(mockUserService)
	h := NewUserHandler(svc)

	svc.On("Register", mock.Anything, mock.AnythingOfType("domain.CreateUserRequest")).Return(&domain.UserResponse{
		ID:       1,
		Username: "alice",
	}, nil)

	body := []byte(`{"username":"alice","password":"password123","email":"alice@example.com"}`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := gin.New()
		r.POST("/register", h.Register)

		req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
	}
}
