// Package router 注册 HTTP 路由。
package router

import (
	"net/http"
	"net/http/pprof"

	"github.com/gin-gonic/gin"
	"go-tutorial/phase3/internal/config"
	"go-tutorial/phase3/internal/handler"
	"go-tutorial/phase3/internal/middleware"
	"go-tutorial/phase3/internal/service"
)

// New 创建 Gin 引擎并注册路由
func New(cfg *config.Config, userService service.UserService) *gin.Engine {
	if cfg.App.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// 全局中间件
	r.Use(middleware.Recovery())
	r.Use(middleware.RequestLogger())

	rateLimiter := middleware.NewRateLimiter(10, 20)
	r.Use(rateLimiter.Limit())

	// 配置 JWT
	middleware.SetJWTConfig(cfg.JWT.Secret, cfg.JWT.ExpireHours)

	userHandler := handler.NewUserHandler(userService)
	healthHandler := handler.NewHealthHandler()

	// 健康检查
	r.GET("/health", healthHandler.Check)

	// API v1
	api := r.Group("/api/v1")
	{
		// 公开接口
		api.POST("/auth/register", userHandler.Register)
		api.POST("/auth/login", userHandler.Login)

		// 需要 JWT 认证
		authorized := api.Group("")
		authorized.Use(middleware.JWTAuth())
		{
			authorized.GET("/users", userHandler.List)
			authorized.GET("/users/:id", userHandler.GetByID)
			authorized.POST("/users", userHandler.Create)
			authorized.PUT("/users/:id", userHandler.Update)
			authorized.DELETE("/users/:id", userHandler.Delete)
		}
	}

	// pprof 性能分析路由（生产环境应限制访问）
	registerPprof(r)

	return r
}

func registerPprof(r *gin.Engine) {
	pr := r.Group("/debug/pprof")
	{
		pr.GET("/", func(c *gin.Context) {
			pprof.Index(c.Writer, c.Request)
		})
		pr.GET("/cmdline", func(c *gin.Context) {
			pprof.Cmdline(c.Writer, c.Request)
		})
		pr.GET("/profile", func(c *gin.Context) {
			pprof.Profile(c.Writer, c.Request)
		})
		pr.GET("/symbol", func(c *gin.Context) {
			pprof.Symbol(c.Writer, c.Request)
		})
		pr.POST("/symbol", func(c *gin.Context) {
			pprof.Symbol(c.Writer, c.Request)
		})
		pr.GET("/trace", func(c *gin.Context) {
			pprof.Trace(c.Writer, c.Request)
		})
		pr.GET("/allocs", pprofHandler("allocs"))
		pr.GET("/block", pprofHandler("block"))
		pr.GET("/goroutine", pprofHandler("goroutine"))
		pr.GET("/heap", pprofHandler("heap"))
		pr.GET("/mutex", pprofHandler("mutex"))
		pr.GET("/threadcreate", pprofHandler("threadcreate"))
	}
}

func pprofHandler(name string) gin.HandlerFunc {
	return func(c *gin.Context) {
		handler := pprof.Handler(name)
		handler.ServeHTTP(c.Writer, c.Request)
	}
}

// Server 返回 http.Server
func Server(cfg *config.Config, userService service.UserService) *http.Server {
	r := New(cfg, userService)
	return &http.Server{
		Addr:    cfg.App.Addr(),
		Handler: r,
	}
}
