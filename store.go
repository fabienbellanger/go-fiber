package main

import (
	"fmt"
	"log"

	"github.com/fabienbellanger/go-fiber/stores"
	"github.com/spf13/viper"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
	"github.com/jmoiron/sqlx"
)

type store interface {
	open() error
	close() error

	user() stores.UserStore
}

type dbStore struct {
	DB        *sqlx.DB
	userStore stores.UserStore
}

func (s *dbStore) open() error {
	log.Printf("%s\n", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		viper.GetString("database.username"),
		viper.GetString("database.password"),
		viper.GetString("database.host"),
		viper.GetInt("database.port"),
		viper.GetString("database.name")))
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
	s.userStore = &stores.DBUserStore{DB: db}

	return nil
}

func (s *dbStore) close() error {
	return s.DB.Close()
}

func (s *dbStore) user() stores.UserStore {
	return s.userStore
}
