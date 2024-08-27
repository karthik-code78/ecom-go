package main

import (
	"github.com/karthik-code78/ecom/shared/configure"
	"github.com/karthik-code78/ecom/shared/logging"
	"github.com/karthik-code78/ecom/shared/migration"
	"log"
	"net/http"
	"order-service/handlers"
	"order-service/models"
	"order-service/router"
	"os"
)

//var locaEnvPath =

func main() {
	log.Printf(os.Getwd())
	logging.Initializelogger()
	logging.Log.Info("lol")
	db, err := configure.ConnectAndReturnDB()
	if err != nil {
		logging.Log.Fatal("Failed to connect to the migration", err)
	}

	logging.Log.Info("db is", db)

	// Initialize tables
	migration.Migrate(&models.Order{}, &models.OrderProduct{})

	// Database for handlers
	handlers.SetDatabase(db)

	// Init router
	mainRouter := router.InitRouter()

	//Start server
	logging.Log.Info("Server is running on port : 8084")
	err = http.ListenAndServe(":8084", mainRouter)
	if err != nil {
		log.Fatal("Failed to start server", err)
	}
}
