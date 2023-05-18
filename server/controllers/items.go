package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAllExercises(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{
		"message": "get all exercises",
	})
}
