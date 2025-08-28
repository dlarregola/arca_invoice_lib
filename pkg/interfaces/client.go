package interfaces

import (
	"context"
	"time"
)

// Logger es la interfaz para logging
type Logger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
}

// ARCAClientManager es la interfaz principal para el manager multi-tenant
type ARCAClientManager interface {
	// GetClientForCompany obtiene un cliente específico para una empresa
	// La configuración se pasa directamente en tiempo de ejecución para mejor escalabilidad
	GetClientForCompany(ctx context.Context, companyConfig CompanyConfig) (ARCAClient, error)

	// ValidateCompanyConfig valida la configuración de una empresa
	ValidateCompanyConfig(config CompanyConfig) error

	// CleanupInactiveClients limpia el cache de clientes inactivos
	CleanupInactiveClients(maxIdleTime time.Duration)

	// InvalidateClient invalida el cache de un cliente específico
	InvalidateClient(companyID string)

	// GetCacheStats retorna estadísticas del cache
	GetCacheStats() CacheStats
}

// ARCAClient es la interfaz para un cliente de una empresa específica
type ARCAClient interface {
	// WSFE retorna el servicio de facturación nacional
	WSFE() WSFEService

	// WSFEX retorna el servicio de facturación internacional
	WSFEX() WSFEXService

	// GetCompanyInfo retorna información de la empresa
	GetCompanyInfo() CompanyInfo

	// IsHealthy verifica el estado de la conexión
	IsHealthy(ctx context.Context) error

	// Close cierra el cliente y limpia recursos
	Close() error
}

// CompanyConfig es la interfaz para configuración de empresa
type CompanyConfig interface {
	// GetCUIT retorna el CUIT de la empresa
	GetCUIT() string

	// GetCertificate retorna el certificado de la empresa
	GetCertificate() []byte

	// GetPrivateKey retorna la clave privada de la empresa
	GetPrivateKey() []byte

	// GetEnvironment retorna el ambiente ("testing" | "production")
	GetEnvironment() string

	// GetCompanyID retorna el ID único de la empresa
	GetCompanyID() string
}

// CompanyInfo representa información de la empresa
type CompanyInfo struct {
	CompanyID   string `json:"company_id"`
	CUIT        string `json:"cuit"`
	Environment string `json:"environment"`
}

// CacheStats representa estadísticas del cache
type CacheStats struct {
	TotalClients    int           `json:"total_clients"`
	ActiveClients   int           `json:"active_clients"`
	InactiveClients int           `json:"inactive_clients"`
	LastCleanup     time.Time     `json:"last_cleanup"`
	MaxIdleTime     time.Duration `json:"max_idle_time"`
}
