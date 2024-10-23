package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID           int            `json:"id" gorm:"column:id,primaryKey;autoIncrement"`
	Name         string         `json:"name" gorm:"column:name"`
	Email        string         `json:"email" gorm:"column:email"`
	UserContacts []UserContacts `gorm:"foreignKey:UserID"`
}

type Contact struct {
	gorm.Model
	ContactID    int            `json:"id" gorm:"column:contact_id,primaryKey;autoIncrement"`
	ContactName  string         `json:"contact_name" gorm:"column:contact_name"`
	PhoneNumber  string         `json:"phone_number" gorm:"column:phone_number;size:11"`
	UserContacts []UserContacts `gorm:"foreignKey:UserID"`
}

type UserContacts struct {
	gorm.Model
	UserContactsID int     `json:"user_contacts_id" gorm:"column:contacts_id,primaryKey;autoIncrement"`
	IsFavorite     bool    `json:"is_favorite" gorm:"column:is_favorite"`
	ContactID      int     `json:"contact_id" gorm:"column:contact_id"`
	UserID         int     `json:"user_id" gorm:"column:user_id"`
	Contact        Contact `gorm:"foreignKey:ContactID;references:ContactID"`
	User           User    `gorm:"foreignKey:ID;references:ID"`
}
