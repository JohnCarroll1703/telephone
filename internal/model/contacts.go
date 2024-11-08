package model

import (
	"time"
)

type Contact struct {
	ContactID   uint      `gorm:"primaryKey"`
	PhoneNumber string    `gorm:"column:phone_number;size:11"`
	DeletedAt   time.Time `gorm:"column:deleted_at"`
}
