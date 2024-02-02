package controllers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"github.com/trainify/database"
	models "github.com/trainify/models"
	"github.com/trainify/queries"
	"golang.org/x/crypto/bcrypt"
)

// HELPERS

// Creates an empty profile for a new user (where they can update later)
func createEmptyProfile() sql.Result {
	result, insertionErr := database.DB.Exec(
		queries.CreateEmptyProfileQuery,
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)

	if insertionErr != nil {
		log.Fatal("Insertion Error: ")
	}

	return result
}

// Create a token and attaches the user data payload to a signed token
func jwtPayloadHelper(user models.User) string {

	// Create token
	token := jwt.New(jwt.GetSigningMethod("HS256"))

	// Set claims and payload for token
	claims := token.Claims.(jwt.MapClaims)
	claims["data"] = user
	claims["exp"] = time.Now().Add(time.Hour * 3).Unix() // expires in 3 hours

	// Sign token
	dotenvErr := godotenv.Load(".env")
	if dotenvErr != nil {
		log.Fatal("Error Loading .env File: ", dotenvErr)
	}

	secretKey := []byte(os.Getenv("SECRET_KEY"))
	signedToken, tokenSignErr := token.SignedString(secretKey)
	if tokenSignErr != nil {
		log.Fatal("Token Sign Error: ", tokenSignErr)
	}

	// return signed token
	return signedToken
}

// CONTROLLERS

// Registers user with email, password, and username
func Register(c *gin.Context) {

	var newUser models.User

	bindErr := c.BindJSON(&newUser)
	if bindErr != nil {
		log.Fatal("Bind Error: ", bindErr)
	}

	registrationInfo := models.Info{
		Username: newUser.Username,
		Email:    newUser.Email,
		Password: newUser.Password,
	}

	// Need Username, Email and Password to register
	missingRegistrationInfoErr := models.RegistrationInfoProvided(registrationInfo)
	if missingRegistrationInfoErr != nil {
		log.Fatal("Missing Required Registration Info Error: ", missingRegistrationInfoErr)
	}

	// Verify that email is real
	invalidEmailErr := models.ValidateEmail(newUser.Email)
	if invalidEmailErr != nil {
		log.Fatal("Invalid Email Error: ", invalidEmailErr)
	}

	// Hash password
	hash, hashErr := bcrypt.GenerateFromPassword([]byte(newUser.Password), 10)
	if hashErr != nil {
		log.Fatal("Error Hashing Password: ", hashErr)
	}

	hashedPassword := string(hash)
	newUser.Password = hashedPassword

	// Create profile and get profile_id to add to the user
	profileID, idErr := createEmptyProfile().LastInsertId()
	if idErr != nil {
		log.Fatal("Error fetching profile ID: ", idErr)
	}
	newUser.ProfileID = int(profileID)

	// Register user in DB, no user_id needed
	stmt, registerErr := database.DB.Prepare(
		queries.CreateUserQuery,
	)

	if registerErr != nil {
		log.Fatal("Error Registering User: ", registerErr)
	}

	defer stmt.Close()

	result, executionErr := stmt.Exec(
		newUser.Username,
		newUser.Email,
		newUser.Password,
		newUser.ProfileID,
	)

	if executionErr != nil {
		log.Fatal("Execution Error: ", executionErr)
	}

	// Attach userID to user object so we can return in JSON
	userID, idErr := result.LastInsertId()
	if idErr != nil {
		log.Fatal("Error Fetching User ID: ", idErr)
	}
	newUser.UserID = int(userID)

	signedToken := jwtPayloadHelper(newUser)

	// Send back response
	c.IndentedJSON(http.StatusCreated, gin.H{
		"user":  newUser,
		"token": signedToken,
	})
	fmt.Println("Registered user successfully")
}

