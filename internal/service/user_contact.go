package service

import (
	"context"
	"fmt"
	"git.tarlanpayments.kz/pkg/golog"
	"go.opentelemetry.io/otel/trace"
	"telephone/internal/config"
	"telephone/internal/model"
	"telephone/internal/repository"
	"telephone/internal/schema"
)

type UserContactsService struct {
	repos  *repository.Repositories
	logger golog.ContextLogger
	tr     trace.Tracer
	cfg    *config.Config
}

func NewUserContactService(
	repos *repository.Repositories,
	logger golog.ContextLogger,
	tr trace.Tracer,
	cfg *config.Config) *UserContactsService {
	return &UserContactsService{
		repos:  repos,
		logger: logger,
		tr:     tr,
		cfg:    cfg,
	}
}

func (s UserContactsService) Create(ctx context.Context, userId int, phone string) error {
	return s.repos.UserContactRepository.AddUserContact(ctx, userId, phone)
}

func (s UserContactsService) AddUserContact(
	ctx context.Context,
	userID int, phone string, contact *schema.Contact) error {
	var _ error
	_, err := s.repos.ContactRepo.GetContactByPhone(ctx, phone)
	if err != nil {
		_ = s.repos.ContactRepo.CreateContact(ctx, schema.NewCreateContactRequest(contact))
	}

	return s.repos.UserContactRepository.AddUserContact(ctx, userID, phone)
}

func (s UserContactsService) GetUserContacts(ctx context.Context, userID int) ([]model.UserContacts, error) {
	user, err := s.repos.UserRepo.GetUserByID(ctx, userID)
	if err != nil {
		_ = fmt.Errorf("user not found. Impossible to retrieve a list of contacts")
	}

	res := user.UserContacts
	return res, err
}
