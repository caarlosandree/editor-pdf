package usecase

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/editor-pdf/backend/internal/domain"
	"github.com/editor-pdf/backend/internal/dto"
	"github.com/editor-pdf/backend/internal/model"
	"github.com/editor-pdf/backend/pkg/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// DocumentUseCase contém os casos de uso de documentos
type DocumentUseCase struct {
	documentRepo    domain.DocumentRepository
	auditLogRepo    domain.AuditLogRepository
	fileStorage     domain.FileStorage
	pdfProcessor    domain.PDFProcessor
	storageBasePath string
}

// NewDocumentUseCase cria uma nova instância de DocumentUseCase
func NewDocumentUseCase(
	documentRepo domain.DocumentRepository,
	auditLogRepo domain.AuditLogRepository,
	fileStorage domain.FileStorage,
	pdfProcessor domain.PDFProcessor,
	storageBasePath string,
) *DocumentUseCase {
	return &DocumentUseCase{
		documentRepo:    documentRepo,
		auditLogRepo:    auditLogRepo,
		fileStorage:     fileStorage,
		pdfProcessor:    pdfProcessor,
		storageBasePath: storageBasePath,
	}
}

// UploadDocument faz upload de um documento PDF
func (uc *DocumentUseCase) UploadDocument(ctx context.Context, userID uuid.UUID, fileData []byte, filename string) (*dto.DocumentResponse, error) {
	// Valida o PDF usando magic bytes
	if err := uc.pdfProcessor.ValidatePDF(ctx, fileData); err != nil {
		return nil, fmt.Errorf("arquivo PDF inválido: %w", err)
	}

	// Calcula checksum
	hash := sha256.Sum256(fileData)
	checksum := hex.EncodeToString(hash[:])

	// Gera nome único para o arquivo
	fileID := uuid.New()
	ext := filepath.Ext(filename)
	if ext == "" {
		ext = ".pdf"
	}
	safeFilename := fmt.Sprintf("%s%s", fileID.String(), ext)

	// Salva o arquivo no storage temporariamente para extrair informações
	tempPath, err := uc.fileStorage.Save(ctx, fileData, "temp_"+safeFilename)
	if err != nil {
		return nil, fmt.Errorf("erro ao salvar arquivo temporário: %w", err)
	}

	// Extrai informações das páginas
	fullPath := filepath.Join(uc.storageBasePath, tempPath)
	pages, err := uc.pdfProcessor.ExtractPages(ctx, fullPath)
	if err != nil {
		// Se não conseguir extrair páginas, continua com 0
		logger.Logger.Warn("Erro ao extrair páginas do PDF", zap.Error(err))
		pages = []model.Page{}
	}

	// Remove arquivo temporário
	_ = uc.fileStorage.Delete(ctx, tempPath)

	// Salva o arquivo no storage com nome final
	filePath, err := uc.fileStorage.Save(ctx, fileData, safeFilename)
	if err != nil {
		return nil, fmt.Errorf("erro ao salvar arquivo: %w", err)
	}

	// Cria registro no banco
	document := &model.Document{
		ID:        fileID,
		UserID:    userID,
		FilePath:  filePath,
		Checksum:  checksum,
		Version:   1,
		Status:    model.DocumentStatusReady,
		PageCount: len(pages),
	}

	if err := uc.documentRepo.Create(ctx, document); err != nil {
		// Tenta remover o arquivo se falhar ao criar registro
		_ = uc.fileStorage.Delete(ctx, filePath)
		return nil, fmt.Errorf("erro ao criar registro do documento: %w", err)
	}

	// Cria log de auditoria
	uc.createAuditLog(ctx, document.ID, userID, "UPLOAD", map[string]interface{}{
		"filename": filename,
		"size":     len(fileData),
	})

	// Obtém URL do arquivo
	fileURL, _ := uc.fileStorage.GetURL(ctx, filePath)

	return uc.toDocumentResponse(document, fileURL), nil
}

// GetDocument busca um documento por ID
func (uc *DocumentUseCase) GetDocument(ctx context.Context, documentID, userID uuid.UUID) (*dto.DocumentResponse, error) {
	document, err := uc.documentRepo.FindByID(ctx, documentID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar documento: %w", err)
	}

	if document == nil {
		return nil, errors.New("documento não encontrado")
	}

	// Obtém URL do arquivo
	fileURL, _ := uc.fileStorage.GetURL(ctx, document.FilePath)

	return uc.toDocumentResponse(document, fileURL), nil
}

