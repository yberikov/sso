package storage

import (
	"context"
	"errors"
	"sso/internal/domain/models"
)

var (
	ErrUserExists   = errors.New("user already exists")
	ErrUserNotFound = errors.New("user not found")
)

type Auth interface {
	SaveUser(ctx context.Context, email string, passHash []byte) (int64, error)
	GetUser(ctx context.Context, email string) (models.User, error)
	EditUser(ctx context.Context, email string, name string, telephone string, birthDate string) error
}
