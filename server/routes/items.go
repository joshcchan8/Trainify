package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/trainify/controllers"
)

// update route will be used to auto-generate new item (with AI) and update existing one
func SetItemRoutes(group *gin.RouterGroup) {
	group.GET("/", controllers.GetAllItems)
	group.POST("/", controllers.CreateItem)
	group.GET("/:id", controllers.GetItem)
	group.PATCH("/:id", controllers.UpdateItem)
	group.DELETE("/:id", controllers.DeleteItem)
}
