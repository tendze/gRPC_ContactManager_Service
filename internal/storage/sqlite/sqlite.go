package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"gRPC_ContactManagement_Service/internal/domain/models"
	_ "github.com/mattn/go-sqlite3"
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
	const op = "sqlite.SaveContact"

	stmt, err := s.db.Prepare("INSERT INTO contacts(creator_email, name, email, phone) VALUES(?, ?, ?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.ExecContext(ctx, creatorEmail, name, email, phone)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (s *Storage) Contact(
	ctx context.Context,
	creatorEmail, name, email, phone string,
) (models.Contact, error) {
	const op = "sqlite.Contact"

	var query, param string
	if name != "" {
		query = "SELECT name, email, phone FROM contacts WHERE creator_email = ? AND name = ?"
		param = name
	} else if email != "" {
		query = "SELECT name, email, phone FROM contacts WHERE creator_email = ? AND email = ?"
		param = email
	} else {
		query = "SELECT name, email, phone FROM contacts WHERE creator_email = ? AND phone = ?"
		param = phone
	}
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return models.Contact{}, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, creatorEmail, param)
	var contact models.Contact
	err = row.Scan(&contact.Name, &contact.Email, &contact.Phone)
	if err != nil {
		return models.Contact{}, fmt.Errorf("%s: %w", op, err)
	}

	return contact, nil
}

func (s *Storage) DeleteContact(
	ctx context.Context,
	creatorEmail string,
	id int64,
) error {
	const op = "sqlite.DeleteContact"

	stmt, err := s.db.Prepare("DELETE FROM contacts WHERE creator_email = ? AND id = ?")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.ExecContext(ctx, creatorEmail, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
