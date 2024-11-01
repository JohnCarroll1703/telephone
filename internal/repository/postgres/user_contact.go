package postgres

import (
	"context"
	"errors"
	"git.tarlanpayments.kz/pkg/golog"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
	"net/http"
	"telephone/internal/model"
	"telephone/internal/schema"
	"telephone/pkg/terr"
)

type UserContacts struct {
	client *http.Client
	log    golog.ContextLogger
	tr     trace.Tracer
	db     *gorm.DB
}

func (u UserContacts) GetByPhone(ctx context.Context, req schema.AddContactRequest,
) (contactRelation *model.UserContacts, err error) {
	err = u.db.Where("user_id = ?", "contact_id = ?").First(&req).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
		}
		return &model.UserContacts{}, terr.RecordNotFound
	}

	return contactRelation, err
}

func (u UserContacts) GetByUserIDContactID(userID int, contactID int) (
	contactRelation *model.UserContacts,
	err error) {
	err = u.db.Where("contact_id = ? AND user_id = ?", contactID, userID).First(&contactRelation).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &model.UserContacts{}, terr.RecordNotFound
		}
	}

	return contactRelation, err
}

func (u UserContacts) AddContacts(userID int, contactID int) (err error) {
	err = u.db.Create(&model.UserContacts{
		UserID:    userID,
		ContactID: contactID}).
		Error

	if err != nil {
		return terr.ErrDbUnexpected
	}

	return nil
}

func NewUserContacts(
	log golog.ContextLogger,
	tr trace.Tracer,
	db *gorm.DB) *UserContacts {
	return &UserContacts{
		log: log,
		tr:  tr,
		db:  db,
	}
}
