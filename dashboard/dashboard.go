package dashboard

import (
	"fmt"
	"net/http"

	reporthandler "aetherlabs.com/glitch-report/report-handler"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

type Config struct {
	DatabaseURL string
}

func Initialize(cfg Config) error {
	var err error
	db, err = gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	return nil
}

func FetchPendingReports(c *gin.Context) {
	var reports []reporthandler.Report
	result := db.Find(&reports)
	if result.Error != nil {
		fmt.Printf("Failed to fetch reports: %v", result.Error)
	}
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("%v reports found.", len(reports)),
		"reports": reports,
	})
}

func DeleteReport() {
	// Delete image from cloudinary
	// Delete report from db
}

func SendResolvedEmail() {
	// POST handler with issue resolvement
}
