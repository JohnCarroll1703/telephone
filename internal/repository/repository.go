package repository

import (
	"context"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
	"telephone/internal/config"
	"telephone/internal/model"
	"telephone/internal/repository/postgres"
)

type User interface {
	CreateUser(ctx context.Context, user *model.User) (_ *model.User, err error)
	GetUserByID(ctx context.Context, id uint) (*model.User, error)
	GetAllUsers() ([]model.User, error)
	GetAllUsersWithPaginationAndFiltering(limit int, page int,
		sort string,
		filter map[string]interface{},
		direction string) ([]model.User,
		*model.Paginate, error)
}

type Contact interface {
	CreateContact(ctx context.Context, contact *model.Contact) (_ *model.Contact, err error)
	GetContactByID(ctx context.Context, id uint64) (*model.Contact, error)
	GetAllContacts() ([]model.Contact, error)
	GetByPhone(
		ctx context.Context,
		phone string,
	) (*model.Contact, error)
}

type UserContacts interface {
	GetByPhone(ctx context.Context, phone string) (
		contactRelation *model.UserContactRelation, err error)
	AddContacts(userID int, contactID int) (_ *model.UserContactRelation, err error)
	GetByUserIDContactID(userID int, contactID int) (
		contactRelation *model.UserContactRelation,
		err error)
	ListFav(userID uint64) ([]model.Contact, error)
	GetAllRelations() ([]model.UserContactRelation, error)
}

type Repositories struct {
	UserRepo              User
	ContactRepo           Contact
	UserContactRepository UserContacts
}

func NewRepositories(
	cfg *config.Config,
	jaegerTrace trace.Tracer,
	db *gorm.DB) *Repositories {
	userRepo := postgres.NewTelephone(jaegerTrace, db)
	contactRepo := postgres.NewContact(jaegerTrace, db)
	userContactsRepo := postgres.NewUserContacts(jaegerTrace, db)
	return &Repositories{
		UserRepo:              userRepo,
		ContactRepo:           contactRepo,
		UserContactRepository: userContactsRepo,
	}
}
