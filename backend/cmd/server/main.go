// @title Editor PDF API
// @version 1.0
// @description API para edição e processamento de documentos PDF
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@editorpdf.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1


package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/editor-pdf/backend/cmd/server/docs" // Importa docs para registrar Swagger
	"github.com/editor-pdf/backend/internal/config"
	"github.com/editor-pdf/backend/internal/handler"
	"github.com/editor-pdf/backend/internal/infrastructure/pdf"
	"github.com/editor-pdf/backend/internal/infrastructure/storage"
	appMiddleware "github.com/editor-pdf/backend/internal/middleware"
	"github.com/editor-pdf/backend/internal/repository"
	"github.com/editor-pdf/backend/internal/usecase"
	"github.com/editor-pdf/backend/internal/util"
	"github.com/editor-pdf/backend/internal/validator"
	"github.com/editor-pdf/backend/pkg/logger"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.uber.org/zap"
)

func main() {
	// Carrega configurações
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Erro ao carregar configurações: %v\n", err)
		os.Exit(1)
	}

	// Inicializa logger
	if err := logger.InitLogger(cfg.Env); err != nil {
		fmt.Printf("Erro ao inicializar logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Logger.Info("Iniciando servidor",
		zap.String("env", cfg.Env),
		zap.String("address", cfg.GetServerAddress()),
	)

	// Conecta ao banco de dados
	db, err := util.NewDB(cfg.DB.DSN())
	if err != nil {
		logger.Logger.Fatal("Erro ao conectar ao banco de dados", zap.Error(err))
	}
	defer func() {
		if err := util.CloseDB(db); err != nil {
			logger.Logger.Error("Erro ao fechar conexão com banco de dados", zap.Error(err))
		}
	}()

	logger.Logger.Info("Conexão com banco de dados estabelecida")

	// Cria instância do Echo
	e := echo.New()

	// Registra validator customizado
	validator.RegisterCustomValidator(e)

	// Middlewares globais
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID()) // Deve vir antes do LoggingMiddleware para garantir RequestID
	e.Use(appMiddleware.LoggingMiddleware()) // Middleware customizado de logging
	e.Use(appMiddleware.SecurityHeaders())
	e.Use(appMiddleware.SetupCORS(cfg.CORS.AllowedOrigins))

	// Health check endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status":  "ok",
			"service": "editor-pdf-backend",
		})
	})

	// Swagger documentation
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Setup de rotas da API
	setupRoutes(e, db, cfg)

	// Inicia servidor em goroutine
	go func() {
		address := cfg.GetServerAddress()
		if err := e.Start(address); err != nil && err != http.ErrServerClosed {
			logger.Logger.Fatal("Erro ao iniciar servidor", zap.Error(err))
		}
	}()

	logger.Logger.Info("Servidor iniciado", zap.String("address", cfg.GetServerAddress()))

	// Aguarda sinal de interrupção para graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	logger.Logger.Info("Encerrando servidor...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		logger.Logger.Fatal("Erro ao encerrar servidor", zap.Error(err))
	}

	logger.Logger.Info("Servidor encerrado")
}

// setupRoutes configura as rotas da API
func setupRoutes(e *echo.Echo, db *sqlx.DB, cfg *config.Config) {
	// Inicializa FileStorage
	fileStorage, err := storage.NewFileStorage(cfg)
	if err != nil {
		logger.Logger.Fatal("Erro ao inicializar FileStorage", zap.Error(err))
	}

	// Inicializa PDFProcessor
	pdfProcessor, err := pdf.NewPDFCPUProcessor()
	if err != nil {
		logger.Logger.Fatal("Erro ao inicializar PDFProcessor", zap.Error(err))
	}

	// Inicializa Repositories
	documentRepo := repository.NewDocumentRepository(db)
	auditLogRepo := repository.NewAuditLogRepository(db)

	// Inicializa UseCases
	documentUseCase := usecase.NewDocumentUseCase(
		documentRepo,
		auditLogRepo,
		fileStorage,
		pdfProcessor,
		cfg.Storage.Path,
	)
	previewUseCase := usecase.NewPDFPreviewUseCase(
		documentRepo,
		pdfProcessor,
		fileStorage,
		cfg.Storage.Path,
	)

	// Inicializa Handlers
	documentHandler := handler.NewDocumentHandler(
		documentUseCase,
		previewUseCase,
		cfg.Storage.MaxUploadSize,
	)

	// API v1
	v1 := e.Group("/api/v1")
	{
		// Rotas públicas de documentos
		documents := v1.Group("/documents")
		{
			documents.POST("", documentHandler.UploadDocument)
			documents.GET("", documentHandler.ListDocuments)
			documents.GET("/:id", documentHandler.GetDocument)
			documents.POST("/:id/process", documentHandler.ProcessDocument)
			documents.GET("/:id/preview/:page", documentHandler.GeneratePreview)
			documents.DELETE("/:id", documentHandler.DeleteDocument)
		}
	}
}
