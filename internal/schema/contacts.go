package schema

import (
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"telephone/internal/model"
	pb "telephone/internal/proto"
	"time"
)

type Contact struct {
	ContactID   uint64 `json:"id"`
	ContactName string `json:"contact_name"`
	PhoneNumber string `json:"phone_number"`
}

func (c Contact) Validate() error {
	if c.PhoneNumber == "" {
		return status.Error(codes.InvalidArgument, "телефон не может быть пустым")
	}

	if err := validator.New().Struct(c); err != nil {
		return err
	}

	return nil
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
	Id          uint   `json:"id"`
	PhoneNumber string `json:"phone_number"`
}

func (r AddContactRequest) Validate() error {
	if r.PhoneNumber == "" {
		return status.Error(codes.InvalidArgument, "а что ты собрался добавлять?")
	}
	//if r.Id == 0 {
	//	return status.Error(codes.InvalidArgument, "ээмм... а кому добавлять номер собрался то?")
	//}

	if err := validator.New().Struct(r); err != nil {
		return err
	}

	return nil
}

type AddContactResponse struct {
	ContactID   uint64 `json:"id"`
	PhoneNumber string `json:"phone_number"`
}

func NewCreateContactRequest(req *Contact) *model.Contact {
	return &model.Contact{
		ContactID:   uint(req.ContactID),
		PhoneNumber: req.PhoneNumber,
	}
}

func NewFromProtoToModelAddContactRequest(req *pb.AddContactRequest) *AddContactRequest {
	return &AddContactRequest{
		PhoneNumber: req.Phone,
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
