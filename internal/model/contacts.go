package model

import (
	"time"
)

type User struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"column:name"`
	Email     string    `json:"email" gorm:"column:email"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	DeletedAt time.Time `json:"deleted_at" gorm:"column:deleted_at"`
}

type Contact struct {
	ContactID   int       `json:"id" gorm:"primaryKey"`
	PhoneNumber string    `json:"phone_number" gorm:"column:phone_number;size:11"`
	DeletedAt   time.Time `json:"deleted_at" gorm:"column:deleted_at"`
}

type UserContactRelation struct {
	UserContactsID int  `json:"user_contacts_id" gorm:"primaryKey"`
	IsFavorite     bool `json:"is_favorite" gorm:"column:is_favorite"`
	ContactID      int  `json:"contact_id" gorm:"column:contact_id"`
	UserID         int  `json:"id" gorm:"column:id"`
}
