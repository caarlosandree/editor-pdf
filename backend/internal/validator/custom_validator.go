package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// CustomValidator é um validator customizado para Echo
type CustomValidator struct {
	validator *validator.Validate
}

// NewCustomValidator cria uma nova instância de CustomValidator
func NewCustomValidator() *CustomValidator {
	return &CustomValidator{
		validator: validator.New(),
	}
}

// Validate valida uma struct
func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return err
	}
	return nil
}

// RegisterCustomValidator registra o validator customizado no Echo
func RegisterCustomValidator(e *echo.Echo) {
	e.Validator = NewCustomValidator()
}
