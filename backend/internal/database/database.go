package database

import (
	"log"

	"github.com/kendall1978/project_salary/backend/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect(databaseURL string) error {
	var err error
	DB, err = gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		return err
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)

	log.Println("Database connected successfully")
	return nil
}

func Migrate() error {
	return DB.AutoMigrate(
		&models.User{},
		&models.UserMFA{},
		&models.RecoveryCode{},
		&models.UserCompensation{},
		&models.MFAChallenge{},
		&models.Region{},
		&models.SalarySnapshot{},
		&models.ScrapeLog{},
		// Portfolio models
		&models.BlogPost{},
		&models.HomeSection{},
		&models.SocialLink{},
		&models.SiteSetting{},
	)
}
