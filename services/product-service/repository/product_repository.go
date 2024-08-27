package repository

import (
	"gorm.io/gorm"
	"product-service/models"
)

type ProductRepository interface {
	Create(product *models.Product) error
	FindAll(queryParams map[string]string) ([]models.Product, error)
	FindById(id uint) (*models.Product, error)
	FindByIds(ids []uint) ([]models.Product, error)
	FindByIdAndQuantity(id uint, qty uint) (*models.Product, error)
	Update(product *models.Product) error
	Delete(product *models.Product) error
}

type ProductRepositoryImpl struct {
	db *gorm.DB
}

func (p *ProductRepositoryImpl) Create(product *models.Product) error {
	return p.db.Create(product).Error
}

func (p *ProductRepositoryImpl) FindAll(queryParams map[string]string) ([]models.Product, error) {
	var products []models.Product
	filterQuery := p.db

	if name, ok := queryParams["name"]; ok && name != "" {
		filterQuery = filterQuery.Where("name LIKE ?", "%"+name+"%")
	}

	if description, ok := queryParams["description"]; ok && description != "" {
		filterQuery = filterQuery.Where("description LIKE ?", "%"+description+"%")
	}

	if minPrice, ok := queryParams["minPrice"]; ok && minPrice != "" {
		filterQuery = filterQuery.Where("price >= ?", minPrice)
	}
	if maxPrice, ok := queryParams["maxPrice"]; ok && maxPrice != "" {
		filterQuery = filterQuery.Where("price <= ?", maxPrice)
	}

	if minQuantity, ok := queryParams["minQuantity"]; ok && minQuantity != "" {
		filterQuery = filterQuery.Where("quantity >= ?", minQuantity)
	}
	if maxQuantity, ok := queryParams["maxQuantity"]; ok && maxQuantity != "" {
		filterQuery = filterQuery.Where("quantity <= ?", maxQuantity)
	}

	if sortBy, ok := queryParams["sortBy"]; ok && sortBy != "" {
		sortOrder := "asc"
		if sortDir, ok := queryParams["sortDir"]; ok && sortDir == "desc" {
			sortOrder = "desc"
		}
		filterQuery = filterQuery.Order(sortBy + " " + sortOrder)
	}

	err := filterQuery.Find(&products).Error
	return products, err
}

func (p *ProductRepositoryImpl) FindById(id uint) (*models.Product, error) {
	var product models.Product
	err := p.db.First(&product, id).Error
	return &product, err
}

func (p *ProductRepositoryImpl) FindByIds(ids []uint) ([]models.Product, error) {
	var products []models.Product
	err := p.db.Where("id IN (?)", ids).Find(&products).Error
	return products, err
}

func (p *ProductRepositoryImpl) FindByIdAndQuantity(id uint, qty uint) (*models.Product, error) {
	var product models.Product
	err := p.db.Where("id = ? AND quantity >= ?", id, qty).First(&product).Error
	return &product, err
}

func (p *ProductRepositoryImpl) Update(product *models.Product) error {
	return p.db.Save(product).Error
}

func (p *ProductRepositoryImpl) Delete(product *models.Product) error {
	return p.db.Delete(product).Error
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &ProductRepositoryImpl{db}
}
