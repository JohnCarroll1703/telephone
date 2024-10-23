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
		telephoneService: NewTelephoneService(deps.Repos, deps.Logger, deps.JeagerTracer, deps.Cgf),
		contactService:   NewContactService(deps.Repos, deps.Logger, deps.JeagerTracer, deps.Cgf),
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
	telephoneService Telephone
	contactService   Contact
}

type Telephone interface {
	CreateUser(ctx context.Context, user *schema.User) error
	GetUserByID(ctx context.Context, id int) (*model.User, error)
	UpdateUser(ctx context.Context, user *schema.User) error
	DeleteUser(ctx context.Context, id int) error
	GetAllUsers(ctx context.Context, tr trace.Tracer,
		funcName string) ([]*model.User, error)
}

type Contact interface {
	GetContacts() ([]*model.Contact, error)
	CreateContact(ctx context.Context, contact *schema.Contact) error
	GetContactByID(ctx context.Context, id int) (*model.Contact, error)
	UpdateContact(ctx context.Context, id int) (*model.Contact, error)
	DeleteContact(ctx context.Context, id int) error
}

type UserContacts interface {
	AddUserContact(ctx context.Context,
		userID int, contact model.Contact,
		isFav bool,
		userContact model.UserContacts) error
	RemoveUserContact(ctx context.Context, userID int, contactID int) error
}
