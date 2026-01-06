package handler

import (
	"io"
	"net/http"
	"strconv"

	"github.com/editor-pdf/backend/internal/dto"
	"github.com/editor-pdf/backend/internal/usecase"
	"github.com/editor-pdf/backend/pkg/response"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// DefaultUserID é o UUID padrão usado quando não há autenticação
var DefaultUserID = uuid.MustParse("00000000-0000-0000-0000-000000000000")

// DocumentHandler contém os handlers de documentos
type DocumentHandler struct {
	documentUseCase  *usecase.DocumentUseCase
	previewUseCase   *usecase.PDFPreviewUseCase
	maxUploadSize    int64
}

// NewDocumentHandler cria uma nova instância de DocumentHandler
func NewDocumentHandler(
	documentUseCase *usecase.DocumentUseCase,
	previewUseCase *usecase.PDFPreviewUseCase,
	maxUploadSize int64,
) *DocumentHandler {
	return &DocumentHandler{
		documentUseCase: documentUseCase,
		previewUseCase:   previewUseCase,
		maxUploadSize:   maxUploadSize,
	}
}

// UploadDocument faz upload de um documento PDF
// @Summary Faz upload de um documento PDF
// @Description Faz upload de um arquivo PDF e cria um registro no banco
// @Tags documents
// @Security Bearer
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Arquivo PDF"
// @Success 201 {object} dto.UploadDocumentResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 413 {object} response.ErrorResponse
// @Router /api/v1/documents [post]
func (h *DocumentHandler) UploadDocument(c echo.Context) error {
	// Usa DefaultUserID quando não há autenticação
	userUUID := DefaultUserID

	// Obtém o arquivo do form
	file, err := c.FormFile("file")
	if err != nil {
		return response.ErrorBadRequest(c, err, "arquivo não fornecido")
	}

	// Valida tamanho do arquivo
	if file.Size > h.maxUploadSize {
		return response.Error(c, http.StatusRequestEntityTooLarge, nil, "arquivo muito grande")
	}

	// Abre o arquivo
	src, err := file.Open()
	if err != nil {
		return response.ErrorBadRequest(c, err, "erro ao abrir arquivo")
	}
	defer src.Close()

	// Lê o conteúdo do arquivo
	fileData, err := io.ReadAll(src)
	if err != nil {
		return response.ErrorBadRequest(c, err, "erro ao ler arquivo")
	}

	// Valida tamanho novamente após ler
	if int64(len(fileData)) > h.maxUploadSize {
		return response.Error(c, http.StatusRequestEntityTooLarge, nil, "arquivo muito grande")
	}

	// Valida MIME type
	contentType := file.Header.Get("Content-Type")
	if contentType != "" && contentType != "application/pdf" {
		return response.ErrorBadRequest(c, nil, "tipo de arquivo inválido (esperado: application/pdf)")
	}

	// Valida magic bytes (PDF deve começar com %PDF)
	if len(fileData) < 4 || string(fileData[0:4]) != "%PDF" {
		return response.ErrorBadRequest(c, nil, "arquivo não é um PDF válido (magic bytes inválidos)")
	}

	// Faz upload do documento (validação adicional será feita no UseCase)
	document, err := h.documentUseCase.UploadDocument(c.Request().Context(), userUUID, fileData, file.Filename)
	if err != nil {
		return response.ErrorInternalServer(c, err, "erro ao fazer upload do documento")
	}

	return response.SuccessCreated(c, dto.UploadDocumentResponse{
		Document: *document,
		Message:  "Documento enviado com sucesso",
	}, "Documento enviado com sucesso")
}

// ListDocuments lista documentos do usuário
// @Summary Lista documentos do usuário
// @Description Retorna uma lista paginada de documentos do usuário autenticado
// @Tags documents
// @Security Bearer
// @Produce json
// @Param limit query int false "Limite de resultados" default(20)
// @Param offset query int false "Offset para paginação" default(0)
// @Success 200 {object} dto.DocumentListResponse
// @Failure 401 {object} response.ErrorResponse
// @Router /api/v1/documents [get]
func (h *DocumentHandler) ListDocuments(c echo.Context) error {
	// Usa DefaultUserID quando não há autenticação
	userUUID := DefaultUserID

	// Parse de query parameters
	limit := 20
	offset := 0

	if limitStr := c.QueryParam("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	if offsetStr := c.QueryParam("offset"); offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	// Lista documentos
	documents, err := h.documentUseCase.ListDocuments(c.Request().Context(), userUUID, limit, offset)
	if err != nil {
		return response.ErrorInternalServer(c, err, "erro ao listar documentos")
	}

	return response.SuccessOK(c, documents)
}

// GetDocument busca um documento por ID
// @Summary Busca um documento por ID
// @Description Retorna os detalhes de um documento específico
// @Tags documents
// @Security Bearer
// @Produce json
// @Param id path string true "ID do documento"
// @Success 200 {object} dto.DocumentResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Router /api/v1/documents/{id} [get]
func (h *DocumentHandler) GetDocument(c echo.Context) error {
	// Usa DefaultUserID quando não há autenticação
	userUUID := DefaultUserID

	documentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.ErrorBadRequest(c, err, "ID de documento inválido")
	}

	// Busca documento
	document, err := h.documentUseCase.GetDocument(c.Request().Context(), documentID, userUUID)
	if err != nil {
		if err.Error() == "documento não encontrado" {
			return response.ErrorNotFound(c, err, "documento não encontrado")
		}
		if err.Error() == "acesso negado" {
			return response.ErrorForbidden(c, err, "acesso negado")
		}
		return response.ErrorInternalServer(c, err, "erro ao buscar documento")
	}

	return response.SuccessOK(c, document)
}

// ProcessDocument processa edições em um documento
// @Summary Processa edições em um documento
// @Description Aplica edições (texto, imagens, etc.) em um documento PDF
// @Tags documents
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path string true "ID do documento"
// @Param request body dto.ProcessDocumentRequest true "Instruções de edição"
// @Success 200 {object} dto.ProcessDocumentResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Router /api/v1/documents/{id}/process [post]
func (h *DocumentHandler) ProcessDocument(c echo.Context) error {
	// Usa DefaultUserID quando não há autenticação
	userUUID := DefaultUserID

	documentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.ErrorBadRequest(c, err, "ID de documento inválido")
	}

	var req dto.ProcessDocumentRequest
	if err := c.Bind(&req); err != nil {
		return response.ErrorBadRequest(c, err, "dados inválidos")
	}

	if err := c.Validate(&req); err != nil {
		return response.ErrorBadRequest(c, err, "validação falhou")
	}

	// Processa documento
	document, err := h.documentUseCase.ProcessDocument(c.Request().Context(), documentID, userUUID, req.Instructions)
	if err != nil {
		if err.Error() == "documento não encontrado" {
			return response.ErrorNotFound(c, err, "documento não encontrado")
		}
		if err.Error() == "acesso negado" {
			return response.ErrorForbidden(c, err, "acesso negado")
		}
		return response.ErrorInternalServer(c, err, "erro ao processar documento")
	}

	return response.SuccessOK(c, dto.ProcessDocumentResponse{
		Document: *document,
		Message:  "Documento processado com sucesso",
	})
}

