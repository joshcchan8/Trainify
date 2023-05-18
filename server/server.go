package main

import (
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Connect to DB
	db, err := sql.Open("mysql", "trainee:8787@tcp(localhost:3306)/database")
	if err != nil {
		return
	}
	defer db.Close()
	fmt.Println("Connected to DB")

	// Router
	router := gin.Default()
	router.Run("localhost:3000")
}
