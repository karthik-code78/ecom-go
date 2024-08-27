package models

type Cart struct {
	ID     uint          `json:"id" gorm:"unique;primaryKey;autoIncrement"`
	Name   string        `json:"name"`
	UserID uint          `json:"userId" gorm:"unique;not null"`
	Value  float64       `json:"value"`
	Items  []CartProduct `json:"items" gorm:"foreignKey:CartID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
