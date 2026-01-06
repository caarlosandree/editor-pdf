package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config contém todas as configurações da aplicação
type Config struct {
	Server  ServerConfig  `mapstructure:"server"`
	DB      DBConfig      `mapstructure:"db"`
	JWT     JWTConfig     `mapstructure:"jwt"`
	CORS    CORSConfig    `mapstructure:"cors"`
	Storage StorageConfig `mapstructure:"storage"`
	Env     string        `mapstructure:"env"`
}

// ServerConfig contém configurações do servidor
type ServerConfig struct {
	Port string `mapstructure:"port"`
	Host string `mapstructure:"host"`
}

// DBConfig contém configurações do banco de dados
type DBConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	SSLMode  string `mapstructure:"sslmode"`
}

// JWTConfig contém configurações JWT
type JWTConfig struct {
	Secret     string `mapstructure:"secret"`
	Expiration string `mapstructure:"expiration"`
}

// CORSConfig contém configurações CORS
type CORSConfig struct {
	AllowedOrigins []string `mapstructure:"allowed_origins"`
}

// StorageConfig contém configurações de armazenamento de arquivos
type StorageConfig struct {
	Path         string `mapstructure:"path"`
	MaxUploadSize int64  `mapstructure:"max_upload_size"` // em bytes
}

// DSN retorna a string de conexão do PostgreSQL
func (c *DBConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		c.Host,
		c.User,
		c.Password,
		c.Name,
		c.Port,
		c.SSLMode,
	)
}

// LoadConfig carrega as configurações do arquivo .env.local, .env e variáveis de ambiente
func LoadConfig() (*Config, error) {
	viper.SetConfigType("env")
	
	// Permite que variáveis de ambiente sobrescrevam valores do arquivo
	viper.AutomaticEnv()

	// Mapeia variáveis de ambiente para a estrutura
	viper.SetDefault("SERVER_PORT", "8080")
	viper.SetDefault("SERVER_HOST", "localhost")
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_USER", "postgres")
	viper.SetDefault("DB_PASSWORD", "postgres")
	viper.SetDefault("DB_NAME", "editor_pdf")
	viper.SetDefault("DB_SSLMODE", "disable")
	viper.SetDefault("JWT_SECRET", "")
	viper.SetDefault("JWT_EXPIRATION", "24h")
	viper.SetDefault("CORS_ALLOWED_ORIGINS", "http://localhost:3000")
	viper.SetDefault("STORAGE_PATH", "./storage")
	viper.SetDefault("STORAGE_MAX_UPLOAD_SIZE", 104857600) // 100MB em bytes
	viper.SetDefault("ENV", "development")

	// Tenta ler primeiro o arquivo .env.local (prioridade maior)
	viper.SetConfigName(".env.local")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./..")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("erro ao ler arquivo de configuração .env.local: %w", err)
		}
		// Se .env.local não existir, tenta ler .env
		viper.SetConfigName(".env")
		if err := viper.MergeInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				return nil, fmt.Errorf("erro ao ler arquivo de configuração .env: %w", err)
			}
		}
	}

	var config Config

	// Mapeia variáveis de ambiente para a estrutura
	config.Server.Port = viper.GetString("SERVER_PORT")
	config.Server.Host = viper.GetString("SERVER_HOST")
	config.DB.Host = viper.GetString("DB_HOST")
	config.DB.Port = viper.GetString("DB_PORT")
	config.DB.User = viper.GetString("DB_USER")
	config.DB.Password = viper.GetString("DB_PASSWORD")
	config.DB.Name = viper.GetString("DB_NAME")
	config.DB.SSLMode = viper.GetString("DB_SSLMODE")
	config.JWT.Secret = viper.GetString("JWT_SECRET")
	config.JWT.Expiration = viper.GetString("JWT_EXPIRATION")
	config.Storage.Path = viper.GetString("STORAGE_PATH")
	config.Storage.MaxUploadSize = viper.GetInt64("STORAGE_MAX_UPLOAD_SIZE")
	config.Env = viper.GetString("ENV")

	// Parse CORS allowed origins
	corsOrigins := viper.GetString("CORS_ALLOWED_ORIGINS")
	if corsOrigins != "" {
		// Split por vírgula e remove espaços
		origins := strings.Split(corsOrigins, ",")
		config.CORS.AllowedOrigins = make([]string, 0, len(origins))
		for _, origin := range origins {
			origin = strings.TrimSpace(origin)
			if origin != "" {
				config.CORS.AllowedOrigins = append(config.CORS.AllowedOrigins, origin)
			}
		}
	}

	// Valida configurações obrigatórias
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("configuração inválida: %w", err)
	}

	return &config, nil
}

// validateConfig valida as configurações obrigatórias
func validateConfig(cfg *Config) error {
	if cfg.DB.Host == "" {
		return fmt.Errorf("DB_HOST é obrigatório")
	}
	if cfg.DB.User == "" {
		return fmt.Errorf("DB_USER é obrigatório")
	}
	if cfg.DB.Name == "" {
		return fmt.Errorf("DB_NAME é obrigatório")
	}
	if cfg.Storage.Path == "" {
		return fmt.Errorf("STORAGE_PATH é obrigatório")
	}
	if cfg.JWT.Secret == "" {
		return fmt.Errorf("JWT_SECRET é obrigatório")
	}
	return nil
}

// IsDevelopment retorna true se o ambiente é de desenvolvimento
func (c *Config) IsDevelopment() bool {
	return c.Env == "development"
}

// IsProduction retorna true se o ambiente é de produção
func (c *Config) IsProduction() bool {
	return c.Env == "production"
}

// GetServerAddress retorna o endereço completo do servidor
func (c *Config) GetServerAddress() string {
	return fmt.Sprintf("%s:%s", c.Server.Host, c.Server.Port)
}
