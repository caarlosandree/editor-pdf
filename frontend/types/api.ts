// Tipos base para respostas da API

export interface ApiResponse<T> {
  success: boolean
  data?: T
  message?: string
}

export interface ApiError {
  success: boolean
  error: string
  message?: string
}

export interface PaginatedResponse<T> {
  data: T[]
  page: number
  limit: number
  total: number
  totalPages: number
}
