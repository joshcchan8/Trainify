package main

import (
	"github.com/trainify/database"
	"github.com/trainify/routes"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Connect to DB
	database.ConnectDB()

	// Router
	router := gin.Default()

	// Item routes (private)

	itemRoutes := router.Group("/items")
	{
		routes.SetItemRoutes(itemRoutes)
	}

	router.Run("localhost:3000")
}
