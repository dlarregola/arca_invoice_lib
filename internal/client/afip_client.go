package client

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/dlarregola/arca_invoice_lib/internal/services/auth"
	"github.com/dlarregola/arca_invoice_lib/internal/services/wsfe"
	"github.com/dlarregola/arca_invoice_lib/internal/services/wsfex"
	"github.com/dlarregola/arca_invoice_lib/internal/shared"
	"github.com/dlarregola/arca_invoice_lib/pkg/interfaces"
)

// arcaClient es la implementación privada del cliente ARCA
type arcaClient struct {
	companyConfig interfaces.CompanyConfig
	config        *shared.InternalConfig
	wsfeService   interfaces.WSFEService
	wsfexService  interfaces.WSFEXService
	authService   interfaces.AuthService
	httpClient    *http.Client
	logger        interfaces.Logger
	mutex         sync.RWMutex
	closed        bool
}

// WSFE retorna el servicio de facturación nacional
func (c *arcaClient) WSFE() interfaces.WSFEService {
	return c.wsfeService
}

// WSFEX retorna el servicio de facturación internacional
func (c *arcaClient) WSFEX() interfaces.WSFEXService {
	return c.wsfexService
}

// GetCompanyInfo retorna información de la empresa
func (c *arcaClient) GetCompanyInfo() interfaces.CompanyInfo {
	return interfaces.CompanyInfo{
		CompanyID:   c.companyConfig.GetCompanyID(),
		CUIT:        c.companyConfig.GetCUIT(),
		Environment: c.companyConfig.GetEnvironment(),
	}
}

// IsHealthy verifica el estado de la conexión
func (c *arcaClient) IsHealthy(ctx context.Context) error {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if c.closed {
		return fmt.Errorf("client is closed")
	}

	// Intentar obtener un token para verificar la conexión
	_, err := c.authService.GetToken(ctx, "wsfe")
	if err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}

	return nil
}

// Close cierra el cliente y limpia recursos
func (c *arcaClient) Close() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.closed {
		return nil
	}

	c.closed = true

	// Limpiar cache de autenticación
	c.authService.ClearCache()

	// Cerrar HTTP client si es necesario
	if c.httpClient != nil {
		// El http.Client no tiene un método Close(), pero podríamos
		// implementar un transport personalizado si es necesario
	}

	c.logger.Infof("Client closed for company %s", c.companyConfig.GetCompanyID())
	return nil
}

// initializeServices inicializa los servicios del cliente
func (c *arcaClient) initializeServices() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Crear HTTP client
	c.httpClient = &http.Client{
		Timeout: c.config.Timeout,
	}

	// Crear servicio de autenticación
	c.authService = auth.NewAuthService(c.config, c.logger)

	// Crear servicio WSFE
	wsfeService, err := wsfe.NewWSFEService(c.authService, c.logger)
	if err != nil {
		return fmt.Errorf("failed to create WSFE service: %w", err)
	}
	c.wsfeService = wsfeService

	// Crear servicio WSFEX
	wsfexService, err := wsfex.NewWSFEXService(c.authService, c.logger)
	if err != nil {
		return fmt.Errorf("failed to create WSFEX service: %w", err)
	}
	c.wsfexService = wsfexService

	c.logger.Infof("Services initialized for company %s", c.companyConfig.GetCompanyID())
	return nil
}

// getBaseURL retorna la URL base según el environment
func (c *arcaClient) getBaseURL() string {
	switch c.config.Environment {
	case "testing":
		return "https://wswhomo.afip.gov.ar"
	case "production":
		return "https://servicios1.afip.gov.ar"
	default:
		return "https://wswhomo.afip.gov.ar"
	}
}

// getWSAAURL retorna la URL del servicio WSAA
func (c *arcaClient) getWSAAURL() string {
	return c.getBaseURL() + "/ws/services/LoginCms"
}

// getWSFEURL retorna la URL del servicio WSFEv1
func (c *arcaClient) getWSFEURL() string {
	return c.getBaseURL() + "/wsfev1/service.asmx"
}

// getWSFEXURL retorna la URL del servicio WSFEXv1
func (c *arcaClient) getWSFEXURL() string {
	return c.getBaseURL() + "/wsfexv1/service.asmx"
}
