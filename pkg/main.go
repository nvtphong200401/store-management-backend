package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/nvtphong200401/store-management/pkg/db"
	"github.com/nvtphong200401/store-management/pkg/handlers"
	"github.com/nvtphong200401/store-management/pkg/registry"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/nvtphong200401/store-management/pkg/docs"
)

//	@title			Swagger Example API
//	@version		1.0
//	@description	This is a sample server celler server.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:8080
//	@BasePath	/api/v1

//	@securityDefinitions.basic	BasicAuth

//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization
//	@description				Description for what is this security definition being used

func init() {
	cmd := exec.Command("bash", "cmd/import.sh")
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
		cmd = exec.Command("bash", "cmd/export.sh")
		err = cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
}

func main() {
	gormDB, err := db.ConnectPostgresDB("D:\\Documents\\store-management-backend\\production.env")
	rd := db.ConnectRedis("D:\\Documents\\store-management-backend\\production.env")
	txStore := db.NewTXStore(gormDB, rd)

	if err != nil {
		log.Panic(err)
		return
	}
	defer txStore.CloseStorage()
	r := registry.NewRegistry(&txStore)

	routersInit := handlers.InitRouter(r.NewAppController())
	txStore.MigrateUp()

	routersInit.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	server := &http.Server{
		Addr:    ":8080",
		Handler: routersInit,
	}

	server.ListenAndServe()

}
