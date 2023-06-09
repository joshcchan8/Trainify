package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	models "github.com/trainify/models"
)

type OpenAIResponse struct {
	ID      string `json:"id"`
	Model   string `json:"model"`
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
		Index        int    `json:"index"`
	} `json:"choices"`
}

// Helper to create a profile string based on the given profile
func createProfileString(profile models.UserProfile) string {
	userProfile := `` + "\n"

	if profile.Age != nil {
		userProfile += (`Age: ` + strconv.Itoa(*profile.Age) + "\n")
	}

	if profile.Weight != nil {
		userProfile += (`Weight: ` + strconv.Itoa(*profile.Weight) + "\n")
	}

	if profile.Height != nil {
		userProfile += (`Height: ` + strconv.Itoa(*profile.Height) + "\n")
	}

	if profile.MaxPushUps != nil {
		userProfile += (`Max push-ups: ` + strconv.Itoa(*profile.MaxPushUps) + "\n")
	}

	if profile.AvgPushUps != nil {
		userProfile += (`Average push-ups: ` + strconv.Itoa(*profile.AvgPushUps) + "\n")
	}

	if profile.MaxPullUps != nil {
		userProfile += (`Max pull-ups: ` + strconv.Itoa(*profile.MaxPullUps) + "\n")
	}

	if profile.AvgPullUps != nil {
		userProfile += (`Average pull-ups: ` + strconv.Itoa(*profile.AvgPullUps) + "\n")
	}

	if profile.MaxSquat != nil {
		userProfile += (`Max squat weight: ` + strconv.Itoa(*profile.MaxSquat) + "\n")
	}

	if profile.AvgSquat != nil {
		userProfile += (`Average squat weight: ` + strconv.Itoa(*profile.AvgSquat) + "\n")
	}

	if profile.MaxBench != nil {
		userProfile += (`Max benchpress weight: ` + strconv.Itoa(*profile.MaxBench) + "\n")
	}

	if profile.AvgBench != nil {
		userProfile += (`Average benchpress weight: ` + strconv.Itoa(*profile.AvgBench) + "\n")
	}

	if profile.CardioLevel != nil {
		userProfile += (`Cardio level (from 1-10): ` + strconv.Itoa(*profile.CardioLevel) + ".\n")
	}

	return userProfile
}

// Helper to make the openAI API request with the given request body and secret API key
func makeOpenAIRequest(secretKey, requestBody string) ([]byte, error) {
	client := &http.Client{}
	url := "https://api.openai.com/v1/chat/completions"
	payload := strings.NewReader(requestBody)

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+secretKey)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// Generates a workout
func getGeneration(c *gin.Context, specs input, profile models.UserProfile) OpenAIResponse {

	var data OpenAIResponse

	// Setup description of role

	roleDescription := `
		Generate a workout that suits the user and fits the provided workout description as best as possible.
		Each list item should have EXACTLY the following format: <exercise #>. <exercise name>: <# of sets> sets of <# of reps> reps at <weight if not bodyweight>.
		Reps and weight should be based on the user's capabilities from user profile (maxes and averages).
		At the end, list estimated calories burned EXACTLY as: Estimated Calories Burned: <# of calories> calories
	`
	sysContent, _ := json.Marshal(roleDescription)

	// Setup inputs and profile as string

	targetedMuscleString := ""
	for i := 0; i < len(specs.TargetedMuscleGroups); i++ {
		targetedMuscleString += (" " + specs.TargetedMuscleGroups[i])
	}

	workoutDescription := `
		Difficulty should be a ` + specs.Difficulty + ` out of 10,
		duration should be ` + strconv.Itoa(specs.Minutes) + ` minutes long, 
		and workout should target:` + targetedMuscleString + `.
	`
	userProfile := createProfileString(profile)

	combined := `
		User profile: ` + userProfile + `
		Workout Description: ` + workoutDescription + `
	`
	userContent, _ := json.Marshal(combined)

	// Create request

	requestBody := fmt.Sprintf(`{
		"model": "gpt-3.5-turbo",
		"messages": [
			{"role": "system", "content": %s},
			{"role": "user", "content": %s}
		],
		"max_tokens": 400
	}`, sysContent, userContent)

	dotenvErr := godotenv.Load(".env")
	if dotenvErr != nil {
		log.Fatal("Error Loading .env File: ", dotenvErr)
	}

	secretKey := os.Getenv("API_KEY")

	response, requestErr := makeOpenAIRequest(secretKey, requestBody)
	if requestErr != nil {
		log.Fatal("OpenAI Request Error:", requestErr)
	}

	jsonConversionErr := json.Unmarshal(response, &data)
	if jsonConversionErr != nil {
		log.Fatal("JSON Conversion Error: ", jsonConversionErr)
	}

	return data
}
