// Package config 使用 Viper 管理应用配置。
//
// 与 Java 对比：
// - Java Spring Boot: application.yml + @ConfigurationProperties
// - Go: Viper + 结构体映射
package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config 是应用配置的总入口
type Config struct {
	App AppConfig `mapstructure:"app"`
	JWT JWTConfig `mapstructure:"jwt"`
	Log LogConfig `mapstructure:"log"`
}

// AppConfig 应用相关配置
type AppConfig struct {
	Name    string `mapstructure:"name"`
	Version string `mapstructure:"version"`
	Port    int    `mapstructure:"port"`
	Mode    string `mapstructure:"mode"`
}

// JWTConfig JWT 相关配置
type JWTConfig struct {
	Secret      string `mapstructure:"secret"`
	ExpireHours int    `mapstructure:"expire_hours"`
}

// LogConfig 日志相关配置
type LogConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

// Load 从配置文件加载配置
func Load(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.AutomaticEnv() // 允许环境变量覆盖配置文件

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config failed: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config failed: %w", err)
	}

	return &cfg, nil
}

// Addr 返回监听地址
func (a *AppConfig) Addr() string {
	return fmt.Sprintf(":%d", a.Port)
}
