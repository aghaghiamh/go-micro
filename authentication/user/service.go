package user

import (
	"errors"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo Repository
}

func New(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (svc *Service) Authenticate(req AuthRequest) (AuthResponse, error) {
	u, uErr := svc.repo.GetByEmail(req.Email)
	if uErr != nil {

		return AuthResponse{}, fmt.Errorf("email is not correct")
	}

	isMatched, passErr := svc.passwordMatches(req.Password, u.HashedPassword)
	if passErr != nil || !isMatched {

		return AuthResponse{}, fmt.Errorf("password is not correct")
	}

	return AuthResponse{}, nil
}

func (svc Service) passwordMatches(providedPass, savedPass string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(savedPass), []byte(providedPass))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {

			return false, nil
		}

		log.Println(err)
		return false, err
	}

	return true, nil
}
