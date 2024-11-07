package service

import (
	"context"
	"errors"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
	"telephone/internal/config"
	"telephone/internal/model"
	"telephone/internal/repository"
	"telephone/pkg/terr"

	"telephone/internal/schema"
)

type UserContactsService struct {
	repos *repository.Repositories
	tr    trace.Tracer
	cfg   *config.Config
}

func NewUserContactService(
	repos *repository.Repositories,
	tr trace.Tracer,
	cfg *config.Config) *UserContactsService {
	return &UserContactsService{
		repos: repos,
		tr:    tr,
		cfg:   cfg,
	}
}

func (s UserContactsService) AddContacts(
	ctx context.Context,
	userID uint64,
	request *schema.AddContactRequest,
) (*schema.AddContactResponse, error) {
	var (
		findContact *model.Contact
		relation    *model.UserContactRelation
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
		ContactID:   uint64(relation.ContactID),
		PhoneNumber: findContact.PhoneNumber,
	}, nil
}

func (s UserContactsService) ListFav(ctx context.Context, userID int,
) ([]model.Contact, error) { // listFav типо
	user, err := s.repos.UserRepo.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, terr.RecordNotFound
		}
		return nil, err
	}

	res, err := s.repos.UserContactRepository.ListFav(uint64(user.ID))
	return res, err
}

func (s UserContactsService) GetByUserIDContactID(userID int, contactID int) (_ *model.UserContactRelation, err error) {
	data, err := s.repos.UserContactRepository.GetByUserIDContactID(userID, contactID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, terr.RecordNotFound
		}
		return nil, err
	}
	return &model.UserContactRelation{
		UserContactsID: data.UserContactsID,
		UserID:         data.UserID,
		ContactID:      data.ContactID,
		IsFavorite:     true,
	}, nil
}

func (s UserContactsService) GetAllRelations() (
	[]model.UserContactRelation, error) {
	return s.repos.UserContactRepository.GetAllRelations()
}
