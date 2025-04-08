package dashboard

import (
	"fmt"
	"net/http"

	"aetherlabs.com/glitch-report/contactUs"

	"github.com/gin-gonic/gin"
)

func FetchPendingInquiries(c *gin.Context) {
	var inquiries []contactUs.ContactUs
	result := db.Find(&inquiries)
	if result.Error != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to fetch inquiries: %v", result.Error))
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":   fmt.Sprintf("%v inquiries found.", len(inquiries)),
		"inquiries": inquiries,
	})
}

func DeleteInquiry(c *gin.Context) {
	var delete Request
	err := c.ShouldBindJSON(&delete)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Inavlid Request Format",
		})
	}
	result := deleteReportByUID(delete.UID)
	if result.Error != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to delete inquiry: %v", result.Error))
		return
	}
	if result.RowsAffected == 0 {
		c.String(http.StatusInternalServerError, "No inquiries found by the given UID")
		return
	}
}
