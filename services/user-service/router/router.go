package router

import (
	"github.com/gorilla/mux"
	"github.com/karthik-code78/ecom/shared/logging"
	"user-service/handlers"
)

func InitRouter() *mux.Router {
	router := mux.NewRouter()

	// middlewares
	router.Use(logging.Logger)

	// Books router
	usersRouter := router.PathPrefix("/users").Subrouter()
	//usersRouter.Use(auth.Authenticate)

	// Books routes - CRUD
	usersRouter.HandleFunc("", handlers.GetAllUsers).Methods("GET")
	usersRouter.HandleFunc("/{id}", handlers.GetUser).Methods("GET")
	usersRouter.HandleFunc("", handlers.CreateUser).Methods("POST")
	usersRouter.HandleFunc("/{id}", handlers.UpdateUser).Methods("PUT")
	usersRouter.HandleFunc("/{id}", handlers.DeleteUser).Methods("DELETE")

	return router
}
