package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/K1ender/moqchat/internal/entity/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Session interface {
	FindSessionByToken(ctx context.Context, token string) (model.Session, error)
	CreateSession(ctx context.Context, session model.Session) (uuid.UUID, error)
	DeleteSession(ctx context.Context, id uuid.UUID) error
	UpdateExpiresAt(ctx context.Context, id uuid.UUID, expiresAt time.Time) error
	FindSessionByID(ctx context.Context, id uuid.UUID) (model.Session, error)
}

type SessionPostgres struct {
	conn *pgxpool.Pool
}

func NewSessionPostgres(conn *pgxpool.Pool) Session {
	return &SessionPostgres{conn: conn}
}

// CreateSession implements [Session].
func (s *SessionPostgres) CreateSession(ctx context.Context, session model.Session) (uuid.UUID, error) {
	tx, err := s.conn.Begin(ctx)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer safeRollback(ctx, tx)

	query := `INSERT INTO sessions(user_id, token, expires_at) VALUES ($1, $2, $3) RETURNING id`

	var id uuid.UUID
	err = tx.QueryRow(ctx, query, session.UserID, session.Token, session.ExpiresAt).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create session: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return id, nil
}

// DeleteSession implements [Session].
func (s *SessionPostgres) DeleteSession(ctx context.Context, id uuid.UUID) error {
	tx, err := s.conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer safeRollback(ctx, tx)

	query := `DELETE FROM sessions WHERE id = $1`

	_, err = tx.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// FindSessionByToken implements [Session].
func (s *SessionPostgres) FindSessionByToken(ctx context.Context, token string) (model.Session, error) {
	tx, err := s.conn.Begin(ctx)
	if err != nil {
		return model.Session{}, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer safeRollback(ctx, tx)

	query := `SELECT id, user_id, token, expires_at, created_at FROM sessions WHERE token = $1`

	var session model.Session
	err = tx.QueryRow(ctx, query, token).
		Scan(
			&session.ID,
			&session.UserID,
			&session.Token,
			&session.ExpiresAt,
			&session.CreatedAt,
		)
	if err != nil {
		return model.Session{}, fmt.Errorf("failed to find session: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return model.Session{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return session, nil
}

// UpdateExpiresAt implements [Session].
func (s *SessionPostgres) UpdateExpiresAt(ctx context.Context, id uuid.UUID, expiresAt time.Time) error {
	tx, err := s.conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer safeRollback(ctx, tx)

	query := `UPDATE sessions SET expires_at = $1 WHERE id = $2`

	_, err = tx.Exec(ctx, query, expiresAt, id)
	if err != nil {
		return fmt.Errorf("failed to update session: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// FindSessionByID implements [Session].
func (s *SessionPostgres) FindSessionByID(ctx context.Context, id uuid.UUID) (model.Session, error) {
	tx, err := s.conn.Begin(ctx)
	if err != nil {
		return model.Session{}, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer safeRollback(ctx, tx)

	query := `SELECT id, user_id, token, expires_at, created_at FROM sessions WHERE id = $1`

	var session model.Session
	err = tx.QueryRow(ctx, query, id).
		Scan(&session.ID, &session.UserID, &session.Token, &session.ExpiresAt, &session.CreatedAt)
	if err != nil {
		return model.Session{}, fmt.Errorf("failed to find session: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return model.Session{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return session, nil
}
