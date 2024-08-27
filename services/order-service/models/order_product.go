package models

type OrderProduct struct {
	ID        uint `json:"-" gorm:"primaryKey;autoIncrement"`
	OrderID   uint `json:"orderID" gorm:"not null"`
	ProductID uint `json:"productID" gorm:"not null"`
	Quantity  uint `json:"quantity" gorm:"not null"`
}
