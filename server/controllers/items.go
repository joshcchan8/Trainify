package controllers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/trainify/database"
	"github.com/trainify/middleware"
	models "github.com/trainify/models"
)

// HELPERS
type input struct {
	Difficulty           string
	Minutes              int
	TargetedMuscleGroups []string
}

// Helper to setup inputs (specs and profile) for the API call
func setupInputs(item models.Item, userID int) (input, models.UserProfile) {

	// Generate Workout
	specifications := input{
		Difficulty:           item.Difficulty,
		Minutes:              item.Minutes,
		TargetedMuscleGroups: item.TargetedMuscleGroups,
	}

	var profileID int
	var profile models.UserProfile

	row1 := database.DB.QueryRow(`SELECT profile_id FROM users WHERE user_id=?`, userID)
	scanErr1 := row1.Scan(&profileID)
	if scanErr1 != nil {
		log.Fatal("Scan Error: ", scanErr1)
	}

	row2 := database.DB.QueryRow(`SELECT * FROM user_profiles WHERE profile_id=?`, profileID)
	scanErr2 := row2.Scan(
		&profile.ProfileID,
		&profile.Age,
		&profile.Weight,
		&profile.Height,
		&profile.MaxPushUps,
		&profile.AvgPushUps,
		&profile.MaxPullUps,
		&profile.AvgPullUps,
		&profile.MaxSquat,
		&profile.AvgSquat,
		&profile.MaxBench,
		&profile.AvgBench,
		&profile.CardioLevel,
	)

	if scanErr2 != nil {
		log.Fatal("Scan Error: ", scanErr2)
	}

	return specifications, profile
}

// Helper to extract the description and calories burned from the content string
func extractData(content string) (string, int) {

	var desc strings.Builder
	var cals int

	scanner := bufio.NewScanner(strings.NewReader(content))
	calPattern := `Estimated Calories Burned: (\d+) calories`
	descPattern := `^\d+`

	regex := regexp.MustCompile(calPattern)

	for scanner.Scan() {
		line := scanner.Text()

		if cals == 0 {
			calMatch := regex.FindStringSubmatch(line)
			if len(calMatch) > 1 {
				cals, _ = strconv.Atoi(calMatch[1])
			}
		}

		if regexp.MustCompile(descPattern).MatchString(line) {
			desc.WriteString(line)
			desc.WriteString("\n")
		}
	}

	if scanner.Err() != nil {
		log.Fatal("Error occurred during content scanning: ", scanner.Err())
	}

	if cals == 0 {
		log.Fatal("No Calories Burned Provided")
	}

	return desc.String(), cals
}

// Helper to get the userID from the payload
func unloadPayload(c *gin.Context) int {

	// Get the custom request from the context
	customReq, ok := c.Get("custom_request")
	if !ok {
		log.Fatal("Request Error: could not get custom request")
		return 0
	}

	req, ok := customReq.(*middleware.CustomRequest)
	if !ok {
		log.Fatal("Request Error: could not convert custom request")
	}

	// Get payload
	payload := req.Payload
	if payload == nil {
		log.Fatal("Payload Error: could not get payload from request")
	}

	// Get userID that created the item
	userID := int(payload["data"].(map[string]interface{})["user_id"].(float64))
	return userID
}

// Helper to detect invalid muscle group errors
func muscleGroupsValidator(muscleGroups []string) {
	muscleGroupErr := models.ValidateMuscleGroups(muscleGroups)
	if muscleGroupErr != nil {
		log.Fatal("Muscle Group Error: ", muscleGroupErr)
	}
}

// Helper to convert from string array to json data
func itemToJson(itemsArray []string) []byte {
	jsonData, jsonConversionErr := json.Marshal(itemsArray)
	if jsonConversionErr != nil {
		log.Fatal("JSON Conversion Error: ", jsonConversionErr)
	}
	return jsonData
}

// Helper to convert from json string to string array
func jsonToItem(jsonString string, itemAddress *[]string) {
	jsonConversionErr := json.Unmarshal([]byte(jsonString), itemAddress)
	if jsonConversionErr != nil {
		log.Fatal("JSON Conversion Error: ", jsonConversionErr)
	}
}

// Helper for GetItem(), UpdateItem(), and DeleteItem() function
func getItemById(id string, userID int) models.Item {
	var item models.Item
	var targetedMuscleString string

	row := database.DB.QueryRow("SELECT * FROM items WHERE item_id=? AND created_by=?", id, userID)
	scanErr := row.Scan(
		&item.ItemID,
		&item.ItemName,
		&item.Difficulty,
		&item.Minutes,
		&item.CaloriesBurned,
		&targetedMuscleString,
		&item.WorkoutDescription,
		&item.CreatedBy,
	)

	if scanErr != nil {
		log.Fatal("Scan Error: ", scanErr)
	}

	jsonToItem(targetedMuscleString, &item.TargetedMuscleGroups)
	return item
}

// CONTROLLERS

// Gets all the items from the DB for a user
func GetAllItems(c *gin.Context) {

	userID := unloadPayload(c)
	rows, queryErr := database.DB.Query("SELECT * FROM items WHERE created_by=?", userID)

	if queryErr != nil {
		log.Fatal("Query Error: ", queryErr)
	}

	var items []models.Item
	var targetedMuscleString string

	for rows.Next() {
		var item models.Item
		scanErr := rows.Scan(
			&item.ItemID,
			&item.ItemName,
			&item.Difficulty,
			&item.Minutes,
			&item.CaloriesBurned,
			&targetedMuscleString,
			&item.WorkoutDescription,
			&item.CreatedBy,
		)

		if scanErr != nil {
			log.Fatal("Scan Error: ", scanErr)
		}

		jsonToItem(targetedMuscleString, &item.TargetedMuscleGroups)
		items = append(items, item)
	}

	c.IndentedJSON(http.StatusOK, gin.H{"items": items})
	fmt.Print("Items read successfully")
}

