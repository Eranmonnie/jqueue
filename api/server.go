package api

import (
	"fmt"
	"github/Eranmonnie/jqueue/config"
	"github/Eranmonnie/jqueue/database"
	"log"
	"net/http"
	"time"
)

type ServerConfig struct {
	Port      string
	RedisPort string
	Db        *database.DatabaseConfig
}

func StartServer(sConfig ServerConfig) error {
	//post to create a job
	// and get a job
	//and probably filters

	s := &http.Server{
		Addr:           sConfig.Port,
		Handler:        nil,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if sConfig.RedisPort != "" {
		log.Println("Redis configuration provided, initializing Redis connection...")
		config.InitRedis(sConfig.RedisPort)
	} else {
		return fmt.Errorf("no Redis configuration provided important for queue")
	}

	if sConfig.Db != nil {
		log.Println("Database configuration provided, initializing database connection...")
		database, err := database.NewGormDB(*sConfig.Db)
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

	log.Println("Server starting on port", sConfig.Port)
	return s.ListenAndServe()
}
