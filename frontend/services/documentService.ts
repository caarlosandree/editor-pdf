// Serviço para operações com documentos

import api from '@/lib/axios'
import type {
  Document,
  DocumentListResponse,
  UploadDocumentResponse,
  ProcessDocumentResponse,
} from '@/types/document'
import type { ProcessDocumentRequest } from '@/types/edit'
import type { ApiResponse } from '@/types/api'

export interface ListDocumentsParams {
  limit?: number
  offset?: number
}

// Helper para extrair data da resposta da API
function extractData<T>(response: ApiResponse<T>): T {
  if (response.success && response.data) {
    return response.data
  }
  throw new Error(response.message || 'Erro ao processar resposta da API')
}

export const documentService = {
  // Lista documentos do usuário
  async listDocuments(
    params?: ListDocumentsParams
  ): Promise<DocumentListResponse> {
    const response = await api.get<ApiResponse<DocumentListResponse>>('/documents', {
      params,
    })
    return extractData(response.data)
  },

  // Busca um documento por ID
  async getDocument(id: string): Promise<Document> {
    const response = await api.get<ApiResponse<Document>>(`/documents/${id}`)
    return extractData(response.data)
  },

  // Faz upload de um documento PDF
  async uploadDocument(file: File): Promise<UploadDocumentResponse> {
    const formData = new FormData()
    formData.append('file', file)

    const response = await api.post<ApiResponse<UploadDocumentResponse>>(
      '/documents',
      formData,
      {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      }
    )
    return extractData(response.data)
  },

  // Processa edições em um documento
  async processDocument(
    id: string,
    request: ProcessDocumentRequest
  ): Promise<ProcessDocumentResponse> {
    const response = await api.post<ApiResponse<ProcessDocumentResponse>>(
      `/documents/${id}/process`,
      request
    )
    return extractData(response.data)
  },

  // Gera preview de uma página do documento (retorna URL da imagem)
  getPreviewUrl(id: string, page: number): string {
    const baseURL = api.defaults.baseURL || ''
    return `${baseURL}/documents/${id}/preview/${page}`
  },

  // Deleta um documento
  async deleteDocument(id: string): Promise<void> {
    await api.delete(`/documents/${id}`)
  },
}
