# Guía de Uso - AFIP Go Library

## Índice

1. [Instalación](#instalación)
2. [Configuración Inicial](#configuración-inicial)
3. [Uso Básico](#uso-básico)
4. [Patrón Multi-Tenant](#patrón-multi-tenant)
5. [Servicios Disponibles](#servicios-disponibles)
6. [Manejo de Errores](#manejo-de-errores)
7. [Configuración Avanzada](#configuración-avanzada)
8. [Ejemplos Prácticos](#ejemplos-prácticos)
9. [Mejores Prácticas](#mejores-prácticas)
10. [Troubleshooting](#troubleshooting)

## Instalación

### Requisitos Previos

- Go 1.19 o superior
- Certificado X.509 de AFIP
- Clave privada correspondiente
- CUIT habilitado en AFIP

### Instalación de la Librería

```bash
go get github.com/your-org/arca_invoice_lib
```

### Importar la Librería

```go
import (
    "arca_invoice_lib/pkg/factory"
    "arca_invoice_lib/pkg/interfaces"
    "arca_invoice_lib/pkg/models"
)
```

## Configuración Inicial

### 1. Obtener Certificados AFIP

Para usar la librería necesitas obtener certificados de AFIP:

1. **Registrarse en AFIP**: Crear cuenta en [AFIP](https://www.afip.gob.ar)
2. **Solicitar Certificado**: Generar certificado X.509 en el sistema de AFIP
3. **Descargar Archivos**: Obtener el certificado (.crt) y clave privada (.key)

### 2. Configuración de Empresa

Cada empresa debe implementar la interfaz `CompanyConfig`:

```go
type CompanyConfiguration struct {
    CompanyID   string
    CUIT        string
    Certificate []byte
    PrivateKey  []byte
    Environment string // "testing" o "production"
}

// Implementar métodos de la interfaz
func (c *CompanyConfiguration) GetCUIT() string { return c.CUIT }
func (c *CompanyConfiguration) GetCertificate() []byte { return c.Certificate }
func (c *CompanyConfiguration) GetPrivateKey() []byte { return c.PrivateKey }
func (c *CompanyConfiguration) GetEnvironment() string { return c.Environment }
func (c *CompanyConfiguration) GetCompanyID() string { return c.CompanyID }
```

### 3. Cargar Certificados

```go
func loadCertificates(certPath, keyPath string) ([]byte, []byte, error) {
    // Cargar certificado
    certData, err := os.ReadFile(certPath)
    if err != nil {
        return nil, nil, fmt.Errorf("failed to read certificate: %w", err)
    }

    // Cargar clave privada
    keyData, err := os.ReadFile(keyPath)
    if err != nil {
        return nil, nil, fmt.Errorf("failed to read private key: %w", err)
    }

    return certData, keyData, nil
}
```

## Uso Básico

### 1. Crear Manager

```go
package main

import (
    "context"
    "log"
    "time"
    
    "arca_invoice_lib/pkg/factory"
    "arca_invoice_lib/pkg/interfaces"
)

func main() {
    // 1. Crear factory
    factory := factory.NewClientManagerFactory()

    // 2. Configurar manager
    manager := factory.CreateManager(factory.ManagerConfig{
        ClientCacheSize:   100,                    // Máximo 100 clientes en cache
        ClientIdleTimeout: 30 * time.Minute,       // Timeout de inactividad
        HTTPTimeout:       30 * time.Second,       // Timeout HTTP
        MaxRetryAttempts:  3,                      // Reintentos
        Logger:            &MyLogger{},            // Logger personalizado
    })

    // 3. Crear configuración de empresa
    companyConfig := &CompanyConfiguration{
        CompanyID:   "empresa-001",
        CUIT:        "20-12345678-9",
        Certificate: certData,  // []byte del certificado
        PrivateKey:  keyData,   // []byte de la clave privada
        Environment: "testing", // "testing" o "production"
    }

    // 4. Obtener cliente
    ctx := context.Background()
    client, err := manager.GetClientForCompany(ctx, companyConfig)
    if err != nil {
        log.Fatal("Failed to get client:", err)
    }

    // 5. Usar servicios
    // ... (ver ejemplos de servicios)
}
```

### 2. Logger Personalizado

```go
type MyLogger struct{}

func (l *MyLogger) Debug(args ...interface{}) {
    log.Printf("[DEBUG] %v", args...)
}

func (l *MyLogger) Debugf(format string, args ...interface{}) {
    log.Printf("[DEBUG] "+format, args...)
}

func (l *MyLogger) Info(args ...interface{}) {
    log.Printf("[INFO] %v", args...)
}

func (l *MyLogger) Infof(format string, args ...interface{}) {
    log.Printf("[INFO] "+format, args...)
}

func (l *MyLogger) Warn(args ...interface{}) {
    log.Printf("[WARN] %v", args...)
}

func (l *MyLogger) Warnf(format string, args ...interface{}) {
    log.Printf("[WARN] "+format, args...)
}

func (l *MyLogger) Error(args ...interface{}) {
    log.Printf("[ERROR] %v", args...)
}

func (l *MyLogger) Errorf(format string, args ...interface{}) {
    log.Printf("[ERROR] "+format, args...)
}
```

## Patrón Multi-Tenant

### 1. Uso en Servicios

```go
// InvoiceService maneja facturas para múltiples empresas
type InvoiceService struct {
    afipManager interfaces.AFIPClientManager
}

func NewInvoiceService(afipManager interfaces.AFIPClientManager) *InvoiceService {
    return &InvoiceService{
        afipManager: afipManager,
    }
}

func (s *InvoiceService) CreateInvoice(ctx context.Context, companyConfig interfaces.CompanyConfig, invoiceData *models.Invoice) (*models.AuthorizationResponse, error) {
    // Obtener cliente específico de la empresa
    client, err := s.afipManager.GetClientForCompany(ctx, companyConfig)
    if err != nil {
        return nil, fmt.Errorf("failed to get AFIP client: %w", err)
    }
    
    // Usar servicio de facturación nacional
    response, err := client.WSFE().AuthorizeInvoice(ctx, invoiceData)
    if err != nil {
        return nil, fmt.Errorf("failed to authorize invoice: %w", err)
    }
    
    return response, nil
}

func (s *InvoiceService) CreateExportInvoice(ctx context.Context, companyConfig interfaces.CompanyConfig, exportInvoice *models.ExportInvoice) (*models.ExportAuthResponse, error) {
    // Obtener cliente específico de la empresa
    client, err := s.afipManager.GetClientForCompany(ctx, companyConfig)
    if err != nil {
        return nil, fmt.Errorf("failed to get AFIP client: %w", err)
    }
    
    // Usar servicio de facturación internacional
    response, err := client.WSFEX().AuthorizeExportInvoice(ctx, exportInvoice)
    if err != nil {
        return nil, fmt.Errorf("failed to authorize export invoice: %w", err)
    }
    
    return response, nil
}
```

### 2. Gestión de Cache

```go
func (s *InvoiceService) ManageCache() {
    // Obtener estadísticas del cache
    stats := s.afipManager.GetCacheStats()
    log.Printf("Cache Stats: Total=%d, Active=%d, Inactive=%d", 
        stats.TotalClients, stats.ActiveClients, stats.InactiveClients)

    // Limpiar clientes inactivos (más de 5 minutos)
    s.afipManager.CleanupInactiveClients(5 * time.Minute)

    // Invalidar cliente específico
    s.afipManager.InvalidateClient("empresa-001")
}
```

### 3. Configuración Dinámica

```go
// DatabaseCompanyConfigProvider obtiene configuraciones desde base de datos
type DatabaseCompanyConfigProvider struct {
    db *sql.DB
}

func (p *DatabaseCompanyConfigProvider) GetCompanyConfig(ctx context.Context, companyID string) (interfaces.CompanyConfig, error) {
    // Consultar base de datos
    query := `SELECT cuit, certificate, private_key, environment 
              FROM companies WHERE company_id = ?`
    
    var cuit, environment string
    var certData, keyData []byte
    
    err := p.db.QueryRowContext(ctx, query, companyID).Scan(&cuit, &certData, &keyData, &environment)
    if err != nil {
        return nil, fmt.Errorf("failed to get company config: %w", err)
    }
    
    return &CompanyConfiguration{
        CompanyID:   companyID,
        CUIT:        cuit,
        Certificate: certData,
        PrivateKey:  keyData,
        Environment: environment,
    }, nil
}
```

## Servicios Disponibles

### 1. Facturación Nacional (WSFE)

#### Autorizar Factura

```go
func (s *InvoiceService) CreateNationalInvoice(ctx context.Context, companyConfig interfaces.CompanyConfig) error {
    client, err := s.afipManager.GetClientForCompany(ctx, companyConfig)
    if err != nil {
        return err
    }

    // Crear factura
    invoice := &models.Invoice{
        InvoiceBase: models.InvoiceBase{
            InvoiceType:   models.InvoiceTypeA,
            PointOfSale:   1,
            DateFrom:      time.Now(),
            DateTo:        time.Now(),
            ConceptType:   models.ConceptTypeProducts,
            CurrencyType:  models.CurrencyTypePES,
            Amount:        1000.0,
            TaxAmount:     210.0,
            TotalAmount:   1210.0,
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

    // Autorizar factura
    response, err := client.WSFE().AuthorizeInvoice(ctx, invoice)
    if err != nil {
        return fmt.Errorf("failed to authorize invoice: %w", err)
    }

    log.Printf("Invoice authorized: CAE=%s, Number=%d", response.CAE, response.InvoiceNumber)
    return nil
}
```

#### Consultar Factura

```go
func (s *InvoiceService) QueryInvoice(ctx context.Context, companyConfig interfaces.CompanyConfig) error {
    client, err := s.afipManager.GetClientForCompany(ctx, companyConfig)
    if err != nil {
        return err
    }

    // Crear consulta
    query := &models.InvoiceQuery{
        InvoiceType:   models.InvoiceTypeA,
        PointOfSale:   1,
        InvoiceNumber: 1001,
        DateFrom:      time.Now().AddDate(0, 0, -30),
        DateTo:        time.Now(),
    }

    // Consultar factura
    invoice, err := client.WSFE().QueryInvoice(ctx, query)
    if err != nil {
        return fmt.Errorf("failed to query invoice: %w", err)
    }

    log.Printf("Invoice found: Number=%d, Total=%f", invoice.InvoiceNumber, invoice.TotalAmount)
    return nil
}
```

#### Obtener Último Comprobante

```go
func (s *InvoiceService) GetLastInvoice(ctx context.Context, companyConfig interfaces.CompanyConfig) error {
    client, err := s.afipManager.GetClientForCompany(ctx, companyConfig)
    if err != nil {
        return err
    }

    // Obtener último comprobante
    lastInvoice, err := client.WSFE().GetLastInvoice(ctx, models.InvoiceTypeA, 1)
    if err != nil {
        return fmt.Errorf("failed to get last invoice: %w", err)
    }

    log.Printf("Last invoice: Type=%d, Number=%d, Date=%s", 
        lastInvoice.InvoiceType, lastInvoice.InvoiceNumber, lastInvoice.Date)
    return nil
}
```

### 2. Facturación Internacional (WSFEX)

#### Autorizar Factura de Exportación

```go
func (s *InvoiceService) CreateExportInvoice(ctx context.Context, companyConfig interfaces.CompanyConfig) error {
    client, err := s.afipManager.GetClientForCompany(ctx, companyConfig)
    if err != nil {
        return err
    }

    // Crear factura de exportación
    exportInvoice := &models.ExportInvoice{
        InvoiceBase: models.InvoiceBase{
            InvoiceType:   models.InvoiceTypeE,
            PointOfSale:   1,
            DateFrom:      time.Now(),
            DateTo:        time.Now(),
            ConceptType:   models.ConceptTypeProducts,
            CurrencyType:  models.CurrencyTypeUSD,
            CurrencyRate:  100.0, // Tipo de cambio
            Amount:        1000.0,
            TaxAmount:     0.0, // Exportación no paga IVA
            TotalAmount:   1000.0,
            Items: []models.Item{
                {
                    Description: "Producto de exportación",
                    Quantity:    1,
                    UnitPrice:   1000.0,
                    TotalPrice:  1000.0,
                    Country:     "US",
                },
            },
        },
        Destination:     "Estados Unidos",
        DestinationCode: "US",
        ExportDate:      time.Now(),
        ExportType:      "Definitiva",
    }

    // Autorizar factura de exportación
    response, err := client.WSFEX().AuthorizeExportInvoice(ctx, exportInvoice)
    if err != nil {
        return fmt.Errorf("failed to authorize export invoice: %w", err)
    }

    log.Printf("Export invoice authorized: CAE=%s, Number=%d", response.CAE, response.InvoiceNumber)
    return nil
}
```

#### Consultar Factura de Exportación

```go
func (s *InvoiceService) QueryExportInvoice(ctx context.Context, companyConfig interfaces.CompanyConfig) error {
    client, err := s.afipManager.GetClientForCompany(ctx, companyConfig)
    if err != nil {
        return err
    }

    // Crear consulta
    query := &models.ExportInvoiceQuery{
        InvoiceType:   models.InvoiceTypeE,
        PointOfSale:   1,
        InvoiceNumber: 1001,
        DateFrom:      time.Now().AddDate(0, 0, -30),
        DateTo:        time.Now(),
    }

    // Consultar factura de exportación
    invoice, err := client.WSFEX().QueryExportInvoice(ctx, query)
    if err != nil {
        return fmt.Errorf("failed to query export invoice: %w", err)
    }

    log.Printf("Export invoice found: Number=%d, Destination=%s", 
        invoice.InvoiceNumber, invoice.Destination)
    return nil
}
```

### 3. Servicios de Consulta

#### Obtener Parámetros del Sistema

```go
func (s *InvoiceService) GetSystemParameters(ctx context.Context, companyConfig interfaces.CompanyConfig) error {
    client, err := s.afipManager.GetClientForCompany(ctx, companyConfig)
    if err != nil {
        return err
    }

    // Obtener parámetros
    params, err := client.WSFE().GetSystemParameters(ctx)
    if err != nil {
        return fmt.Errorf("failed to get system parameters: %w", err)
    }

    log.Printf("System parameters: DocumentTypes=%d, InvoiceTypes=%d", 
        len(params.DocumentTypes), len(params.InvoiceTypes))
    return nil
}
```

#### Obtener Tipos de Moneda

```go
func (s *InvoiceService) GetCurrencyTypes(ctx context.Context, companyConfig interfaces.CompanyConfig) error {
    client, err := s.afipManager.GetClientForCompany(ctx, companyConfig)
    if err != nil {
        return err
    }

    // Obtener tipos de moneda
    currencies, err := client.WSFEX().GetCurrencyTypes(ctx)
    if err != nil {
        return fmt.Errorf("failed to get currency types: %w", err)
    }

    for _, currency := range currencies {
        log.Printf("Currency: %s - %s (Active: %t)", 
            currency.ID, currency.Description, currency.Active)
    }
    return nil
}
```

## Manejo de Errores

### 1. Tipos de Errores

```go
func (s *InvoiceService) handleErrors(err error) {
    // Verificar tipo de error
    switch {
    case errors.Is(err, &errors.CompanyConfigError{}):
        var configErr *errors.CompanyConfigError
        if errors.As(err, &configErr) {
            log.Printf("Configuration error for company %s: %s", 
                configErr.CompanyID, configErr.Message)
        }
        
    case errors.Is(err, &errors.AuthenticationError{}):
        var authErr *errors.AuthenticationError
        if errors.As(err, &authErr) {
            log.Printf("Authentication error for service %s: %s", 
                authErr.Service, authErr.Message)
        }
        
    case errors.Is(err, &errors.InvoiceError{}):
        var invoiceErr *errors.InvoiceError
        if errors.As(err, &invoiceErr) {
            log.Printf("Invoice error: Type=%s, POS=%d, Message=%s", 
                invoiceErr.InvoiceType, invoiceErr.PointOfSale, invoiceErr.Message)
        }
        
    default:
        log.Printf("Unexpected error: %v", err)
    }
}
```

### 2. Retry Strategy

```go
func (s *InvoiceService) createInvoiceWithRetry(ctx context.Context, companyConfig interfaces.CompanyConfig, invoice *models.Invoice) (*models.AuthorizationResponse, error) {
    var lastErr error
    maxRetries := 3
    
    for attempt := 1; attempt <= maxRetries; attempt++ {
        response, err := s.CreateInvoice(ctx, companyConfig, invoice)
        if err == nil {
            return response, nil
        }
        
        lastErr = err
        
        // Verificar si es un error retryable
        if !isRetryableError(err) {
            return nil, err
        }
        
        // Esperar antes del siguiente intento
        if attempt < maxRetries {
            backoff := time.Duration(attempt) * time.Second
            log.Printf("Retry attempt %d/%d in %v", attempt, maxRetries, backoff)
            time.Sleep(backoff)
        }
    }
    
    return nil, fmt.Errorf("failed after %d attempts: %w", maxRetries, lastErr)
}

func isRetryableError(err error) bool {
    // Errores de red, timeouts, etc.
    return strings.Contains(err.Error(), "timeout") ||
           strings.Contains(err.Error(), "connection") ||
           strings.Contains(err.Error(), "temporary")
}
```

## Configuración Avanzada

### 1. Configuración de Timeouts

```go
func createManagerWithCustomTimeouts() interfaces.AFIPClientManager {
    factory := factory.NewClientManagerFactory()
    
    return factory.CreateManager(factory.ManagerConfig{
        ClientCacheSize:   200,                    // Más clientes en cache
        ClientIdleTimeout: 60 * time.Minute,       // Timeout más largo
        HTTPTimeout:       60 * time.Second,       // Timeout HTTP más largo
        MaxRetryAttempts:  5,                      // Más reintentos
        Logger:            &DetailedLogger{},
    })
}
```

### 2. Configuración de Cache

```go
func configureCache(manager interfaces.AFIPClientManager) {
    // Obtener estadísticas iniciales
    stats := manager.GetCacheStats()
    log.Printf("Initial cache stats: %+v", stats)
    
    // Configurar limpieza periódica
    go func() {
        ticker := time.NewTicker(10 * time.Minute)
        defer ticker.Stop()
        
        for range ticker.C {
            manager.CleanupInactiveClients(15 * time.Minute)
            
            stats := manager.GetCacheStats()
            log.Printf("Cache cleanup completed: %+v", stats)
        }
    }()
}
```

### 3. Health Check

```go
func (s *InvoiceService) HealthCheck(ctx context.Context, companyConfig interfaces.CompanyConfig) error {
    client, err := s.afipManager.GetClientForCompany(ctx, companyConfig)
    if err != nil {
        return fmt.Errorf("failed to get client: %w", err)
    }
    
    // Verificar salud del cliente
    err = client.IsHealthy(ctx)
    if err != nil {
        return fmt.Errorf("client health check failed: %w", err)
    }
    
    // Obtener información de la empresa
    companyInfo := client.GetCompanyInfo()
    log.Printf("Company health check passed: %s (%s)", 
        companyInfo.CompanyID, companyInfo.Environment)
    
    return nil
}
```

## Ejemplos Prácticos

### 1. Servicio Completo de Facturación

```go
// InvoiceService completo con todos los métodos
type InvoiceService struct {
    afipManager interfaces.AFIPClientManager
    logger      interfaces.Logger
}

func NewInvoiceService(afipManager interfaces.AFIPClientManager, logger interfaces.Logger) *InvoiceService {
    return &InvoiceService{
        afipManager: afipManager,
        logger:      logger,
    }
}

// CreateInvoice crea una factura nacional
func (s *InvoiceService) CreateInvoice(ctx context.Context, companyConfig interfaces.CompanyConfig, invoiceData *models.Invoice) (*models.AuthorizationResponse, error) {
    s.logger.Infof("Creating invoice for company %s", companyConfig.GetCompanyID())
    
    client, err := s.afipManager.GetClientForCompany(ctx, companyConfig)
    if err != nil {
        s.logger.Errorf("Failed to get client: %v", err)
        return nil, fmt.Errorf("failed to get AFIP client: %w", err)
    }
    
    response, err := client.WSFE().AuthorizeInvoice(ctx, invoiceData)
    if err != nil {
        s.logger.Errorf("Failed to authorize invoice: %v", err)
        return nil, fmt.Errorf("failed to authorize invoice: %w", err)
    }
    
    s.logger.Infof("Invoice authorized successfully: CAE=%s, Number=%d", response.CAE, response.InvoiceNumber)
    return response, nil
}

// CreateExportInvoice crea una factura de exportación
func (s *InvoiceService) CreateExportInvoice(ctx context.Context, companyConfig interfaces.CompanyConfig, exportInvoice *models.ExportInvoice) (*models.ExportAuthResponse, error) {
    s.logger.Infof("Creating export invoice for company %s", companyConfig.GetCompanyID())
    
    client, err := s.afipManager.GetClientForCompany(ctx, companyConfig)
    if err != nil {
        s.logger.Errorf("Failed to get client: %v", err)
        return nil, fmt.Errorf("failed to get AFIP client: %w", err)
    }
    
    response, err := client.WSFEX().AuthorizeExportInvoice(ctx, exportInvoice)
    if err != nil {
        s.logger.Errorf("Failed to authorize export invoice: %v", err)
        return nil, fmt.Errorf("failed to authorize export invoice: %w", err)
    }
    
    s.logger.Infof("Export invoice authorized successfully: CAE=%s, Number=%d", response.CAE, response.InvoiceNumber)
    return response, nil
}

// QueryInvoice consulta una factura
func (s *InvoiceService) QueryInvoice(ctx context.Context, companyConfig interfaces.CompanyConfig, query *models.InvoiceQuery) (*models.Invoice, error) {
    s.logger.Infof("Querying invoice for company %s", companyConfig.GetCompanyID())
    
    client, err := s.afipManager.GetClientForCompany(ctx, companyConfig)
    if err != nil {
        s.logger.Errorf("Failed to get client: %v", err)
        return nil, fmt.Errorf("failed to get AFIP client: %w", err)
    }
    
    invoice, err := client.WSFE().QueryInvoice(ctx, query)
    if err != nil {
        s.logger.Errorf("Failed to query invoice: %v", err)
        return nil, fmt.Errorf("failed to query invoice: %w", err)
    }
    
    s.logger.Infof("Invoice queried successfully: Number=%d", invoice.InvoiceNumber)
    return invoice, nil
}
```

### 2. Aplicación Web con HTTP Handler

```go
// HTTPHandler maneja requests HTTP para facturación
type HTTPHandler struct {
    invoiceService *InvoiceService
}

func NewHTTPHandler(invoiceService *InvoiceService) *HTTPHandler {
    return &HTTPHandler{
        invoiceService: invoiceService,
    }
}

// CreateInvoiceHandler maneja requests POST para crear facturas
func (h *HTTPHandler) CreateInvoiceHandler(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    
    // Parsear request
    var request struct {
        CompanyID string           `json:"company_id"`
        Invoice   *models.Invoice `json:"invoice"`
    }
    
    if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }
    
    // Obtener configuración de empresa (desde base de datos, etc.)
    companyConfig, err := h.getCompanyConfig(ctx, request.CompanyID)
    if err != nil {
        http.Error(w, "Failed to get company config", http.StatusInternalServerError)
        return
    }
    
    // Crear factura
    response, err := h.invoiceService.CreateInvoice(ctx, companyConfig, request.Invoice)
    if err != nil {
        http.Error(w, fmt.Sprintf("Failed to create invoice: %v", err), http.StatusInternalServerError)
        return
    }
    
    // Retornar respuesta
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func (h *HTTPHandler) getCompanyConfig(ctx context.Context, companyID string) (interfaces.CompanyConfig, error) {
    // Implementar lógica para obtener configuración desde base de datos
    // Este es un ejemplo simplificado
    return &CompanyConfiguration{
        CompanyID:   companyID,
        CUIT:        "20-12345678-9",
        Certificate: []byte("certificado..."),
        PrivateKey:  []byte("clave privada..."),
        Environment: "testing",
    }, nil
}
```

### 3. Worker para Procesamiento en Lote

```go
// InvoiceWorker procesa facturas en lote
type InvoiceWorker struct {
    invoiceService *InvoiceService
    queue          chan InvoiceJob
    workers        int
}

type InvoiceJob struct {
    CompanyConfig interfaces.CompanyConfig
    Invoice       *models.Invoice
    Result        chan InvoiceResult
}

type InvoiceResult struct {
    Response *models.AuthorizationResponse
    Error    error
}

func NewInvoiceWorker(invoiceService *InvoiceService, workers int) *InvoiceWorker {
    worker := &InvoiceWorker{
        invoiceService: invoiceService,
        queue:          make(chan InvoiceJob, 100),
        workers:        workers,
    }
    
    // Iniciar workers
    for i := 0; i < workers; i++ {
        go worker.processJobs()
    }
    
    return worker
}

func (w *InvoiceWorker) processJobs() {
    for job := range w.queue {
        ctx := context.Background()
        
        response, err := w.invoiceService.CreateInvoice(ctx, job.CompanyConfig, job.Invoice)
        
        job.Result <- InvoiceResult{
            Response: response,
            Error:    err,
        }
    }
}

func (w *InvoiceWorker) SubmitJob(companyConfig interfaces.CompanyConfig, invoice *models.Invoice) chan InvoiceResult {
    result := make(chan InvoiceResult, 1)
    
    w.queue <- InvoiceJob{
        CompanyConfig: companyConfig,
        Invoice:       invoice,
        Result:        result,
    }
    
    return result
}
```

## Mejores Prácticas

### 1. Gestión de Configuraciones

```go
// ConfigManager maneja configuraciones de empresas
type ConfigManager struct {
    cache map[string]interfaces.CompanyConfig
    mutex sync.RWMutex
    db    *sql.DB
}

func (cm *ConfigManager) GetCompanyConfig(ctx context.Context, companyID string) (interfaces.CompanyConfig, error) {
    // Verificar cache primero
    cm.mutex.RLock()
    if config, exists := cm.cache[companyID]; exists {
        cm.mutex.RUnlock()
        return config, nil
    }
    cm.mutex.RUnlock()
    
    // Obtener de base de datos
    config, err := cm.loadFromDatabase(ctx, companyID)
    if err != nil {
        return nil, err
    }
    
    // Guardar en cache
    cm.mutex.Lock()
    cm.cache[companyID] = config
    cm.mutex.Unlock()
    
    return config, nil
}
```

### 2. Manejo de Context

```go
func (s *InvoiceService) CreateInvoiceWithTimeout(ctx context.Context, companyConfig interfaces.CompanyConfig, invoice *models.Invoice) (*models.AuthorizationResponse, error) {
    // Crear contexto con timeout
    ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
    defer cancel()
    
    // Crear factura
    response, err := s.CreateInvoice(ctx, companyConfig, invoice)
    if err != nil {
        return nil, err
    }
    
    return response, nil
}
```

### 3. Logging Estructurado

```go
type StructuredLogger struct {
    logger *log.Logger
}

func (l *StructuredLogger) log(level, message string, fields map[string]interface{}) {
    logEntry := map[string]interface{}{
        "timestamp": time.Now().Format(time.RFC3339),
        "level":     level,
        "message":   message,
    }
    
    for key, value := range fields {
        logEntry[key] = value
    }
    
    jsonData, _ := json.Marshal(logEntry)
    l.logger.Printf("%s", jsonData)
}

func (l *StructuredLogger) Info(args ...interface{}) {
    l.log("INFO", fmt.Sprint(args...), nil)
}

func (l *StructuredLogger) Infof(format string, args ...interface{}) {
    l.log("INFO", fmt.Sprintf(format, args...), nil)
}
```

### 4. Validación de Datos

```go
func validateInvoice(invoice *models.Invoice) error {
    if invoice == nil {
        return errors.New("invoice cannot be nil")
    }
    
    if invoice.InvoiceType <= 0 {
        return errors.New("invalid invoice type")
    }
    
    if invoice.PointOfSale <= 0 {
        return errors.New("invalid point of sale")
    }
    
    if len(invoice.Items) == 0 {
        return errors.New("invoice must have at least one item")
    }
    
    for i, item := range invoice.Items {
        if item.Description == "" {
            return fmt.Errorf("item %d: description cannot be empty", i)
        }
        
        if item.Quantity <= 0 {
            return fmt.Errorf("item %d: quantity must be positive", i)
        }
        
        if item.UnitPrice <= 0 {
            return fmt.Errorf("item %d: unit price must be positive", i)
        }
    }
    
    return nil
}
```

## Troubleshooting

### 1. Errores Comunes

#### Error de Autenticación
```
Error: Authentication failed for service wsfe
```
**Solución**: Verificar que el certificado y clave privada sean válidos y correspondan al CUIT.

#### Error de Timeout
```
Error: Request timeout after 30s
```
**Solución**: Aumentar el timeout en la configuración del manager.

#### Error de Cache
```
Error: Cache is full
```
**Solución**: Aumentar `ClientCacheSize` o implementar limpieza más frecuente.

### 2. Debugging

#### Habilitar Logging Detallado
```go
type DebugLogger struct{}

func (l *DebugLogger) Debug(args ...interface{}) {
    log.Printf("[DEBUG] %v", args...)
}

// Usar en configuración
manager := factory.CreateManager(factory.ManagerConfig{
    Logger: &DebugLogger{},
    // ... otras configuraciones
})
```

#### Verificar Estado del Cache
```go
stats := manager.GetCacheStats()
log.Printf("Cache Stats: %+v", stats)
```

#### Verificar Salud del Cliente
```go
err := client.IsHealthy(ctx)
if err != nil {
    log.Printf("Client health check failed: %v", err)
}
```

### 3. Performance

#### Optimizaciones Recomendadas

1. **Cache Size**: Ajustar según el número de empresas activas
2. **Idle Timeout**: Balancear entre memoria y performance
3. **HTTP Timeout**: Ajustar según latencia de red
4. **Retry Attempts**: Configurar según tolerancia a fallos

#### Monitoreo

```go
func monitorPerformance(manager interfaces.AFIPClientManager) {
    ticker := time.NewTicker(5 * time.Minute)
    defer ticker.Stop()
    
    for range ticker.C {
        stats := manager.GetCacheStats()
        log.Printf("Performance Stats: %+v", stats)
        
        // Alertar si hay muchos clientes inactivos
        if stats.InactiveClients > stats.TotalClients/2 {
            log.Printf("Warning: High number of inactive clients")
        }
    }
}
```
