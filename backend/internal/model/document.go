package model

import (
	"time"

	"github.com/google/uuid"
)

// DocumentStatus representa o status de um documento
type DocumentStatus string

const (
	DocumentStatusProcessing DocumentStatus = "PROCESSING"
	DocumentStatusReady      DocumentStatus = "READY"
	DocumentStatusError      DocumentStatus = "ERROR"
)

// Document representa um documento PDF
type Document struct {
	ID        uuid.UUID      `db:"id"`
	UserID    uuid.UUID      `db:"user_id"`
	FilePath  string         `db:"file_path"`
	Checksum  string         `db:"checksum"`
	Version   int            `db:"version"`
	Status    DocumentStatus `db:"status"`
	PageCount int            `db:"page_count"`
	CreatedAt time.Time      `db:"created_at"`
	UpdatedAt time.Time      `db:"updated_at"`
}
