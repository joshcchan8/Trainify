package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/trainify/controllers"
)

func SetItemRoutes(group *gin.RouterGroup) {
	group.GET("/", controllers.GetAllExercises)
}
