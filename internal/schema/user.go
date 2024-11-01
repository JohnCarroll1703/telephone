package schema

import (
	"telephone/internal/model"
	"time"
)

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func NewUpdateUserRequest(req *User) *model.User {
	return &model.User{
		ID:        req.ID,
		Name:      req.Name,
		Email:     req.Email,
		CreatedAt: time.Time{},
	}
}

func NewCreateUserRequest(req *User) *model.User {
	return &model.User{
		ID:    req.ID,
		Name:  req.Name,
		Email: req.Email,
	}
}
