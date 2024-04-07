package initializer

import (
	"database/sql"
	"fmt"
	"os"
)

var DB *sql.DB

func ConnectToDatabase() {
	dsn := os.Getenv("DB_URL")
	DB, err := sql.Open("sqlserver", dsn)
	if err != nil {
		fmt.Println("Ошибка подключения к SQL Server: ", err.Error())
		return
	}
	defer DB.Close()
}
