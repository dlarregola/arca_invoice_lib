package interfaces

import (
	"github.com/dlarregola/arca_invoice_lib/pkg/models"
	"context"
)

// WSFEService es la interfaz para el servicio de facturación nacional
type WSFEService interface {
	// AuthorizeInvoice autoriza un comprobante
	AuthorizeInvoice(ctx context.Context, invoice *models.Invoice) (*models.AuthorizationResponse, error)

	// QueryInvoice consulta un comprobante
	QueryInvoice(ctx context.Context, query *models.InvoiceQuery) (*models.Invoice, error)

	// GetLastAuthorizedInvoice obtiene el último comprobante autorizado
	GetLastAuthorizedInvoice(ctx context.Context, pointOfSale int, invoiceType int) (*models.LastInvoiceResponse, error)

	// QueryCAEA consulta un CAEA
	QueryCAEA(ctx context.Context, caea string) (*models.CAEAResponse, error)

	// GetDocumentTypes obtiene los tipos de documento disponibles
	GetDocumentTypes(ctx context.Context) ([]models.DocumentType, error)

	// GetCurrencies obtiene las monedas disponibles
	GetCurrencies(ctx context.Context) ([]models.Currency, error)

	// GetConceptTypes obtiene los tipos de concepto disponibles
	GetConceptTypes(ctx context.Context) ([]models.ConceptType, error)

	// GetInvoiceTypes obtiene los tipos de comprobante disponibles
	GetInvoiceTypes(ctx context.Context) ([]models.InvoiceType, error)
}
