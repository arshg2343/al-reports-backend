package main

import (
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	// "github.com/joho/godotenv"

	"aetherlabs.com/glitch-report/dashboard"
	reporthandler "aetherlabs.com/glitch-report/report-handler"
)

func main() {
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatalf("Initialization of env error: %v", err)
	// }

	err := reporthandler.Initialize(reporthandler.Config{
		Cloudinary: reporthandler.CloudinaryConfig{
			CloudURL: os.Getenv("CLOUD_URL"),
		},
		DatabaseURL: os.Getenv("DATABASE_URL"),
	})
	if err != nil {
		log.Fatalf("Initialization error: %v", err)
	}

	err = dashboard.Initialize(dashboard.Config{DatabaseURL: os.Getenv("DATABASE_URL")})
	if err != nil {
		log.Fatalf("Initialization error: %v", err)
	}

	server := gin.Default()

	server.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	server.POST("/report-glitch", reporthandler.HandleReport)
	server.GET("/active-reports", dashboard.FetchPendingReports)
	server.POST("/delete-report", dashboard.DeleteReport)
	// Handle resolving and sending email
	server.POST("/resolve-report")

	server.Run(":8080")
}
