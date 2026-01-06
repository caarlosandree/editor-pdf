import axios from 'axios'

const api = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1',
  headers: {
    'Content-Type': 'application/json',
  },
})

// Interceptor para adicionar headers (se necessário no futuro)
api.interceptors.request.use(
  (config) => {
    // Não adiciona token JWT (autenticação não é necessária)
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Interceptor para tratamento de erros
api.interceptors.response.use(
  (response) => response,
  (error) => {
    // Apenas rejeita o erro sem redirecionar
    // (autenticação não é necessária neste projeto)
    return Promise.reject(error)
  }
)

export default api
