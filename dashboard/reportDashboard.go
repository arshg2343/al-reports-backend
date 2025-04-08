package dashboard

import (
	"fmt"
	"net/http"

	reporthandler "aetherlabs.com/glitch-report/report-handler"
	"github.com/gin-gonic/gin"
)

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
	var delete Request
	err := c.ShouldBindJSON(&delete)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Inavlid Request Format",
		})
	}
	result := deleteReportByUID(delete.UID)
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

func ResolveReport(c *gin.Context) {
	var resolve Request
	err := c.ShouldBindJSON(&resolve)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Inavlid Request Format",
		})
	}
	sendResolvedEmail()
	deleteReportByUID(resolve.UID)
}

func sendResolvedEmail() {
	// send issue resolved mail
}
