package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/trainify/controllers"
	"github.com/trainify/middleware"
)

func SetUserRoutes(group *gin.RouterGroup) {
	group.POST("/register", controllers.Register)
	group.POST("/login", controllers.Login)
	group.PATCH("/update", middleware.AuthenticationMiddleware(), controllers.UpdateUser)
	group.DELETE("/delete", middleware.AuthenticationMiddleware(), controllers.DeleteUser)
}
