package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-tutorial/phase3/internal/config"
	"go-tutorial/phase3/internal/domain"
	"go-tutorial/phase3/internal/logger"
	"go-tutorial/phase3/internal/repository"
	"go-tutorial/phase3/internal/router"
	"go-tutorial/phase3/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupIntegrationServer(t *testing.T) *httptest.Server {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	err = db.AutoMigrate(&domain.User{})
	require.NoError(t, err)

	cfg := &config.Config{
		App: config.AppConfig{Mode: "test"},
		JWT: config.JWTConfig{Secret: "test-secret", ExpireHours: 24},
		Log: config.LogConfig{Level: "error"},
	}

	logger.Init(cfg.Log.Level, cfg.Log.Format)

	repo := repository.NewUserRepository(db)
	svc := service.NewUserService(repo)
	r := router.New(cfg, svc)

	return httptest.NewServer(r)
}

func TestIntegration_RegisterAndLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	server := setupIntegrationServer(t)
	defer server.Close()

	client := server.Client()

	// 注册
	registerBody := `{"username":"alice","password":"password123","email":"alice@example.com"}`
	resp, err := client.Post(server.URL+"/api/v1/auth/register", "application/json", bytes.NewBufferString(registerBody))
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	// 登录
	loginBody := `{"username":"alice","password":"password123"}`
	resp, err = client.Post(server.URL+"/api/v1/auth/login", "application/json", bytes.NewBufferString(loginBody))
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var loginResp domain.LoginResponse
	err = json.NewDecoder(resp.Body).Decode(&loginResp)
	require.NoError(t, err)
	assert.NotEmpty(t, loginResp.Token)
	assert.Equal(t, "alice", loginResp.User.Username)
}

func TestIntegration_Health(t *testing.T) {
	gin.SetMode(gin.TestMode)
	server := setupIntegrationServer(t)
	defer server.Close()

	resp, err := server.Client().Get(server.URL + "/health")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
