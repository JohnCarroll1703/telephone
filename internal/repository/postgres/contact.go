package postgres

import (
	"context"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
	"net/http"
	"telephone/internal/model"
	"telephone/pkg/terr"
)

type Contact struct {
	client   *http.Client
	tr       trace.Tracer
	postgres *pgx.Conn
	db       *gorm.DB
}

func (c Contact) CreateContact(ctx context.Context, contact *model.Contact) (_ *model.Contact, err error) {
	if err = c.db.WithContext(ctx).Create(&contact).Error; err != nil {
		return nil, terr.ErrDbUnexpected
	}
	return nil, nil
}

func (c Contact) GetContactByID(ctx context.Context, id uint64) (*model.Contact, error) {
	var contact model.Contact
	err := c.db.First(&contact, id).Error
	return &contact, err
}

func (c Contact) GetAllContacts() ([]model.Contact, error) {
	var contacts []model.Contact
	err := c.db.Model(&model.Contact{}).Find(&contacts).Error
	return contacts, err
}

func (c Contact) GetByPhone(
	ctx context.Context,
	phone string,
) (*model.Contact, error) {
	var contact *model.Contact
	err := c.db.Where("phone_number = ?", phone).Find(&contact).Error
	return contact, err
}

func NewContact(
	tr trace.Tracer,
	db *gorm.DB) *Contact {
	return &Contact{
		tr: tr,
		db: db,
	}
}
