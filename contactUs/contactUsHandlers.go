package contactUs

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ContactUs struct {
	Name string `json:"name"`
	Email string `json:"email"`
	Contact string `json:"contact"`
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

	return nil
}

func HandleContactUs(c *gin.Context) {
	var contactUs ContactUs
	if err := c.ShouldBindJSON(&contactUs); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}



}

func (contactUs *ContactUs) saveToDB() {
	
}