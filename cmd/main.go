package cmd

import (
	"github/Eranmonnie/jqueue/api"
	"github/Eranmonnie/jqueue/database"
)

func main() {
	// Initialize the server configuration
	sConfig := api.ServerConfig{
		Port:      ":8080",
		RedisPort: "6379",
		Db: &database.DatabaseConfig{
			Host:     "localhost",
			User:     "user",
			Password: "password",
			Name:     "dbname",
			Port:     "5432",
			Type:     "postgres",
			Options: map[string]string{
				"sslmode": "disable",
			},
		},
	}
	// Start the server
	if err := api.StartServer(sConfig); err != nil {
		panic(err)
	}
	// Start the worker

	//want to do it ni a way that it can be run in a different process without user specifically hard coding it

}
