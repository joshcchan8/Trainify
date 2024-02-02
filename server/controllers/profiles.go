package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/trainify/database"
	models "github.com/trainify/models"
	"github.com/trainify/queries"
)

// Users only have a single profile

// HELPERS

func GetProfileByUserID(c *gin.Context) models.UserProfile {
	var profile models.UserProfile
	userID := unloadPayload(c)

	row1 := database.DB.QueryRow(queries.GetProfileIdFromUserIdQuery, userID)
	scanErr1 := row1.Scan(&profile.ProfileID)
	if scanErr1 != nil {
		log.Fatal("Scan Error: ", scanErr1)
	}

	row2 := database.DB.QueryRow(queries.GetProfileQuery, profile.ProfileID)
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

	return profile
}

// CONTROLLERS

// Gets the profile of the user currently logged in
func GetProfile(c *gin.Context) {

	var profile models.UserProfile = GetProfileByUserID(c)

	c.IndentedJSON(http.StatusOK, gin.H{
		"profile": profile,
	})
	fmt.Println("Profile read successfully")
}

// Updates the profile of the user currently logged in
func UpdateProfile(c *gin.Context) {

	var updatedProfile models.UserProfile = GetProfileByUserID(c)

	// Fix: binding data to pointers currently, need to convert JSON to addresses
	bindErr := c.ShouldBindJSON(&updatedProfile)
	if bindErr != nil {
		log.Fatal("Bind Error: ", bindErr)
	}

	stmt, updateErr := database.DB.Prepare(
		queries.UpdateProfileQuery,
	)

	if updateErr != nil {
		log.Fatal("Insertion error: ", updateErr)
	}

	defer stmt.Close()

	_, executionErr := stmt.Exec(
		updatedProfile.Age,
		updatedProfile.Weight,
		updatedProfile.Height,
		updatedProfile.MaxPushUps,
		updatedProfile.AvgPushUps,
		updatedProfile.MaxPullUps,
		updatedProfile.AvgPullUps,
		updatedProfile.MaxSquat,
		updatedProfile.AvgSquat,
		updatedProfile.MaxBench,
		updatedProfile.AvgBench,
		updatedProfile.CardioLevel,
		updatedProfile.ProfileID,
	)

	if executionErr != nil {
		log.Fatal("Execution Error: ", executionErr)
	}

	c.IndentedJSON(http.StatusOK, gin.H{"updated": updatedProfile})
	fmt.Println("Profile updated successfully")
}
