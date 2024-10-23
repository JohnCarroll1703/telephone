package schema

import "telephone/internal/model"

type Contact struct {
	ContactID   int    `json:"id"`
	ContactName string `json:"contact_name"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
}

func NewCreateContactRequest(req *Contact) *model.Contact {
	return &model.Contact{
		ContactID:   req.ContactID,
		ContactName: req.ContactName,
		PhoneNumber: req.PhoneNumber,
	}
}

func NewUpdateContactRequest(req *Contact) *model.Contact {
	return &model.Contact{
		ContactID:   req.ContactID,
		ContactName: req.ContactName,
		PhoneNumber: req.PhoneNumber,
	}
}
