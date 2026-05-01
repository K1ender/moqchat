package repository

import (
	"context"
	"fmt"

	"github.com/K1ender/moqchat/internal/entity/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type User interface {
	CreateUser(ctx context.Context, user model.User) (uuid.UUID, error)
	FindUserByEmail(ctx context.Context, email string) (model.User, error)
}

type UserPostgres struct {
	conn *pgxpool.Pool
}

func NewUserPostgres(conn *pgxpool.Pool) User {
	return &UserPostgres{
		conn: conn,
	}
}

// CreateUser implements [User].
func (r *UserPostgres) CreateUser(ctx context.Context, user model.User) (uuid.UUID, error) {
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer safeRollback(ctx, tx)

	query := `INSERT INTO users(username, email, password) VALUES ($1, $2, $3) RETURNING id`

	var id uuid.UUID
	err = tx.QueryRow(ctx, query, user.Username, user.Email, user.Password).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create user: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return id, nil
}

// FindUserByEmail implements [User].
func (r *UserPostgres) FindUserByEmail(ctx context.Context, email string) (model.User, error) {
	var user model.User
	query := `SELECT id, username, email, password, created_at FROM users WHERE email = $1`

	err := r.conn.QueryRow(ctx, query, email).
		Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		return model.User{}, fmt.Errorf("failed to find user: %w", err)
	}

	return user, nil
}
