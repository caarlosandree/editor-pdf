package pdf

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"

	"github.com/editor-pdf/backend/internal/domain"
	appModel "github.com/editor-pdf/backend/internal/model"
	"github.com/editor-pdf/backend/pkg/logger"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	pdfcpuModel "github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/unidoc/unipdf/v3/common"
	"github.com/unidoc/unipdf/v3/contentstream"
	"github.com/unidoc/unipdf/v3/core"
	"github.com/unidoc/unipdf/v3/model"
	"github.com/unidoc/unipdf/v3/render"
	"go.uber.org/zap"
)

// PDFCPUProcessor implementa PDFProcessor usando a biblioteca pdfcpu
type PDFCPUProcessor struct{}

// NewPDFCPUProcessor cria uma nova instância de PDFCPUProcessor
func NewPDFCPUProcessor() (domain.PDFProcessor, error) {
	return &PDFCPUProcessor{}, nil
}

// ValidatePDF valida se um arquivo é um PDF válido usando magic bytes
func (p *PDFCPUProcessor) ValidatePDF(ctx context.Context, data []byte) error {
	// Magic bytes de PDF: %PDF-1.
	if len(data) < 8 {
		return fmt.Errorf("arquivo muito pequeno para ser um PDF")
	}

	magicBytes := string(data[:8])
	if magicBytes[:4] != "%PDF" {
		return fmt.Errorf("arquivo não é um PDF válido (magic bytes inválidos)")
	}

	// Validação adicional tentando ler o PDF
	tempFile, err := os.CreateTemp("", "validate_*.pdf")
	if err != nil {
		return fmt.Errorf("erro ao criar arquivo temporário: %w", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	if _, err := tempFile.Write(data); err != nil {
		return fmt.Errorf("erro ao escrever arquivo temporário: %w", err)
	}

	// Tenta ler o PDF com pdfcpu (validação básica)
	if _, err := api.ReadContextFile(tempFile.Name()); err != nil {
		return fmt.Errorf("PDF inválido: %w", err)
	}

	return nil
}

// ExtractPages extrai informações sobre as páginas de um PDF
func (p *PDFCPUProcessor) ExtractPages(ctx context.Context, filePath string) ([]appModel.Page, error) {
	ctxFile, err := api.ReadContextFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler PDF: %w", err)
	}

	pageDims, err := ctxFile.PageDims()
	if err != nil {
		return nil, fmt.Errorf("erro ao obter dimensões das páginas: %w", err)
	}

	pages := make([]appModel.Page, 0, len(pageDims))
	for i, pageDim := range pageDims {
		pages = append(pages, appModel.Page{
			Number: i + 1,
			Width:  pageDim.Width,
			Height: pageDim.Height,
		})
	}

	return pages, nil
}

// AddText adiciona texto a uma página específica do PDF
func (p *PDFCPUProcessor) AddText(ctx context.Context, filePath string, pageNum int, x, y float64, text string, fontSize float64) error {
	// Desabilita logs do unipdf para evitar poluição
	common.SetLogger(common.NewConsoleLogger(common.LogLevelError))

	// Carrega o PDF
	reader, file, err := model.NewPdfReaderFromFile(filePath, nil)
	if err != nil {
		return fmt.Errorf("erro ao carregar PDF: %w", err)
	}
	defer file.Close()

	// Valida número da página
	numPages, err := reader.GetNumPages()
	if err != nil {
		return fmt.Errorf("erro ao obter número de páginas: %w", err)
	}

	if pageNum < 1 || pageNum > numPages {
		return fmt.Errorf("página inválida: %d (PDF tem %d páginas)", pageNum, numPages)
	}

	// Obtém a página (GetPage usa índice baseado em 1)
	page, err := reader.GetPage(pageNum)
	if err != nil {
		return fmt.Errorf("erro ao obter página %d: %w", pageNum, err)
	}

	// Obtém dimensões da página
	pageRect, err := page.GetMediaBox()
	if err != nil {
		return fmt.Errorf("erro ao obter dimensões da página: %w", err)
	}

	// Converte coordenadas: y=0 é no bottom no PDF, então precisamos inverter
	// Se y é fornecido do topo, converter: y_pdf = pageHeight - y
	pageHeight := pageRect.Height()
	yPdf := pageHeight - y

	// Obtém o content stream existente (retorna []string)
	contentStreams, err := page.GetContentStreams()
	if err != nil {
		return fmt.Errorf("erro ao obter content stream: %w", err)
	}

	// Cria um novo content stream com o texto
	contentCreator := contentstream.NewContentCreator()
	contentCreator.Add_q() // Save graphics state

	// Configura fonte (Helvetica por padrão)
	contentCreator.Add_BT()               // Begin text object
	contentCreator.Add_Tf("F1", fontSize) // Set font and size
	contentCreator.Add_Td(x, yPdf)        // Move to position

	// Adiciona o texto (precisa ser PdfObjectString)
	textObj := core.MakeString(text)
	contentCreator.Add_Tj(*textObj)
	contentCreator.Add_ET() // End text object
	contentCreator.Add_Q()  // Restore graphics state

	// Adiciona o novo conteúdo ao content stream existente
	newContent := string(contentCreator.Bytes())
	contentStreams = append(contentStreams, newContent)

	// Atualiza o content stream da página
	if err := page.SetContentStreams(contentStreams, core.NewFlateEncoder()); err != nil {
		return fmt.Errorf("erro ao atualizar content stream: %w", err)
	}

	// Adiciona fonte Helvetica se não existir
	resources := page.Resources
	if resources == nil {
		resources = model.NewPdfPageResources()
		page.Resources = resources
	}

	// Obtém ou cria o Font dictionary
	var fontDict *core.PdfObjectDictionary
	if resources.Font == nil {
		fontDict = core.MakeDict()
		resources.Font = fontDict
	} else {
		// Tenta fazer type assertion
		if dict, ok := resources.Font.(*core.PdfObjectDictionary); ok {
			fontDict = dict
		} else {
			// Se não for dict, cria um novo
			fontDict = core.MakeDict()
			resources.Font = fontDict
		}
	}

	// Cria fonte Helvetica se não existir
	if fontDict.Get("F1") == nil {
		helveticaFontDict := core.MakeDict()
		helveticaFontDict.Set("Type", core.MakeName("Font"))
		helveticaFontDict.Set("Subtype", core.MakeName("Type1"))
		helveticaFontDict.Set("BaseFont", core.MakeName("Helvetica"))

		fontObj := core.MakeIndirectObject(helveticaFontDict)
		fontDict.Set("F1", fontObj)
	}

	// Cria um novo writer e copia todas as páginas
	writer := model.NewPdfWriter()

	// Copia todas as páginas, modificando apenas a página especificada
	for i := 1; i <= numPages; i++ {
		var pageToAdd *model.PdfPage
		if i == pageNum {
			// Usa a página modificada
			pageToAdd = page
		} else {
			// Copia a página original
			otherPage, err := reader.GetPage(i)
			if err != nil {
				return fmt.Errorf("erro ao obter página %d: %w", i, err)
			}
			pageToAdd = otherPage
		}

		if err := writer.AddPage(pageToAdd); err != nil {
			return fmt.Errorf("erro ao adicionar página %d: %w", i, err)
		}
	}

	// Salva em arquivo temporário primeiro
	tempFile, err := os.CreateTemp("", "pdf_edit_*.pdf")
	if err != nil {
		return fmt.Errorf("erro ao criar arquivo temporário: %w", err)
	}
	tempPath := tempFile.Name()
	tempFile.Close()
	defer os.Remove(tempPath)

	if err := writer.WriteToFile(tempPath); err != nil {
		return fmt.Errorf("erro ao salvar PDF temporário: %w", err)
	}

	// Substitui o arquivo original pelo modificado
	if err := os.Rename(tempPath, filePath); err != nil {
		return fmt.Errorf("erro ao substituir arquivo: %w", err)
	}

	logger.Logger.Debug("Texto adicionado ao PDF",
		zap.String("file", filePath),
		zap.Int("page", pageNum),
		zap.String("text", text),
		zap.Float64("x", x),
		zap.Float64("y", y),
	)

	return nil
}

// AddImage adiciona uma imagem a uma página específica do PDF
func (p *PDFCPUProcessor) AddImage(ctx context.Context, filePath string, pageNum int, x, y, width, height float64, imagePath string) error {
	// Desabilita logs do unipdf para evitar poluição
	common.SetLogger(common.NewConsoleLogger(common.LogLevelError))

	// Valida se o arquivo de imagem existe
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		return fmt.Errorf("arquivo de imagem não encontrado: %s", imagePath)
	}

	// Carrega o PDF
	reader, file, err := model.NewPdfReaderFromFile(filePath, nil)
	if err != nil {
		return fmt.Errorf("erro ao carregar PDF: %w", err)
	}
	defer file.Close()

	// Valida número da página
	numPages, err := reader.GetNumPages()
	if err != nil {
		return fmt.Errorf("erro ao obter número de páginas: %w", err)
	}

	if pageNum < 1 || pageNum > numPages {
		return fmt.Errorf("página inválida: %d (PDF tem %d páginas)", pageNum, numPages)
	}

	// Obtém a página
	page, err := reader.GetPage(pageNum)
	if err != nil {
		return fmt.Errorf("erro ao obter página %d: %w", pageNum, err)
	}

	// Obtém dimensões da página
	pageRect, err := page.GetMediaBox()
	if err != nil {
		return fmt.Errorf("erro ao obter dimensões da página: %w", err)
	}

	// Converte coordenadas: y=0 é no bottom no PDF
	pageHeight := pageRect.Height()
	yPdf := pageHeight - y - height // Ajusta para que y seja do topo

	// Carrega a imagem usando image.Decode
	imgFile, err := os.Open(imagePath)
	if err != nil {
		return fmt.Errorf("erro ao abrir arquivo de imagem: %w", err)
	}
	defer imgFile.Close()

	goImg, _, err := image.Decode(imgFile)
	if err != nil {
		return fmt.Errorf("erro ao decodificar imagem: %w", err)
	}

	// Converte image.Image do Go para model.Image do unipdf
	bounds := goImg.Bounds()
	imgWidth := int64(bounds.Dx())
	imgHeight := int64(bounds.Dy())

	// Extrai os dados da imagem
	var imgData []byte
	switch img := goImg.(type) {
	case *image.RGBA:
		imgData = img.Pix
	case *image.NRGBA:
		imgData = img.Pix
	default:
		// Converte para RGBA se necessário
		rgba := image.NewRGBA(bounds)
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				rgba.Set(x, y, goImg.At(x, y))
			}
		}
		imgData = rgba.Pix
	}

	// Cria model.Image
	pdfImg := &model.Image{
		Width:            imgWidth,
		Height:           imgHeight,
		BitsPerComponent: 8,
		ColorComponents:  3, // RGB
		Data:             imgData,
	}

	// Cria um XObject de imagem
	// Usa DeviceRGB como colorspace e FlateEncoder
	colorspace := model.NewPdfColorspaceDeviceRGB()
	encoder := core.NewFlateEncoder()
	ximg, err := model.NewXObjectImageFromImage(pdfImg, colorspace, encoder)
	if err != nil {
		return fmt.Errorf("erro ao criar XObject de imagem: %w", err)
	}

	// Adiciona a imagem aos recursos da página
	resources := page.Resources
	if resources == nil {
		resources = model.NewPdfPageResources()
		page.Resources = resources
	}

	// Obtém ou cria o XObject dictionary
	var xObjectDict *core.PdfObjectDictionary
	if resources.XObject == nil {
		xObjectDict = core.MakeDict()
		resources.XObject = xObjectDict
	} else {
		if dict, ok := resources.XObject.(*core.PdfObjectDictionary); ok {
			xObjectDict = dict
		} else {
			xObjectDict = core.MakeDict()
			resources.XObject = xObjectDict
		}
	}

	// Gera um nome único para a imagem
	imageName := fmt.Sprintf("Img%d", len(xObjectDict.Keys())+1)
	imageNameObj := core.MakeName(imageName)
	xObjectDict.Set(*imageNameObj, ximg.ToPdfObject())

	// Obtém o content stream existente
	contentStreams, err := page.GetContentStreams()
	if err != nil {
		return fmt.Errorf("erro ao obter content stream: %w", err)
	}

	// Cria um novo content stream com a imagem
	contentCreator := contentstream.NewContentCreator()
	contentCreator.Add_q() // Save graphics state

	// Posiciona e dimensiona a imagem
	// Matrix: [width 0 0 height x y] cm
	contentCreator.Add_cm(width, 0, 0, height, x, yPdf) // cm = concat matrix
	contentCreator.Add_Do(*imageNameObj)                // Do = draw object
	contentCreator.Add_Q()                              // Restore graphics state

	// Adiciona o novo conteúdo ao content stream existente
	newContent := string(contentCreator.Bytes())
	contentStreams = append(contentStreams, newContent)

	// Atualiza o content stream da página
	if err := page.SetContentStreams(contentStreams, core.NewFlateEncoder()); err != nil {
		return fmt.Errorf("erro ao atualizar content stream: %w", err)
	}

	// Cria um novo writer e copia todas as páginas
	writer := model.NewPdfWriter()

	// Copia todas as páginas, modificando apenas a página especificada
	for i := 1; i <= numPages; i++ {
		var pageToAdd *model.PdfPage
		if i == pageNum {
			pageToAdd = page
		} else {
			otherPage, err := reader.GetPage(i)
			if err != nil {
				return fmt.Errorf("erro ao obter página %d: %w", i, err)
			}
			pageToAdd = otherPage
		}

		if err := writer.AddPage(pageToAdd); err != nil {
			return fmt.Errorf("erro ao adicionar página %d: %w", i, err)
		}
	}

	// Salva em arquivo temporário primeiro
	tempFile, err := os.CreateTemp("", "pdf_edit_*.pdf")
	if err != nil {
		return fmt.Errorf("erro ao criar arquivo temporário: %w", err)
	}
	tempPath := tempFile.Name()
	tempFile.Close()
	defer os.Remove(tempPath)

	if err := writer.WriteToFile(tempPath); err != nil {
		return fmt.Errorf("erro ao salvar PDF temporário: %w", err)
	}

	// Substitui o arquivo original pelo modificado
	if err := os.Rename(tempPath, filePath); err != nil {
		return fmt.Errorf("erro ao substituir arquivo: %w", err)
	}

	logger.Logger.Debug("Imagem adicionada ao PDF",
		zap.String("file", filePath),
		zap.Int("page", pageNum),
		zap.String("image", imagePath),
		zap.Float64("x", x),
		zap.Float64("y", y),
		zap.Float64("width", width),
		zap.Float64("height", height),
	)

	return nil
}

