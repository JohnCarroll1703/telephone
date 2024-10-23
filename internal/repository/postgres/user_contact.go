package postgres

import (
	"context"
	"errors"
	"fmt"
	"git.tarlanpayments.kz/pkg/golog"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
	"net/http"
	"telephone/internal/model"
)

type UserContacts struct {
	client   *http.Client
	log      golog.ContextLogger
	tr       trace.Tracer
	postgres *pgx.Conn
	db       *gorm.DB
}

func (u UserContacts) AddUserContact(ctx context.Context, userID int,
	contact model.Contact, isFav bool,
	userContact model.UserContacts) error {
	var user model.User
	if err := u.db.First(&user, user.ID).Error; err != nil {
		return fmt.Errorf("user doesn't exist")
	}

	var existingContact model.Contact
	if err := u.db.Where("contact_name = ? AND phone_number = ?", contact.ContactName, contact.PhoneNumber).
		First(&existingContact).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := u.db.Create(&contact).Error; err != nil {
				return fmt.Errorf("Error creating contact: %v", err)
			}
			existingContact = contact
		} else {
			return fmt.Errorf("Error finding contact: %v", err)
		}
	}

	// create the association in UserContacts
	userContact = model.UserContacts{
		UserID:     userID,
		ContactID:  existingContact.ContactID,
		IsFavorite: isFav,
	}

	if err := u.db.Create(&userContact).Error; err != nil {
		return fmt.Errorf("Error adding contact to user: %v", err)
	}

	return nil
}

func (u UserContacts) RemoveUserContact(ctx context.Context, userID int, contactID int) error {
	//TODO implement me
	var userContact model.UserContacts

	// Find the UserContacts entry that matches the userID and contactID
	if err := u.db.Where("user_id = ? AND contact_id = ?", userID, contactID).
		First(&userContact).Error; err != nil {
		return fmt.Errorf("Error finding contact for user: %v", err)
	}

	// тут мы удаляем ассоциацию в таблице
	if err := u.db.Delete(&userContact).Error; err != nil {
		return fmt.Errorf("Error deleting contact for user: %v", err)
	}

	return nil
}

func NewUserContacts(log golog.ContextLogger,
	tr trace.Tracer,
	postgre *pgx.Conn) *UserContacts {
	return &UserContacts{
		log:      log,
		tr:       tr,
		postgres: postgre,
	}
}
