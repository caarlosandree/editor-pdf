package repository

import (
	"context"
	"encoding/json"

	"github.com/editor-pdf/backend/internal/domain"
	"github.com/editor-pdf/backend/internal/model"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// auditLogRepository implementa AuditLogRepository usando sqlx
type auditLogRepository struct {
	db *sqlx.DB
}

// NewAuditLogRepository cria uma nova instância de AuditLogRepository
func NewAuditLogRepository(db *sqlx.DB) domain.AuditLogRepository {
	return &auditLogRepository{db: db}
}

// Create cria um novo log de auditoria
func (r *auditLogRepository) Create(ctx context.Context, log *model.AuditLog) error {
	query := `
		INSERT INTO audit_logs (id, document_id, user_id, action, metadata, created_at)
		VALUES (:id, :document_id, :user_id, :action, :metadata, :created_at)
	`

	if log.ID == uuid.Nil {
		log.ID = uuid.New()
	}

	// Se Metadata for nil, converte para JSON vazio
	if log.Metadata == nil {
		log.Metadata = json.RawMessage("{}")
	}

	_, err := r.db.NamedExecContext(ctx, query, log)
	if err != nil {
		return err
	}

	return nil
}

// FindByDocumentID busca logs de auditoria de um documento
func (r *auditLogRepository) FindByDocumentID(ctx context.Context, documentID uuid.UUID, limit, offset int) ([]*model.AuditLog, int, error) {
	var logs []*model.AuditLog
	query := `
		SELECT id, document_id, user_id, action, metadata, created_at 
		FROM audit_logs 
		WHERE document_id = $1 
		ORDER BY created_at DESC 
		LIMIT $2 OFFSET $3
	`

	err := r.db.SelectContext(ctx, &logs, query, documentID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	// Conta total de logs
	var total int
	countQuery := `SELECT COUNT(*) FROM audit_logs WHERE document_id = $1`
	err = r.db.GetContext(ctx, &total, countQuery, documentID)
	if err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}
