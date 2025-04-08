package contactUs

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ContactUs struct {
	gorm.Model
	UID     string
	Name    string `json:"name"`
	Email   string `json:"email"`
	Contact string `json:"contact"`
	Type    string `json:"type"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

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

	err = db.AutoMigrate(&ContactUs{})
	if err != nil {
		return fmt.Errorf("failed to auto-migrate database: %v", err)
	}

	return nil
}

func HandleContactUs(c *gin.Context) {
	var contactUs ContactUs
	if err := c.ShouldBindJSON(&contactUs); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	contactUs.UID = uuid.New().String()

	err := contactUs.saveToDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save report"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Inquiry saved successfully"})

}

func (contactUs *ContactUs) saveToDB() error {
	result := db.Create(contactUs)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
