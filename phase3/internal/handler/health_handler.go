package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthHandler 健康检查处理器
type HealthHandler struct{}

// NewHealthHandler 创建 HealthHandler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// Check 健康检查
func (h *HealthHandler) Check(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"time":   timeNow(),
	})
}

func timeNow() string {
	return "healthy"
}
