package schema

import (
	"telephone/internal/model"
	pb "telephone/internal/proto"
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
		Name:  req.Name,
		Email: req.Email,
	}
}

func NewFromProtoToModelUserRequest(req *pb.CreateUserRequest) *User {
	return &User{
		ID:    int(req.User.UserId),
		Name:  req.User.Name,
		Email: req.User.Email,
	}
}
