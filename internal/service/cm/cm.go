package cm

import (
	"context"
	"errors"
	"fmt"
	"gRPC_ContactManagement_Service/internal/domain/models"
	"gRPC_ContactManagement_Service/internal/lib/logger/sl"
	"gRPC_ContactManagement_Service/internal/storage"
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
		creatorEmail, name, email, phone string,
	) (uid int64, err error)
}

type ContactProvider interface {
	Contact(
		ctx context.Context,
		creatorEmail, name, email, phone string,
	) (models.Contact, error)
}

type ContactDeleter interface {
	DeleteContact(
		ctx context.Context,
		creatorEmail string,
		id int64,
	) error
}

var (
	ErrContactExists   = errors.New("contact exists")
	ErrContactNotFound = errors.New("contact not found")
)

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
	creatorEmail, name, email, phone string,
) (int64, error) {
	const op = "cm.CreateContact"
	log := cmg.log.With(
		slog.String("op", op),
	)
	log.Info("creating contact")

	uid, err := cmg.contactSaver.SaveContact(ctx, creatorEmail, name, email, phone)
	if err != nil {
		if errors.Is(err, storage.ErrContactExists) {
			log.Warn("contact already exists", sl.Err(err))
			return -1, fmt.Errorf("%s: %w", op, ErrContactExists)
		}
		log.Error("failed to save contact", sl.Err(err))
		return -1, fmt.Errorf("%s: %w", op, err)
	}
	return uid, nil
}

func (cmg *ContactManager) GetContactByName(
	ctx context.Context,
	creatorEmail, name string,
) (models.Contact, error) {
	const op = "cm.GetContactByName"
	log := cmg.log.With(
		slog.String("op", op),
	)
	log.Info("searching for contact", slog.String("name", name))

	contact, err := cmg.contactProvider.Contact(ctx, creatorEmail, name, "", "")
	if err != nil {
		if errors.Is(err, storage.ErrContactNotFound) {
			return models.Contact{}, fmt.Errorf("%s: %w", op, ErrContactNotFound)
		}
		return models.Contact{}, fmt.Errorf("%s: %w", op, err)
	}
	return contact, nil
}

func (cmg *ContactManager) GetContactByEmail(
	ctx context.Context,
	creatorEmail, email string,
) (models.Contact, error) {
	const op = "cm.GetContactByEmail"
	log := cmg.log.With(
		slog.String("op", op),
	)
	log.Info("searching for contact", slog.String("email", email))

	contact, err := cmg.contactProvider.Contact(ctx, creatorEmail, "", email, "")
	if err != nil {
		if errors.Is(err, storage.ErrContactNotFound) {
			return models.Contact{}, fmt.Errorf("%s: %w", op, ErrContactNotFound)
		}
		return models.Contact{}, fmt.Errorf("%s: %w", op, err)
	}
	return contact, nil
}

func (cmg *ContactManager) GetContactByPhone(
	ctx context.Context,
	creatorEmail, phone string,
) (models.Contact, error) {
	const op = "cm.GetContactByPhone"
	log := cmg.log.With(
		slog.String("op", op),
	)
	log.Info("searching for contact", slog.String("phone", phone))

	contact, err := cmg.contactProvider.Contact(ctx, creatorEmail, "", "", phone)
	if err != nil {
		if errors.Is(err, storage.ErrContactNotFound) {
			return models.Contact{}, fmt.Errorf("%s: %w", op, ErrContactNotFound)
		}
		return models.Contact{}, fmt.Errorf("%s: %w", op, err)
	}
	return contact, nil
}

func (cmg *ContactManager) DeleteContact(
	ctx context.Context,
	creatorEmail string,
	id int64,
) error {
	const op = "cm.DeleteContact"
	log := cmg.log.With(
		slog.String("op", op),
	)
	log.Info("trying to delete contact")

	err := cmg.contactDeleter.DeleteContact(ctx, creatorEmail, id)
	if err != nil {
		if errors.Is(err, storage.ErrContactNotFound) {
			return ErrContactNotFound
		}
		return err
	}
	return nil
}
