package main

import (
	"github.com/trainify/database"
	"github.com/trainify/middleware"
	"github.com/trainify/routes"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Connect to DB
	database.ConnectDB()

	// Router
	router := gin.Default()
	router.Use(middleware.CorsMiddleware())

	// User routes (public)
	userRoutes := router.Group("/users")
	routes.SetUserRoutes(userRoutes)

	// Item routes (private)
	itemRoutes := router.Group("/items")
	itemRoutes.Use(middleware.AuthenticationMiddleware())
	routes.SetItemRoutes(itemRoutes)

	// Profile routes (private)
	profileRoutes := router.Group("/profiles")
	profileRoutes.Use(middleware.AuthenticationMiddleware())
	routes.SetProfileRoutes(profileRoutes)

	router.Run("localhost:8000")
}
