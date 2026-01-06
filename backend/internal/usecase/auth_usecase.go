package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/editor-pdf/backend/internal/config"
	"github.com/editor-pdf/backend/internal/domain"
	"github.com/editor-pdf/backend/internal/dto"
	"github.com/editor-pdf/backend/internal/model"
	"github.com/editor-pdf/backend/pkg/logger"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// AuthUseCase contém os casos de uso de autenticação
type AuthUseCase struct {
	userRepo domain.UserRepository
	config   *config.Config
}

// NewAuthUseCase cria uma nova instância de AuthUseCase
func NewAuthUseCase(userRepo domain.UserRepository, cfg *config.Config) *AuthUseCase {
	return &AuthUseCase{
		userRepo: userRepo,
		config:   cfg,
	}
}

// Register registra um novo usuário
func (uc *AuthUseCase) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.AuthResponse, error) {
	// Verifica se o usuário já existe
	existingUser, err := uc.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		logger.Logger.Error("Erro ao buscar usuário por email", zap.Error(err))
		return nil, fmt.Errorf("erro ao verificar email: %w", err)
	}

	if existingUser != nil {
		return nil, errors.New("email já está em uso")
	}

	// Hash da senha
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Logger.Error("Erro ao gerar hash da senha", zap.Error(err))
		return nil, fmt.Errorf("erro ao processar senha: %w", err)
	}

	// Cria o usuário
	user := &model.User{
		ID:           uuid.New(),
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Name:         req.Name,
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		logger.Logger.Error("Erro ao criar usuário", zap.Error(err))
		return nil, fmt.Errorf("erro ao criar usuário: %w", err)
	}

	// Gera token JWT
	token, err := uc.generateToken(user.ID.String(), user.Email)
	if err != nil {
		return nil, fmt.Errorf("erro ao gerar token: %w", err)
	}

	return &dto.AuthResponse{
		Token: token,
		User: dto.UserResponse{
			ID:    user.ID.String(),
			Email: user.Email,
			Name:  user.Name,
		},
	}, nil
}

// Login autentica um usuário
func (uc *AuthUseCase) Login(ctx context.Context, req *dto.LoginRequest) (*dto.AuthResponse, error) {
	// Busca o usuário por email
	user, err := uc.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		logger.Logger.Error("Erro ao buscar usuário por email", zap.Error(err))
		return nil, fmt.Errorf("credenciais inválidas")
	}

	if user == nil {
		return nil, errors.New("credenciais inválidas")
	}

	// Verifica a senha
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("credenciais inválidas")
	}

	// Gera token JWT
	token, err := uc.generateToken(user.ID.String(), user.Email)
	if err != nil {
		return nil, fmt.Errorf("erro ao gerar token: %w", err)
	}

	return &dto.AuthResponse{
		Token: token,
		User: dto.UserResponse{
			ID:    user.ID.String(),
			Email: user.Email,
			Name:  user.Name,
		},
	}, nil
}

// ValidateToken valida um token JWT e retorna os dados do usuário
func (uc *AuthUseCase) ValidateToken(ctx context.Context, tokenString string) (*dto.UserResponse, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verifica o método de assinatura
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de assinatura inesperado: %v", token.Header["alg"])
		}
		return []byte(uc.config.JWT.Secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("token inválido: %w", err)
	}

	if !token.Valid {
		return nil, errors.New("token inválido")
	}

	// Extrai claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("erro ao extrair claims do token")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return nil, errors.New("user_id não encontrado no token")
	}

	// Busca o usuário no banco para garantir que ainda existe
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("user_id inválido: %w", err)
	}

	user, err := uc.userRepo.FindByID(ctx, userUUID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar usuário: %w", err)
	}

	if user == nil {
		return nil, errors.New("usuário não encontrado")
	}

	return &dto.UserResponse{
		ID:    user.ID.String(),
		Email: user.Email,
		Name:  user.Name,
	}, nil
}

// generateToken gera um token JWT
func (uc *AuthUseCase) generateToken(userID, email string) (string, error) {
	// Parse da duração de expiração
	expiration, err := time.ParseDuration(uc.config.JWT.Expiration)
	if err != nil {
		// Default para 24h se não conseguir parsear
		expiration = 24 * time.Hour
	}

	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     time.Now().Add(expiration).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(uc.config.JWT.Secret))
	if err != nil {
		return "", fmt.Errorf("erro ao assinar token: %w", err)
	}

	return tokenString, nil
}
