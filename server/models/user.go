package model

import (
	"time"
)

type user struct {
	ID       string    `json:"id"`
	Username string    `json:"username"`
	DOB      time.Time `json:"dob"`
}
