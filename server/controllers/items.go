package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/trainify/database"
	model "github.com/trainify/models"
)

func GetAllExercises(c *gin.Context) {

	rows, err := database.DB.Query("SELECT * FROM items")

	if err != nil {
		c.IndentedJSON(http.StatusOK, gin.H{
			"message": "get all exercises failed",
		})
		return
	}

	var data []model.Item

	for rows.Next() {
		var item model.Item
		err := rows.Scan(
			&item.ItemID,
			&item.ItemType,
			&item.ItemName,
			&item.Difficulty,
			&item.Minutes,
			&item.CaloriesBurned,
			&item.TargetedMuscleGroup,
		)

		if err != nil {
			log.Fatal(err)
			return
		}

		fmt.Println(item)
		data = append(data, item)
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"data": data,
	})
}
