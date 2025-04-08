package dashboard

import (
	"fmt"

	reporthandler "aetherlabs.com/glitch-report/report-handler"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

type Config struct {
	DatabaseURL string
}

type Request struct {
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

func deleteReportByUID(uid string) *gorm.DB {
	result := db.Where("uid = ?", uid).Delete(&reporthandler.Report{})
	return result
}
