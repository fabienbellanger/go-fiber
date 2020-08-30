package main

import (
	"github.com/jmoiron/sqlx"
)

type userStore interface {
	getUser() User
}

type dbUserStore struct {
	DB *sqlx.DB
}

func (s *dbUserStore) getUser() User {
	return User{
		ID:        1,
		Lastname:  "Bellanger",
		Firstname: "Fabien",
	}
}
