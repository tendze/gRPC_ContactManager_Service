package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"gRPC_ContactManagement_Service/internal/domain/models"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "sqlite.New"
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &Storage{db: db}, nil
}

func (s *Storage) SaveContact(
	ctx context.Context,
	creatorEmail, name, email, phone string,
) (uid int64, err error) {
	panic("implement me")
}

func (s *Storage) Contact(
	ctx context.Context,
	creatorEmail, name, email, phone string,
) (models.Contact, error) {
	panic("implement me")
}

func (s *Storage) DeleteContact(
	ctx context.Context,
	creatorEmail string,
	id int64,
) error {
	panic("implement me")
}
