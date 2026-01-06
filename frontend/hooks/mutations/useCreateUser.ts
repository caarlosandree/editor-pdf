import { useMutation, useQueryClient } from '@tanstack/react-query'
import { createUser } from '@/services/userService'
import type { CreateUserData } from '@/types/user'

export const useCreateUser = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (data: CreateUserData) => createUser(data),
    onSuccess: () => {
      // Invalida cache para refetch automático
      queryClient.invalidateQueries({ queryKey: ['users'] })
    },
    onError: (error) => {
      console.error('Erro ao criar usuário:', error)
    },
  })
}