// ListDocuments lista documentos de um usuário
func (uc *DocumentUseCase) ListDocuments(ctx context.Context, userID uuid.UUID, limit, offset int) (*dto.DocumentListResponse, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	documents, total, err := uc.documentRepo.FindByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar documentos: %w", err)
	}

	responses := make([]dto.DocumentResponse, 0, len(documents))
	for _, doc := range documents {
		fileURL, _ := uc.fileStorage.GetURL(ctx, doc.FilePath)
		responses = append(responses, *uc.toDocumentResponse(doc, fileURL))
	}

	return &dto.DocumentListResponse{
		Documents: responses,
		Total:     total,
		Limit:     limit,
		Offset:    offset,
	}, nil
}

// ProcessDocument processa edições em um documento
func (uc *DocumentUseCase) ProcessDocument(ctx context.Context, documentID, userID uuid.UUID, instructions []dto.EditInstruction) (*dto.DocumentResponse, error) {
	// Busca o documento
	document, err := uc.documentRepo.FindByID(ctx, documentID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar documento: %w", err)
	}

	if document == nil {
		return nil, errors.New("documento não encontrado")
	}

	// Lê o PDF original
	pdfData, err := uc.fileStorage.Read(ctx, document.FilePath)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler PDF original: %w", err)
	}

	// Valida o PDF
	if err := uc.pdfProcessor.ValidatePDF(ctx, pdfData); err != nil {
		return nil, fmt.Errorf("PDF inválido: %w", err)
	}

	// Salva PDF temporário para processamento
	tempInputPath := filepath.Join(uc.storageBasePath, fmt.Sprintf("temp_%s_input.pdf", documentID.String()))
	tempPath, err := uc.fileStorage.Save(ctx, pdfData, filepath.Base(tempInputPath))
	if err != nil {
		return nil, fmt.Errorf("erro ao salvar PDF temporário: %w", err)
	}
	defer uc.fileStorage.Delete(ctx, tempPath)

	// Gera caminho de saída
	newVersion := document.Version + 1
	outputFilename := fmt.Sprintf("%s_v%d.pdf", documentID.String(), newVersion)
	outputPath, err := uc.fileStorage.Save(ctx, pdfData, outputFilename)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar novo arquivo: %w", err)
	}

	// Obtém caminho completo do arquivo temporário e de saída
	fullTempPath := filepath.Join(uc.storageBasePath, tempPath)
	fullOutputPath := filepath.Join(uc.storageBasePath, outputPath)

	// Processa cada edição sequencialmente
	for i, instruction := range instructions {
		switch instruction.Type {
		case "text":
			// Valida campos obrigatórios
			if instruction.Content == "" {
				_ = uc.fileStorage.Delete(ctx, outputPath)
				return nil, fmt.Errorf("edição %d: conteúdo de texto não pode ser vazio", i+1)
			}

			fontSize := 12.0 // Tamanho padrão
			if instruction.FontSize != nil && *instruction.FontSize > 0 {
				fontSize = *instruction.FontSize
			}

			// Aplica a edição de texto
			if err := uc.pdfProcessor.AddText(ctx, fullTempPath, instruction.Page, instruction.X, instruction.Y, instruction.Content, fontSize); err != nil {
				_ = uc.fileStorage.Delete(ctx, outputPath)
				return nil, fmt.Errorf("erro ao adicionar texto na edição %d: %w", i+1, err)
			}

		case "image":
			// Valida campos obrigatórios
			if instruction.Content == "" {
				_ = uc.fileStorage.Delete(ctx, outputPath)
				return nil, fmt.Errorf("edição %d: caminho da imagem não pode ser vazio", i+1)
			}
			if instruction.Width == nil || *instruction.Width <= 0 {
				_ = uc.fileStorage.Delete(ctx, outputPath)
				return nil, fmt.Errorf("edição %d: largura da imagem deve ser maior que zero", i+1)
			}
			if instruction.Height == nil || *instruction.Height <= 0 {
				_ = uc.fileStorage.Delete(ctx, outputPath)
				return nil, fmt.Errorf("edição %d: altura da imagem deve ser maior que zero", i+1)
			}

			// Resolve caminho da imagem (pode ser relativo ou absoluto)
			imagePath := instruction.Content
			if !filepath.IsAbs(imagePath) {
				imagePath = filepath.Join(uc.storageBasePath, imagePath)
			}

			// Aplica a edição de imagem
			if err := uc.pdfProcessor.AddImage(ctx, fullTempPath, instruction.Page, instruction.X, instruction.Y, *instruction.Width, *instruction.Height, imagePath); err != nil {
				_ = uc.fileStorage.Delete(ctx, outputPath)
				return nil, fmt.Errorf("erro ao adicionar imagem na edição %d: %w", i+1, err)
			}

		case "drawing":
			_ = uc.fileStorage.Delete(ctx, outputPath)
			return nil, fmt.Errorf("edição %d: tipo 'drawing' ainda não está implementado", i+1)

		default:
			_ = uc.fileStorage.Delete(ctx, outputPath)
			return nil, fmt.Errorf("edição %d: tipo de edição desconhecido: %s", i+1, instruction.Type)
		}
	}

	// Copia o arquivo processado para o caminho de saída
	processedData, err := os.ReadFile(fullTempPath)
	if err != nil {
		_ = uc.fileStorage.Delete(ctx, outputPath)
		return nil, fmt.Errorf("erro ao ler PDF processado: %w", err)
	}

	if err := os.WriteFile(fullOutputPath, processedData, 0644); err != nil {
		_ = uc.fileStorage.Delete(ctx, outputPath)
		return nil, fmt.Errorf("erro ao salvar PDF processado: %w", err)
	}

	logger.Logger.Info("Edições processadas com sucesso",
		zap.String("document_id", documentID.String()),
		zap.Int("instructions_count", len(instructions)),
		zap.Int("new_version", newVersion),
	)

	// Incrementa versão
	if err := uc.documentRepo.IncrementVersion(ctx, documentID); err != nil {
		_ = uc.fileStorage.Delete(ctx, outputPath)
		return nil, fmt.Errorf("erro ao incrementar versão: %w", err)
	}

	// Atualiza caminho do arquivo
	document.FilePath = outputPath
	document.Version = newVersion
	if err := uc.documentRepo.Update(ctx, document); err != nil {
		_ = uc.fileStorage.Delete(ctx, outputPath)
		return nil, fmt.Errorf("erro ao atualizar documento: %w", err)
	}

	// Cria log de auditoria
	uc.createAuditLog(ctx, documentID, userID, "PROCESS", map[string]interface{}{
		"instructions_count": len(instructions),
		"new_version":        newVersion,
	})

	// Obtém URL do arquivo
	fileURL, _ := uc.fileStorage.GetURL(ctx, document.FilePath)

	return uc.toDocumentResponse(document, fileURL), nil
}

