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

// Gets all the items from the DB for a user
func GetAllItems(c *gin.Context) {

	rows, queryErr := database.DB.Query("SELECT * FROM items")

	if queryErr != nil {
		log.Fatal("Query Error: ", queryErr)
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
			log.Fatal("JSON Conversion Error: ", err)
		}

		items = append(items, item)
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"items": items,
	})

	fmt.Print("Items read successfully")
}

// Creates a new item in the DB for a specific user
func CreateItem(c *gin.Context) {

	var newItem model.Item

	createItemErr := c.BindJSON(&newItem)
	if createItemErr != nil {
		log.Fatal("Creation error: ", createItemErr)
	}

	stmt, insertionErr := database.DB.Prepare("INSERT INTO items VALUES (?, ?, ?, ?, ?, ?, ?)")
	if insertionErr != nil {
		log.Fatal("Insertion error: ", insertionErr)
	}

	defer stmt.Close()

	jsonData, jsonConversionErr := json.Marshal(newItem.TargetedMuscleGroup)
	if jsonConversionErr != nil {
		log.Fatal("JSON Conversion Error: ", jsonConversionErr)
	}

	_, executionErr := stmt.Exec(
		newItem.ItemID,
		newItem.ItemType,
		newItem.ItemName,
		newItem.Difficulty,
		newItem.Minutes,
		newItem.CaloriesBurned,
		jsonData,
	)

	if executionErr != nil {
		log.Fatal("Execution Error: ", executionErr)
	}

	fmt.Println("Item created successfully")
}
