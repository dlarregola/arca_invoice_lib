package client

import (
	"github.com/dlarregola/arca_invoice_lib/internal/shared"
	"github.com/dlarregola/arca_invoice_lib/pkg/errors"
	"github.com/dlarregola/arca_invoice_lib/pkg/interfaces"
	"context"
	"fmt"
	"sync"
	"time"
)

// ManagerConfig representa la configuración del manager
type ManagerConfig struct {
	// Configuración de cache
	ClientCacheSize   int
	ClientIdleTimeout time.Duration

	// Configuración de red
	HTTPTimeout      time.Duration
	MaxRetryAttempts int

	// Logging
	Logger Logger
}

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

// clientManager es la implementación privada del manager multi-tenant
type clientManager struct {
	clientCache  map[string]*cachedClient
	cacheMutex   sync.RWMutex
	config       ManagerConfig
	lastCleanup  time.Time
	cleanupMutex sync.Mutex
}

// cachedClient representa un cliente en cache
type cachedClient struct {
	client    interfaces.AFIPClient
	lastUsed  time.Time
	companyID string
	createdAt time.Time
}

// internalConfig representa la configuración interna del cliente
type internalConfig = shared.InternalConfig

// newClientManager crea una nueva instancia del manager
func NewClientManager(config ManagerConfig) interfaces.AFIPClientManager {
	return &clientManager{
		clientCache: make(map[string]*cachedClient),
		config:      config,
		lastCleanup: time.Now(),
	}
}

