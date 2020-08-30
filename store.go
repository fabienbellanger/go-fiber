package main

import (
	"fmt"
	"log"

	"github.com/spf13/viper"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
	"github.com/jmoiron/sqlx"
)

type store interface {
	open() error
	close() error

	user() userStore
}

type dbStore struct {
	DB        *sqlx.DB
	userStore userStore
}

func (s *dbStore) open() error {
	db, err := sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		viper.GetString("database.username"),
		viper.GetString("database.password"),
		viper.GetString("database.host"),
		viper.GetInt("database.port"),
		viper.GetString("database.name")))
	if err != nil {
		return err
	}
	log.Printf("Connected to DB %s", viper.GetString("database.name"))
	s.DB = db

	// User store
	// ----------
	s.userStore = &dbUserStore{DB: db}

	return nil
}

func (s *dbStore) close() error {
	return s.DB.Close()
}

func (s *dbStore) user() userStore {
	return s.userStore
}
