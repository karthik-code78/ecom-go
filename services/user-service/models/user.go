package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Firstname string `json:"firstname" gorm:"not null"`
	Lastname  string `json:"lastname" gorm:"not null"`
	Email     string `json:"email" gorm:"unique;not null"`
	Password  string `json:"password" gorm:"type:binary(60);not null"`
}

type UserResponse struct {
	ID        uint   `json:"id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}
