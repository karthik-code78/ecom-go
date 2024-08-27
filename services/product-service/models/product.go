package models

import "gorm.io/gorm"

// Model for Product

type Product struct {
	gorm.Model
	Name        string  `json:"name" gorm:"not null"`
	Description string  `json:"description"`
	Price       float64 `json:"price" gorm:"type:Decimal;not null"`
	Quantity    uint    `json:"quantity" gorm:"not null"`
}

type ProductIds struct {
	Ids []uint `json:"ids"`
}
