package domain

import "context"

// FileStorage define a interface para armazenamento de arquivos
type FileStorage interface {
	// Save salva um arquivo e retorna o caminho relativo
	Save(ctx context.Context, data []byte, filename string) (string, error)

	// Read lê um arquivo do storage
	Read(ctx context.Context, filePath string) ([]byte, error)

	// Delete remove um arquivo do storage
	Delete(ctx context.Context, filePath string) error

	// Exists verifica se um arquivo existe no storage
	Exists(ctx context.Context, filePath string) (bool, error)

	// GetURL retorna a URL completa para acessar o arquivo
	GetURL(ctx context.Context, filePath string) (string, error)
}
