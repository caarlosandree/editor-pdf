package response

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// SuccessResponse representa uma resposta de sucesso
type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

// ErrorResponse representa uma resposta de erro
type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// Success retorna uma resposta de sucesso
func Success(c echo.Context, statusCode int, data interface{}, message string) error {
	return c.JSON(statusCode, SuccessResponse{
		Success: true,
		Data:    data,
		Message: message,
	})
}

// SuccessOK retorna uma resposta de sucesso com status 200
func SuccessOK(c echo.Context, data interface{}) error {
	return Success(c, http.StatusOK, data, "")
}

// SuccessCreated retorna uma resposta de sucesso com status 201
func SuccessCreated(c echo.Context, data interface{}, message string) error {
	return Success(c, http.StatusCreated, data, message)
}

// Error retorna uma resposta de erro
func Error(c echo.Context, statusCode int, err error, message string) error {
	errorMsg := ""
	if err != nil {
		errorMsg = err.Error()
	}
	return c.JSON(statusCode, ErrorResponse{
		Success: false,
		Error:   errorMsg,
		Message: message,
	})
}

// ErrorBadRequest retorna uma resposta de erro com status 400
func ErrorBadRequest(c echo.Context, err error, message string) error {
	return Error(c, http.StatusBadRequest, err, message)
}

// ErrorUnauthorized retorna uma resposta de erro com status 401
func ErrorUnauthorized(c echo.Context, err error, message string) error {
	return Error(c, http.StatusUnauthorized, err, message)
}

// ErrorForbidden retorna uma resposta de erro com status 403
func ErrorForbidden(c echo.Context, err error, message string) error {
	return Error(c, http.StatusForbidden, err, message)
}

// ErrorNotFound retorna uma resposta de erro com status 404
func ErrorNotFound(c echo.Context, err error, message string) error {
	return Error(c, http.StatusNotFound, err, message)
}

// ErrorInternalServer retorna uma resposta de erro com status 500
func ErrorInternalServer(c echo.Context, err error, message string) error {
	return Error(c, http.StatusInternalServerError, err, message)
}
