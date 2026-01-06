import { useQuery } from '@tanstack/react-query'
import { getUsers, getUserById } from '@/services/userService'
import type { User } from '@/types/user'

export const useUsers = () => {
  return useQuery<User[]>({
    queryKey: ['users'],
    queryFn: getUsers,
    staleTime: 5 * 60 * 1000, // 5 minutos
  })
}

export const useUser = (userId: number) => {
  return useQuery<User>({
    queryKey: ['users', userId],
    queryFn: () => getUserById(userId),
    enabled: !!userId, // Só executa se userId existir
  })
}
