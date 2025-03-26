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

type DeleteRequest struct {
	UID string `json:"uid"`
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
		c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to fetch reports: %v", result.Error))
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("%v reports found.", len(reports)),
		"reports": reports,
	})
}

func DeleteReport(c *gin.Context) {
	var delete DeleteRequest
	err := c.ShouldBindJSON(&delete)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Inavlid Request Format",
		})
	}
	result := db.Where("uid = ?", delete.UID).Delete(&reporthandler.Report{})
	if result.Error != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to delete report: %v", result.Error))
		return
	}
	if result.RowsAffected == 0 {
		c.String(http.StatusInternalServerError, "No reports found by the given UID")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Report deleted successfully.",
	})
}

func SendResolvedEmail() {
	// POST handler with issue resolvement
}
