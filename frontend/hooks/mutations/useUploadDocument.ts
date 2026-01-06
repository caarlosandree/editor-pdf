// Hook para upload de documentos

import { useMutation, useQueryClient } from '@tanstack/react-query'
import { documentService } from '@/services/documentService'
import type { UploadDocumentResponse } from '@/types/document'
import { toast } from 'sonner'

export const useUploadDocument = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (file: File) => documentService.uploadDocument(file),
    onSuccess: (data: UploadDocumentResponse) => {
      // Invalida cache de documentos para refetch
      queryClient.invalidateQueries({ queryKey: ['documents'] })
      toast.success(data.message || 'Documento enviado com sucesso')
    },
    onError: (error: Error) => {
      toast.error('Erro ao enviar documento', {
        description: error.message || 'Tente novamente',
      })
    },
  })
}
