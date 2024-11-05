package model

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	ID        int       `json:"id" gorm:"column:id,primaryKey;autoIncrement"`
	Name      string    `json:"name" gorm:"column:name"`
	Email     string    `json:"email" gorm:"column:email"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
}

type Contact struct {
	gorm.Model
	ContactID   int    `json:"id" gorm:"column:contact_id,primaryKey"`
	ContactName string `json:"contact_name" gorm:"column:contact_name"`
	PhoneNumber string `json:"phone_number" gorm:"column:phone_number;size:11"`
}

type UserContactRelation struct {
	gorm.Model
	UserContactsID int  `json:"user_contacts_id" gorm:"column:contacts_id,primaryKey"`
	IsFavorite     bool `json:"is_favorite" gorm:"column:is_favorite"`
	ContactID      int  `json:"contact_id" gorm:"column:contact_id"`
	UserID         int  `json:"id" gorm:"column:id"`
}
