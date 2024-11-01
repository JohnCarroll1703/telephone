package service

import (
	"context"
	"errors"
	"git.tarlanpayments.kz/pkg/golog"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
	"telephone/internal/config"
	"telephone/internal/model"
	"telephone/internal/repository"
	"telephone/pkg/terr"

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

func (s UserContactsService) AddContacts(
	ctx context.Context,
	userID uint64,
	request *schema.AddContactRequest,
) (*schema.AddContactResponse, error) {
	var (
		findContact *model.Contact
		relation    *model.UserContacts
	)

	findContact, err := s.repos.ContactRepo.GetByPhone(ctx, request.PhoneNumber)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		findContact, err = s.repos.ContactRepo.CreateContact(ctx, findContact)
		if err != nil {
			return nil, err
		}
	}

	relation, err = s.repos.UserContactRepository.GetByUserIDContactID(int(userID), findContact.ContactID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		relation, err = s.repos.UserContactRepository.
			AddContacts(int(userID), findContact.ContactID)
		if err != nil {
			return nil, err
		}
	}

	return &schema.AddContactResponse{
		ID:          uint64(relation.ContactID),
		ContactName: findContact.ContactName,
		PhoneNumber: findContact.PhoneNumber,
	}, nil
}

func (s UserContactsService) GetUserContacts(ctx context.Context, userID int) ([]model.UserContacts, error) {
	user, err := s.repos.UserRepo.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []model.UserContacts{}, terr.RecordNotFound
		}
	}

	res := user.UserContacts
	return res, err
}

func (s UserContactsService) GetByUserIDContactID(userID int, contactID int) (_ *model.UserContacts, err error) {
	data, err := s.repos.UserContactRepository.GetByUserIDContactID(userID, contactID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &model.UserContacts{}, terr.RecordNotFound
		}
	}
	return &model.UserContacts{
		UserContactsID: data.UserContactsID,
		UserID:         data.UserID,
		ContactID:      data.ContactID,
		IsFavorite:     data.IsFavorite,
		Contact:        data.Contact,
		User:           data.User,
	}, nil
}
