package router

import (
	"github.com/gorilla/mux"
	"github.com/karthik-code78/ecom/shared/logging"
	"order-service/handlers"
)

func InitRouter() *mux.Router {
	router := mux.NewRouter()

	// middlewares
	router.Use(logging.Logger)

	// Books router
	productsRouter := router.PathPrefix("/orders").Subrouter()
	//productsRouter.Use(auth.Authenticate)

	// Books routes - CRUD
	productsRouter.HandleFunc("", handlers.GetAllOrders).Methods("GET")
	productsRouter.HandleFunc("", handlers.CreateOrder).Methods("POST")

	return router
}
