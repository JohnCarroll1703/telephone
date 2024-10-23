package postgres

import (
	"context"
	"fmt"
	"git.tarlanpayments.kz/pkg/golog"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
	"net/http"
	"telephone/internal/model"
)

type Contact struct {
	client   *http.Client
	log      golog.ContextLogger
	tr       trace.Tracer
	postgres *pgx.Conn
	db       *gorm.DB
}

func (c Contact) CreateContact(ctx context.Context, contact *model.Contact) error {
	//TODO implement me zaebal
	return c.db.Create(&contact).Error
	panic("implement me")
}

func (c Contact) GetContactByID(ctx context.Context, id int) (*model.Contact, error) {
	//TODO implement me
	var contact model.Contact
	err := c.db.First(&contact, id).Error
	return &contact, err
}

func (c Contact) GetAllContacts() ([]*model.Contact, error) {
	var db *gorm.DB
	var contacts []*model.Contact
	err := db.Model(&model.Contact{}).Find(&contacts).Error
	return contacts, err
}

func (c Contact) UpdateContact(ctx context.Context, id int, updatedContact *model.Contact) error {
	//TODO implement me
	var contact model.Contact
	if err := c.db.First(&id).Error; err != nil {
		return fmt.Errorf("user not found tipo")
	}

	contact.ContactID = updatedContact.ContactID
	contact.ContactName = updatedContact.ContactName
	return c.db.Save(&updatedContact).Error
}

func (c Contact) DeleteContact(ctx context.Context, id int) error {
	//TODO implement me
	panic("implement me")
}

func NewContact(
	log golog.ContextLogger,
	tr trace.Tracer,
	postgre *pgx.Conn) *Contact {
	return &Contact{
		log:      log,
		tr:       tr,
		postgres: postgre,
	}
}
