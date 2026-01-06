// Serviço para operações de autenticação

import api from '@/lib/axios'
import type { User } from '@/types/user'
import type { ApiResponse } from '@/types/api'

export interface RegisterRequest {
  name: string
  email: string
  password: string
}

export interface LoginRequest {
  email: string
  password: string
}

export interface AuthResponse {
  token: string
  user: User
}

// Helper para extrair data da resposta da API
function extractData<T>(response: ApiResponse<T>): T {
  if (response.success && response.data) {
    return response.data
  }
  throw new Error(response.message || 'Erro ao processar resposta da API')
}

export const authService = {
  // Registra um novo usuário
  async register(data: RegisterRequest): Promise<AuthResponse> {
    const response = await api.post<ApiResponse<AuthResponse>>('/auth/register', data)
    return extractData(response.data)
  },

  // Faz login
  async login(data: LoginRequest): Promise<AuthResponse> {
    const response = await api.post<ApiResponse<AuthResponse>>('/auth/login', data)
    return extractData(response.data)
  },

  // Busca dados do usuário autenticado
  async getMe(): Promise<User> {
    const response = await api.get<ApiResponse<User>>('/auth/me')
    return extractData(response.data)
  },
}
