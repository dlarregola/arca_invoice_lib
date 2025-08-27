package wsfex

import (
	"github.com/dlarregola/arca_invoice_lib/pkg/interfaces"
	"github.com/dlarregola/arca_invoice_lib/pkg/models"
	"context"
	"fmt"
	"time"
)

// wsfexService es la implementación privada del servicio WSFEX
type wsfexService struct {
	authService interfaces.AuthService
	logger      interfaces.Logger
}

// newWSFEXService crea un nuevo servicio WSFEX
func newWSFEXService(authService interfaces.AuthService, logger interfaces.Logger) (interfaces.WSFEXService, error) {
	return &wsfexService{
		authService: authService,
		logger:      logger,
	}, nil
}

// AuthorizeExportInvoice autoriza un comprobante de exportación
func (s *wsfexService) AuthorizeExportInvoice(ctx context.Context, invoice *models.ExportInvoice) (*models.ExportAuthResponse, error) {
	// Validar factura de exportación
	if err := s.validateExportInvoice(invoice); err != nil {
		return nil, fmt.Errorf("invalid export invoice: %w", err)
	}

	// Obtener token de autenticación
	_, err := s.authService.GetToken(ctx, "wsfex")
	if err != nil {
		return nil, fmt.Errorf("auth failed: %w", err)
	}

	// TODO: Implementar llamada SOAP real
	s.logger.Infof("Authorizing export invoice %d for point of sale %d", invoice.InvoiceNumber, invoice.PointOfSale)

	return &models.ExportAuthResponse{
		AuthorizationResponse: models.AuthorizationResponse{
			CAE:               "12345678901234",
			CAEExpirationDate: invoice.DateTo.AddDate(0, 1, 0),
			InvoiceNumber:     invoice.InvoiceNumber,
			PointOfSale:       invoice.PointOfSale,
			InvoiceType:       invoice.InvoiceType,
			AuthorizationDate: invoice.DateFrom,
			Status:            "A",
			Message:           "Autorizado",
		},
		ExportType: invoice.ExportType,
	}, nil
}

// QueryExportInvoice consulta un comprobante de exportación
func (s *wsfexService) QueryExportInvoice(ctx context.Context, query *models.ExportInvoiceQuery) (*models.ExportInvoice, error) {
	// Obtener token de autenticación
	_, err := s.authService.GetToken(ctx, "wsfex")
	if err != nil {
		return nil, fmt.Errorf("auth failed: %w", err)
	}

	// TODO: Implementar consulta SOAP real
	s.logger.Infof("Querying export invoice %d for point of sale %d", query.InvoiceNumber, query.PointOfSale)

	// Retornar factura de exportación simulada
	return &models.ExportInvoice{
		InvoiceBase: models.InvoiceBase{
			InvoiceType:   query.InvoiceType,
			PointOfSale:   query.PointOfSale,
			InvoiceNumber: query.InvoiceNumber,
			DateFrom:      query.DateFrom,
			DateTo:        query.DateTo,
		},
		Destination:     "Estados Unidos",
		DestinationCode: "US",
		ExportDate:      time.Now(),
		ExportType:      "Definitiva",
	}, nil
}

// GetExportDestinations obtiene los destinos de exportación disponibles
func (s *wsfexService) GetExportDestinations(ctx context.Context) ([]models.Destination, error) {
	// Obtener token de autenticación
	_, err := s.authService.GetToken(ctx, "wsfex")
	if err != nil {
		return nil, fmt.Errorf("auth failed: %w", err)
	}

	// TODO: Implementar consulta SOAP real
	s.logger.Info("Getting export destinations")

	return []models.Destination{
		{ID: "US", Description: "Estados Unidos", Active: true},
		{ID: "BR", Description: "Brasil", Active: true},
		{ID: "CL", Description: "Chile", Active: true},
		{ID: "UY", Description: "Uruguay", Active: true},
	}, nil
}

// GetCurrencies obtiene las monedas disponibles
func (s *wsfexService) GetCurrencies(ctx context.Context) ([]models.Currency, error) {
	// Obtener token de autenticación
	_, err := s.authService.GetToken(ctx, "wsfex")
	if err != nil {
		return nil, fmt.Errorf("auth failed: %w", err)
	}

	// TODO: Implementar consulta SOAP real
	s.logger.Info("Getting currencies")

	return []models.Currency{
		{ID: "USD", Description: "Dólar Estadounidense", Active: true},
		{ID: "EUR", Description: "Euro", Active: true},
		{ID: "BRL", Description: "Real Brasileño", Active: true},
	}, nil
}

// GetUnitTypes obtiene los tipos de unidad disponibles
func (s *wsfexService) GetUnitTypes(ctx context.Context) ([]models.UnitType, error) {
	// Obtener token de autenticación
	_, err := s.authService.GetToken(ctx, "wsfex")
	if err != nil {
		return nil, fmt.Errorf("auth failed: %w", err)
	}

	// TODO: Implementar consulta SOAP real
	s.logger.Info("Getting unit types")

	return []models.UnitType{
		{ID: "UN", Description: "Unidad", Active: true},
		{ID: "KG", Description: "Kilogramo", Active: true},
		{ID: "M", Description: "Metro", Active: true},
		{ID: "L", Description: "Litro", Active: true},
	}, nil
}

// validateExportInvoice valida los datos de una factura de exportación
func (s *wsfexService) validateExportInvoice(invoice *models.ExportInvoice) error {
	if invoice == nil {
		return fmt.Errorf("export invoice cannot be nil")
	}

	if invoice.InvoiceNumber <= 0 {
		return fmt.Errorf("invoice number must be greater than 0")
	}

	if invoice.PointOfSale <= 0 {
		return fmt.Errorf("point of sale must be greater than 0")
	}

	if invoice.Amount <= 0 {
		return fmt.Errorf("amount must be greater than 0")
	}

	if len(invoice.Items) == 0 {
		return fmt.Errorf("export invoice must have at least one item")
	}

	if invoice.Destination == "" {
		return fmt.Errorf("destination cannot be empty")
	}

	if invoice.ExportType == "" {
		return fmt.Errorf("export type cannot be empty")
	}

	return nil
}

// NewWSFEXService crea un nuevo servicio WSFEX
func NewWSFEXService(authService interfaces.AuthService, logger interfaces.Logger) (interfaces.WSFEXService, error) {
	return newWSFEXService(authService, logger)
}
