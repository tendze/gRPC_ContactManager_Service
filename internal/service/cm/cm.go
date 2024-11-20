package cm

import (
	"context"
	"gRPC_ContactManagement_Service/internal/domain/models"
	"log/slog"
)

type ContactManager struct {
	log             *slog.Logger
	contactSaver    ContactSaver
	contactProvider ContactProvider
	contactDeleter  ContactDeleter
}

type ContactSaver interface {
	SaveContact(
		ctx context.Context,
		name, email, phone string,
	) (uid int64, err error)
}

type ContactProvider interface {
	Contact(
		ctx context.Context,
		name, email, phone string,
	) (models.Contact, error)
}

type ContactDeleter interface {
	Delete(
		ctx context.Context,
		name, email, phone string,
	) error
}

func New(
	log *slog.Logger,
	saver ContactSaver,
	provider ContactProvider,
	deleter ContactDeleter,
) *ContactManager {
	return &ContactManager{
		log:             log,
		contactSaver:    saver,
		contactProvider: provider,
		contactDeleter:  deleter,
	}
}

func (cmg *ContactManager) CreateContact(
	ctx context.Context,
	name, email, phone string,
) (int64, error) {
	panic("implement me")
}

func (cmg *ContactManager) GetContactByName(
	ctx context.Context,
	name string,
) (models.Contact, error) {
	panic("implement me")
}

func (cmg *ContactManager) GetContactByEmail(
	ctx context.Context,
	email string,
) (models.Contact, error) {
	panic("implement me")
}

func (cmg *ContactManager) GetContactByPhone(
	ctx context.Context,
	email string,
) (models.Contact, error) {
	panic("implement me")
}

func (cmg *ContactManager) DeleteContact(
	ctx context.Context,
	id int64,
) error {
	panic("implement me")
}
