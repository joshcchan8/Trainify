package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var DB *sql.DB

func ConnectDB() {

	dotenvErr := godotenv.Load(".env")
	if dotenvErr != nil {
		log.Fatal("Error Loading .env File: ", dotenvErr)
	}

	connectionString := os.Getenv("CONNECTION_STRING")

	db, databaseErr := sql.Open("mysql", connectionString)
	if databaseErr != nil {
		log.Fatal("Database connection error")
	}

	db.SetConnMaxIdleTime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal("database ping error", pingErr)
	}

	DB = db
	fmt.Println("Connected to DB")
}
