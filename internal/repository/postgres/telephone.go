package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
	"math"
	"net/http"
	"strings"
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

func (t Telephone) GetUserByID(ctx context.Context, id uint) (user *model.User, err error) {
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

func (t Telephone) GetAllUsersWithPaginationAndFiltering(limit int, page int,
	sort string, filter map[string]interface{}, direction string) ([]model.User, *model.Paginate, error) {
	var users []model.User
	var pagination model.Paginate
	query := t.db

	if filter["name"] != "" {
		query = query.Where("name = ?", filter["name"].(string))
	}

	pagination.Limit = limit
	pagination.Page = page

	if strings.EqualFold(direction, "asc") {
		direction = "asc"
	}

	if strings.EqualFold(direction, "desc") {
		direction = "desc"
	}

	query = query.Order(fmt.Sprintf("%s %s", direction, sort))

	err := t.db.Scopes(t.Paginate(users, &pagination, query,
		int64(len(users)))).Find(&users).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, terr.RecordNotFound
		}
		return nil, nil, err
	}

	return users, &pagination, nil
}

func (t Telephone) Paginate(value interface{}, pagination *model.Paginate,
	db *gorm.DB, currRecord int64) func(db *gorm.DB) *gorm.DB {
	var totalRecords int64
	db.Model(value).Count(&totalRecords)

	pagination.TotalRecords = totalRecords
	pagination.TotalPage = int(math.Ceil(float64(totalRecords)) / float64(pagination.GetLimit()))
	pagination.Records = int64(pagination.Limit*(pagination.Page-1)) + currRecord
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit())
	}
}
