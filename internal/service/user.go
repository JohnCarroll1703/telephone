package service

import (
	"context"
	"git.tarlanpayments.kz/pkg/golog"
	"go.opentelemetry.io/otel/trace"
	"telephone/internal/config"
	"telephone/internal/model"
	"telephone/internal/repository"
	"telephone/internal/schema"
	"telephone/pkg/tracing"
)

type TelephoneService struct {
	repos  *repository.Repositories
	logger golog.ContextLogger
	tr     trace.Tracer
	cfg    *config.Config
}

func (t TelephoneService) CreateUser(ctx context.Context, user *schema.User) error {
	return t.repos.UserRepo.CreateUser(ctx, schema.NewCreateUserRequest(user))
}

func (t TelephoneService) GetUserByID(ctx context.Context, id int) (*model.User, error) {
	return t.repos.UserRepo.GetUserByID(ctx, id)
}

func (t TelephoneService) UpdateUser(ctx context.Context, user *schema.User) error {
	return t.repos.UserRepo.UpdateUser(ctx, schema.NewUpdateUserRequest(user))
}

func (t TelephoneService) DeleteUser(ctx context.Context, id int) error {
	return t.repos.UserRepo.DeleteUser(ctx, id)
}

func (t TelephoneService) GetAllUsers(ctx context.Context, tr trace.Tracer,
	funcName string) ([]*model.User, error) {
	ctx, span := tracing.CreateSpan(ctx, tr, "GetAllUsers_Func")
	defer span.End()
	return t.repos.UserRepo.GetAllUsers()
}

func NewTelephoneService(
	repos *repository.Repositories,
	logger golog.ContextLogger,
	tr trace.Tracer,
	cfg *config.Config,
) *TelephoneService {
	return &TelephoneService{
		repos:  repos,
		logger: logger,
		tr:     tr,
		cfg:    cfg,
	}
}
