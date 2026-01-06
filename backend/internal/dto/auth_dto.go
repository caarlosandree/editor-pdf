package dto

// RegisterRequest representa a requisição de registro
// @Description Dados necessários para registro de um novo usuário
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email" example:"usuario@example.com"`
	Password string `json:"password" validate:"required,min=8" example:"senhaSegura123"`
	Name     string `json:"name" validate:"required,min=3,max=255" example:"João Silva"`
}

// LoginRequest representa a requisição de login
// @Description Credenciais de autenticação do usuário
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email" example:"usuario@example.com"`
	Password string `json:"password" validate:"required" example:"senhaSegura123"`
}

// AuthResponse representa a resposta de autenticação
// @Description Resposta contendo token JWT e dados do usuário autenticado
type AuthResponse struct {
	Token string      `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	User  UserResponse `json:"user"`
}

// UserResponse representa os dados do usuário na resposta
// @Description Dados do usuário retornados nas respostas da API
type UserResponse struct {
	ID    string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Email string `json:"email" example:"usuario@example.com"`
	Name  string `json:"name" example:"João Silva"`
}
