package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/nvtphong200401/store-management/pkg/handlers"
	"github.com/nvtphong200401/store-management/pkg/models"
)

func init() {
	models.SetUp()
}

func main() {
	cmd := exec.Command("sh", "cmd/import.sh")
	// Redirect the command's standard output and error to the current process's standard output and error
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	// Register a signal handler for interrupt signals
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-signalChannel
		fmt.Println("\nStopping Golang project...")
		// Execute your shell script here
		cmd = exec.Command("sh", "cmd/export.sh")
		err = cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()

	routersInit := handlers.InitRouter()
	server := &http.Server{
		Addr:    ":8080",
		Handler: routersInit,
	}
	server.ListenAndServe()
	models.CloseDB()

}
