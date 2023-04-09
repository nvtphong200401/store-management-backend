package main

import (
	"net/http"

	"github.com/nvtphong200401/store-management/pkg/handlers"
	"github.com/nvtphong200401/store-management/pkg/models"
)

func init() {
	models.SetUp()
}

func main() {
	routersInit := handlers.InitRouter()
	server := &http.Server{
		Addr:    ":8080",
		Handler: routersInit,
	}
	server.ListenAndServe()
	models.CloseDB()
}
