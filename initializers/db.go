package initializers

import (
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"log"
	"os"
)

var DB *gorm.DB

func ConnectToDatabase() {
	var err error
	dsn := os.Getenv("DB_URL")
	DB, err = gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Ошибка подключения к SQL Server: %s", err.Error())
		return
	}
	log.Printf("Connected to Database: %s", DB.Name())
}
