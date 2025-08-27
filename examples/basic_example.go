package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dlarregola/arca_invoice_lib/internal/client"
	"github.com/dlarregola/arca_invoice_lib/pkg/factory"
	"github.com/dlarregola/arca_invoice_lib/pkg/interfaces"
	"github.com/dlarregola/arca_invoice_lib/pkg/models"
)

// CompanyConfiguration implementa la interfaz CompanyConfig
type CompanyConfiguration struct {
	CompanyID   string
	CUIT        string
	Certificate []byte
	PrivateKey  []byte
	Environment string
}

// Implementar métodos de la interfaz
func (c *CompanyConfiguration) GetCUIT() string                    { return c.CUIT }
func (c *CompanyConfiguration) GetCertificate() []byte             { return c.Certificate }
func (c *CompanyConfiguration) GetPrivateKey() []byte              { return c.PrivateKey }
func (c *CompanyConfiguration) GetEnvironment() string             { return c.Environment }
func (c *CompanyConfiguration) GetCompanyID() string               { return c.CompanyID }

// MyLogger implementa la interfaz Logger
type MyLogger struct{}

func (l *MyLogger) Debug(args ...interface{})                 { log.Printf("[DEBUG] %v", args...) }
func (l *MyLogger) Debugf(format string, args ...interface{}) { log.Printf("[DEBUG] "+format, args...) }
func (l *MyLogger) Info(args ...interface{})                  { log.Printf("[INFO] %v", args...) }
func (l *MyLogger) Infof(format string, args ...interface{})  { log.Printf("[INFO] "+format, args...) }
func (l *MyLogger) Warn(args ...interface{})                  { log.Printf("[WARN] %v", args...) }
func (l *MyLogger) Warnf(format string, args ...interface{})  { log.Printf("[WARN] "+format, args...) }
func (l *MyLogger) Error(args ...interface{})                 { log.Printf("[ERROR] %v", args...) }
func (l *MyLogger) Errorf(format string, args ...interface{}) { log.Printf("[ERROR] "+format, args...) }

func main() {
	// 1. Cargar certificados (en un caso real, estos vendrían de archivos o base de datos)
	certData, keyData, err := loadCertificates("cert.crt", "key.key")
	if err != nil {
		log.Fatal("Failed to load certificates:", err)
	}

	// 2. Crear factory
	factory := factory.NewClientManagerFactory()

	// 3. Configurar manager
	manager := factory.CreateManager(client.ManagerConfig{
		ClientCacheSize:   100,                    // Máximo 100 clientes en cache
		ClientIdleTimeout: 30 * time.Minute,       // Timeout de inactividad
		HTTPTimeout:       30 * time.Second,       // Timeout HTTP
		MaxRetryAttempts:  3,                      // Reintentos
		Logger:            &MyLogger{},            // Logger personalizado
	})

	// 4. Crear configuración de empresa
	companyConfig := &CompanyConfiguration{
		CompanyID:   "empresa-001",
		CUIT:        "20-12345678-9",
		Certificate: certData,
		PrivateKey:  keyData,
		Environment: "testing", // "testing" o "production"
	}

	// 5. Crear contexto
	ctx := context.Background()

	// 6. Obtener cliente
	client, err := manager.GetClientForCompany(ctx, companyConfig)
	if err != nil {
		log.Fatal("Failed to get client:", err)
	}

	// 7. Verificar salud del cliente
	if err := client.IsHealthy(ctx); err != nil {
		log.Printf("Client health check failed: %v", err)
	} else {
		log.Println("Client is healthy")
	}

	// 8. Obtener información de la empresa
	companyInfo := client.GetCompanyInfo()
	log.Printf("Company Info: %s (%s) - %s", companyInfo.CompanyID, companyInfo.CUIT, companyInfo.Environment)

	// 9. Ejemplo: Crear una factura nacional
	if err := createNationalInvoice(ctx, client); err != nil {
		log.Printf("Failed to create national invoice: %v", err)
	}

	// 10. Ejemplo: Crear una factura de exportación
	if err := createExportInvoice(ctx, client); err != nil {
		log.Printf("Failed to create export invoice: %v", err)
	}

	// 11. Ejemplo: Consultar factura
	if err := queryInvoice(ctx, client); err != nil {
		log.Printf("Failed to query invoice: %v", err)
	}

	// 12. Ejemplo: Obtener parámetros del sistema
	if err := getSystemParameters(ctx, client); err != nil {
		log.Printf("Failed to get system parameters: %v", err)
	}

	// 13. Mostrar estadísticas del cache
	stats := manager.GetCacheStats()
	log.Printf("Cache Stats: Total=%d, Active=%d, Inactive=%d", 
		stats.TotalClients, stats.ActiveClients, stats.InactiveClients)

	// 14. Cerrar cliente
	if err := client.Close(); err != nil {
		log.Printf("Failed to close client: %v", err)
	}
}

// loadCertificates carga certificados desde archivos
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

// createNationalInvoice crea una factura nacional de ejemplo
func createNationalInvoice(ctx context.Context, client interfaces.AFIPClient) error {
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

// createExportInvoice crea una factura de exportación de ejemplo
func createExportInvoice(ctx context.Context, client interfaces.AFIPClient) error {
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

// queryInvoice consulta una factura de ejemplo
func queryInvoice(ctx context.Context, client interfaces.AFIPClient) error {
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

// getSystemParameters obtiene parámetros del sistema
func getSystemParameters(ctx context.Context, client interfaces.AFIPClient) error {
	// Obtener tipos de documento
	docTypes, err := client.WSFE().GetDocumentTypes(ctx)
	if err != nil {
		return fmt.Errorf("failed to get document types: %w", err)
	}

	// Obtener tipos de factura
	invoiceTypes, err := client.WSFE().GetInvoiceTypes(ctx)
	if err != nil {
		return fmt.Errorf("failed to get invoice types: %w", err)
	}

	log.Printf("System parameters: DocumentTypes=%d, InvoiceTypes=%d", 
		len(docTypes), len(invoiceTypes))
	return nil
}
