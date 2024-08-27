package repository

import (
	"gorm.io/gorm"
	"user-service/models"
)

type UserRepository interface {
	Create(user *models.User) error
	FindAll(queryParams map[string]string) ([]models.User, error)
	FindById(id uint) (*models.User, error)
	FindByEmailId(emailId string) (*models.User, error)
	Update(user *models.User) error
	Delete(user *models.User) error
}

type UserRepositoryImpl struct {
	db *gorm.DB
}

func (u UserRepositoryImpl) Create(user *models.User) error {
	return u.db.Create(user).Error
}

func (u UserRepositoryImpl) FindAll(queryParams map[string]string) ([]models.User, error) {
	var users []models.User
	filterQuery := u.db

	// Filtering/Search by user name
	if name, ok := queryParams["firstname"]; ok && name != "" {
		filterQuery = filterQuery.Where("firstname LIKE ?", "%"+name+"%")
	}
	if lastname, ok := queryParams["lastname"]; ok && lastname != "" {
		filterQuery = filterQuery.Where("lastname LIKE ?", "%"+lastname+"%")
	}

	// Sort by a specific field
	if sortBy, ok := queryParams["sortBy"]; ok && sortBy != "" {
		sortOrder := "asc"
		if sortDir, ok := queryParams["sortDir"]; ok && sortDir == "desc" {
			sortOrder = "desc"
		}
		filterQuery = filterQuery.Order(sortBy + " " + sortOrder)
	}

	err := filterQuery.Find(&users).Error
	return users, err
}

func (u UserRepositoryImpl) FindById(id uint) (*models.User, error) {
	var user models.User
	err := u.db.First(&user, id).Error
	return &user, err
}

func (u UserRepositoryImpl) FindByEmailId(emailId string) (*models.User, error) {
	var user models.User
	err := u.db.Where("email_id = ?", emailId).First(&user).Error
	return &user, err
}

func (u UserRepositoryImpl) Update(user *models.User) error {
	return u.db.Save(user).Error
}

func (u UserRepositoryImpl) Delete(user *models.User) error {
	return u.db.Delete(user).Error
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{db}
}
