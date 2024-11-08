package service

import (
	"context"
	"go.opentelemetry.io/otel/trace"
	"telephone/internal/config"
	"telephone/internal/model"
	"telephone/internal/repository"
	"telephone/internal/schema"
)

type TelephoneService struct {
	repos *repository.Repositories
	tr    trace.Tracer
	cfg   *config.Config
}

func (t TelephoneService) CreateUser(ctx context.Context, user *schema.User) (_ *model.User, err error) {
	return t.repos.UserRepo.CreateUser(ctx, schema.NewCreateUserRequest(user))
}

func (t TelephoneService) GetUserByID(ctx context.Context, id uint) (schema.User, error) {
	data, err := t.repos.UserRepo.GetUserByID(ctx, id)
	if err != nil {
		return schema.User{}, err
	}
	res := schema.User{
		ID:    data.ID,
		Name:  data.Name,
		Email: data.Email,
	}
	return res, err
}

func (t TelephoneService) GetAllUsers(ctx context.Context,
) ([]model.User, error) {
	return t.repos.UserRepo.GetAllUsers()
}

func (t TelephoneService) GetAllUsersWithPaginationAndFiltering(
	limit int, page int,
	sort string,
	filter map[string]interface{},
	direction string,
) ([]model.User, *model.Paginate, error) {
	resp, pagination, err := t.repos.UserRepo.GetAllUsersWithPaginationAndFiltering(
		limit, page,
		sort,
		filter,
		direction)

	return resp, pagination, err
}

func NewTelephoneService(
	repos *repository.Repositories,
	tr trace.Tracer,
	cfg *config.Config,
) *TelephoneService {
	return &TelephoneService{
		repos: repos,
		tr:    tr,
		cfg:   cfg,
	}
}
