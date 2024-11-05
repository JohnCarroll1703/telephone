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
	userContactService UserContactRelation
}

type Telephone interface {
	CreateUser(ctx context.Context, user *schema.User) error
	GetUserByID(ctx context.Context, id int) (schema.User, error)
	GetAllUsers(ctx context.Context, tr trace.Tracer,
		funcName string) ([]model.User, error)
}

type Contact interface {
	GetContacts() ([]model.Contact, error)
	CreateContact(ctx context.Context, contact *schema.Contact) (_ *model.Contact, err error)
	GetContactByID(ctx context.Context, id uint64) (schema.ContactResponse, error)
}

type UserContactRelation interface {
	AddContacts(
		ctx context.Context,
		userID uint64,
		request *schema.AddContactRequest,
	) (*schema.AddContactResponse, error)
	ListFav(ctx context.Context, userID int,
	) (*model.UserContactRelation, error)
}
