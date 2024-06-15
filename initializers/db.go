package initializers

import (
	"errors"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"webApp/model"
	"webApp/model/entity"
)

var DB *gorm.DB

func ConnectToDatabase() {
	var err error
	dsn := os.Getenv("DB_URL")
	DB, err = gorm.Open(sqlserver.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		log.Printf("Failed to connect to database. %s\n", err.Error())
		return
	}

	log.Printf("Connected to Database: %s", DB.Name())

	err = DB.AutoMigrate(&model.User{}, &entity.Apichats{})
	if err != nil {
		log.Printf("Failed to execute migrate Database: %s\n", err.Error())
		return
	}

	result := DB.Where(&entity.Apichats{Name: "gigaChat"}).First(&entity.Apichats{})
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		DB.Create(&entity.Apichats{Name: "gigaChat"})
	}
}
