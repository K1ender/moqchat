package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/K1ender/moqchat/internal/entity/model"
	"github.com/K1ender/moqchat/internal/repository"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrWrongPassword = errors.New("wrong password")
)

type Auth interface {
	Login(ctx context.Context, email string, password string) (uuid.UUID, error)
	Register(ctx context.Context, username string, email string, password string) (uuid.UUID, error)
}

type AuthUsecase struct {
	userRepo repository.User
}

func NewAuthUsecase(userRepo repository.User) Auth {
	return &AuthUsecase{userRepo: userRepo}
}

// Login implements [Auth].
func (a *AuthUsecase) Login(ctx context.Context, email string, password string) (uuid.UUID, error) {
	user, err := a.userRepo.FindUserByEmail(ctx, email)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to find user: %w", err)
	}

	if comparePassword(user.Password, []byte(password)) != nil {
		return uuid.Nil, ErrWrongPassword
	}

	return user.ID, nil
}

// Register implements [Auth].
func (a *AuthUsecase) Register(ctx context.Context, username string, email string, password string) (uuid.UUID, error) {
	passwordHash, err := hashPassword([]byte(password))
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := model.User{Username: username, Email: email, Password: passwordHash}
	id, err := a.userRepo.CreateUser(ctx, user)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create user: %w", err)
	}

	return id, nil
}

func hashPassword(password []byte) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	return hash, nil
}

func comparePassword(hashedPassword, password []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, password)
}
