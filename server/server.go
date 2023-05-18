package main

import (
	"github.com/gin-gonic/gin"
	// "database/sql"
	// "github.com/go-sql-driver/myysql"
)

func main() {
	router := gin.Default()
	router.Run("localhost:3000")
}
