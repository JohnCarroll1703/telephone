package service

import (
	"context"
	"git.tarlanpayments.kz/pkg/golog"
	"git.tarlanpayments.kz/processing/events/goevents/producer"
	"go.opentelemetry.io/otel/trace"
	"telephone/internal/config"
	"telephone/internal/model"
	"telephone/internal/repository"
	"telephone/internal/schema"
)

func NewServices(deps Deps) *Services {
	return &Services{
		telephoneService:   NewTelephoneService(deps.Repos, deps.Logger, deps.JeagerTracer, deps.Cgf),
		contactService:     NewContactService(deps.Repos, deps.Logger, deps.JeagerTracer, deps.Cgf),
		userContactService: NewUserContactService(deps.Repos, deps.Logger, deps.JeagerTracer, deps.Cgf),
	}
}

type Deps struct {
	Repos        *repository.Repositories
	Cgf          *config.Config
	Logger       golog.ContextLogger
	Producer     producer.Producer
	JeagerTracer trace.Tracer
}

type Services struct {
	telephoneService   Telephone
	contactService     Contact
	userContactService UserContacts
}

type Telephone interface {
	CreateUser(ctx context.Context, user *schema.User) error
	GetUserByID(ctx context.Context, id int) (schema.User, error)
	GetAllUsers(ctx context.Context, tr trace.Tracer,
		funcName string) ([]model.User, error)
}

type Contact interface {
	GetContacts() ([]model.Contact, error)
	CreateContact(ctx context.Context, contact *schema.Contact) error
	GetContactByID(ctx context.Context, id int) (*model.Contact, error)
}

type UserContacts interface {
	AddUserContact(
		ctx context.Context,
		userID int, phone string, contact *schema.Contact) error
	Create(ctx context.Context, userId int, phone string) error
	GetUserContacts(ctx context.Context, userID int) ([]model.UserContacts, error)
}
