package stores

import (
	"github.com/fabienbellanger/go-fiber/models"
	"github.com/jmoiron/sqlx"
)

type UserStore interface {
	GetUser() models.User
}

type DBUserStore struct {
	DB *sqlx.DB
}

func (s *DBUserStore) GetUser() models.User {
	return models.User{
		ID:        1,
		Lastname:  "Doe",
		Firstname: "John",
	}
}
