package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/dlarregola/arca_invoice_lib/pkg/factory"
	"github.com/dlarregola/arca_invoice_lib/pkg/interfaces"
	"github.com/dlarregola/arca_invoice_lib/pkg/models"
)

// AdvancedCompanyConfiguration implementa la interfaz CompanyConfig
type AdvancedCompanyConfiguration struct {
	CompanyID   string
	CUIT        string
	Certificate []byte
	PrivateKey  []byte
	Environment string
}

// Implementar métodos de la interfaz
func (c *AdvancedCompanyConfiguration) GetCUIT() string        { return c.CUIT }
func (c *AdvancedCompanyConfiguration) GetCertificate() []byte { return c.Certificate }
func (c *AdvancedCompanyConfiguration) GetPrivateKey() []byte  { return c.PrivateKey }
func (c *AdvancedCompanyConfiguration) GetEnvironment() string { return c.Environment }
func (c *AdvancedCompanyConfiguration) GetCompanyID() string   { return c.CompanyID }

// AdvancedLogger implementa logging estructurado
type AdvancedLogger struct{}

func (l *AdvancedLogger) log(level, message string, fields map[string]interface{}) {
	log.Printf("[%s] %s - %v", level, message, fields)
}

func (l *AdvancedLogger) Debug(args ...interface{}) { l.log("DEBUG", fmt.Sprint(args...), nil) }
func (l *AdvancedLogger) Debugf(format string, args ...interface{}) {
	l.log("DEBUG", fmt.Sprintf(format, args...), nil)
}
func (l *AdvancedLogger) Info(args ...interface{}) { l.log("INFO", fmt.Sprint(args...), nil) }
func (l *AdvancedLogger) Infof(format string, args ...interface{}) {
	l.log("INFO", fmt.Sprintf(format, args...), nil)
}
func (l *AdvancedLogger) Warn(args ...interface{}) { l.log("WARN", fmt.Sprint(args...), nil) }
func (l *AdvancedLogger) Warnf(format string, args ...interface{}) {
	l.log("WARN", fmt.Sprintf(format, args...), nil)
}
func (l *AdvancedLogger) Error(args ...interface{}) { l.log("ERROR", fmt.Sprint(args...), nil) }
func (l *AdvancedLogger) Errorf(format string, args ...interface{}) {
	l.log("ERROR", fmt.Sprintf(format, args...), nil)
}

// DatabaseCompanyConfigProvider obtiene configuraciones desde base de datos
type DatabaseCompanyConfigProvider struct {
	db *sql.DB
}

func NewDatabaseCompanyConfigProvider(db *sql.DB) *DatabaseCompanyConfigProvider {
	return &DatabaseCompanyConfigProvider{db: db}
}

func (p *DatabaseCompanyConfigProvider) GetCompanyConfig(ctx context.Context, companyID string) (interfaces.CompanyConfig, error) {
	// Consultar base de datos
	query := `SELECT cuit, certificate, private_key, environment 
              FROM companies WHERE company_id = ? AND active = true`

	var cuit, environment string
	var certData, keyData []byte

	err := p.db.QueryRowContext(ctx, query, companyID).Scan(&cuit, &certData, &keyData, &environment)
	if err != nil {
		return nil, fmt.Errorf("failed to get company config: %w", err)
	}

	return &AdvancedCompanyConfiguration{
		CompanyID:   companyID,
		CUIT:        cuit,
		Certificate: certData,
		PrivateKey:  keyData,
		Environment: environment,
	}, nil
}

// AdvancedInvoiceService maneja facturas para múltiples empresas
type AdvancedInvoiceService struct {
	arcaManager    interfaces.ARCAClientManager
	logger         interfaces.Logger
	configProvider *DatabaseCompanyConfigProvider
}

func NewAdvancedInvoiceService(arcaManager interfaces.ARCAClientManager, logger interfaces.Logger, configProvider *DatabaseCompanyConfigProvider) *AdvancedInvoiceService {
	return &AdvancedInvoiceService{
		arcaManager:    arcaManager,
		logger:         logger,
		configProvider: configProvider,
	}
}

