package models

import (
	"strings"
)

// User represents a user in database.
type User struct {
	ID        int    `json:"id"`
	Lastname  string `json:"lastname"`
	Firstname string `json:"firstname"`
}

// Fullname returns user fullname.
func (u User) Fullname() string {
	return strings.TrimSpace(u.Firstname + " " + u.Lastname)
}
