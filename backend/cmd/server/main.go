package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kendall1978/project_salary/backend/internal/config"
	"github.com/kendall1978/project_salary/backend/internal/crypto"
	"github.com/kendall1978/project_salary/backend/internal/database"
	"github.com/kendall1978/project_salary/backend/internal/handlers"
	"github.com/kendall1978/project_salary/backend/internal/routes"
	"github.com/kendall1978/project_salary/backend/internal/scraper"
	"github.com/robfig/cron/v3"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if err := crypto.SetKey(cfg.EncryptionKey); err != nil {
		log.Fatalf("Invalid encryption key: %v", err)
	}

	handlers.SetJWTSecret(cfg.JWTSecret)

	if err := database.Connect(cfg.DatabaseURL); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := database.Migrate(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	os.MkdirAll("uploads", 0755)

	routes.Setup(r)

	// Monthly OES data check — 1st of every month at midnight
	cronScheduler := cron.New()
	cronScheduler.AddFunc("0 0 1 * *", func() {
		log.Println("Running scheduled OES data check...")
		if err := scraper.FetchOESData(false, ""); err != nil {
			log.Printf("Scheduled OES import failed: %v", err)
		}
	})
	cronScheduler.Start()
	defer cronScheduler.Stop()

	log.Printf("Server starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
