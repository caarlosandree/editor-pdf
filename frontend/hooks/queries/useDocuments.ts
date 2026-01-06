// Hook para listar documentos

import { useQuery } from '@tanstack/react-query'
import { documentService, type ListDocumentsParams } from '@/services/documentService'
import type { DocumentListResponse } from '@/types/document'

export const useDocuments = (params?: ListDocumentsParams) => {
  return useQuery<DocumentListResponse>({
    queryKey: ['documents', params],
    queryFn: () => documentService.listDocuments(params),
    staleTime: 30 * 1000, // 30 segundos
  })
}
