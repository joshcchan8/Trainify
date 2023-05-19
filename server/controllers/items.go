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
		scanErr := rows.Scan(
			&item.ItemID,
			&item.ItemType,
			&item.ItemName,
			&item.Difficulty,
			&item.Minutes,
			&item.CaloriesBurned,
			&targetedMuscleString,
		)

		if scanErr != nil {
			log.Fatal("Scan Error: ", scanErr)
		}

		jsonConversionErr := json.Unmarshal([]byte(targetedMuscleString), &item.TargetedMuscleGroup)
		if jsonConversionErr != nil {
			log.Fatal("JSON Conversion Error: ", jsonConversionErr)
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

	bindErr := c.BindJSON(&newItem)
	if bindErr != nil {
		log.Fatal("Bind Error: ", bindErr)
	}

	muscleGroupErr := model.ValidateMuscleGroups(newItem.TargetedMuscleGroup)
	if muscleGroupErr != nil {
		log.Fatal("Muscle Group Error: ", muscleGroupErr)
	}

	// item_id auto-generated
	stmt, insertionErr := database.DB.Prepare(
		`INSERT INTO items (item_type, item_name, difficulty, minutes, calories_burned, targeted_muscle_groups) 
		VALUES (?, ?, ?, ?, ?, ?)`,
	)
	if insertionErr != nil {
		log.Fatal("Insertion Error: ", insertionErr)
	}

	defer stmt.Close()

	jsonData, jsonConversionErr := json.Marshal(newItem.TargetedMuscleGroup)
	if jsonConversionErr != nil {
		log.Fatal("JSON Conversion Error: ", jsonConversionErr)
	}

	_, executionErr := stmt.Exec(
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

	c.IndentedJSON(http.StatusCreated, gin.H{
		"item": newItem,
	})
	fmt.Println("Item created successfully")
}

// Helper for GetItem() function
func getItemById(id string) model.Item {
	var item model.Item
	var targetedMuscleString string

	row := database.DB.QueryRow("SELECT * FROM items WHERE item_id=?", id)
	scanErr := row.Scan(
		&item.ItemID,
		&item.ItemType,
		&item.ItemName,
		&item.Difficulty,
		&item.Minutes,
		&item.CaloriesBurned,
		&targetedMuscleString,
	)

	if scanErr != nil {
		log.Fatal("Scan Error: ", scanErr)
	}

	jsonConversionErr := json.Unmarshal([]byte(targetedMuscleString), &item.TargetedMuscleGroup)
	if jsonConversionErr != nil {
		log.Fatal("JSON Conversion Error: ", jsonConversionErr)
	}

	return item
}

// Gets item with specific ID from DB
func GetItem(c *gin.Context) {
	id := c.Param("id")
	item := getItemById(id)

	c.IndentedJSON(http.StatusOK, gin.H{
		"item": item,
	})
	fmt.Println("Item read successfully")
}

// Updates item with specific ID in DB with new specified data
func UpdateItem(c *gin.Context) {
	id := c.Param("id")
	updatedItem := getItemById(id)

	bindErr := c.BindJSON(&updatedItem)
	if bindErr != nil {
		log.Fatal("Update Error: ", bindErr)
	}

	muscleGroupErr := model.ValidateMuscleGroups(updatedItem.TargetedMuscleGroup)
	if muscleGroupErr != nil {
		log.Fatal("Muscle Group Error: ", muscleGroupErr)
	}

	// cannot update item_id or item_type (workout cannot become schedule, and vice versa)
	stmt, updateErr := database.DB.Prepare(
		`UPDATE items
		SET item_name=?,
			difficulty=?,
			minutes=?,
			calories_burned=?,
			targeted_muscle_groups=?
		WHERE item_id=?`,
	)
	if updateErr != nil {
		log.Fatal("Insertion error: ", updateErr)
	}

	defer stmt.Close()

	jsonData, jsonConversionErr := json.Marshal(updatedItem.TargetedMuscleGroup)
	if jsonConversionErr != nil {
		log.Fatal("JSON Conversion Error: ", jsonConversionErr)
	}

	_, executionErr := stmt.Exec(
		updatedItem.ItemName,
		updatedItem.Difficulty,
		updatedItem.Minutes,
		updatedItem.CaloriesBurned,
		jsonData,
		updatedItem.ItemID,
	)

	if executionErr != nil {
		log.Fatal("Execution Error: ", executionErr)
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"updated": updatedItem,
	})
	fmt.Println("Item updated successfully")
}

func DeleteItem(c *gin.Context) {
	id := c.Param("id")
	deletedItem := getItemById(id)

	stmt, deleteErr := database.DB.Prepare(
		`DELETE FROM items
		WHERE item_id=?`,
	)
	if deleteErr != nil {
		log.Fatal("Deletion Error: ", deleteErr)
	}

	defer stmt.Close()

	_, executionErr := stmt.Exec(id)

	if executionErr != nil {
		log.Fatal("Execution Error: ", executionErr)
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"deleted": deletedItem,
	})
	fmt.Println("Item deleted successfully")
}
