// Hook para processar edições em documentos

import { useMutation, useQueryClient } from '@tanstack/react-query'
import { documentService } from '@/services/documentService'
import type { ProcessDocumentRequest } from '@/types/edit'
import type { ProcessDocumentResponse } from '@/types/document'
import { toast } from 'sonner'

export const useProcessDocument = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({
      id,
      request,
    }: {
      id: string
      request: ProcessDocumentRequest
    }) => documentService.processDocument(id, request),
    onSuccess: (data: ProcessDocumentResponse, variables) => {
      // Invalida cache do documento específico e lista
      queryClient.invalidateQueries({ queryKey: ['documents', variables.id] })
      queryClient.invalidateQueries({ queryKey: ['documents'] })
      toast.success(data.message || 'Documento processado com sucesso')
    },
    onError: (error: Error) => {
      toast.error('Erro ao processar documento', {
        description: error.message || 'Tente novamente',
      })
    },
  })
}
