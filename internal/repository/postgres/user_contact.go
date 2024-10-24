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

type UserContacts struct {
	client   *http.Client
	log      golog.ContextLogger
	tr       trace.Tracer
	postgres *pgx.Conn
	db       *gorm.DB
}

func (u UserContacts) GetByContact(ctx context.Context, phone string) (*model.UserContacts, error) {
	var userContact model.UserContacts
	err := u.db.Where("user_id = ?", "contact_id = ?").Find(&phone).Error
	return &userContact, err
}

func (u UserContacts) AddUserContact(ctx context.Context,
	userID int, phone string) error {

	contact, err := u.GetByContact(ctx, phone)
	if err != nil {

	}
	return nil
}

func NewUserContacts(log golog.ContextLogger,
	tr trace.Tracer,
	postgre *pgx.Conn) *UserContacts {
	return &UserContacts{
		log:      log,
		tr:       tr,
		postgres: postgre,
	}
}
