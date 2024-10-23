package service

import (
	"context"
	"git.tarlanpayments.kz/pkg/golog"
	"go.opentelemetry.io/otel/trace"
	"telephone/internal/config"
	"telephone/internal/model"
	"telephone/internal/repository"
	"telephone/internal/schema"
)

type ContactService struct {
	repos  *repository.Repositories
	logger golog.ContextLogger
	tr     trace.Tracer
	cfg    *config.Config
}

func (c ContactService) GetContacts() ([]*model.Contact, error) {
	return c.repos.ContactRepo.GetAllContacts()
}

func (c ContactService) CreateContact(ctx context.Context, contact *schema.Contact) error {
	return c.repos.ContactRepo.CreateContact(ctx, schema.NewCreateContactRequest(contact))
}

func (c ContactService) GetContactByID(ctx context.Context, id int) (*model.Contact, error) {
	return c.repos.ContactRepo.GetContactByID(ctx, id)
}

func (c ContactService) UpdateContact(ctx context.Context, id int) (*model.Contact, error) {
	//TODO implement me

	panic("implement me")
}

func (c ContactService) DeleteContact(ctx context.Context, id int) error {
	return c.repos.ContactRepo.DeleteContact(ctx, id)
}

func NewContactService(
	repos *repository.Repositories,
	log golog.ContextLogger,
	tr trace.Tracer,
	cfg *config.Config,
) *ContactService {
	return &ContactService{
		repos:  repos,
		logger: log,
		tr:     tr,
		cfg:    cfg,
	}
}
