package model

type itemType string
type muscle string

const (
	Workout  itemType = "workout"
	Schedule itemType = "schedule"
)

const (
	Chest      muscle = "chest"
	Back       muscle = "back"
	Shoulders  muscle = "shoulders"
	Triceps    muscle = "triceps"
	Biceps     muscle = "biceps"
	UpperAbs   muscle = "upper abs"
	LowerAbs   muscle = "lower abs"
	Obliques   muscle = "obliques"
	Quadriceps muscle = "quadriceps"
	Hamstrings muscle = "hamstrings"
	Hips       muscle = "hips"
	Glutes     muscle = "glutes"
	Calves     muscle = "calves"
)

type Item struct {
	ItemID              string   `json:"item_id"`
	ItemType            itemType `json:"item_type"`
	ItemName            string   `json:"item_name"`
	Difficulty          string   `json:"difficulty"`
	Minutes             int      `json:"minutes"`
	CaloriesBurned      int      `json:"calories_burned"`
	TargetedMuscleGroup []muscle `json:"targeted_muscle_group"`
}
