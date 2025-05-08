package api

import (
	"github/Eranmonnie/jqueue/database"
	"log"
	"net/http"
	"time"
)

type ServerConfig struct {
	Port string
	Db   *database.DatabaseConfig
}

func StartServer(config ServerConfig) error {
	//post to create a job
	// and get a job
	//and probably filters

	s := &http.Server{
		Addr:           config.Port,
		Handler:        nil,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if config.Db != nil {
		log.Println("Database configuration provided, initializing database connection...")
		database, err := database.NewGormDB(*config.Db)
		if err != nil {
			log.Fatalf("Failed to create database connection: %v", err)
		}
		if err := database.Connect(); err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}
		defer func() {
			if err := database.Close(); err != nil {
				log.Fatalf("Failed to close database connection: %v", err)
			}
		}()
	} else {
		log.Println("No database configuration provided")
	}

	log.Println("Server starting on port", config.Port)
	return s.ListenAndServe()
}
