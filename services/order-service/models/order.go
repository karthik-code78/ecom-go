package models

type Order struct {
	ID     uint `json:"-" gorm:"primaryKey;autoIncrement"`
	CartID uint `json:"cartID" gorm:"default:0; not null"`
	UserID uint `json:"userID" gorm:"not null"`
	// 0 - Created, 1 - Accepted, 2 - Denied
	ApprovalStatus int            `json:"isApproved" gorm:"default:0; check:approval_status >= 0 AND approval_status < 3; not null"`
	OrderValue     float64        `json:"orderValue" gorm:"default:0.0; not null"`
	Items          []OrderProduct `json:"items" gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE;"`
}
