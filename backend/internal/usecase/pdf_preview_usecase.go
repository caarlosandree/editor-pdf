package usecase

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/editor-pdf/backend/internal/domain"
	"github.com/editor-pdf/backend/pkg/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// PDFPreviewUseCase contém os casos de uso para preview de PDF
type PDFPreviewUseCase struct {
	documentRepo    domain.DocumentRepository
	pdfProcessor    domain.PDFProcessor
	fileStorage     domain.FileStorage
	storageBasePath string
}

// NewPDFPreviewUseCase cria uma nova instância de PDFPreviewUseCase
func NewPDFPreviewUseCase(
	documentRepo domain.DocumentRepository,
	pdfProcessor domain.PDFProcessor,
	fileStorage domain.FileStorage,
	storageBasePath string,
) *PDFPreviewUseCase {
	return &PDFPreviewUseCase{
		documentRepo:    documentRepo,
		pdfProcessor:    pdfProcessor,
		fileStorage:     fileStorage,
		storageBasePath: storageBasePath,
	}
}

// GeneratePreview gera uma preview (imagem) de uma página específica do PDF
func (uc *PDFPreviewUseCase) GeneratePreview(ctx context.Context, documentID, userID uuid.UUID, pageNum int) ([]byte, error) {
	// Busca o documento
	document, err := uc.documentRepo.FindByID(ctx, documentID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar documento: %w", err)
	}

	if document == nil {
		return nil, errors.New("documento não encontrado")
	}

	// Valida número da página
	if pageNum < 1 || pageNum > document.PageCount {
		return nil, fmt.Errorf("página inválida: %d (documento tem %d páginas)", pageNum, document.PageCount)
	}

	// Obtém o caminho completo do arquivo PDF
	fullFilePath := filepath.Join(uc.storageBasePath, document.FilePath)

	// Gera a preview usando o PDFProcessor
	previewBytes, err := uc.pdfProcessor.GeneratePreview(ctx, fullFilePath, pageNum)
	if err != nil {
		return nil, fmt.Errorf("erro ao gerar preview: %w", err)
	}

	logger.Logger.Debug("Preview gerado com sucesso",
		zap.String("document_id", documentID.String()),
		zap.Int("page", pageNum),
		zap.Int("size_bytes", len(previewBytes)),
	)

	return previewBytes, nil
}
