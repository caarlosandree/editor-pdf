package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// AuditLog representa um log de auditoria (anotação de ações)
type AuditLog struct {
	ID         uuid.UUID       `db:"id"`
	DocumentID uuid.UUID       `db:"document_id"`
	UserID     uuid.UUID       `db:"user_id"`
	Action     string          `db:"action"`
	Metadata   json.RawMessage `db:"metadata"`
	CreatedAt  time.Time       `db:"created_at"`
}
