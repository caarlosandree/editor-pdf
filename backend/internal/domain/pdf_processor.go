package domain

import (
	"context"

	"github.com/editor-pdf/backend/internal/model"
)

// PDFProcessor define a interface para processamento de arquivos PDF
type PDFProcessor interface {
	// ExtractPages extrai informações sobre as páginas de um PDF
	ExtractPages(ctx context.Context, filePath string) ([]model.Page, error)

	// AddText adiciona texto a uma página específica do PDF
	// Coordenadas em PDF points (72 DPI)
	AddText(ctx context.Context, filePath string, pageNum int, x, y float64, text string, fontSize float64) error

	// AddImage adiciona uma imagem a uma página específica do PDF
	// Coordenadas e dimensões em PDF points (72 DPI)
	AddImage(ctx context.Context, filePath string, pageNum int, x, y, width, height float64, imagePath string) error

	// MergePDFs mescla múltiplos PDFs em um único arquivo
	MergePDFs(ctx context.Context, outputPath string, inputPaths []string) error

	// GeneratePreview gera uma preview (imagem) de uma página específica do PDF
	GeneratePreview(ctx context.Context, filePath string, pageNum int) ([]byte, error)

	// ValidatePDF valida se um arquivo é um PDF válido usando magic bytes
	ValidatePDF(ctx context.Context, data []byte) error
}
