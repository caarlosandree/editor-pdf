import api from '@/lib/axios'
import type { User, CreateUserData, UpdateUserData } from '@/types/user'

export const getUsers = async (): Promise<User[]> => {
  const response = await api.get<User[]>('/users')
  return response.data
}

export const getUserById = async (id: number): Promise<User> => {
  const response = await api.get<User>(`/users/${id}`)
  return response.data
}

export const createUser = async (data: CreateUserData): Promise<User> => {
  const response = await api.post<User>('/users', data)
  return response.data
}

export const updateUser = async (
  id: number,
  data: UpdateUserData
): Promise<User> => {
  const response = await api.put<User>(`/users/${id}`, data)
  return response.data
}

export const deleteUser = async (id: number): Promise<void> => {
  await api.delete(`/users/${id}`)
}
