package main

import (
	"log"

	"github.com/aimerneige/auto-you-koma/internal/config"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig("config.example.yaml")
	if err != nil {
		log.Fatalf("Fail to load configuration: %v", err)
	}

	// 设定运行模式
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 初始化 Gin
	r := gin.Default()

	// 提供一个健康检查 API
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"message": "Auto Yon Koma API is running",
		})
	})

	log.Printf("Server starting on port %d...", cfg.Server.Port)
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
}
