package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/editor-pdf/backend/internal/domain"
	"github.com/editor-pdf/backend/internal/model"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// userRepository implementa UserRepository usando sqlx
type userRepository struct {
	db *sqlx.DB
}

// NewUserRepository cria uma nova instância de UserRepository
func NewUserRepository(db *sqlx.DB) domain.UserRepository {
	return &userRepository{db: db}
}

// Create cria um novo usuário
func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	query := `
		INSERT INTO users (id, email, password_hash, name, created_at, updated_at)
		VALUES (:id, :email, :password_hash, :name, :created_at, :updated_at)
	`

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}

	_, err := r.db.NamedExecContext(ctx, query, user)
	if err != nil {
		return err
	}

	return nil
}

// FindByID busca um usuário por ID
func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	var user model.User
	query := `SELECT id, email, password_hash, name, created_at, updated_at FROM users WHERE id = $1`

	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

// FindByEmail busca um usuário por email
func (r *userRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	query := `SELECT id, email, password_hash, name, created_at, updated_at FROM users WHERE email = $1`

	err := r.db.GetContext(ctx, &user, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

// Update atualiza um usuário
func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	query := `
		UPDATE users 
		SET email = :email, password_hash = :password_hash, name = :name, updated_at = :updated_at
		WHERE id = :id
	`

	user.UpdatedAt = time.Now()

	_, err := r.db.NamedExecContext(ctx, query, user)
	return err
}

// Delete remove um usuário
func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
