package schema

import "gorm.io/gorm"

type UserContacts struct {
	gorm.Model
	UserContactsID int     `json:"user_contacts_id" gorm:"column:contacts_id,primaryKey;autoIncrement"`
	IsFavorite     bool    `json:"is_favorite" gorm:"column:is_favorite"`
	ContactID      int     `json:"contact_id" gorm:"column:contact_id"`
	UserID         int     `json:"user_id" gorm:"column:user_id"`
	Contact        Contact `gorm:"foreignKey:ContactID;references:ContactID"`
	User           User    `gorm:"foreignKey:ID;references:ID"`
}
