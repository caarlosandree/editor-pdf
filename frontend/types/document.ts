// Tipos relacionados a documentos PDF

export interface Document {
  id: string
  user_id: string
  file_path: string
  file_url: string
  checksum: string
  version: number
  status: 'uploaded' | 'processing' | 'processed' | 'error'
  page_count: number
  created_at: string
  updated_at: string
}

export interface DocumentListResponse {
  documents: Document[]
  total: number
  limit: number
  offset: number
}

export interface UploadDocumentResponse {
  document: Document
  message: string
}

export interface ProcessDocumentResponse {
  document: Document
  message: string
}
