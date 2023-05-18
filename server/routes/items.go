package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/joshchan4444/Trainify/controllers"
)

func SetItemRoutes(group *gin.RouterGroup) {
	group.GET("/", controllers.GetAllExercises)
}
