package main

import (
	"log"

	"github.com/aimerneige/auto-you-koma/internal/compositor"
	"github.com/aimerneige/auto-you-koma/internal/config"
	"github.com/aimerneige/auto-you-koma/internal/handler"
	"github.com/aimerneige/auto-you-koma/internal/llm"
	"github.com/aimerneige/auto-you-koma/internal/middleware"
	"github.com/aimerneige/auto-you-koma/internal/model"
	sqlite_repo "github.com/aimerneige/auto-you-koma/internal/repository/sqlite"
	"github.com/aimerneige/auto-you-koma/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	cfg, err := config.LoadConfig("config.example.yaml")
	if err != nil {
		log.Fatalf("Fail to load configuration: %v", err)
	}

	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize DB
	db, err := gorm.Open(sqlite.Open(cfg.Database.SQLite.Path), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// Auto Migrate models
	err = db.AutoMigrate(
		&model.User{},
		&model.Character{},
		&model.CharacterVariant{},
		&model.CharacterImage{},
		&model.CharacterGroup{},
		&model.CharacterGroupMember{},
		&model.Project{},
		&model.Script{},
		&model.Generation{},
		&model.Series{},
		&model.CharacterState{},
	)
	if err != nil {
		log.Fatalf("failed to auto migrate: %v", err)
	}

	userRepo := sqlite_repo.NewUserRepository(db)
	authSvc := service.NewAuthService(userRepo, cfg.Auth)
	authHandler := handler.NewAuthHandler(authSvc)

	charRepo := sqlite_repo.NewCharacterRepository(db)
	charSvc := service.NewCharacterService(charRepo, cfg.Storage)
	charHandler := handler.NewCharacterHandler(charSvc)

	// Phase 1.3 LLM and Scripts
	openAIGen := llm.NewOpenAIGenerator(cfg.LLM.Text.OpenAI)
	scriptRepo := sqlite_repo.NewScriptRepository(db)
	scriptSvc := service.NewScriptService(scriptRepo, openAIGen)
	scriptHandler := handler.NewScriptHandler(scriptSvc)

	// Phase 1.5 Image Generation and Compositing
	imageGen := llm.NewGenericImageGenerator(cfg.LLM.Image)
	comp := compositor.NewCompositor(cfg.Storage.BasePath)
	genRepo := sqlite_repo.NewGenerationRepository(db)
	genSvc := service.NewGenerationService(genRepo, scriptRepo, imageGen, comp)
	genHandler := handler.NewGenerationHandler(genSvc)

	// Phase 2.1 Continuity
	seriesRepo := sqlite_repo.NewSeriesRepository(db)
	stateRepo := sqlite_repo.NewCharacterStateRepository(db)
	contSvc := service.NewContinuityService(stateRepo, scriptRepo, openAIGen)
	seriesHandler := handler.NewSeriesHandler(seriesRepo, contSvc)

	r := gin.Default()

	// Serve Static Files for Assets
	r.Static("/assets", cfg.Storage.BasePath)

	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "message": "Auto Yon Koma API is running"})
	})

	apiV1 := r.Group("/api/v1")
	{
		authRoutes := apiV1.Group("/auth")
		{
			authRoutes.POST("/register", authHandler.Register)
			authRoutes.POST("/login", authHandler.Login)
			
			// Protected routes
			protected := apiV1.Group("")
			protected.Use(middleware.AuthMiddleware(cfg.Auth))
			{
				authProt := protected.Group("/auth")
				authProt.POST("/2fa/setup", authHandler.Setup2FA)
				authProt.POST("/2fa/verify", authHandler.Verify2FA)
				authProt.GET("/me", authHandler.Me)

				charRoutes := protected.Group("/characters")
				{
					charRoutes.GET("", charHandler.List)
					charRoutes.POST("", charHandler.Create)
					charRoutes.GET("/:id", charHandler.Get)
					charRoutes.POST("/:id/images", charHandler.UploadImage)
					charRoutes.POST("/:id/variants", charHandler.AddVariant)
				}

				scriptRoutes := protected.Group("/scripts")
				{
					scriptRoutes.GET("", scriptHandler.ListByProject)
					scriptRoutes.POST("", scriptHandler.Create)
					scriptRoutes.GET("/:id", scriptHandler.Get)
					scriptRoutes.PUT("/:id", scriptHandler.Update)
					scriptRoutes.POST("/generate", scriptHandler.GenerateStream)
					scriptRoutes.POST("/:id/parse", scriptHandler.Parse)
					scriptRoutes.PUT("/:id/panels/update", scriptHandler.UpdatePanel)
					scriptRoutes.POST("/:id/panels/regenerate", scriptHandler.RegeneratePanel)
				}

				genRoutes := protected.Group("/generations")
				{
					genRoutes.POST("", genHandler.Start)
					genRoutes.GET("/:id", genHandler.Get)
				}

				seriesRoutes := protected.Group("/series")
				{
					seriesRoutes.POST("", seriesHandler.Create)
					seriesRoutes.GET("", seriesHandler.List)
					seriesRoutes.POST("/:id/continuity/:scriptId/:characterId", seriesHandler.SythesizeMemory)
				}
			}
		}
	}

	log.Printf("Server starting on port %d...", cfg.Server.Port)
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
}
