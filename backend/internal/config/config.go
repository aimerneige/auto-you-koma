package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Queue    QueueConfig    `yaml:"queue"`
	LLM      LLMConfig      `yaml:"llm"`
	Auth     AuthConfig     `yaml:"auth"`
	Storage  StorageConfig  `yaml:"storage"`
	Quota    QuotaConfig    `yaml:"quota"`
}

type ServerConfig struct {
	Port int    `yaml:"port"`
	Mode string `yaml:"mode"`
}

type DatabaseConfig struct {
	Driver string       `yaml:"driver"`
	SQLite SQLiteConfig `yaml:"sqlite"`
}

type SQLiteConfig struct {
	Path string `yaml:"path"`
}

type QueueConfig struct {
	Driver string `yaml:"driver"`
}

type LLMConfig struct {
	Text  TextLLMConfig  `yaml:"text"`
	Image ImageLLMConfig `yaml:"image"`
}

type TextLLMConfig struct {
	Provider string       `yaml:"provider"`
	Gemini   GeminiConfig `yaml:"gemini"`
}

type GeminiConfig struct {
	APIKey string `yaml:"api_key"`
	Model  string `yaml:"model"`
}

type ImageLLMConfig struct {
	Provider   string           `yaml:"provider"`
	NanoBanana NanoBananaConfig `yaml:"nanobanana"`
}

type NanoBananaConfig struct {
	APIKey string `yaml:"api_key"`
}

type AuthConfig struct {
	JWTSecret   string `yaml:"jwt_secret"`
	TokenExpiry string `yaml:"token_expiry"`
}

type StorageConfig struct {
	BasePath string `yaml:"base_path"`
}

type QuotaConfig struct {
	DefaultDailyLimit int `yaml:"default_daily_limit"`
}

// LoadConfig loads the app configuration from yaml
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Printf("Failed to read config file: %v", err)
		return nil, err
	}
	
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Printf("Failed to unmarshal config yaml: %v", err)
		return nil, err
	}
	
	return &cfg, nil
}
