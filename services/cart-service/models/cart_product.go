package models

type CartProduct struct {
	ID        uint    `json:"-" gorm:"primaryKey;auto_increment"`
	CartID    uint    `json:"cartId"`
	Cart      Cart    `json:"-" gorm:"foreignKey:CartID;references:ID"`
	ProductId uint    `json:"productId"`
	Quantity  uint    `json:"quantity"`
	Value     float64 `json:"-"`
}
