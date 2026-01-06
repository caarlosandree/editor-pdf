// Hook para deletar documentos

import { useMutation, useQueryClient } from '@tanstack/react-query'
import { documentService } from '@/services/documentService'
import { toast } from 'sonner'

export const useDeleteDocument = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (id: string) => documentService.deleteDocument(id),
    onSuccess: () => {
      // Invalida cache de documentos para refetch
      queryClient.invalidateQueries({ queryKey: ['documents'] })
      toast.success('Documento deletado com sucesso')
    },
    onError: (error: Error) => {
      toast.error('Erro ao deletar documento', {
        description: error.message || 'Tente novamente',
      })
    },
  })
}
