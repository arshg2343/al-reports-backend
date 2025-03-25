package main

import (
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	reporthandler "aetherlabs.com/glitch-report/report-handler"
)

func main() {

	err := reporthandler.Initialize(reporthandler.Config{
		Cloudinary: reporthandler.CloudinaryConfig{
			CloudURL: os.Getenv("CLOUD_URL"),
		},
		DatabaseURL: os.Getenv("DB_URL"),
	})
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

	server.Run(":8080")
}
