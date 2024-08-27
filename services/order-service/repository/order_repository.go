package repository

import (
	"gorm.io/gorm"
	"order-service/models"
)

type OrderRepository interface {
	FindAll() ([]models.Order, error)
	FindById(id int) (*models.Order, error)
	FindByUserId(userId int) (*models.Order, error)
	Create(order *models.Order) error
	Update(order *models.Order) (*models.Order, error)
	Delete(order *models.Order) error
}

type OrderRepositoryImpl struct {
	db *gorm.DB
}

func (o OrderRepositoryImpl) FindAll() ([]models.Order, error) {
	var orders []models.Order
	err := o.db.Find(&orders).Error
	return orders, err
}

func (o OrderRepositoryImpl) FindById(id int) (*models.Order, error) {
	var order *models.Order
	err := o.db.First(&order, id).Error
	return order, err
}

func (o OrderRepositoryImpl) FindByUserId(userId int) (*models.Order, error) {
	var order *models.Order
	err := o.db.Where("user_id = ?", userId).First(&order).Error
	return order, err
}

func (o OrderRepositoryImpl) Create(order *models.Order) error {
	err := o.db.Create(&order).Error
	return err
}

func (o OrderRepositoryImpl) Update(order *models.Order) (*models.Order, error) {
	err := o.db.Save(&order).Error
	return order, err
}

func (o OrderRepositoryImpl) Delete(order *models.Order) error {
	err := o.db.Delete(&order).Error
	return err
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return OrderRepositoryImpl{db: db}
}
