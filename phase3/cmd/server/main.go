// Package main 是 Phase 3 程序的入口。
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-tutorial/phase3/internal/config"
	"go-tutorial/phase3/internal/domain"
	"go-tutorial/phase3/internal/logger"
	"go-tutorial/phase3/internal/repository"
	"go-tutorial/phase3/internal/router"
	"go-tutorial/phase3/internal/service"

	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// 加载配置
	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "load config failed: %v\n", err)
		os.Exit(1)
	}

	// 初始化日志
	if err := logger.Init(cfg.Log.Level, cfg.Log.Format); err != nil {
		fmt.Fprintf(os.Stderr, "init logger failed: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	// 初始化数据库
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		logger.L().Fatal("open database failed", zap.Error(err))
	}

	// 自动迁移
	if err := db.AutoMigrate(&domain.User{}); err != nil {
		logger.L().Fatal("auto migrate failed", zap.Error(err))
	}

	// 依赖注入
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)

	// 创建 HTTP 服务
	srv := router.Server(cfg, userService)

	// 优雅关闭
	go func() {
		logger.L().Info("server starting", zap.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.L().Fatal("server listen failed", zap.Error(err))
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.L().Info("server shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.L().Error("server shutdown failed", zap.Error(err))
	}

	logger.L().Info("server exited")
}
