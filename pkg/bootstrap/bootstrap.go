package bootstrap

import (
	"fmt"
	"log"
	"os"

	"github.com/tzincker/gocourse_domain/domain"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitLogger() *log.Logger {
	return log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
}

func DBConnection() (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_NAME"))

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if os.Getenv("DATABASE_IS_DEBUG") == "true" {
		db = db.Debug()
	}

	if os.Getenv("DATABASE_MIGRATE") == "true" {
		if err := db.AutoMigrate(&domain.Enrollment{}); err != nil {
			return nil, err
		}
	}

	if err != nil {
		return nil, err
	}

	return db, nil
}

func Address() string {
	return fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT"))
}
