package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/trainify/controllers"
)

func SetProfileRoutes(group *gin.RouterGroup) {
	group.GET("/", controllers.GetProfile)
	group.PATCH("/", controllers.UpdateProfile)
}