// GetClientForCompany obtiene un cliente específico para una empresa
func (m *clientManager) GetClientForCompany(ctx context.Context, companyConfig interfaces.CompanyConfig) (interfaces.AFIPClient, error) {
	// Validar configuración
	if err := m.ValidateCompanyConfig(companyConfig); err != nil {
		return nil, fmt.Errorf("invalid company config: %w", err)
	}

	companyID := companyConfig.GetCompanyID()

	// Verificar cache primero
	if client := m.getCachedClient(companyID); client != nil {
		return client, nil
	}

	// Crear nuevo cliente
	client, err := m.createNewClient(companyConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	// Guardar en cache
	m.cacheClient(companyID, client)

	return client, nil
}

// ValidateCompanyConfig valida la configuración de una empresa
func (m *clientManager) ValidateCompanyConfig(config interfaces.CompanyConfig) error {
	if config == nil {
		return errors.NewCompanyConfigError("", "config", "configuration cannot be nil")
	}

	companyID := config.GetCompanyID()
	if companyID == "" {
		return errors.NewCompanyConfigError(companyID, "company_id", "company ID cannot be empty")
	}

	if config.GetCUIT() == "" {
		return errors.NewCompanyConfigError(companyID, "cuit", "CUIT cannot be empty")
	}

	if len(config.GetCertificate()) == 0 {
		return errors.NewCompanyConfigError(companyID, "certificate", "certificate cannot be empty")
	}

	if len(config.GetPrivateKey()) == 0 {
		return errors.NewCompanyConfigError(companyID, "private_key", "private key cannot be empty")
	}

	env := config.GetEnvironment()
	if env != "testing" && env != "production" {
		return errors.NewCompanyConfigError(companyID, "environment", "environment must be 'testing' or 'production'")
	}

	return nil
}

// CleanupInactiveClients limpia el cache de clientes inactivos
func (m *clientManager) CleanupInactiveClients(maxIdleTime time.Duration) {
	m.cleanupMutex.Lock()
	defer m.cleanupMutex.Unlock()

	now := time.Now()
	cutoff := now.Add(-maxIdleTime)

	m.cacheMutex.Lock()
	defer m.cacheMutex.Unlock()

	var toRemove []string
	for companyID, cached := range m.clientCache {
		if cached.lastUsed.Before(cutoff) {
			toRemove = append(toRemove, companyID)
		}
	}

	for _, companyID := range toRemove {
		if cached, exists := m.clientCache[companyID]; exists {
			// Cerrar cliente antes de remover
			if err := cached.client.Close(); err != nil {
				m.config.Logger.Warnf("Error closing client for company %s: %v", companyID, err)
			}
			delete(m.clientCache, companyID)
			m.config.Logger.Infof("Removed inactive client for company %s", companyID)
		}
	}

	m.lastCleanup = now
}

// InvalidateClient invalida el cache de un cliente específico
func (m *clientManager) InvalidateClient(companyID string) {
	m.cacheMutex.Lock()
	defer m.cacheMutex.Unlock()

	if cached, exists := m.clientCache[companyID]; exists {
		// Cerrar cliente antes de remover
		if err := cached.client.Close(); err != nil {
			m.config.Logger.Warnf("Error closing client for company %s: %v", companyID, err)
		}
		delete(m.clientCache, companyID)
		m.config.Logger.Infof("Invalidated client for company %s", companyID)
	}
}

// GetCacheStats retorna estadísticas del cache
func (m *clientManager) GetCacheStats() interfaces.CacheStats {
	m.cacheMutex.RLock()
	defer m.cacheMutex.RUnlock()

	now := time.Now()
	cutoff := now.Add(-m.config.ClientIdleTimeout)

	activeCount := 0
	inactiveCount := 0

	for _, cached := range m.clientCache {
		if cached.lastUsed.After(cutoff) {
			activeCount++
		} else {
			inactiveCount++
		}
	}

	return interfaces.CacheStats{
		TotalClients:    len(m.clientCache),
		ActiveClients:   activeCount,
		InactiveClients: inactiveCount,
		LastCleanup:     m.lastCleanup,
		MaxIdleTime:     m.config.ClientIdleTimeout,
	}
}

// getCachedClient obtiene un cliente del cache
func (m *clientManager) getCachedClient(companyID string) interfaces.AFIPClient {
	m.cacheMutex.RLock()
	defer m.cacheMutex.RUnlock()

	cached, exists := m.clientCache[companyID]
	if !exists {
		return nil
	}

	// Verificar si el cliente aún es válido
	if time.Since(cached.lastUsed) > m.config.ClientIdleTimeout {
		// Cliente expirado, remover del cache
		m.cacheMutex.RUnlock()
		m.cacheMutex.Lock()
		delete(m.clientCache, companyID)
		m.cacheMutex.Unlock()
		m.cacheMutex.RLock()
		return nil
	}

	// Actualizar último uso
	cached.lastUsed = time.Now()
	return cached.client
}

// cacheClient guarda un cliente en el cache
func (m *clientManager) cacheClient(companyID string, client interfaces.AFIPClient) {
	m.cacheMutex.Lock()
	defer m.cacheMutex.Unlock()

	// Verificar límite de cache
	if len(m.clientCache) >= m.config.ClientCacheSize {
		// Remover el cliente más antiguo
		var oldestCompanyID string
		var oldestTime time.Time
		for id, cached := range m.clientCache {
			if oldestCompanyID == "" || cached.lastUsed.Before(oldestTime) {
				oldestCompanyID = id
				oldestTime = cached.lastUsed
			}
		}
		if oldestCompanyID != "" {
			if cached, exists := m.clientCache[oldestCompanyID]; exists {
				if err := cached.client.Close(); err != nil {
					m.config.Logger.Warnf("Error closing old client for company %s: %v", oldestCompanyID, err)
				}
				delete(m.clientCache, oldestCompanyID)
			}
		}
	}

	m.clientCache[companyID] = &cachedClient{
		client:    client,
		lastUsed:  time.Now(),
		companyID: companyID,
		createdAt: time.Now(),
	}
}

// createNewClient crea un nuevo cliente AFIP
func (m *clientManager) createNewClient(config interfaces.CompanyConfig) (interfaces.AFIPClient, error) {
	// Crear configuración interna
	internalConfig := &internalConfig{
		CUIT:          config.GetCUIT(),
		Certificate:   config.GetCertificate(),
		PrivateKey:    config.GetPrivateKey(),
		Environment:   config.GetEnvironment(),
		Timeout:       m.config.HTTPTimeout,
		RetryAttempts: m.config.MaxRetryAttempts,
	}

	// Crear cliente interno
	client := &afipClient{
		companyConfig: config,
		config:        internalConfig,
		logger:        m.config.Logger,
	}

	// Inicializar servicios
	if err := client.initializeServices(); err != nil {
		return nil, fmt.Errorf("failed to initialize services: %w", err)
	}

	return client, nil
}
