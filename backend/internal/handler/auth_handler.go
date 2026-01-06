package handler

import (
	"net/http"

	"github.com/editor-pdf/backend/internal/dto"
	"github.com/editor-pdf/backend/internal/middleware"
	"github.com/editor-pdf/backend/internal/usecase"
	"github.com/editor-pdf/backend/pkg/response"
	"github.com/labstack/echo/v4"
)

// AuthHandler contém os handlers de autenticação
type AuthHandler struct {
	authUseCase *usecase.AuthUseCase
}

// NewAuthHandler cria uma nova instância de AuthHandler
func NewAuthHandler(authUseCase *usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{
		authUseCase: authUseCase,
	}
}

// Register registra um novo usuário
// @Summary Registra um novo usuário
// @Description Cria uma nova conta de usuário
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Dados de registro"
// @Success 201 {object} dto.AuthResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 409 {object} response.ErrorResponse
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c echo.Context) error {
	var req dto.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return response.ErrorBadRequest(c, err, "dados inválidos")
	}

	if err := c.Validate(&req); err != nil {
		return response.ErrorBadRequest(c, err, "validação falhou")
	}

	authResponse, err := h.authUseCase.Register(c.Request().Context(), &req)
	if err != nil {
		if err.Error() == "email já está em uso" {
			return response.Error(c, http.StatusConflict, err, "email já está em uso")
		}
		return response.ErrorInternalServer(c, err, "erro ao registrar usuário")
	}

	return response.SuccessCreated(c, authResponse, "Usuário registrado com sucesso")
}

// Login autentica um usuário
// @Summary Autentica um usuário
// @Description Realiza login e retorna token JWT
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Credenciais de login"
// @Success 200 {object} dto.AuthResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c echo.Context) error {
	var req dto.LoginRequest
	if err := c.Bind(&req); err != nil {
		return response.ErrorBadRequest(c, err, "dados inválidos")
	}

	if err := c.Validate(&req); err != nil {
		return response.ErrorBadRequest(c, err, "validação falhou")
	}

	authResponse, err := h.authUseCase.Login(c.Request().Context(), &req)
	if err != nil {
		if err.Error() == "credenciais inválidas" {
			return response.ErrorUnauthorized(c, err, "credenciais inválidas")
		}
		return response.ErrorInternalServer(c, err, "erro ao fazer login")
	}

	return response.SuccessOK(c, authResponse)
}

// Me retorna os dados do usuário autenticado
// @Summary Retorna dados do usuário autenticado
// @Description Retorna os dados do usuário baseado no token JWT
// @Tags auth
// @Security Bearer
// @Produce json
// @Success 200 {object} dto.UserResponse
// @Failure 401 {object} response.ErrorResponse
// @Router /api/v1/auth/me [get]
func (h *AuthHandler) Me(c echo.Context) error {
	user := middleware.GetUser(c)
	if user == nil {
		return response.ErrorUnauthorized(c, nil, "usuário não autenticado")
	}

	return response.SuccessOK(c, user)
}