// DeleteDocument remove um documento
func (uc *DocumentUseCase) DeleteDocument(ctx context.Context, documentID, userID uuid.UUID) error {
	// Busca o documento
	document, err := uc.documentRepo.FindByID(ctx, documentID)
	if err != nil {
		return fmt.Errorf("erro ao buscar documento: %w", err)
	}

	if document == nil {
		return errors.New("documento não encontrado")
	}

	// Remove o arquivo do storage
	if err := uc.fileStorage.Delete(ctx, document.FilePath); err != nil {
		logger.Logger.Warn("Erro ao deletar arquivo do storage", zap.Error(err))
	}

	// Remove do banco
	if err := uc.documentRepo.Delete(ctx, documentID); err != nil {
		return fmt.Errorf("erro ao deletar documento: %w", err)
	}

	// Cria log de auditoria
	uc.createAuditLog(ctx, documentID, userID, "DELETE", nil)

	return nil
}

// toDocumentResponse converte model.Document para dto.DocumentResponse
func (uc *DocumentUseCase) toDocumentResponse(doc *model.Document, fileURL string) *dto.DocumentResponse {
	return &dto.DocumentResponse{
		ID:        doc.ID.String(),
		UserID:    doc.UserID.String(),
		FilePath:  doc.FilePath,
		FileURL:   fileURL,
		Checksum:  doc.Checksum,
		Version:   doc.Version,
		Status:    string(doc.Status),
		PageCount: doc.PageCount,
		CreatedAt: doc.CreatedAt,
		UpdatedAt: doc.UpdatedAt,
	}
}

// createAuditLog cria um log de auditoria
func (uc *DocumentUseCase) createAuditLog(ctx context.Context, documentID, userID uuid.UUID, action string, metadata map[string]interface{}) {
	log := &model.AuditLog{
		ID:         uuid.New(),
		DocumentID: documentID,
		UserID:     userID,
		Action:     action,
		CreatedAt:  time.Now(),
	}

	if metadata != nil {
		// Serializa metadata para JSON
		metadataJSON, err := json.Marshal(metadata)
		if err == nil {
			log.Metadata = metadataJSON
		}
	}

	if err := uc.auditLogRepo.Create(ctx, log); err != nil {
		logger.Logger.Warn("Erro ao criar log de auditoria", zap.Error(err))
	}
}
