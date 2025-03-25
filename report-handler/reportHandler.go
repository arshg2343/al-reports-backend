package reporthandler

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	cld *cloudinary.Cloudinary
	db  *gorm.DB
)

// CloudinaryConfig holds Cloudinary connection details
type CloudinaryConfig struct {
	CloudURL string
}

// Report struct with gorm tags for database mapping
type Report struct {
	gorm.Model
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

// Config holds the configuration for the package
type Config struct {
	Cloudinary   CloudinaryConfig
	DatabasePath string
}

// Initialize initializes the package with given configuration
func Initialize(cfg Config) error {
	// Initialize Cloudinary
	var err error
	cld, err = cloudinary.NewFromURL(cfg.Cloudinary.CloudURL)
	if err != nil {
		return fmt.Errorf("failed to initialize Cloudinary: %v", err)
	}

	// Ensure database directory exists
	dbDir := "database"
	if err := os.MkdirAll(dbDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create database directory: %v", err)
	}

	// Use provided database path or default
	dbPath := cfg.DatabasePath
	if dbPath == "" {
		dbPath = dbDir + "/glitch_reports.db"
	}

	// Initialize database
	db, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	// Auto migrate the Report model
	err = db.AutoMigrate(&Report{})
	if err != nil {
		return fmt.Errorf("failed to auto-migrate database: %v", err)
	}

	return nil
}

// New creates a new Report instance with validation
func New(email, username, deviceType, browserInfo, glitchType, glitchLocation, glitchDescription, stepsToReproduce, urgency string) (*Report, error) {
	// Validate input fields
	if email == "" || username == "" || deviceType == "" || browserInfo == "" ||
		glitchType == "" || glitchLocation == "" || glitchDescription == "" ||
		stepsToReproduce == "" || urgency == "" {
		return nil, errors.New("invalid input: all fields are required")
	}

	return &Report{
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

// HandleReport processes the incoming glitch report
func HandleReport(c *gin.Context) {
	// Check if Cloudinary and DB are initialized
	if cld == nil || db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Service not initialized"})
		return
	}

	// Extract form data
	email := c.PostForm("email")
	username := c.PostForm("username")
	deviceType := c.PostForm("deviceType")
	browserInfo := c.PostForm("browserInfo")
	glitchType := c.PostForm("glitchType")
	glitchLocation := c.PostForm("glitchLocation")
	glitchDescription := c.PostForm("glitchDescription")
	stepsToReproduce := c.PostForm("stepsToReproduce")
	urgency := c.PostForm("urgency")

	// Create new report
	report, err := New(email, username, deviceType, browserInfo, glitchType,
		glitchLocation, glitchDescription, stepsToReproduce, urgency)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Handle file upload (optional)
	var imageURL string
	file, err := c.FormFile("attachment")
	if err == nil && file != nil {
		// If an image is uploaded
		imageURL, err = uploadToCloudinary(file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image"})
			return
		}
		report.AttachmentURL = imageURL
	}

	// Save report to database
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

// uploadToCloudinary handles file upload to Cloudinary
func uploadToCloudinary(file *multipart.FileHeader) (string, error) {
	// Check if Cloudinary is initialized
	if cld == nil {
		return "", errors.New("cloudinary not initialized")
	}

	// Open the file
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// Upload to Cloudinary
	uploadResult, err := cld.Upload.Upload(context.TODO(), src, uploader.UploadParams{
		Folder: "glitch_reports", // Optional: organize uploads in a specific folder
	})
	if err != nil {
		return "", err
	}

	return uploadResult.SecureURL, nil
}

// GetAllReports retrieves all reports from the database
func GetAllReports() ([]Report, error) {
	var reports []Report
	result := db.Find(&reports)
	if result.Error != nil {
		return nil, result.Error
	}
	return reports, nil
}
