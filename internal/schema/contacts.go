package schema

import (
	"telephone/internal/model"
	"time"
)

type Contact struct {
	ContactID   uint64 `json:"id"`
	ContactName string `json:"contact_name"`
	PhoneNumber string `json:"phone_number"`
}

type ContactRequest struct {
	ContactName string `json:"contact_name"`
	PhoneNumber string `json:"phone_number"`
}

type ContactResponse struct {
	ContactID   uint64 `json:"id"`
	ContactName string `json:"contact_name"`
	PhoneNumber string `json:"phone_number"`
}

type AddContactRequest struct {
	ContactName string `json:"contact_name"`
	PhoneNumber string `json:"phone_number"`
}

type AddContactResponse struct {
	ID          uint64 `json:"id"`
	ContactName string `json:"contact_name"`
	PhoneNumber string `json:"phone_number"`
}

func NewCreateContactRequest(req *Contact) *model.Contact {
	return &model.Contact{
		ContactName: req.ContactName,
		PhoneNumber: req.PhoneNumber,
	}
}

func NewContactRequest(req *ContactRequest) *model.Contact {
	return &model.Contact{
		ContactName: req.ContactName,
		PhoneNumber: req.PhoneNumber,
	}
}

func NewUpdateContactRequest(req *AddContactRequest) *model.Contact {
	return &model.Contact{
		ContactName: req.ContactName,
		PhoneNumber: req.PhoneNumber,
	}
}

type ContactFilter struct {
	CreatedAt   time.Time `json:"created_at" form:"created_at" time_format:"2003-03-17"`
	PhoneNumber string    `json:"phone_number"`
	SortBy      string    `json:"sort_by" validate:"omitempty"`
}
