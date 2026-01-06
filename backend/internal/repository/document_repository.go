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

// documentRepository implementa DocumentRepository usando sqlx
type documentRepository struct {
	db *sqlx.DB
}

// NewDocumentRepository cria uma nova instância de DocumentRepository
func NewDocumentRepository(db *sqlx.DB) domain.DocumentRepository {
	return &documentRepository{db: db}
}

// Create cria um novo documento
func (r *documentRepository) Create(ctx context.Context, document *model.Document) error {
	query := `
		INSERT INTO documents (id, user_id, file_path, checksum, version, status, page_count, created_at, updated_at)
		VALUES (:id, :user_id, :file_path, :checksum, :version, :status, :page_count, :created_at, :updated_at)
	`

	now := time.Now()
	document.CreatedAt = now
	document.UpdatedAt = now

	if document.ID == uuid.Nil {
		document.ID = uuid.New()
	}

	if document.Version == 0 {
		document.Version = 1
	}

	if document.Status == "" {
		document.Status = model.DocumentStatusReady
	}

	_, err := r.db.NamedExecContext(ctx, query, document)
	if err != nil {
		return err
	}

	return nil
}

// FindByID busca um documento por ID
func (r *documentRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Document, error) {
	var document model.Document
	query := `
		SELECT id, user_id, file_path, checksum, version, status, page_count, created_at, updated_at 
		FROM documents 
		WHERE id = $1
	`

	err := r.db.GetContext(ctx, &document, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &document, nil
}

// FindByUserID busca todos os documentos de um usuário
// NOTA: Como não há autenticação, lista todos os documentos (ignora userID)
func (r *documentRepository) FindByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*model.Document, int, error) {
	var documents []*model.Document
	query := `
		SELECT id, user_id, file_path, checksum, version, status, page_count, created_at, updated_at 
		FROM documents 
		ORDER BY created_at DESC 
		LIMIT $1 OFFSET $2
	`

	err := r.db.SelectContext(ctx, &documents, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	// Conta total de documentos
	var total int
	countQuery := `SELECT COUNT(*) FROM documents`
	err = r.db.GetContext(ctx, &total, countQuery)
	if err != nil {
		return nil, 0, err
	}

	return documents, total, nil
}

// Update atualiza um documento
func (r *documentRepository) Update(ctx context.Context, document *model.Document) error {
	query := `
		UPDATE documents 
		SET file_path = :file_path, checksum = :checksum, version = :version, 
		    status = :status, page_count = :page_count, updated_at = :updated_at
		WHERE id = :id
	`

	document.UpdatedAt = time.Now()

	_, err := r.db.NamedExecContext(ctx, query, document)
	return err
}

// Delete remove um documento
func (r *documentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM documents WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// IncrementVersion incrementa a versão de um documento
func (r *documentRepository) IncrementVersion(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE documents 
		SET version = version + 1, updated_at = $1
		WHERE id = $2
	`

	_, err := r.db.ExecContext(ctx, query, time.Now(), id)
	return err
}
