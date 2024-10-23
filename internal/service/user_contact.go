package service

import (
	"context"
	"git.tarlanpayments.kz/pkg/golog"
	"go.opentelemetry.io/otel/trace"
	"telephone/internal/config"
	"telephone/internal/model"
	"telephone/internal/repository"
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

func (s UserContactsService) AddUserContact(ctx context.Context,
	userID int, contact model.Contact,
	isFav bool,
	userContact model.UserContacts) error {
	return s.repos.UserContactRepository.AddUserContact(ctx, userID, contact, isFav, userContact)
}

func (s UserContactsService) RemoveUserContact(ctx context.Context, userID int, contactID int) error {
	return s.repos.UserContactRepository.RemoveUserContact(ctx, userID, contactID)
}
