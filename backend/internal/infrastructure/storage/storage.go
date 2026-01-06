package storage

import (
	"github.com/editor-pdf/backend/internal/config"
	"github.com/editor-pdf/backend/internal/domain"
)

// NewFileStorage cria uma instância de FileStorage baseada na configuração
// Por enquanto, apenas suporta LocalStorage, mas está preparado para futuras implementações (S3, etc.)
func NewFileStorage(cfg *config.Config) (domain.FileStorage, error) {
	baseURL := cfg.Server.Host + ":" + cfg.Server.Port + "/files"
	return NewLocalStorage(cfg.Storage.Path, baseURL)
}
