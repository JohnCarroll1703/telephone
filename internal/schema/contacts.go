package schema

import (
	"telephone/internal/model"
	pb "telephone/internal/proto"
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
	ContactID   uint64 `json:"id"`
	PhoneNumber string `json:"phone_number"`
}

func NewCreateContactRequest(req *Contact) *model.Contact {
	return &model.Contact{
		ContactID:   int(req.ContactID),
		PhoneNumber: req.PhoneNumber,
	}
}

func NewFromProtoToModelCreateContactRequest(req *pb.CreateContactRequest) *Contact {
	return &Contact{
		PhoneNumber: req.Contact.Phone,
	}
}

type ContactFilter struct {
	CreatedAt   time.Time `json:"created_at" form:"created_at" time_format:"2003-03-17"`
	PhoneNumber string    `json:"phone_number"`
	SortBy      string    `json:"sort_by" validate:"omitempty"`
}
