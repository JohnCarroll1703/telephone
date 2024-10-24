package postgres

import (
	"context"
	"git.tarlanpayments.kz/pkg/golog"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
	"net/http"
	"telephone/internal/model"
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

func (t Telephone) GetAllUsers() ([]model.User, error) {
	var db *gorm.DB
	var users []model.User
	err := db.Model(&model.User{}).Find(&users).Error
	return users, err
}

func (t Telephone) GetUserByID(ctx context.Context, id int) (model.User, error) {
	var user model.User
	res := t.db.Find(&user, id)
	if res.Error != nil {
		return model.User{}, res.Error
	}

	if res.RowsAffected == 0 {
		return model.User{}, res.Error
	}

	return res, nil
}

func (t Telephone) CreateUser(
	ctx context.Context,
	user *model.User) error {
	return t.db.Create(&user).Error
}
