package postgres

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
	"net/http"
	"telephone/internal/model"
	"telephone/pkg/terr"
	"time"
)

type Telephone struct {
	client   *http.Client
	tr       trace.Tracer
	postgres *pgx.Conn
	db       *gorm.DB
}

func NewTelephone(tr trace.Tracer, db *gorm.DB) *Telephone {
	return &Telephone{
		client: &http.Client{
			Timeout: time.Minute,
		},
		tr: tr,
		db: db,
	}
}

func (t Telephone) GetAllUsers() ([]model.User, error) {
	var users []model.User
	err := t.db.Model(&model.User{}).Find(&users).Error
	return users, err
}

func (t Telephone) GetUserByID(ctx context.Context, id int) (user *model.User, err error) {
	err = t.db.Find(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, terr.RecordNotFound
		}
		return &model.User{}, err
	}

	return user, nil
}

func (t Telephone) CreateUser(
	ctx context.Context,
	user *model.User) (_ *model.User, err error) {
	if err = t.db.WithContext(ctx).Create(&user).Error; err != nil {
		return nil, terr.ErrDbUnexpected
	}

	return nil, nil
}