// Creates a new item in the DB for a specific user
func CreateItem(c *gin.Context) {

	var newItem models.Item
	userID := unloadPayload(c)

	bindErr := c.BindJSON(&newItem)
	if bindErr != nil {
		log.Fatal("Bind Error: ", bindErr)
	}

	muscleGroupsValidator(newItem.TargetedMuscleGroups)

	// Get inputs
	specifications, profile := setupInputs(newItem, userID)

	// OpenAI Response:
	data := getGeneration(c, specifications, profile)
	newItem.WorkoutDescription, newItem.CaloriesBurned = extractData(data.Choices[0].Message.Content)

	// item_id auto-generated
	stmt, insertionErr := database.DB.Prepare(
		`INSERT INTO items (item_name, difficulty, minutes, calories_burned, targeted_muscle_groups, workout_description, created_by) 
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
	)
	if insertionErr != nil {
		log.Fatal("Insertion Error: ", insertionErr)
	}

	defer stmt.Close()

	jsonData := itemToJson(newItem.TargetedMuscleGroups)
	result, executionErr := stmt.Exec(
		newItem.ItemName,
		newItem.Difficulty,
		newItem.Minutes,
		newItem.CaloriesBurned,
		jsonData,
		newItem.WorkoutDescription,
		userID,
	)

	if executionErr != nil {
		log.Fatal("Execution Error: ", executionErr)
	}

	// Attach userID to user object so we can return in JSON
	itemID, idErr := result.LastInsertId()
	if idErr != nil {
		log.Fatal("Error Fetching Item ID: ", idErr)
	}

	newItem.ItemID = int(itemID)
	newItem.CreatedBy = userID

	c.IndentedJSON(http.StatusCreated, gin.H{"item": newItem})
	fmt.Println("Item created successfully")
}

// Gets item with specific ID from DB
func GetItem(c *gin.Context) {
	id := c.Param("id")
	userID := unloadPayload(c)
	item := getItemById(id, userID)

	c.IndentedJSON(http.StatusOK, gin.H{"item": item})
	fmt.Println("Item read successfully")
}

// Updates item_name, difficulty, minutes, calories_burned, and targeted muscle groups
func UpdateItem(c *gin.Context) {
	id := c.Param("id")
	userID := unloadPayload(c)
	updatedItem := getItemById(id, userID)

	bindErr := c.ShouldBindJSON(&updatedItem)
	if bindErr != nil {
		log.Fatal("Bind Error: ", bindErr)
	}

	muscleGroupsValidator(updatedItem.TargetedMuscleGroups)

	// cannot update item_id
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

	jsonData := itemToJson(updatedItem.TargetedMuscleGroups)
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

	c.IndentedJSON(http.StatusOK, gin.H{"updated": updatedItem})
	fmt.Println("Item updated successfully")
}

// Regenerates workout_description and calories_burned, leaves the rest the same. (do not accept workout_description and calories_burn as params)
func RegenerateItem(c *gin.Context) {

	id := c.Param("id")
	userID := unloadPayload(c)
	regeneratedItem := getItemById(id, userID)

	bindErr := c.ShouldBindJSON(&regeneratedItem)
	if bindErr != nil {
		log.Fatal("Bind Error: ", bindErr)
	}

	muscleGroupsValidator(regeneratedItem.TargetedMuscleGroups)

	// Get inputs
	specifications, profile := setupInputs(regeneratedItem, userID)

	// OpenAI Response:
	data := getGeneration(c, specifications, profile)
	regeneratedItem.WorkoutDescription, regeneratedItem.CaloriesBurned = extractData(data.Choices[0].Message.Content)

	// item_id auto-generated
	stmt, updateErr := database.DB.Prepare(
		`UPDATE items 
		SET item_name=?,
			difficulty=?,
			minutes=?,
			calories_burned=?,
			targeted_muscle_groups=?, 
			workout_description=?,
			created_by=? 
		WHERE item_id=?`,
	)
	if updateErr != nil {
		log.Fatal("Insertion Error: ", updateErr)
	}

	defer stmt.Close()

	jsonData := itemToJson(regeneratedItem.TargetedMuscleGroups)
	_, executionErr := stmt.Exec(
		regeneratedItem.ItemName,
		regeneratedItem.Difficulty,
		regeneratedItem.Minutes,
		regeneratedItem.CaloriesBurned,
		jsonData,
		regeneratedItem.WorkoutDescription,
		userID,
		regeneratedItem.ItemID,
	)

	if executionErr != nil {
		log.Fatal("Execution Error: ", executionErr)
	}

	c.IndentedJSON(http.StatusOK, gin.H{"regenerated": regeneratedItem})
	fmt.Println("Item regenerated successfully")
}

func DeleteItem(c *gin.Context) {
	id := c.Param("id")
	userID := unloadPayload(c)
	deletedItem := getItemById(id, userID)

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

	c.IndentedJSON(http.StatusOK, gin.H{"deleted": deletedItem})
	fmt.Println("Item deleted successfully")
}