// MergePDFs mescla múltiplos PDFs em um único arquivo
func (p *PDFCPUProcessor) MergePDFs(ctx context.Context, outputPath string, inputPaths []string) error {
	// Cria diretório de saída se não existir
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório: %w", err)
	}

	// Usa pdfcpu para mesclar PDFs
	config := pdfcpuModel.NewDefaultConfiguration()
	if err := api.MergeCreateFile(inputPaths, outputPath, false, config); err != nil {
		return fmt.Errorf("erro ao mesclar PDFs: %w", err)
	}

	logger.Logger.Debug("PDFs mesclados com sucesso",
		zap.String("output", outputPath),
		zap.Int("count", len(inputPaths)),
	)

	return nil
}

// GeneratePreview gera uma preview (imagem) de uma página específica do PDF
func (p *PDFCPUProcessor) GeneratePreview(ctx context.Context, filePath string, pageNum int) ([]byte, error) {
	// Desabilita logs do unipdf para evitar poluição
	common.SetLogger(common.NewConsoleLogger(common.LogLevelError))

	// Carrega o PDF
	reader, file, err := model.NewPdfReaderFromFile(filePath, nil)
	if err != nil {
		return nil, fmt.Errorf("erro ao carregar PDF: %w", err)
	}
	defer file.Close()

	// Valida número da página
	numPages, err := reader.GetNumPages()
	if err != nil {
		return nil, fmt.Errorf("erro ao obter número de páginas: %w", err)
	}

	if pageNum < 1 || pageNum > numPages {
		return nil, fmt.Errorf("página inválida: %d (PDF tem %d páginas)", pageNum, numPages)
	}

	// Obtém a página
	page, err := reader.GetPage(pageNum)
	if err != nil {
		return nil, fmt.Errorf("erro ao obter página %d: %w", pageNum, err)
	}

	// Obtém dimensões da página para calcular largura de saída
	pageRect, err := page.GetMediaBox()
	if err != nil {
		return nil, fmt.Errorf("erro ao obter dimensões da página: %w", err)
	}

	// Calcula largura de saída para 150 DPI (72 points = 1 inch)
	// 150 DPI = 150 pixels por inch = 150/72 pixels por point
	dpi := 150.0
	pointsPerInch := 72.0
	pixelsPerPoint := dpi / pointsPerInch
	outputWidth := int(pageRect.Width() * pixelsPerPoint)

	// Cria um renderer
	device := render.NewImageDevice()
	device.OutputWidth = outputWidth

	// Renderiza a página
	img, err := device.Render(page)
	if err != nil {
		return nil, fmt.Errorf("erro ao renderizar página: %w", err)
	}

	// Converte para bytes PNG
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, fmt.Errorf("erro ao codificar imagem PNG: %w", err)
	}

	logger.Logger.Debug("Preview gerado com sucesso",
		zap.String("file", filePath),
		zap.Int("page", pageNum),
		zap.Int("size_bytes", buf.Len()),
	)

	return buf.Bytes(), nil
}

