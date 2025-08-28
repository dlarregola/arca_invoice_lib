package client

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// ARCAClient representa el cliente principal de ARCA
type ARCAClient struct {
	config      *Config
	auth        *WSAAAuth
	wsfe        interface{}
	wsfex       interface{}
	logger      interface{}
	loggerMutex sync.RWMutex
}

// NewARCAClient crea un nuevo cliente ARCA
func NewARCAClient(config Config) (*ARCAClient, error) {
	// Validar configuración
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	// Crear logger básico
	logger := &basicLogger{}

	// Crear autenticador
	auth := NewWSAAAuth(&config, logger)

	client := &ARCAClient{
		config: &config,
		auth:   auth,
		logger: logger,
	}

	return client, nil
}

// WSFE retorna el servicio de facturación nacional
func (c *ARCAClient) WSFE() interface{} {
	return c.wsfe
}

// WSFEX retorna el servicio de facturación internacional
func (c *ARCAClient) WSFEX() interface{} {
	return c.wsfex
}

// GetConfig retorna la configuración del cliente
func (c *ARCAClient) GetConfig() *Config {
	return c.config
}

// SetLogger establece un logger personalizado
func (c *ARCAClient) SetLogger(logger interface{}) {
	c.loggerMutex.Lock()
	defer c.loggerMutex.Unlock()
	c.logger = logger
}

// GetLogger retorna el logger actual
func (c *ARCAClient) GetLogger() interface{} {
	c.loggerMutex.RLock()
	defer c.loggerMutex.RUnlock()
	return c.logger
}

// ClearAuthCache limpia el cache de autenticación
func (c *ARCAClient) ClearAuthCache() {
	c.auth.ClearCache()
}

// GetAuthCacheSize retorna el tamaño del cache de autenticación
func (c *ARCAClient) GetAuthCacheSize() int {
	return c.auth.GetCacheSize()
}

// TestConnection prueba la conexión con ARCA
func (c *ARCAClient) TestConnection(ctx context.Context) error {
	// Intentar obtener un ticket de acceso para el servicio de testing
	_, err := c.auth.GetAccessTicket(ctx, "wsfe")
	if err != nil {
		return fmt.Errorf("connection test failed: %w", err)
	}

	return nil
}

// GetSystemStatus obtiene el estado del sistema ARCA
func (c *ARCAClient) GetSystemStatus(ctx context.Context) (*SystemStatus, error) {
	// Por ahora retornamos un status básico
	// En una implementación real, consultaríamos los parámetros del sistema
	return &SystemStatus{
		Status:     "ok",
		Message:    "System is operational",
		Timestamp:  time.Now(),
		LastUpdate: time.Now(),
	}, nil
}

// SystemStatus representa el estado del sistema ARCA
type SystemStatus struct {
	Status     string    `json:"status"`
	Message    string    `json:"message"`
	Timestamp  time.Time `json:"timestamp"`
	LastUpdate time.Time `json:"last_update,omitempty"`
}

// basicLogger implementa un logger básico
type basicLogger struct{}

func (l *basicLogger) Debug(args ...interface{})                 {}
func (l *basicLogger) Debugf(format string, args ...interface{}) {}
func (l *basicLogger) Info(args ...interface{})                  {}
func (l *basicLogger) Infof(format string, args ...interface{})  {}
func (l *basicLogger) Warn(args ...interface{})                  {}
func (l *basicLogger) Warnf(format string, args ...interface{})  {}
func (l *basicLogger) Error(args ...interface{})                 {}
func (l *basicLogger) Errorf(format string, args ...interface{}) {}
func (l *basicLogger) GetLevel() interface{}                     { return "info" }
