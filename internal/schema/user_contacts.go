package schema

import (
	"gorm.io/gorm"
	pb "telephone/internal/proto"
)

type UserContactRelation struct {
	gorm.Model
	UserContactsID int     `json:"user_contacts_id" gorm:"column:contacts_id,primaryKey"`
	IsFavorite     bool    `json:"is_favorite" gorm:"column:is_favorite"`
	ContactID      int     `json:"contact_id" gorm:"column:contact_id"`
	UserID         int     `json:"user_id" gorm:"column:user_id"`
	Contact        Contact `gorm:"foreignKey:ContactID;references:ContactID"`
	User           User    `gorm:"foreignKey:ID;references:ID"`
}

func NewFromProtoToModelRelationRequest(req *pb.CreateUserContactRelationRequest) *UserContactRelation {
	return &UserContactRelation{
		UserContactsID: int(req.UserContact.UserContactsId),
		IsFavorite:     true,
		ContactID:      int(req.UserContact.ContactId),
		UserID:         int(req.UserContact.UserId),
		Contact: Contact{
			ContactID: uint64(req.UserContact.ContactId),
		},
		User: User{
			ID: int(req.UserContact.UserId),
		},
	}
}
