package model

import (
	"errors"
)

type itemType string

const (
	Workout  itemType = "workout"
	Schedule itemType = "schedule"
)

type Item struct {
	ItemID              int      `json:"item_id"`
	ItemType            itemType `json:"item_type"`
	ItemName            string   `json:"item_name"`
	Difficulty          string   `json:"difficulty"`
	Minutes             int      `json:"minutes"`
	CaloriesBurned      int      `json:"calories_burned"`
	TargetedMuscleGroup []string `json:"targeted_muscle_group"`
}

// validates whether selected muscle groups are part of the valid list
func ValidateMuscleGroups(groups []string) error {

	muscleGroups := [13]string{
		"chest", "back", "shoulders", "triceps",
		"biceps", "upper abs", "lower abs", "obliques",
		"quadriceps", "hamstrings", "hips", "glutes", "calves",
	}
	muscleGroupErr := errors.New("invalid muscle group")

	// checks that given muscle groups match those in the array
	for _, muscleGroup := range groups {
		valid := false
		for _, element := range muscleGroups {
			if muscleGroup == element {
				valid = true
				break
			}
		}
		if !valid {
			return muscleGroupErr
		}
	}
	return nil
}
