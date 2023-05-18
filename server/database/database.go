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
	db, err := sql.Open("mysql", "trainee:8787@tcp(localhost:3306)/trainify")
	if err != nil {
		log.Fatal("Database connection error")
	}

	db.SetConnMaxIdleTime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	DB = db
	fmt.Println("Connected to DB")
}
