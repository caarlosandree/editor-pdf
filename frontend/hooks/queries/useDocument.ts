// Hook para buscar um documento específico

import { useQuery } from '@tanstack/react-query'
import { documentService } from '@/services/documentService'
import type { Document } from '@/types/document'

export const useDocument = (id: string | undefined) => {
  return useQuery<Document>({
    queryKey: ['documents', id],
    queryFn: () => documentService.getDocument(id!),
    enabled: !!id,
    staleTime: 60 * 1000, // 1 minuto
  })
}
