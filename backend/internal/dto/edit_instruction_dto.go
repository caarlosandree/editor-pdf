package dto

// EditInstruction representa uma instrução de edição de PDF
// @Description Instrução individual para editar um documento PDF (adicionar texto, imagem ou desenho)
type EditInstruction struct {
	Type     string                 `json:"type" validate:"required,oneof=text image drawing" example:"text" enums:"text,image,drawing"`
	Page     int                    `json:"page" validate:"required,min=1" example:"1"`
	X        float64                `json:"x" validate:"required" example:"100.5"`
	Y        float64                `json:"y" validate:"required" example:"200.5"`
	Width    *float64               `json:"width,omitempty" example:"150.0"`
	Height   *float64               `json:"height,omitempty" example:"50.0"`
	Content  string                 `json:"content,omitempty" example:"Texto a ser inserido"`
	FontSize *float64               `json:"fontSize,omitempty" example:"12.0"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// ProcessDocumentRequest representa a requisição para processar edições em um documento
// @Description Requisição contendo lista de instruções de edição a serem aplicadas no documento
type ProcessDocumentRequest struct {
	Instructions []EditInstruction `json:"instructions" validate:"required,min=1,dive"`
}

// ProcessDocumentResponse representa a resposta após processar edições
// @Description Resposta após processamento bem-sucedido das edições no documento
type ProcessDocumentResponse struct {
	Document DocumentResponse `json:"document"`
	Message  string           `json:"message" example:"Documento processado com sucesso"`
}
