package error

import (
	"errors"
	"fmt"
)

// Erros customizados da aplicação
var (
	ErrNotFound      = errors.New("recurso não encontrado")
	ErrInvalidInput  = errors.New("entrada inválida")
	ErrUnauthorized  = errors.New("não autorizado")
	ErrForbidden     = errors.New("acesso negado")
	ErrInternalError = errors.New("erro interno do servidor")
)

// AppError representa um erro da aplicação com contexto adicional
type AppError struct {
	Err     error
	Message string
	Code    string
}

// Error implementa a interface error
func (e *AppError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	if e.Err != nil {
		return e.Err.Error()
	}
	return "erro desconhecido"
}

// Unwrap retorna o erro original
func (e *AppError) Unwrap() error {
	return e.Err
}

// NewAppError cria um novo AppError
func NewAppError(err error, message string, code string) *AppError {
	return &AppError{
		Err:     err,
		Message: message,
		Code:    code,
	}
}

// WrapError envolve um erro com contexto adicional
func WrapError(err error, message string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", message, err)
}
