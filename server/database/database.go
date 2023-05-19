package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func ConnectDB() {
	db, databaseErr := sql.Open("mysql", "trainee:8787@tcp(localhost:3306)/trainify")
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
