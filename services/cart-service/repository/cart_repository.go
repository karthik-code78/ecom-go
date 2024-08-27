package repository

import (
	"cart-service/models"
	"gorm.io/gorm"
)

type CartRepository interface {
	FindAll() ([]models.Cart, error)
	FindById(id uint) (*models.Cart, error)
	Create(cart *models.Cart) error
	Update(cart *models.Cart) error
	Delete(cart *models.Cart) error
}

type CartRepositoryImpl struct {
	db *gorm.DB
}

func (c *CartRepositoryImpl) FindAll() ([]models.Cart, error) {
	var carts []models.Cart
	err := c.db.Preload("Items").Find(&carts).Error
	return carts, err
}

func (c *CartRepositoryImpl) FindById(id uint) (*models.Cart, error) {
	var cart models.Cart
	err := c.db.Preload("Items").First(&cart, id).Error
	return &cart, err
}

func (c *CartRepositoryImpl) Create(cart *models.Cart) error {
	return c.db.Create(cart).Error
}

func (c *CartRepositoryImpl) Update(cart *models.Cart) error {
	return c.db.Save(cart).Error
}

func (c *CartRepositoryImpl) Delete(cart *models.Cart) error {
	return c.db.Delete(cart).Error
}

func NewCartRepository(db *gorm.DB) CartRepository {
	return &CartRepositoryImpl{db}
}
