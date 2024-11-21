package storage

import "errors"

var (
	ErrContactExists   = errors.New("contact exists")
	ErrContactNotFound = errors.New("contact not found")
)
