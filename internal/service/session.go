package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/K1ender/moqchat/internal/entity/model"
	"github.com/K1ender/moqchat/internal/repository"
	"github.com/google/uuid"
)

const (
	SessionDuration     = 30 * 24 * time.Hour
	HalfSessionDuration = SessionDuration / 2
	RandomStringLength  = 32
)

type Session interface {
	Create(ctx context.Context, userID uuid.UUID) (string, error)
	GetUserIDFromToken(ctx context.Context, token string) (uuid.UUID, error)
	Get(ctx context.Context, id uuid.UUID) (model.Session, error)
	Extend(ctx context.Context, id uuid.UUID) error
}

type SessionUsecase struct {
	sessionRepo repository.Session
}

func NewSessionUsecase(sessionRepo repository.Session) Session {
	return &SessionUsecase{sessionRepo: sessionRepo}
}

var (
	ErrSessionExpired = errors.New("session expired")
)

// Create implements [Session].
func (s *SessionUsecase) Create(ctx context.Context, userID uuid.UUID) (string, error) {
	token, err := cryptoRandomString(RandomStringLength)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	tokenHash := hashToken([]byte(token))

	session := model.Session{UserID: userID, Token: tokenHash, ExpiresAt: time.Now().Add(SessionDuration)}
	_, err = s.sessionRepo.CreateSession(ctx, session)
	if err != nil {
		return "", fmt.Errorf("failed to create session: %w", err)
	}

	return token, nil
}

// Extend implements [Session].
func (s *SessionUsecase) Extend(ctx context.Context, id uuid.UUID) error {
	session, err := s.sessionRepo.FindSessionByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to find session: %w", err)
	}

	err = s.sessionRepo.UpdateExpiresAt(ctx, id, session.ExpiresAt.Add(SessionDuration))
	if err != nil {
		return fmt.Errorf("failed to extend session: %w", err)
	}

	return nil
}

// Get implements [Session].
func (s *SessionUsecase) Get(ctx context.Context, id uuid.UUID) (model.Session, error) {
	session, err := s.sessionRepo.FindSessionByID(ctx, id)
	if err != nil {
		return model.Session{}, fmt.Errorf("failed to get session: %w", err)
	}

	return session, nil
}

// GetUserIDFromToken implements [Session].
func (s *SessionUsecase) GetUserIDFromToken(ctx context.Context, token string) (uuid.UUID, error) {
	tokenHash := hashToken([]byte(token))

	session, err := s.sessionRepo.FindSessionByToken(ctx, tokenHash)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to get session: %w", err)
	}

	if session.ExpiresAt.Before(time.Now()) {
		return uuid.Nil, ErrSessionExpired
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
	hashedToken := sha256.Sum256(token)
	tokenHash := hex.EncodeToString(hashedToken[:])
	return tokenHash
}
