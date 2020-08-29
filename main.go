package main

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

func main() {
	// Configuration initialization
	// ----------------------------
	if err := initConfig(); err != nil {
		log.Fatalln(err)
	}

	// Server creation
	// ---------------
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

// initConfig initializes configuration from config.toml file.
func initConfig() error {
	viper.SetConfigFile("config.toml")
	return viper.ReadInConfig()
}

// run launches a server instance.
func run() error {
	server := newServer()
	log.Printf("Server in %s mode\n", server.mode)

	// Database initialization
	// -----------------------
	server.store = &dbStore{}
	err := server.store.open()
	if err != nil {
		return err
	}
	defer server.store.close()

	// HTTP server initialization
	// --------------------------
	err = server.router.Listen(fmt.Sprintf("%v:%v",
		viper.GetString("server.host"),
		viper.GetString("server.port")))
	if err != nil {
		return err
	}

	return nil
}