// CreateInvoice crea una factura para una empresa específica
func (s *AdvancedInvoiceService) CreateInvoice(ctx context.Context, companyID string, invoiceData *models.Invoice) (*models.AuthorizationResponse, error) {
	s.logger.Infof("Creating invoice for company %s", companyID)

	// Obtener configuración de empresa
	companyConfig, err := s.configProvider.GetCompanyConfig(ctx, companyID)
	if err != nil {
		s.logger.Errorf("Failed to get company config: %v", err)
		return nil, fmt.Errorf("failed to get company config: %w", err)
	}

	// Obtener cliente específico de la empresa
	client, err := s.arcaManager.GetClientForCompany(ctx, companyConfig)
	if err != nil {
		s.logger.Errorf("Failed to get ARCA client: %v", err)
		return nil, fmt.Errorf("failed to get ARCA client: %w", err)
	}

	// Autorizar factura
	response, err := client.WSFE().AuthorizeInvoice(ctx, invoiceData)
	if err != nil {
		s.logger.Errorf("Failed to authorize invoice: %v", err)
		return nil, fmt.Errorf("failed to authorize invoice: %w", err)
	}

	s.logger.Infof("Invoice authorized successfully: CAE=%s, Number=%d", response.CAE, response.InvoiceNumber)
	return response, nil
}

// CreateInvoiceBatch procesa múltiples facturas en paralelo
func (s *AdvancedInvoiceService) CreateInvoiceBatch(ctx context.Context, jobs []AdvancedInvoiceJob) []AdvancedInvoiceResult {
	results := make([]AdvancedInvoiceResult, len(jobs))
	var wg sync.WaitGroup

	// Procesar trabajos en paralelo
	for i, job := range jobs {
		wg.Add(1)
		go func(index int, job AdvancedInvoiceJob) {
			defer wg.Done()

			// Crear contexto con timeout para cada trabajo
			jobCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
			defer cancel()

			response, err := s.CreateInvoice(jobCtx, job.CompanyID, job.Invoice)
			results[index] = AdvancedInvoiceResult{
				CompanyID: job.CompanyID,
				Response:  response,
				Error:     err,
			}
		}(i, job)
	}

	wg.Wait()
	return results
}

// AdvancedInvoiceJob representa un trabajo de facturación
type AdvancedInvoiceJob struct {
	CompanyID string
	Invoice   *models.Invoice
}

// AdvancedInvoiceResult representa el resultado de un trabajo
type AdvancedInvoiceResult struct {
	CompanyID string
	Response  *models.AuthorizationResponse
	Error     error
}

// AdvancedCacheManager maneja el cache de configuraciones
type AdvancedCacheManager struct {
	cache map[string]interfaces.CompanyConfig
	mutex sync.RWMutex
	ttl   time.Duration
}

func NewAdvancedCacheManager(ttl time.Duration) *AdvancedCacheManager {
	cm := &AdvancedCacheManager{
		cache: make(map[string]interfaces.CompanyConfig),
		ttl:   ttl,
	}

	// Iniciar limpieza periódica
	go cm.cleanupPeriodic()

	return cm
}

func (cm *AdvancedCacheManager) Get(companyID string) (interfaces.CompanyConfig, bool) {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	config, exists := cm.cache[companyID]
	return config, exists
}

func (cm *AdvancedCacheManager) Set(companyID string, config interfaces.CompanyConfig) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	cm.cache[companyID] = config
}

func (cm *AdvancedCacheManager) cleanupPeriodic() {
	ticker := time.NewTicker(cm.ttl)
	defer ticker.Stop()

	for range ticker.C {
		cm.mutex.Lock()
		cm.cache = make(map[string]interfaces.CompanyConfig)
		cm.mutex.Unlock()
	}
}

// PerformanceMonitor monitorea el rendimiento
type PerformanceMonitor struct {
	arcaManager interfaces.ARCAClientManager
	logger      interfaces.Logger
}

func NewPerformanceMonitor(arcaManager interfaces.ARCAClientManager, logger interfaces.Logger) *PerformanceMonitor {
	pm := &PerformanceMonitor{
		arcaManager: arcaManager,
		logger:      logger,
	}

	// Iniciar monitoreo periódico
	go pm.monitorPeriodic()

	return pm
}

func (pm *PerformanceMonitor) monitorPeriodic() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		stats := pm.arcaManager.GetCacheStats()
		pm.logger.Infof("Performance Stats: %+v", stats)

		// Alertar si hay muchos clientes inactivos
		if stats.InactiveClients > stats.TotalClients/2 {
			pm.logger.Warnf("High number of inactive clients: %d/%d", stats.InactiveClients, stats.TotalClients)
		}

		// Limpiar clientes inactivos si es necesario
		if stats.InactiveClients > 10 {
			pm.arcaManager.CleanupInactiveClients(15 * time.Minute)
			pm.logger.Infof("Cleaned up inactive clients")
		}
	}
}

