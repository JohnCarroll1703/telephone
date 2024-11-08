package schema

import (
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"telephone/internal/model"
	pb "telephone/internal/proto"
	"time"
)

type User struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func (u User) Validate() error {
	if u.Name == "" {
		return status.Error(codes.InvalidArgument, "имя обязательно для создания пользователя!")
	}
	if err := validator.New().Struct(u); err != nil {
		return err
	}

	return nil
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
		ID:    uint(req.User.UserId),
		Name:  req.User.Name,
		Email: req.User.Email,
	}
}
