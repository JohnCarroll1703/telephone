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

func (u UserContacts) ListFav(userID uint64) ([]model.Contact, error) {
	var listContacts []model.Contact
	err := u.db.Table("user_contact_relations").
		Select("contact_id").
		Where("id = ?", userID).
		Find(&listContacts).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, terr.RecordNotFound
		}
	}

	return listContacts, err
}

func (u UserContacts) GetByPhone(ctx context.Context, phone string) (contactRelation *model.UserContactRelation, err error) {

	err = u.db.Where("user_id = ?", "contact_id = ?").First(&phone).Error
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
	err = u.db.Where("contact_id = ? AND id = ?", contactID, userID).First(&contactRelation).Error

	return contactRelation, err
}

func (u UserContacts) AddContacts(userID int, contactID int) (_ *model.UserContactRelation, err error) {
	err = u.db.Create(&model.UserContactRelation{
		UserID:     userID,
		IsFavorite: true,
		ContactID:  contactID}).
		Error

	if err != nil {
		return nil, terr.ErrDbUnexpected
	}

	return nil, nil
}

func (u UserContacts) GetAllRelations() ([]model.UserContactRelation, error) {
	var relations []model.UserContactRelation
	err := u.db.Model(&model.UserContactRelation{}).Find(&relations).Error
	return relations, err
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
