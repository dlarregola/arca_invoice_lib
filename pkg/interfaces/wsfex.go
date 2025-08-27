package interfaces

import (
	"github.com/dlarregola/arca_invoice_lib/pkg/models"
	"context"
)

// WSFEXService es la interfaz para el servicio de facturación internacional
type WSFEXService interface {
	// AuthorizeExportInvoice autoriza un comprobante de exportación
	AuthorizeExportInvoice(ctx context.Context, invoice *models.ExportInvoice) (*models.ExportAuthResponse, error)

	// QueryExportInvoice consulta un comprobante de exportación
	QueryExportInvoice(ctx context.Context, query *models.ExportInvoiceQuery) (*models.ExportInvoice, error)

	// GetExportDestinations obtiene los destinos de exportación disponibles
	GetExportDestinations(ctx context.Context) ([]models.Destination, error)

	// GetCurrencies obtiene las monedas disponibles
	GetCurrencies(ctx context.Context) ([]models.Currency, error)

	// GetUnitTypes obtiene los tipos de unidad disponibles
	GetUnitTypes(ctx context.Context) ([]models.UnitType, error)
}
