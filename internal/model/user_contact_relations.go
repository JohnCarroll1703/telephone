package model

type UserContactRelation struct {
	UserContactsID uint `gorm:"primaryKey"`
	IsFavorite     bool `gorm:"column:is_fav"`
	ContactID      int  `gorm:"column:contact_id"`
	UserID         int  `gorm:"column:id"`
}
