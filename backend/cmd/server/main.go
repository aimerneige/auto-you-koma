package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/aimerneige/auto-you-koma/internal/config"
	"github.com/aimerneige/auto-you-koma/internal/handler"
	"github.com/aimerneige/auto-you-koma/internal/llm"
	"github.com/aimerneige/auto-you-koma/internal/repository"
	"github.com/aimerneige/auto-you-koma/internal/service"
	"github.com/aimerneige/auto-you-koma/internal/storage"
)

func main() {
	// Initialize database
	if err := config.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Create default Gin engine
	r := gin.Default()

	// Configure CORS
	r.Use(cors.Default())

	// Initialize image generator (using mock for now)
	imageGenerator := &llm.MockImageGenerator{}

	// Initialize storage
	storage := storage.NewStorage(".")

	// Initialize repositories
	charRepo := repository.NewCharacterRepository(config.GetDB())

	// Initialize services
	charSvc := service.NewCharacterService(charRepo, imageGenerator)

	// Initialize handlers
	charHandler := handler.NewCharacterHandler(charSvc)
	imageHandler := handler.NewImageHandler(storage)

	// Register API routes
	api := r.Group("/api/v1")
	charHandler.RegisterRoutes(api)
	imageHandler.RegisterRoutes(api)

	// Simple health check route
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// Start server
	log.Println("Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}