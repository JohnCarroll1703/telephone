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

func (c ContactService) GetContacts() ([]model.Contact, error) {
	return c.repos.ContactRepo.GetAllContacts()
}

func (c ContactService) CreateContact(ctx context.Context, contact *schema.Contact) (_ *model.Contact, err error) {
	return c.repos.ContactRepo.CreateContact(ctx, schema.NewCreateContactRequest(contact))
}

func (c ContactService) GetContactByID(ctx context.Context, id uint64) (schema.ContactResponse, error) {
	data, err := c.repos.ContactRepo.GetContactByID(ctx, id)
	if err != nil {
		return schema.ContactResponse{}, err
	}
	res := schema.ContactResponse{
		ContactID:   uint64(data.ContactID),
		ContactName: data.ContactName,
		PhoneNumber: data.PhoneNumber,
	}

	return res, err
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
