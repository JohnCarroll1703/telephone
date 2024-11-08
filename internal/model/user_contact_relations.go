package model

type UserContactRelation struct {
	UserContactsID uint `gorm:"primaryKey"`
	IsFavorite     bool `gorm:"column:is_fav"`
	ContactID      uint `gorm:"column:contact_id"`
	UserID         uint `gorm:"column:id"`
}
