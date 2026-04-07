package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"go-tutorial/phase3/internal/logger"
	"go.uber.org/zap"
)

// Recovery 自定义错误恢复中间件
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.L().Error("panic recovered",
					zap.Any("error", err),
					zap.String("stack", string(debug.Stack())),
				)

				c.JSON(http.StatusInternalServerError, gin.H{
					"error": fmt.Sprintf("internal server error: %v", err),
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}
