// Tipos relacionados a usuários

export interface User {
  id: string // UUID
  name: string
  email: string
  createdAt?: string
  updatedAt?: string
}

export interface CreateUserData {
  name: string
  email: string
  password: string
}

export interface UpdateUserData {
  name?: string
  email?: string
}
