package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/K1ender/moqchat/internal/entity/model"
	"github.com/K1ender/moqchat/internal/repository"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	SessionDuration     = 30 * 24 * time.Hour
	HalfSessionDuration = SessionDuration / 2
	RandomStringLength  = 32
)

type Auth interface {
	Login(ctx context.Context, email string, password string) (uuid.UUID, error)
	Register(ctx context.Context, username string, email string, password string) (uuid.UUID, error)

	CreateSession(ctx context.Context, userID uuid.UUID) (string, error)
	GetUserIDFromToken(ctx context.Context, token string) (uuid.UUID, error)
	GetSession(ctx context.Context, id uuid.UUID) (model.Session, error)
	ExtendSession(ctx context.Context, id uuid.UUID) error
}

type AuthUsecase struct {
	userRepo    repository.User
	sessionRepo repository.Session
}

func NewAuthUsecase(userRepo repository.User, sessionRepo repository.Session) Auth {
	return &AuthUsecase{userRepo: userRepo, sessionRepo: sessionRepo}
}

// CreateSession implements [Auth].
func (a *AuthUsecase) CreateSession(ctx context.Context, userID uuid.UUID) (string, error) {
	token, err := cryptoRandomString(RandomStringLength)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	tokenHash := hashToken([]byte(token))

	session := model.Session{UserID: userID, Token: tokenHash, ExpiresAt: time.Now().Add(SessionDuration)}
	_, err = a.sessionRepo.CreateSession(ctx, session)
	if err != nil {
		return "", fmt.Errorf("failed to create session: %w", err)
	}

	return token, nil
}

// ExtendSession implements [Auth].
func (a *AuthUsecase) ExtendSession(ctx context.Context, id uuid.UUID) error {
	err := a.sessionRepo.UpdateExpiresAt(ctx, id, time.Now().Add(HalfSessionDuration))
	if err != nil {
		return fmt.Errorf("failed to extend session: %w", err)
	}

	return nil
}

// GetSession implements [Auth].
func (a *AuthUsecase) GetSession(ctx context.Context, id uuid.UUID) (model.Session, error) {
	session, err := a.sessionRepo.FindSessionByID(ctx, id)
	if err != nil {
		return model.Session{}, fmt.Errorf("failed to get session: %w", err)
	}

	return session, nil
}

// Login implements [Auth].
// TODO: decide whether Login should create a session
// keeping it here simplifies API, but separating gives more control (e.g. stateless auth, different session strategies)
func (a *AuthUsecase) Login(ctx context.Context, email string, password string) (uuid.UUID, error) {
	user, err := a.userRepo.FindUserByEmail(ctx, email)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to find user: %w", err)
	}

	if comparePassword(user.Password, []byte(password)) != nil {
		return uuid.Nil, fmt.Errorf("invalid password")
	}

	return user.ID, nil
}

// Register implements [Auth].
// TODO: decide whether Register should create a session
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

// GetUserIDFromToken implements [Auth].
func (a *AuthUsecase) GetUserIDFromToken(ctx context.Context, token string) (uuid.UUID, error) {
	tokenHash := hashToken([]byte(token))

	session, err := a.sessionRepo.FindSessionByToken(ctx, tokenHash)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to get session: %w", err)
	}

	return session.UserID, nil
}

func cryptoRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random string: %w", err)
	}

	return hex.EncodeToString(bytes), nil
}

func hashToken(token []byte) string {
	hashedToken := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(hashedToken[:])
	return tokenHash
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
