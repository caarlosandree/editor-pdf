package dto

import "time"

// DocumentResponse representa a resposta de um documento
// @Description Informações completas de um documento PDF
type DocumentResponse struct {
	ID        string    `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	UserID    string    `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	FilePath  string    `json:"file_path" example:"/storage/documents/550e8400-e29b-41d4-a716-446655440000.pdf"`
	FileURL   string    `json:"file_url" example:"/api/v1/documents/550e8400-e29b-41d4-a716-446655440000/file"`
	Checksum  string    `json:"checksum" example:"a1b2c3d4e5f6..."`
	Version   int       `json:"version" example:"1"`
	Status    string    `json:"status" example:"processed" enums:"uploaded,processing,processed,error"`
	PageCount int       `json:"page_count" example:"10"`
	CreatedAt time.Time `json:"created_at" example:"2024-01-15T10:30:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2024-01-15T10:30:00Z"`
}

// DocumentListResponse representa a resposta de uma lista de documentos
// @Description Lista paginada de documentos do usuário
type DocumentListResponse struct {
	Documents []DocumentResponse `json:"documents"`
	Total     int                `json:"total" example:"50"`
	Limit     int                `json:"limit" example:"20"`
	Offset    int                `json:"offset" example:"0"`
}

// UploadDocumentResponse representa a resposta após upload de documento
// @Description Resposta após upload bem-sucedido de um documento PDF
type UploadDocumentResponse struct {
	Document DocumentResponse `json:"document"`
	Message  string           `json:"message" example:"Documento enviado com sucesso"`
}
