package postgres

import (
	"context"
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
	return c.db.Create(&contact).Error
}

func (c Contact) GetContactByID(ctx context.Context, id int) (*model.Contact, error) {
	var contact model.Contact
	err := c.db.First(&contact, id).Error
	return &contact, err
}

func (c Contact) GetAllContacts() ([]model.Contact, error) {
	var db *gorm.DB
	var contacts []model.Contact
	err := db.Model(&model.Contact{}).Find(&contacts).Error
	return contacts, err
}

func (c Contact) GetContactByPhone(ctx context.Context,
	phone string) (*model.Contact, error) {
	var contact model.Contact
	err := c.db.Where("contact_id = ?").Find(&phone).Error
	return &contact, err
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
