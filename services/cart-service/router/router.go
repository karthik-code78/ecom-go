package router

import (
	"cart-service/handlers"
	"github.com/gorilla/mux"
	"github.com/karthik-code78/ecom/shared/logging"
)

func InitRouter() *mux.Router {
	router := mux.NewRouter()

	// middlewares
	router.Use(logging.Logger)

	// cart router
	cartRouter := router.PathPrefix("/cart").Subrouter()
	//cartRouter.Use(auth.Authenticate)

	// cart routes - CRUD
	cartRouter.HandleFunc("", handlers.GetAllCarts).Methods("GET")
	cartRouter.HandleFunc("", handlers.CreateCart).Methods("POST")
	cartRouter.HandleFunc("/{id}", handlers.GetCartByID).Methods("GET")
	cartRouter.HandleFunc("/addProducts", handlers.AddProductsToCart).Methods("POST")
	cartRouter.HandleFunc("/withUserBody", handlers.CreateCartFromUser).Methods("POST")
	//cartsRouter.HandleFunc("/{id}", handlers.GetAdmin).Methods("GET")
	//cartsRouter.HandleFunc("", handlers.CreateAdmin).Methods("POST")
	//cartsRouter.HandleFunc("/{id}", handlers.UpdateAdmin).Methods("PUT")
	//cartsRouter.HandleFunc("/{id}", handlers.DeleteAdmin).Methods("DELETE")

	return router
}
