package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/joshchan4444/Trainify/controllers"
)

func SetRoutes(router *gin.Engine) {
	router.GET("/", controllers.GetAllExercises)
}
