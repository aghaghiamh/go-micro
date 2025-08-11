package user

import (
	"auth/domain"
)

type Repository interface {
	GetByEmail(email string) (*domain.User, error)
}
