package handlers

import (
	"cart-service/repository"
	"github.com/karthik-code78/ecom/shared/logging"
	"gorm.io/gorm"
)

var db *gorm.DB

var CartRepo repository.CartRepository
var CartProductsRepo repository.CartProductRepository

func SetDatabase(database *gorm.DB) {
	logging.Log.Info(database)
	db = database

	CartRepo = repository.NewCartRepository(db)
	CartProductsRepo = repository.NewCartProductRepository(db)
}