// Logs in a user based on email and password
func Login(c *gin.Context) {

	var loginInfo models.Info
	var user models.User

	bindErr := c.BindJSON(&loginInfo)
	if bindErr != nil {
		log.Fatal("Bind Error :", bindErr)
	}

	// Need Email and Password to login
	missingLoginInfoErr := models.LoginInfoProvided(loginInfo)
	if missingLoginInfoErr != nil {
		log.Fatal("Missing Required Login Info Error: ", missingLoginInfoErr)
	}

	// Retrieve Password from DB
	row := database.DB.QueryRow(queries.GetUserByEmailQuery, loginInfo.Email)
	scanErr := row.Scan(
		&user.UserID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.ProfileID,
	)

	if scanErr != nil {
		log.Fatal("Scan Error: ", scanErr)
	}

	// Verify Password
	invalidPasswordErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginInfo.Password))
	if invalidPasswordErr != nil {
		log.Fatal("Invalid Password Error: ", invalidPasswordErr)
	}

	signedToken := jwtPayloadHelper(user)

	c.IndentedJSON(http.StatusOK, gin.H{
		"user":  user,
		"token": signedToken,
	})
	fmt.Println("Logged in user successfully")
}

// Updates a users password/username (but not email, user_id, or profile_id)
func UpdateUser(c *gin.Context) {

	var updatedUser models.User
	userID := unloadPayload(c)

	bindErr := c.BindJSON(&updatedUser)
	if bindErr != nil {
		log.Fatal("Bind Error :", bindErr)
	}

	// Hash password
	hash, hashErr := bcrypt.GenerateFromPassword([]byte(updatedUser.Password), 10)
	if hashErr != nil {
		log.Fatal("Error Hashing Password: ", hashErr)
	}

	hashedPassword := string(hash)
	updatedUser.Password = hashedPassword

	// Update logged in user info with updated user info
	_, executionErr := database.DB.Exec(
		queries.UpdateUserQuery,
		updatedUser.Username,
		updatedUser.Password,
		userID,
	)

	if executionErr != nil {
		log.Fatal("Execution Error: ", executionErr)
	}

	// Get remaining data of user for JSON response
	row := database.DB.QueryRow(queries.GetUserQuery, userID)
	scanErr := row.Scan(
		&updatedUser.UserID,
		&updatedUser.Username,
		&updatedUser.Email,
		&updatedUser.Password,
		&updatedUser.ProfileID,
	)

	if scanErr != nil {
		log.Fatal("Scan Error: ", scanErr)
	}

	c.IndentedJSON(http.StatusOK, gin.H{"updated": updatedUser})
	fmt.Println("User updated successfully")
}

// Deletes a user from the DB
func DeleteUser(c *gin.Context) {

	var user models.User
	userID := unloadPayload(c)

	// Retrieve user info of logged in user from DB
	row := database.DB.QueryRow(queries.GetUserQuery, userID)
	scanErr := row.Scan(
		&user.UserID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.ProfileID,
	)

	if scanErr != nil {
		log.Fatal("Scan Error: ", scanErr)
	}

	// Prepare to delete associated profile and items
	deleteProfile, profileErr := database.DB.Prepare(queries.DeleteProfileQuery)
	deleteItems, itemsErr := database.DB.Prepare(queries.DeleteAllItemsQuery)
	deleteUser, userErr := database.DB.Prepare(queries.DeleteUserQuery)

	if profileErr != nil || itemsErr != nil || userErr != nil {
		log.Fatal("Statement Preparation Error")
	}

	defer deleteProfile.Close()
	defer deleteItems.Close()
	defer deleteUser.Close()

	// Delete items, then profiles, then the user from the DB

	_, itemExecutionErr := deleteItems.Exec(userID)
	_, userExecutionErr := deleteUser.Exec(userID)
	_, profileExecutionErr := deleteProfile.Exec(user.ProfileID)

	if itemExecutionErr != nil || userExecutionErr != nil || profileExecutionErr != nil {
		log.Fatal("Delete User Execution Error")
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"deleted_user": user,
	})
	fmt.Println("Deleted user successfully")
}
