package model

import (
	"time"
)

type User struct {
	ID       string    `json:"id"`
	Username string    `json:"username"`
	DOB      time.Time `json:"dob"`
}
