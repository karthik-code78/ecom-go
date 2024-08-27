package router

import (
	"github.com/gorilla/mux"
	"github.com/karthik-code78/ecom/shared/logging"
	"product-service/handlers"
)

func InitRouter() *mux.Router {
	router := mux.NewRouter()

	// middlewares
	router.Use(logging.Logger)

	// Books router
	productsRouter := router.PathPrefix("/products").Subrouter()
	//productsRouter.Use(auth.Authenticate)

	// Books routes - CRUD
	productsRouter.HandleFunc("", handlers.GetAllProducts).Methods("GET")
	productsRouter.HandleFunc("/{id}", handlers.GetProduct).Methods("GET")
	productsRouter.HandleFunc("/{id}/{qty}", handlers.GetProductByQuantity).Methods("GET")
	productsRouter.HandleFunc("", handlers.CreateProduct).Methods("POST")
	productsRouter.HandleFunc("/getByIds", handlers.GetProductsByIds).Methods("POST")
	productsRouter.HandleFunc("/{id}", handlers.UpdateProduct).Methods("PUT")
	productsRouter.HandleFunc("/{id}", handlers.DeleteProduct).Methods("DELETE")
	productsRouter.HandleFunc("/updateQty", handlers.UpdateProductQuantity).Methods("POST")

	return router
}
