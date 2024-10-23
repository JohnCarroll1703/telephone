package schema

import "telephone/internal/model"

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func NewUpdateUserRequest(req *User) *model.User {
	return &model.User{
		ID:    req.ID,
		Name:  req.Name,
		Email: req.Email,
	}
}

func NewCreateUserRequest(req *User) *model.User {
	return &model.User{
		ID:    req.ID,
		Name:  req.Name,
		Email: req.Email,
	}
}
