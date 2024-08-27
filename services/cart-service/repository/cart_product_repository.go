package repository

import (
	"cart-service/models"
	"gorm.io/gorm"
)

type CartProductRepository interface {
	Create(cartProduct *models.CartProduct) error
	FindAllByCartId(cartId int) ([]models.CartProduct, error)
	FindByCartIdAndProductId(cartId uint, productId uint) (*models.CartProduct, error)
	FindByProductIds(productIds []uint) ([]models.CartProduct, error)
	Update(cartProduct *models.CartProduct) error
	UpdateByProductQty(cartProduct *models.CartProduct, productQty int) error
}

type CartProductRepositoryImpl struct {
	db *gorm.DB
}

func (c CartProductRepositoryImpl) Create(cartProduct *models.CartProduct) error {
	return c.db.Create(&cartProduct).Error
}

func (c CartProductRepositoryImpl) FindAllByCartId(cartId int) ([]models.CartProduct, error) {
	var cartProducts []models.CartProduct
	err := c.db.Find(&cartProducts, "CartID=?", cartId).Error
	return cartProducts, err
}

func (c CartProductRepositoryImpl) FindByCartIdAndProductId(cartId uint, productId uint) (*models.CartProduct, error) {
	var cartProduct models.CartProduct
	err := c.db.First(&cartProduct, "CartID=? AND ProductID=?", cartId, productId).Error
	return &cartProduct, err
}

func (c CartProductRepositoryImpl) FindByProductIds(productIds []uint) ([]models.CartProduct, error) {
	var cartProducts []models.CartProduct
	err := c.db.Find(&cartProducts).Where("ProductID IN ?", productIds).Error
	return cartProducts, err
}

func (c CartProductRepositoryImpl) Update(cartProduct *models.CartProduct) error {
	return c.db.Save(&cartProduct).Error
}

func (c CartProductRepositoryImpl) UpdateByProductQty(cartProduct *models.CartProduct, productQty int) error {
	return c.db.Model(&cartProduct).Update("AvailableQty", productQty).Error
}

func NewCartProductRepository(db *gorm.DB) CartProductRepository {
	return &CartProductRepositoryImpl{db}
}
