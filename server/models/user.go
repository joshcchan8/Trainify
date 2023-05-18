package model

import (
	"time"
)

type User struct {
	UserID   string    `json:"user_id"`
	Username string    `json:"username"`
	DOB      time.Time `json:"dob"`
}
