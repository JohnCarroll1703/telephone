package schema

import (
	"github.com/go-playground/validator/v10"
	"time"
)

const (
	DefaultOrder = "id ASC"
)

type (
	UserAndContactFilter struct {
		CreatedAt   time.Time `json:"created_at" form:"created_at" time_format:"2003-03-17"`
		PhoneNumber string    `json:"phone_number"`
		Name        string    `json:"name"`
		SortBy      string    `json:"sort_by" validate:"omitempty"`
	}
)

func (uc *UserAndContactFilter) Validate() error {
	return validator.New().Struct(uc)
}
