package factory

import (
	"time"

	"github.com/dlarregola/arca_invoice_lib/internal/client"
	"github.com/dlarregola/arca_invoice_lib/pkg/interfaces"
)

// ClientManagerFactory es la interfaz para crear managers
type ClientManagerFactory interface {
	CreateManager() interfaces.ARCAClientManager
}

// clientManagerFactory es la implementación privada del factory
type clientManagerFactory struct {
	config client.ManagerConfig
}

// NewClientManagerFactory crea una nueva instancia del factor
// Add config params to the factory to override the default values
func NewClientManagerFactory(cacheSize int, idleTimeout time.Duration, httpTimeout time.Duration, maxRetryAttempts int, logger interfaces.Logger) ClientManagerFactory {
	config := client.ManagerConfig{
		ClientCacheSize:   cacheSize,
		ClientIdleTimeout: idleTimeout,
		HTTPTimeout:       httpTimeout,
		MaxRetryAttempts:  maxRetryAttempts,
		Logger:            logger,
	}
	return &clientManagerFactory{config: config}
}

// CreateManager crea un nuevo manager con la configuración especificada
func (f *clientManagerFactory) CreateManager() interfaces.ARCAClientManager {
	// Configurar valores por defecto
	if f.config.ClientCacheSize <= 0 {
		f.config.ClientCacheSize = 100
	}
	if f.config.ClientIdleTimeout <= 0 {
		f.config.ClientIdleTimeout = 30 * time.Minute
	}
	if f.config.HTTPTimeout <= 0 {
		f.config.HTTPTimeout = 30 * time.Second
	}
	if f.config.MaxRetryAttempts <= 0 {
		f.config.MaxRetryAttempts = 3
	}
	if f.config.Logger == nil {
		f.config.Logger = &noopLogger{}
	}

	// Crear manager
	return createClientManager(f.config)
}

// noopLogger es un logger que no hace nada
type noopLogger struct{}

func (l *noopLogger) Debug(args ...interface{})                 {}
func (l *noopLogger) Debugf(format string, args ...interface{}) {}
func (l *noopLogger) Info(args ...interface{})                  {}
func (l *noopLogger) Infof(format string, args ...interface{})  {}
func (l *noopLogger) Warn(args ...interface{})                  {}
func (l *noopLogger) Warnf(format string, args ...interface{})  {}
func (l *noopLogger) Error(args ...interface{})                 {}
func (l *noopLogger) Errorf(format string, args ...interface{}) {}

// createClientManager crea una nueva instancia del manager
// Esta función debe ser implementada en el paquete internal/client
func createClientManager(config client.ManagerConfig) interfaces.ARCAClientManager {
	// Esta es una función stub que será reemplazada por la implementación real
	// en el paquete internal/client
	return client.NewClientManager(config)
}
