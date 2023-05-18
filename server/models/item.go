package model

type workoutType string
type muscle string

const (
	Workout  workoutType = "workout"
	Schedule workoutType = "schedule"
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
	Type                workoutType `json:"type"`
	Name                string      `json:"name"`
	TargetedMuscleGroup []muscle    `json:"targeted muscle group"`
	Difficulty          int         `json:"difficulty"`
	Minutes             int         `json:"minutes"`
	CaloriesBurned      int         `json:"calories burned"`
}
