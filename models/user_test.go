package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFullname(t *testing.T) {
	u := User{
		ID:        1,
		Lastname:  "Bellanger",
		Firstname: "Fabien",
	}
	got := u.Fullname()
	expected := "Fabien Bellanger"
	assert.Equal(t, expected, got)
}

func TestFullnameWithoutLastname(t *testing.T) {
	u := User{
		ID:        1,
		Lastname:  "",
		Firstname: "Fabien",
	}
	got := u.Fullname()
	expected := "Fabien"
	assert.Equal(t, expected, got)
}

func TestFullnameEmpty(t *testing.T) {
	u := User{
		ID:        1,
		Lastname:  "",
		Firstname: "",
	}
	got := u.Fullname()
	expected := ""
	assert.Equal(t, expected, got)
}
