package test

import (
	"fmt"
	"github/Eranmonnie/jqueue/jobs"
	"github/Eranmonnie/jqueue/worker"
	"testing"
)

func Hello(payload interface{}) error {
	fmt.Printf("%v!\n", payload)
	return nil
}

func TestSystem(T *testing.T) {
	// Initialize the server configuration
	// serverConfig := api.ServerConfig{
	// 	Port:      ":8080",
	// 	RedisPort: "6379",
	// 	Db: &database.DatabaseConfig{
	// 		Host:     "localhost",
	// 		Port:     "5432",
	// 		User:     "user",
	// 		Password: "password",
	// 		Name:     "jqueue_db",
	// 	},
	// }

	// // Start the server
	// if err := api.StartServer(serverConfig); err != nil {
	// 	T.Fatalf("Failed to start server: %v", err)
	// }

	// job := jobs.NewJob(1, "test_job", 1, 3, "Test payload")
	// if err := Hello(job); err != nil {
	// 	T.Fatalf("Failed to process job: %v", err)
	// }

	// Register the job handler
	if err := worker.Registry.RegisterHandler("test_job", Hello); err != nil {
		T.Fatalf("Failed to register job handler: %v", err)
	}

	//create job
	job := jobs.NewJob(1, "test_job", 1, 3, "Test payload")
	handler, err := worker.Registry.GetHandler(job.Name)
	if err != nil {
		T.Fatalf("Failed to get job handler: %v", err)
	}
	//testing without redis now cuz of no docker image yet lol
	if err := handler(job.Payload); err != nil {
		T.Fatalf("Failed to process job: %v", err)
	}
}
