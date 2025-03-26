package reporthandler

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	cld *cloudinary.Cloudinary
	db  *gorm.DB
)

type CloudinaryConfig struct {
	CloudURL string
}

// Add uid to Reports.

// Send confirmation mail after receiving report.

type Report struct {
	gorm.Model
	UID               string `gorm:"primaryKey;not null;unique"`
	Email             string `gorm:"not null"`
	Username          string `gorm:"not null"`
	DeviceType        string `gorm:"not null"`
	BrowserInfo       string `gorm:"not null"`
	GlitchType        string `gorm:"not null"`
	GlitchLocation    string `gorm:"not null"`
	GlitchDescription string `gorm:"not null"`
	StepsToReproduce  string `gorm:"not null"`
	AttachmentURL     string
	Urgency           string `gorm:"not null"`
}

type Config struct {
	Cloudinary  CloudinaryConfig
	DatabaseURL string
}

func Initialize(cfg Config) error {
	var err error
	cld, err = cloudinary.NewFromURL(cfg.Cloudinary.CloudURL)
	if err != nil {
		return fmt.Errorf("failed to initialize Cloudinary: %v", err)
	}

	db, err = gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	err = db.AutoMigrate(&Report{})
	if err != nil {
		return fmt.Errorf("failed to auto-migrate database: %v", err)
	}

	return nil
}

func New(email, username, deviceType, browserInfo, glitchType, glitchLocation, glitchDescription, stepsToReproduce, urgency string) (*Report, error) {
	if email == "" || username == "" || deviceType == "" || browserInfo == "" ||
		glitchType == "" || glitchLocation == "" || glitchDescription == "" ||
		stepsToReproduce == "" || urgency == "" {
		return nil, errors.New("invalid input: all fields are required")
	}

	return &Report{
		UID:               uuid.New().String(),
		Email:             email,
		Username:          username,
		DeviceType:        deviceType,
		BrowserInfo:       browserInfo,
		GlitchType:        glitchType,
		GlitchLocation:    glitchLocation,
		GlitchDescription: glitchDescription,
		StepsToReproduce:  stepsToReproduce,
		Urgency:           urgency,
	}, nil
}

func HandleReport(c *gin.Context) {
	if cld == nil || db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Service not initialized"})
		return
	}

	email := c.PostForm("email")
	username := c.PostForm("username")
	deviceType := c.PostForm("deviceType")
	browserInfo := c.PostForm("browserInfo")
	glitchType := c.PostForm("glitchType")
	glitchLocation := c.PostForm("glitchLocation")
	glitchDescription := c.PostForm("glitchDescription")
	stepsToReproduce := c.PostForm("stepsToReproduce")
	urgency := c.PostForm("urgency")

	report, err := New(email, username, deviceType, browserInfo, glitchType,
		glitchLocation, glitchDescription, stepsToReproduce, urgency)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var imageURL string
	file, err := c.FormFile("attachment")
	if err == nil && file != nil {
		imageURL, err = uploadToCloudinary(file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image"})
			return
		}
		report.AttachmentURL = imageURL
	}

	result := db.Create(report)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save report"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Report submitted successfully",
		"report":  report,
	})
}

func uploadToCloudinary(file *multipart.FileHeader) (string, error) {
	if cld == nil {
		return "", errors.New("cloudinary not initialized")
	}

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	uploadResult, err := cld.Upload.Upload(context.TODO(), src, uploader.UploadParams{
		Folder: "glitch_reports",
	})
	if err != nil {
		return "", err
	}

	return uploadResult.SecureURL, nil
}
