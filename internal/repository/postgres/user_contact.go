package postgres

import (
	"context"
	"errors"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"telephone/internal/model"
	"telephone/internal/schema"
	"telephone/pkg/terr"
)

const (
	sortBySeparator string = ", "
)

type UserContacts struct {
	client *http.Client
	tr     trace.Tracer
	db     *gorm.DB
}

func (u UserContacts) ListFav(userID int) (contactRelation *model.UserContactRelation, err error) {
	err = u.db.Where("id = ?", userID).Find(&model.UserContactRelation{}).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, terr.RecordNotFound
		}
		return nil, err
	}

	return contactRelation, nil
}

func (u UserContacts) GetByPhone(ctx context.Context, req schema.ContactRequest) (contactRelation *model.UserContactRelation, err error) {

	err = u.db.Where("user_id = ?", "contact_id = ?").First(&req).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
		}
		return nil, terr.RecordNotFound
	}

	return contactRelation, err
}

func (u UserContacts) GetByUserIDContactID(userID int, contactID int) (
	contactRelation *model.UserContactRelation,
	err error) {
	err = u.db.Where("contact_id = ? AND user_id = ?", contactID, userID).First(&contactRelation).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &model.UserContactRelation{}, terr.RecordNotFound
		}
		return nil, err
	}

	return contactRelation, err
}

func (u UserContacts) AddContacts(userID int, contactID int) (_ *model.UserContactRelation, err error) {
	err = u.db.Create(&model.UserContactRelation{
		UserID:    userID,
		ContactID: contactID}).
		Error

	if err != nil {
		return nil, terr.ErrDbUnexpected
	}

	return nil, nil
}

func (u UserContacts) sortBy(filter schema.UserAndContactFilter) string {
	sortRes := strings.Split(filter.SortBy, ",")
	sort := ""

	for key, value := range sortRes {
		if key%2 == 0 {
			sort += value + " "
		} else {
			sort += value + sortBySeparator
		}
	}

	if sort == " " {
		return schema.DefaultOrder
	}

	sort = sort[:len(sort)-2]

	return sort
}

func NewUserContacts(
	tr trace.Tracer,
	db *gorm.DB) *UserContacts {
	return &UserContacts{
		tr: tr,
		db: db,
	}
}
