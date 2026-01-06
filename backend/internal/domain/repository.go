package domain

import (
	"context"

	"github.com/editor-pdf/backend/internal/model"
	"github.com/google/uuid"
)

// UserRepository define a interface para operações de usuário no banco de dados
type UserRepository interface {
	// Create cria um novo usuário
	Create(ctx context.Context, user *model.User) error

	// FindByID busca um usuário por ID
	FindByID(ctx context.Context, id uuid.UUID) (*model.User, error)

	// FindByEmail busca um usuário por email
	FindByEmail(ctx context.Context, email string) (*model.User, error)

	// Update atualiza um usuário
	Update(ctx context.Context, user *model.User) error

	// Delete remove um usuário
	Delete(ctx context.Context, id uuid.UUID) error
}

// DocumentRepository define a interface para operações de documento no banco de dados
type DocumentRepository interface {
	// Create cria um novo documento
	Create(ctx context.Context, document *model.Document) error

	// FindByID busca um documento por ID
	FindByID(ctx context.Context, id uuid.UUID) (*model.Document, error)

	// FindByUserID busca todos os documentos de um usuário
	FindByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*model.Document, int, error)

	// Update atualiza um documento
	Update(ctx context.Context, document *model.Document) error

	// Delete remove um documento
	Delete(ctx context.Context, id uuid.UUID) error

	// IncrementVersion incrementa a versão de um documento
	IncrementVersion(ctx context.Context, id uuid.UUID) error
}

// AuditLogRepository define a interface para operações de log de auditoria
type AuditLogRepository interface {
	// Create cria um novo log de auditoria
	Create(ctx context.Context, log *model.AuditLog) error

	// FindByDocumentID busca logs de auditoria de um documento
	FindByDocumentID(ctx context.Context, documentID uuid.UUID, limit, offset int) ([]*model.AuditLog, int, error)
}
