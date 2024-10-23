package repository

import (
	"context"
	"git.tarlanpayments.kz/pkg/golog"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel/trace"
	"telephone/internal/config"
	"telephone/internal/model"
	"telephone/internal/repository/postgres"
)

type User interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUserByID(ctx context.Context, id int) (*model.User, error)
	GetAllUsers() ([]*model.User, error)
	UpdateUser(ctx context.Context, user *model.User) error
	DeleteUser(ctx context.Context, id int) error
}

type Contact interface {
	CreateContact(ctx context.Context, contact *model.Contact) error
	GetContactByID(ctx context.Context, id int) (*model.Contact, error)
	GetAllContacts() ([]*model.Contact, error)
	UpdateContact(ctx context.Context, id int, updatedContact *model.Contact) error
	DeleteContact(ctx context.Context, id int) error
}

type UserContacts interface {
	AddUserContact(ctx context.Context, userID int,
		contact model.Contact, isFav bool,
		userContact model.UserContacts) error
	RemoveUserContact(ctx context.Context,
		userID int,
		contactID int) error
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
	postgre *pgx.Conn) *Repositories {
	userRepo := postgres.NewTelephone(logger, jaegerTrace, postgre)
	contactRepo := postgres.NewContact(logger, jaegerTrace, postgre)
	userContactsRepo := postgres.NewUserContacts(logger, jaegerTrace, postgre)
	return &Repositories{
		UserRepo:              userRepo,
		ContactRepo:           contactRepo,
		UserContactRepository: userContactsRepo,
	}
}
