package controllers

import (
	"encoding/json"
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

	var items []model.Item
	var targetedMuscleString string

	for rows.Next() {
		var item model.Item
		err := rows.Scan(
			&item.ItemID,
			&item.ItemType,
			&item.ItemName,
			&item.Difficulty,
			&item.Minutes,
			&item.CaloriesBurned,
			&targetedMuscleString,
		)

		if err != nil {
			log.Fatal(err)
			return
		}

		jsonConversionErr := json.Unmarshal([]byte(targetedMuscleString), &item.TargetedMuscleGroup)
		if jsonConversionErr != nil {
			fmt.Println("Error converting JSON: ", err)
			return
		}

		fmt.Println(item)
		items = append(items, item)
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"items": items,
	})
}
