package main

import (
	"arca_invoice_lib/pkg/client"
	"arca_invoice_lib/pkg/models"
	"arca_invoice_lib/pkg/wsfe"
	"arca_invoice_lib/pkg/wsfex"
	"context"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	// Cargar certificado y clave privada desde archivos
	cert, err := os.ReadFile("certificate.crt")
	if err != nil {
		log.Fatal("Error loading certificate:", err)
	}

	privateKey, err := os.ReadFile("private.key")
	if err != nil {
		log.Fatal("Error loading private key:", err)
	}

	// Configurar cliente AFIP
	config := client.Config{
		Environment:   models.EnvironmentTesting, // Usar ambiente de testing
		CUIT:          "20-12345678-9",           // Reemplazar con tu CUIT
		Certificate:   cert,
		PrivateKey:    privateKey,
		Timeout:       30 * time.Second,
		RetryAttempts: 3,
		LogLevel:      "info",
	}

	// Crear cliente AFIP
	afipClient, err := client.NewAFIPClient(config)
	if err != nil {
		log.Fatal("Error creating AFIP client:", err)
	}

	// Crear contexto
	ctx := context.Background()

	// Probar conexión
	fmt.Println("Testing connection to AFIP...")
	if err := afipClient.TestConnection(ctx); err != nil {
		log.Fatal("Connection test failed:", err)
	}
	fmt.Println("✓ Connection successful")

	// Obtener estado del sistema
	fmt.Println("\nGetting system status...")
	status, err := afipClient.GetSystemStatus(ctx)
	if err != nil {
		log.Printf("Warning: Could not get system status: %v", err)
	} else {
		fmt.Printf("✓ System status: %s - %s\n", status.Status, status.Message)
	}

	// Obtener parámetros del sistema
	fmt.Println("\nGetting system parameters...")
	// Nota: Los servicios WSFE y WSFEX necesitan ser configurados después de crear el cliente
	// Por ahora solo mostramos un mensaje informativo
	fmt.Println("✓ Services need to be configured after client creation")

	// Obtener último comprobante autorizado
	fmt.Println("\nGetting last authorized invoice...")
	// Nota: Los servicios WSFE y WSFEX necesitan ser configurados después de crear el cliente
	fmt.Println("✓ Services need to be configured after client creation")

	// Ejemplo de creación de factura (sin enviar)
	fmt.Println("\nCreating sample invoice...")
	invoice := createSampleInvoice()
	fmt.Printf("✓ Sample invoice created for: %s\n", invoice.NameFrom)

	// Mostrar información del cache de autenticación
	fmt.Printf("\nAuth cache size: %d\n", afipClient.GetAuthCacheSize())

	fmt.Println("\nExample completed successfully!")
}

// createSampleInvoice crea una factura de ejemplo
func createSampleInvoice() *wsfe.Invoice {
	return &wsfe.Invoice{
		InvoiceBase: models.InvoiceBase{
			InvoiceType:   models.InvoiceTypeB,
			PointOfSale:   1,
			InvoiceNumber: 1,
			DateFrom:      time.Now(),
			DateTo:        time.Now(),
			ConceptType:   models.ConceptTypeProducts,
			CurrencyType:  models.CurrencyTypePES,
			CurrencyRate:  1.0,
			Amount:        1000.00,
			TaxAmount:     210.00,
			TotalAmount:   1210.00,
			Items: []models.Item{
				{
					Description: "Producto de ejemplo",
					Quantity:    1,
					UnitPrice:   1000.00,
					TotalPrice:  1000.00,
					Taxes: []models.Tax{
						{
							Type:   models.TaxTypeIVA,
							Rate:   models.TaxRate21,
							Base:   1000.00,
							Amount: 210.00,
						},
					},
				},
			},
			Notes: "Factura de ejemplo - NO ENVIAR A AFIP",
		},
		DocType:       models.DocumentTypeCUIT,
		DocNumber:     "20-12345678-9",
		DocTypeFrom:   models.DocumentTypeDNI,
		DocNumberFrom: "12345678",
		NameFrom:      "Cliente de Ejemplo",
		ServiceFrom:   "Servicio de ejemplo",
	}
}

// Ejemplo de uso con facturación de exportación
func exportExample(afipClient *client.AFIPClient, ctx context.Context) {
	fmt.Println("\n=== Export Invoice Example ===")

	// Crear factura de exportación de ejemplo
	exportInvoice := &wsfex.ExportInvoice{
		InvoiceBase: models.InvoiceBase{
			InvoiceType:   models.InvoiceTypeE,
			PointOfSale:   1,
			InvoiceNumber: 1,
			DateFrom:      time.Now(),
			DateTo:        time.Now(),
			ConceptType:   models.ConceptTypeProducts,
			CurrencyType:  models.CurrencyTypeUSD,
			CurrencyRate:  1.0,
			Amount:        100.00,
			TaxAmount:     0.00,
			TotalAmount:   100.00,
			Items: []models.Item{
				{
					Description: "Producto de exportación",
					Quantity:    1,
					UnitPrice:   100.00,
					TotalPrice:  100.00,
				},
			},
			Notes: "Factura de exportación de ejemplo - NO ENVIAR A AFIP",
		},
		DocType:       models.DocumentTypeCUIT,
		DocNumber:     "20-12345678-9",
		DocTypeFrom:   models.DocumentTypePAS,
		DocNumberFrom: "ABC123456",
		NameFrom:      "Cliente Internacional",
		CountryFrom:   "BR",
		ServiceFrom:   "Servicio de exportación",
	}

	fmt.Printf("✓ Export invoice created for: %s (%s)\n", exportInvoice.NameFrom, exportInvoice.CountryFrom)

	// Obtener parámetros de exportación
	// Nota: Los servicios WSFE y WSFEX necesitan ser configurados después de crear el cliente
	fmt.Println("✓ Export services need to be configured after client creation")
}
