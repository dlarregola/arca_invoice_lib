package factory

import (
	"github.com/dlarregola/arca_invoice_lib/internal/client"
	"github.com/dlarregola/arca_invoice_lib/pkg/interfaces"
	"time"
)

// ClientManagerFactory es la interfaz para crear managers
type ClientManagerFactory interface {
	CreateManager(config client.ManagerConfig) interfaces.AFIPClientManager
}

// clientManagerFactory es la implementación privada del factory
type clientManagerFactory struct{}

// NewClientManagerFactory crea una nueva instancia del factory
func NewClientManagerFactory() ClientManagerFactory {
	return &clientManagerFactory{}
}

// CreateManager crea un nuevo manager con la configuración especificada
func (f *clientManagerFactory) CreateManager(config client.ManagerConfig) interfaces.AFIPClientManager {
	// Configurar valores por defecto
	if config.ClientCacheSize <= 0 {
		config.ClientCacheSize = 100
	}
	if config.ClientIdleTimeout <= 0 {
		config.ClientIdleTimeout = 30 * time.Minute
	}
	if config.HTTPTimeout <= 0 {
		config.HTTPTimeout = 30 * time.Second
	}
	if config.MaxRetryAttempts <= 0 {
		config.MaxRetryAttempts = 3
	}
	if config.Logger == nil {
		config.Logger = &noopLogger{}
	}

	// Crear manager
	return createClientManager(config)
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
func createClientManager(config client.ManagerConfig) interfaces.AFIPClientManager {
	// Esta es una función stub que será reemplazada por la implementación real
	// en el paquete internal/client
	return client.NewClientManager(config)
}
