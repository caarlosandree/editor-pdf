package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/editor-pdf/backend/internal/domain"
	"github.com/editor-pdf/backend/pkg/logger"
	"go.uber.org/zap"
)

// LocalStorage implementa FileStorage usando armazenamento local em disco
type LocalStorage struct {
	basePath string
	baseURL  string
}

// NewLocalStorage cria uma nova instância de LocalStorage
func NewLocalStorage(basePath, baseURL string) (domain.FileStorage, error) {
	// Cria o diretório base se não existir
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("erro ao criar diretório de storage: %w", err)
	}

	return &LocalStorage{
		basePath: basePath,
		baseURL:  baseURL,
	}, nil
}

// Save salva um arquivo e retorna o caminho relativo
func (s *LocalStorage) Save(ctx context.Context, data []byte, filename string) (string, error) {
	// Sanitiza o nome do arquivo
	filename = sanitizeFilename(filename)

	// Gera caminho completo
	fullPath := filepath.Join(s.basePath, filename)

	// Cria diretórios necessários
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("erro ao criar diretório: %w", err)
	}

	// Salva o arquivo
	if err := os.WriteFile(fullPath, data, 0644); err != nil {
		return "", fmt.Errorf("erro ao salvar arquivo: %w", err)
	}

	logger.Logger.Debug("Arquivo salvo",
		zap.String("path", fullPath),
		zap.Int("size", len(data)),
	)

	// Retorna caminho relativo
	return filename, nil
}

// Read lê um arquivo do storage
func (s *LocalStorage) Read(ctx context.Context, filePath string) ([]byte, error) {
	// Sanitiza o caminho para prevenir path traversal
	filePath = sanitizePath(filePath)

	fullPath := filepath.Join(s.basePath, filePath)

	// Verifica se o arquivo existe
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("arquivo não encontrado: %s", filePath)
	}

	data, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler arquivo: %w", err)
	}

	return data, nil
}

// Delete remove um arquivo do storage
func (s *LocalStorage) Delete(ctx context.Context, filePath string) error {
	// Sanitiza o caminho
	filePath = sanitizePath(filePath)

	fullPath := filepath.Join(s.basePath, filePath)

	if err := os.Remove(fullPath); err != nil {
		if os.IsNotExist(err) {
			return nil // Arquivo já não existe, não é erro
		}
		return fmt.Errorf("erro ao deletar arquivo: %w", err)
	}

	logger.Logger.Debug("Arquivo deletado", zap.String("path", fullPath))
	return nil
}

// Exists verifica se um arquivo existe no storage
func (s *LocalStorage) Exists(ctx context.Context, filePath string) (bool, error) {
	// Sanitiza o caminho
	filePath = sanitizePath(filePath)

	fullPath := filepath.Join(s.basePath, filePath)

	_, err := os.Stat(fullPath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, fmt.Errorf("erro ao verificar existência do arquivo: %w", err)
}

// GetURL retorna a URL completa para acessar o arquivo
func (s *LocalStorage) GetURL(ctx context.Context, filePath string) (string, error) {
	// Sanitiza o caminho
	filePath = sanitizePath(filePath)

	// Remove barras iniciais
	filePath = strings.TrimPrefix(filePath, "/")

	// Constrói URL
	url := fmt.Sprintf("%s/%s", strings.TrimSuffix(s.baseURL, "/"), filePath)
	return url, nil
}

// sanitizeFilename remove caracteres perigosos do nome do arquivo
func sanitizeFilename(filename string) string {
	// Remove caracteres perigosos
	filename = strings.ReplaceAll(filename, "..", "")
	filename = strings.ReplaceAll(filename, "/", "_")
	filename = strings.ReplaceAll(filename, "\\", "_")
	filename = strings.TrimSpace(filename)

	return filename
}

// sanitizePath previne path traversal attacks
func sanitizePath(path string) string {
	// Remove ".." e normaliza o caminho
	path = filepath.Clean(path)
	path = strings.TrimPrefix(path, "/")
	return path
}

// WriteFile escreve dados em um arquivo de forma segura
func (s *LocalStorage) WriteFile(ctx context.Context, filePath string, reader io.Reader) error {
	// Sanitiza o caminho
	filePath = sanitizePath(filePath)

	fullPath := filepath.Join(s.basePath, filePath)

	// Cria diretórios necessários
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório: %w", err)
	}

	// Cria o arquivo
	file, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("erro ao criar arquivo: %w", err)
	}
	defer file.Close()

	// Copia dados do reader para o arquivo
	if _, err := io.Copy(file, reader); err != nil {
		return fmt.Errorf("erro ao escrever arquivo: %w", err)
	}

	return nil
}
