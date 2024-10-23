package postgres

import (
	"context"
	"git.tarlanpayments.kz/pkg/golog"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
	"net/http"
	"telephone/internal/model"
	"telephone/internal/repository"
	"time"
)

type Telephone struct {
	client   *http.Client
	log      golog.ContextLogger
	tr       trace.Tracer
	postgres *pgx.Conn
	db       *gorm.DB
}

func NewTelephone(log golog.ContextLogger, tr trace.Tracer, postgres *pgx.Conn) *Telephone {
	return &Telephone{
		client: &http.Client{
			Timeout: time.Minute,
		},
		log:      log,
		tr:       tr,
		postgres: postgres,
	}
}

func (t Telephone) GetAllUsers() ([]*model.User, error) {
	var db *gorm.DB
	var users []*model.User
	err := db.Model(&model.User{}).Find(&users).Error
	return users, err
}

func (t Telephone) DeleteUser(ctx context.Context, id int) error {
	//TODO implement me
	return t.db.Delete(&model.User{}, id).Error
}

func (t Telephone) GetUserByID(ctx context.Context, id int) (*model.User, error) {
	var user model.User
	err := t.db.First(&user, id).Error
	return &user, err
}

func (t Telephone) CreateUser(ctx context.Context, user *model.User) error {
	return t.db.Create(&user).Error
}

func (t Telephone) AddContact(ctx context.Context, userID int, contact model.Contact) error {
	return t.db.Create(&contact).Error
}

func (t Telephone) UpdateUser(ctx context.Context, user *model.User) error {
	return t.db.Save(&user).Error
}

func (t Telephone) UpdateContact(ctx context.Context, contact *model.Contact) error {
	return t.db.Save(&contact).Error
}

func (t Telephone) DeleteContact(ctx context.Context, contactID int) error {
	return t.db.Delete(&contactID).Error
}

func New() repository.User {
	return &Telephone{}
}
