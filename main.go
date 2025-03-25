package main

import (
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	reporthandler "aetherlabs.com/glitch-report/report-handler"
)

func main() {

	err := reporthandler.Initialize(reporthandler.Config{
		Cloudinary: reporthandler.CloudinaryConfig{
			CloudURL: "cloudinary://715997977375351:i-lE-9w2nbNuc0i6l25lwVw-PkE@dl0n4fqcs",
		},
		DatabasePath: "", // Optional: leave empty for default
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
