package repository

import (
	"context"
	"git.tarlanpayments.kz/pkg/golog"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
	"telephone/internal/config"
	"telephone/internal/model"
	"telephone/internal/repository/postgres"
	"telephone/internal/schema"
)

type User interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUserByID(ctx context.Context, id int) (*model.User, error)
	GetAllUsers() ([]model.User, error)
}

type Contact interface {
	CreateContact(ctx context.Context, contact *model.Contact) (_ *model.Contact, err error)
	GetContactByID(ctx context.Context, id uint64) (*model.Contact, error)
	GetAllContacts() ([]model.Contact, error)
	GetByPhone(
		ctx context.Context,
		phone string,
	) (resp *model.Contact, err error)
}

type UserContacts interface {
	GetByPhone(ctx context.Context, req schema.ContactRequest,
	) (contactRelation *model.UserContactRelation, err error)
	AddContacts(userID int, contactID int) (_ *model.UserContactRelation, err error)
	GetByUserIDContactID(userID int, contactID int) (
		contactRelation *model.UserContactRelation,
		err error)
	ListFav(userID int) (contactRelation *model.UserContactRelation, err error)
}

type Repositories struct {
	UserRepo              User
	ContactRepo           Contact
	UserContactRepository UserContacts
}

func NewRepositories(
	cfg *config.Config,
	jaegerTrace trace.Tracer,
	logger golog.ContextLogger,
	db *gorm.DB) *Repositories {
	userRepo := postgres.NewTelephone(logger, jaegerTrace, db)
	contactRepo := postgres.NewContact(logger, jaegerTrace, db)
	userContactsRepo := postgres.NewUserContacts(logger, jaegerTrace, db)
	return &Repositories{
		UserRepo:              userRepo,
		ContactRepo:           contactRepo,
		UserContactRepository: userContactsRepo,
	}
}
