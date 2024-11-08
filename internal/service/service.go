package service

import (
	"context"
	"go.opentelemetry.io/otel/trace"
	"telephone/internal/config"
	"telephone/internal/model"
	"telephone/internal/repository"
	"telephone/internal/schema"
)

func NewServices(deps Deps) *Services {
	return &Services{
		TelephoneService:   NewTelephoneService(deps.Repos, deps.JeagerTracer, deps.Cgf),
		ContactService:     NewContactService(deps.Repos, deps.JeagerTracer, deps.Cgf),
		UserContactService: NewUserContactService(deps.Repos, deps.JeagerTracer, deps.Cgf),
	}
}

type Deps struct {
	Repos        *repository.Repositories
	Cgf          *config.Config
	JeagerTracer trace.Tracer
}

type Services struct {
	TelephoneService   Telephone
	ContactService     Contact
	UserContactService UserContactRelation
}

type Telephone interface {
	CreateUser(ctx context.Context, user *schema.User) (*model.User, error)
	GetUserByID(ctx context.Context, id uint) (schema.User, error)
	GetAllUsers(ctx context.Context) ([]model.User, error)
	GetAllUsersWithPaginationAndFiltering(
		limit int, page int,
		sort string,
		filter map[string]interface{},
		direction string,
	) ([]model.User, *model.Paginate, error)
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
	) ([]model.Contact, error)
	GetAllRelations() (
		[]model.UserContactRelation, error)
}
