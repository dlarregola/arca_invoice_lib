package client

import (
	"time"

	"github.com/dlarregola/arca_invoice_lib/pkg/models"
)

// Config representa la configuración del cliente ARCA
type Config struct {
	// Configuración básica
	Environment models.Environment `json:"environment" yaml:"environment"`
	CUIT        string             `json:"cuit" yaml:"cuit"`
	Certificate []byte             `json:"certificate" yaml:"certificate"`
	PrivateKey  []byte             `json:"private_key" yaml:"private_key"`

	// Configuración de red
	Timeout       time.Duration `json:"timeout" yaml:"timeout"`
	RetryAttempts int           `json:"retry_attempts" yaml:"retry_attempts"`
	RetryDelay    time.Duration `json:"retry_delay" yaml:"retry_delay"`

	// Configuración de logging
	LogLevel     string `json:"log_level" yaml:"log_level"`
	LogRequests  bool   `json:"log_requests" yaml:"log_requests"`
	LogResponses bool   `json:"log_responses" yaml:"log_responses"`

	// Configuración de autenticación
	AuthCacheTTL time.Duration `json:"auth_cache_ttl" yaml:"auth_cache_ttl"`
}

// DefaultConfig retorna una configuración por defecto
func DefaultConfig() Config {
	return Config{
		Environment:   models.EnvironmentTesting,
		Timeout:       30 * time.Second,
		RetryAttempts: 3,
		RetryDelay:    1 * time.Second,
		LogLevel:      "info",
		LogRequests:   false,
		LogResponses:  false,
		AuthCacheTTL:  23 * time.Hour, // Cache por 23 horas (tokens expiran en 24h)
	}
}

// Validate valida la configuración
func (c *Config) Validate() error {
	var errors models.ValidationErrors

	// Validar environment
	if c.Environment != models.EnvironmentTesting && c.Environment != models.EnvironmentProduction {
		errors.Add("environment", "Environment debe ser 'testing' o 'production'", c.Environment)
	}

	// Validar CUIT
	if c.CUIT == "" {
		errors.Add("cuit", "CUIT no puede estar vacío", c.CUIT)
	} else {
		if err := validateCUIT(c.CUIT); err != nil {
			errors.Add("cuit", err.Error(), c.CUIT)
		}
	}

	// Validar certificado
	if len(c.Certificate) == 0 {
		errors.Add("certificate", "Certificado no puede estar vacío", nil)
	}

	// Validar clave privada
	if len(c.PrivateKey) == 0 {
		errors.Add("private_key", "Clave privada no puede estar vacía", nil)
	}

	// Validar timeout
	if c.Timeout <= 0 {
		errors.Add("timeout", "Timeout debe ser mayor a 0", c.Timeout)
	}

	// Validar retry attempts
	if c.RetryAttempts < 0 {
		errors.Add("retry_attempts", "Retry attempts no puede ser negativo", c.RetryAttempts)
	}

	// Validar retry delay
	if c.RetryDelay < 0 {
		errors.Add("retry_delay", "Retry delay no puede ser negativo", c.RetryDelay)
	}

	// Validar auth cache TTL
	if c.AuthCacheTTL <= 0 {
		errors.Add("auth_cache_ttl", "Auth cache TTL debe ser mayor a 0", c.AuthCacheTTL)
	}

	if errors.HasErrors() {
		return errors
	}

	return nil
}

// GetBaseURL retorna la URL base según el environment
func (c *Config) GetBaseURL() string {
	switch c.Environment {
	case models.EnvironmentTesting:
		return "https://wswhomo.afip.gov.ar"
	case models.EnvironmentProduction:
		return "https://servicios1.afip.gov.ar"
	default:
		return "https://wswhomo.afip.gov.ar"
	}
}

// GetWSAAURL retorna la URL del servicio WSAA
func (c *Config) GetWSAAURL() string {
	return c.GetBaseURL() + "/ws/services/LoginCms"
}

// GetWSFEURL retorna la URL del servicio WSFEv1
func (c *Config) GetWSFEURL() string {
	return c.GetBaseURL() + "/wsfev1/service.asmx"
}

// GetWSFEXURL retorna la URL del servicio WSFEXv1
func (c *Config) GetWSFEXURL() string {
	return c.GetBaseURL() + "/wsfexv1/service.asmx"
}

// validateCUIT valida el formato de un CUIT
func validateCUIT(cuit string) error {
	// Importar la función de validación desde utils
	// Esta es una implementación simplificada para evitar dependencias circulares
	if cuit == "" {
		return models.NewValidationError("cuit", "CUIT no puede estar vacío", cuit)
	}

	// Validar formato básico: XX-XXXXXXXX-X
	if len(cuit) != 13 || cuit[2] != '-' || cuit[11] != '-' {
		return models.NewValidationError("cuit", "CUIT debe tener formato XX-XXXXXXXX-X", cuit)
	}

	return nil
}

// WithEnvironment configura el environment
func (c *Config) WithEnvironment(env models.Environment) *Config {
	c.Environment = env
	return c
}

// WithCUIT configura el CUIT
func (c *Config) WithCUIT(cuit string) *Config {
	c.CUIT = cuit
	return c
}

// WithCertificate configura el certificado
func (c *Config) WithCertificate(cert []byte) *Config {
	c.Certificate = cert
	return c
}

// WithPrivateKey configura la clave privada
func (c *Config) WithPrivateKey(key []byte) *Config {
	c.PrivateKey = key
	return c
}

// WithTimeout configura el timeout
func (c *Config) WithTimeout(timeout time.Duration) *Config {
	c.Timeout = timeout
	return c
}

// WithRetryAttempts configura los intentos de retry
func (c *Config) WithRetryAttempts(attempts int) *Config {
	c.RetryAttempts = attempts
	return c
}

// WithRetryDelay configura el delay entre retries
func (c *Config) WithRetryDelay(delay time.Duration) *Config {
	c.RetryDelay = delay
	return c
}

// WithLogLevel configura el nivel de logging
func (c *Config) WithLogLevel(level string) *Config {
	c.LogLevel = level
	return c
}

// WithLogRequests habilita el logging de requests
func (c *Config) WithLogRequests(enabled bool) *Config {
	c.LogRequests = enabled
	return c
}

// WithLogResponses habilita el logging de responses
func (c *Config) WithLogResponses(enabled bool) *Config {
	c.LogResponses = enabled
	return c
}

// WithAuthCacheTTL configura el TTL del cache de autenticación
func (c *Config) WithAuthCacheTTL(ttl time.Duration) *Config {
	c.AuthCacheTTL = ttl
	return c
}
