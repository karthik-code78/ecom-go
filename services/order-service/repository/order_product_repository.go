package repository

import (
	"gorm.io/gorm"
	"order-service/models"
)

type OrderProductRepository interface {
	Create(orderProduct models.OrderProduct) error
	CreateBatch(orderProducts []models.OrderProduct) error
	FindByOrderId(orderId string) ([]models.OrderProduct, error)
}

type OrderProductRepositoryImpl struct {
	db *gorm.DB
}

func (o OrderProductRepositoryImpl) Create(orderProduct models.OrderProduct) error {
	err := o.db.Create(&orderProduct).Error
	return err
}

func (o OrderProductRepositoryImpl) CreateBatch(orderProducts []models.OrderProduct) error {
	err := o.db.Create(&orderProducts).Error
	return err
}

func (o OrderProductRepositoryImpl) FindByOrderId(orderId string) ([]models.OrderProduct, error) {
	var orderProducts []models.OrderProduct
	err := o.db.Where("OrderID = ?", orderId).Find(&orderProducts).Error
	return orderProducts, err
}

func NewOrderProductRepository(db *gorm.DB) OrderProductRepository {
	return &OrderProductRepositoryImpl{db}
}
