package data

import (
	"auth/domain"
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
)

const DB_TIMEOUT = 30

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

type userModel struct {
	id             uuid.UUID
	email          string
	firstname      string
	lastname       string
	hashedPassword string
	isActive       bool
	createdAt      time.Time
	updatedAt      time.Time
}

func userScanner(row *sql.Row, fetchedUser *userModel) error {
	return row.Scan(&fetchedUser.id, &fetchedUser.email, &fetchedUser.firstname, &fetchedUser.lastname,
		&fetchedUser.hashedPassword, &fetchedUser.isActive, &fetchedUser.createdAt, &fetchedUser.updatedAt)
}

func (ur *UserRepo) GetByEmail(email string) (*domain.User, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), DB_TIMEOUT*time.Second)
	defer cancelFunc()

	var fetchedUser userModel

	query := `SELECT * FROM users WHERE email = $1`
	row := ur.db.QueryRowContext(ctx, query, email)

	sErr := userScanner(row, &fetchedUser)
	if sErr != nil {
		if sErr == sql.ErrNoRows {
			return nil, sErr
		}
		log.Println(sErr.Error())
		return nil, sErr
	}

	return fromUserModeltoDomain(fetchedUser), nil
}

func fromUserModeltoDomain(um userModel) *domain.User {
	return &domain.User{
		ID:             um.id,
		Email:          um.email,
		FirstName:      um.firstname,
		LastName:       um.lastname,
		HashedPassword: um.hashedPassword,
		IsActive:       um.isActive,
	}
}