func runAdvancedExample() {
	// 1. Configurar base de datos (simulado)
	db := setupMockDatabase()

	// 2. Crear factory y manager
	factory := factory.NewClientManagerFactory(200, 60*time.Minute, 60*time.Second, 5, &AdvancedLogger{})
	manager := factory.CreateManager()

	// 3. Crear proveedor de configuraciones
	configProvider := NewDatabaseCompanyConfigProvider(db)

	// 4. Crear servicio de facturación
	invoiceService := NewAdvancedInvoiceService(manager, &AdvancedLogger{}, configProvider)

	// 5. Crear monitor de performance
	monitor := NewPerformanceMonitor(manager, &AdvancedLogger{})
	_ = monitor // Usar para evitar warning de variable no utilizada

	// 6. Crear contexto
	ctx := context.Background()

	// 7. Ejemplo: Crear facturas para múltiples empresas
	companies := []string{"empresa-001", "empresa-002", "empresa-003"}

	for _, companyID := range companies {
		// Crear factura de ejemplo
		invoice := createAdvancedSampleInvoice()

		// Crear factura
		response, err := invoiceService.CreateInvoice(ctx, companyID, invoice)
		if err != nil {
			log.Printf("Failed to create invoice for %s: %v", companyID, err)
			continue
		}

		log.Printf("Invoice created for %s: CAE=%s, Number=%d",
			companyID, response.CAE, response.InvoiceNumber)
	}

	// 8. Ejemplo: Procesamiento en lote
	batchJobs := []AdvancedInvoiceJob{
		{CompanyID: "empresa-001", Invoice: createAdvancedSampleInvoice()},
		{CompanyID: "empresa-002", Invoice: createAdvancedSampleInvoice()},
		{CompanyID: "empresa-003", Invoice: createAdvancedSampleInvoice()},
	}

	results := invoiceService.CreateInvoiceBatch(ctx, batchJobs)

	for _, result := range results {
		if result.Error != nil {
			log.Printf("Batch job failed for %s: %v", result.CompanyID, result.Error)
		} else {
			log.Printf("Batch job succeeded for %s: CAE=%s",
				result.CompanyID, result.Response.CAE)
		}
	}

	// 9. Mostrar estadísticas finales
	stats := manager.GetCacheStats()
	log.Printf("Final Cache Stats: Total=%d, Active=%d, Inactive=%d",
		stats.TotalClients, stats.ActiveClients, stats.InactiveClients)

	// 10. Simular trabajo continuo
	log.Println("Starting continuous work simulation...")
	time.Sleep(10 * time.Second)

	// 11. Limpiar recursos
	log.Println("Cleaning up resources...")
	manager.CleanupInactiveClients(0) // Limpiar todos los clientes
}

// setupMockDatabase configura una base de datos simulada
func setupMockDatabase() *sql.DB {
	// En un caso real, esto conectaría a una base de datos real
	// Por ahora, retornamos nil para simular
	return nil
}

// createAdvancedSampleInvoice crea una factura de ejemplo
func createAdvancedSampleInvoice() *models.Invoice {
	return &models.Invoice{
		InvoiceBase: models.InvoiceBase{
			InvoiceType:  models.InvoiceTypeA,
			PointOfSale:  1,
			DateFrom:     time.Now(),
			DateTo:       time.Now(),
			ConceptType:  models.ConceptTypeProducts,
			CurrencyType: models.CurrencyTypePES,
			Amount:       1000.0,
			TaxAmount:    210.0,
			TotalAmount:  1210.0,
			Items: []models.Item{
				{
					Description: "Producto de ejemplo",
					Quantity:    1,
					UnitPrice:   1000.0,
					TotalPrice:  1000.0,
					Taxes: []models.Tax{
						{
							Type:   models.TaxTypeIVA,
							Rate:   models.TaxRate21,
							Base:   1000.0,
							Amount: 210.0,
						},
					},
				},
			},
		},
		DocType:       models.DocumentTypeCUIT,
		DocNumber:     "20-12345678-9",
		DocTypeFrom:   models.DocumentTypeCUIT,
		DocNumberFrom: "20-87654321-0",
		NameFrom:      "Cliente Ejemplo S.A.",
		ServiceFrom:   time.Now(),
	}
}