// GeneratePreview gera preview de uma página do documento
// @Summary Gera preview de uma página
// @Description Retorna uma imagem (preview) de uma página específica do PDF
// @Tags documents
// @Security Bearer
// @Produce image/png
// @Param id path string true "ID do documento"
// @Param page path int true "Número da página"
// @Success 200 {file} binary
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Router /api/v1/documents/{id}/preview/{page} [get]
func (h *DocumentHandler) GeneratePreview(c echo.Context) error {
	// Usa DefaultUserID quando não há autenticação
	userUUID := DefaultUserID

	documentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.ErrorBadRequest(c, err, "ID de documento inválido")
	}

	pageNum, err := strconv.Atoi(c.Param("page"))
	if err != nil || pageNum < 1 {
		return response.ErrorBadRequest(c, err, "número de página inválido")
	}

	// Gera preview
	preview, err := h.previewUseCase.GeneratePreview(c.Request().Context(), documentID, userUUID, pageNum)
	if err != nil {
		if err.Error() == "documento não encontrado" {
			return response.ErrorNotFound(c, err, "documento não encontrado")
		}
		if err.Error() == "acesso negado" {
			return response.ErrorForbidden(c, err, "acesso negado")
		}
		return response.ErrorInternalServer(c, err, "erro ao gerar preview")
	}

	c.Response().Header().Set("Content-Type", "image/png")
	return c.Blob(http.StatusOK, "image/png", preview)
}

// DeleteDocument remove um documento
// @Summary Remove um documento
// @Description Remove um documento e seu arquivo associado
// @Tags documents
// @Security Bearer
// @Produce json
// @Param id path string true "ID do documento"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Router /api/v1/documents/{id} [delete]
func (h *DocumentHandler) DeleteDocument(c echo.Context) error {
	// Usa DefaultUserID quando não há autenticação
	userUUID := DefaultUserID

	documentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.ErrorBadRequest(c, err, "ID de documento inválido")
	}

	// Deleta documento
	if err := h.documentUseCase.DeleteDocument(c.Request().Context(), documentID, userUUID); err != nil {
		if err.Error() == "documento não encontrado" {
			return response.ErrorNotFound(c, err, "documento não encontrado")
		}
		if err.Error() == "acesso negado" {
			return response.ErrorForbidden(c, err, "acesso negado")
		}
		return response.ErrorInternalServer(c, err, "erro ao deletar documento")
	}

	return response.SuccessOK(c, map[string]string{"message": "Documento deletado com sucesso"})
}
