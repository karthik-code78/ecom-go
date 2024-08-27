package handlers

import (
	"gorm.io/gorm"
	"order-service/repository"
)

var db *gorm.DB

var OrdersRepo repository.OrderRepository
var OrderProductsRepo repository.OrderProductRepository

func SetDatabase(database *gorm.DB) {
	db = database

	OrdersRepo = repository.NewOrderRepository(db)
	OrderProductsRepo = repository.NewOrderProductRepository(db)
}