// ProcessEdits processa múltiplas edições em um PDF
// Esta é uma função auxiliar que pode ser usada pelo UseCase
func (p *PDFCPUProcessor) ProcessEdits(ctx context.Context, inputPath, outputPath string, edits []EditInstruction) error {
	// Cria diretório de saída se não existir
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório: %w", err)
	}

	// Copia o arquivo de entrada para um arquivo temporário de trabalho
	tempFile, err := os.CreateTemp("", "pdf_edit_*.pdf")
	if err != nil {
		return fmt.Errorf("erro ao criar arquivo temporário: %w", err)
	}
	tempPath := tempFile.Name()
	tempFile.Close()
	defer os.Remove(tempPath)

	inputData, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("erro ao ler PDF original: %w", err)
	}

	if err := os.WriteFile(tempPath, inputData, 0644); err != nil {
		return fmt.Errorf("erro ao copiar PDF para arquivo temporário: %w", err)
	}

	// Processa cada edição sequencialmente
	for i, edit := range edits {
		switch edit.Type {
		case "text":
			// Valida campos obrigatórios
			if edit.Content == "" {
				return fmt.Errorf("edição %d: conteúdo de texto não pode ser vazio", i+1)
			}
			fontSize := edit.FontSize
			if fontSize <= 0 {
				fontSize = 12.0 // Tamanho padrão
			}

			// Aplica a edição de texto
			if err := p.AddText(ctx, tempPath, edit.Page, edit.X, edit.Y, edit.Content, fontSize); err != nil {
				return fmt.Errorf("erro ao adicionar texto na edição %d: %w", i+1, err)
			}

		case "image":
			// Valida campos obrigatórios
			if edit.Content == "" {
				return fmt.Errorf("edição %d: caminho da imagem não pode ser vazio", i+1)
			}
			if edit.Width <= 0 {
				return fmt.Errorf("edição %d: largura da imagem deve ser maior que zero", i+1)
			}
			if edit.Height <= 0 {
				return fmt.Errorf("edição %d: altura da imagem deve ser maior que zero", i+1)
			}

			// Aplica a edição de imagem
			if err := p.AddImage(ctx, tempPath, edit.Page, edit.X, edit.Y, edit.Width, edit.Height, edit.Content); err != nil {
				return fmt.Errorf("erro ao adicionar imagem na edição %d: %w", i+1, err)
			}

		case "drawing":
			return fmt.Errorf("edição %d: tipo 'drawing' ainda não está implementado", i+1)

		default:
			return fmt.Errorf("edição %d: tipo de edição desconhecido: %s", i+1, edit.Type)
		}
	}

	// Move o arquivo temporário processado para o destino final
	if err := os.Rename(tempPath, outputPath); err != nil {
		// Se rename falhar (diferentes filesystems), copia o arquivo
		outputData, readErr := os.ReadFile(tempPath)
		if readErr != nil {
			return fmt.Errorf("erro ao ler arquivo processado: %w (também falhou ao renomear: %v)", readErr, err)
		}
		if writeErr := os.WriteFile(outputPath, outputData, 0644); writeErr != nil {
			return fmt.Errorf("erro ao salvar PDF processado: %w (também falhou ao renomear: %v)", writeErr, err)
		}
	}

	logger.Logger.Debug("Edições processadas com sucesso",
		zap.String("input", inputPath),
		zap.String("output", outputPath),
		zap.Int("edits_count", len(edits)),
	)

	return nil
}

// EditInstruction representa uma instrução de edição
type EditInstruction struct {
	Type     string                 `json:"type"` // "text", "image", "drawing"
	Page     int                    `json:"page"`
	X        float64                `json:"x"`
	Y        float64                `json:"y"`
	Width    float64                `json:"width,omitempty"`
	Height   float64                `json:"height,omitempty"`
	Content  string                 `json:"content,omitempty"`
	FontSize float64                `json:"fontSize,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// Helper function para converter imagem para bytes
func imageToBytes(img image.Image, format string) ([]byte, error) {
	var buf bytes.Buffer

	switch format {
	case "jpeg", "jpg":
		if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 90}); err != nil {
			return nil, err
		}
	case "png":
		if err := png.Encode(&buf, img); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("formato de imagem não suportado: %s", format)
	}

	return buf.Bytes(), nil
}
