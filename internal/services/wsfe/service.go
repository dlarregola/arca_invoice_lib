package wsfe

import (
	"github.com/dlarregola/arca_invoice_lib/pkg/interfaces"
	"github.com/dlarregola/arca_invoice_lib/pkg/models"
	"context"
	"fmt"
	"time"
)

// wsfeService es la implementación privada del servicio WSFE
type wsfeService struct {
	authService interfaces.AuthService
	logger      interfaces.Logger
}

// newWSFEService crea un nuevo servicio WSFE
func newWSFEService(authService interfaces.AuthService, logger interfaces.Logger) (interfaces.WSFEService, error) {
	return &wsfeService{
		authService: authService,
		logger:      logger,
	}, nil
}

// AuthorizeInvoice autoriza un comprobante
func (s *wsfeService) AuthorizeInvoice(ctx context.Context, invoice *models.Invoice) (*models.AuthorizationResponse, error) {
	// Validar factura
	if err := s.validateInvoice(invoice); err != nil {
		return nil, fmt.Errorf("invalid invoice: %w", err)
	}

	// Obtener token de autenticación
	_, err := s.authService.GetToken(ctx, "wsfe")
	if err != nil {
		return nil, fmt.Errorf("auth failed: %w", err)
	}

	// TODO: Implementar llamada SOAP real
	// Por ahora retornamos una respuesta simulada
	s.logger.Infof("Authorizing invoice %d for point of sale %d", invoice.InvoiceNumber, invoice.PointOfSale)

	return &models.AuthorizationResponse{
		CAE:               "12345678901234",
		CAEExpirationDate: invoice.DateTo.AddDate(0, 1, 0), // 1 mes después
		InvoiceNumber:     invoice.InvoiceNumber,
		PointOfSale:       invoice.PointOfSale,
		InvoiceType:       invoice.InvoiceType,
		AuthorizationDate: invoice.DateFrom,
		Status:            "A",
		Message:           "Autorizado",
	}, nil
}

// QueryInvoice consulta un comprobante
func (s *wsfeService) QueryInvoice(ctx context.Context, query *models.InvoiceQuery) (*models.Invoice, error) {
	// Obtener token de autenticación
	_, err := s.authService.GetToken(ctx, "wsfe")
	if err != nil {
		return nil, fmt.Errorf("auth failed: %w", err)
	}

	// TODO: Implementar consulta SOAP real
	s.logger.Infof("Querying invoice %d for point of sale %d", query.InvoiceNumber, query.PointOfSale)

	// Retornar factura simulada
	return &models.Invoice{
		InvoiceBase: models.InvoiceBase{
			InvoiceType:   query.InvoiceType,
			PointOfSale:   query.PointOfSale,
			InvoiceNumber: query.InvoiceNumber,
			DateFrom:      query.DateFrom,
			DateTo:        query.DateTo,
		},
	}, nil
}

// GetLastAuthorizedInvoice obtiene el último comprobante autorizado
func (s *wsfeService) GetLastAuthorizedInvoice(ctx context.Context, pointOfSale int, invoiceType int) (*models.LastInvoiceResponse, error) {
	// Obtener token de autenticación
	_, err := s.authService.GetToken(ctx, "wsfe")
	if err != nil {
		return nil, fmt.Errorf("auth failed: %w", err)
	}

	// TODO: Implementar consulta SOAP real
	s.logger.Infof("Getting last authorized invoice for point of sale %d, type %d", pointOfSale, invoiceType)

	return &models.LastInvoiceResponse{
		InvoiceType:   models.InvoiceType(invoiceType),
		PointOfSale:   pointOfSale,
		InvoiceNumber: 1000,
		Date:          time.Now(),
	}, nil
}

// QueryCAEA consulta un CAEA
func (s *wsfeService) QueryCAEA(ctx context.Context, caea string) (*models.CAEAResponse, error) {
	// Obtener token de autenticación
	_, err := s.authService.GetToken(ctx, "wsfe")
	if err != nil {
		return nil, fmt.Errorf("auth failed: %w", err)
	}

	// TODO: Implementar consulta SOAP real
	s.logger.Infof("Querying CAEA: %s", caea)

	return &models.CAEAResponse{
		CAEA:           caea,
		ExpirationDate: time.Now().AddDate(0, 1, 0),
		Status:         "A",
		Message:        "CAEA válido",
	}, nil
}

// GetDocumentTypes obtiene los tipos de documento disponibles
func (s *wsfeService) GetDocumentTypes(ctx context.Context) ([]models.DocumentType, error) {
	// Obtener token de autenticación
	_, err := s.authService.GetToken(ctx, "wsfe")
	if err != nil {
		return nil, fmt.Errorf("auth failed: %w", err)
	}

	// TODO: Implementar consulta SOAP real
	s.logger.Info("Getting document types")

	return []models.DocumentType{
		models.DocumentTypeDNI,
		models.DocumentTypeCUIT,
		models.DocumentTypeCUIL,
	}, nil
}

// GetCurrencies obtiene las monedas disponibles
func (s *wsfeService) GetCurrencies(ctx context.Context) ([]models.Currency, error) {
	// Obtener token de autenticación
	_, err := s.authService.GetToken(ctx, "wsfe")
	if err != nil {
		return nil, fmt.Errorf("auth failed: %w", err)
	}

	// TODO: Implementar consulta SOAP real
	s.logger.Info("Getting currencies")

	return []models.Currency{
		{ID: "PES", Description: "Peso Argentino", Active: true},
		{ID: "USD", Description: "Dólar Estadounidense", Active: true},
		{ID: "EUR", Description: "Euro", Active: true},
	}, nil
}

// GetConceptTypes obtiene los tipos de concepto disponibles
func (s *wsfeService) GetConceptTypes(ctx context.Context) ([]models.ConceptType, error) {
	// Obtener token de autenticación
	_, err := s.authService.GetToken(ctx, "wsfe")
	if err != nil {
		return nil, fmt.Errorf("auth failed: %w", err)
	}

	// TODO: Implementar consulta SOAP real
	s.logger.Info("Getting concept types")

	return []models.ConceptType{
		models.ConceptTypeProducts,
		models.ConceptTypeServices,
		models.ConceptTypeMixed,
	}, nil
}

// GetInvoiceTypes obtiene los tipos de comprobante disponibles
func (s *wsfeService) GetInvoiceTypes(ctx context.Context) ([]models.InvoiceType, error) {
	// Obtener token de autenticación
	_, err := s.authService.GetToken(ctx, "wsfe")
	if err != nil {
		return nil, fmt.Errorf("auth failed: %w", err)
	}

	// TODO: Implementar consulta SOAP real
	s.logger.Info("Getting invoice types")

	return []models.InvoiceType{
		models.InvoiceTypeA,
		models.InvoiceTypeB,
		models.InvoiceTypeC,
		models.InvoiceTypeE,
	}, nil
}

// validateInvoice valida los datos de una factura
func (s *wsfeService) validateInvoice(invoice *models.Invoice) error {
	if invoice == nil {
		return fmt.Errorf("invoice cannot be nil")
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
		return fmt.Errorf("invoice must have at least one item")
	}

	return nil
}

// NewWSFEService crea un nuevo servicio WSFE
func NewWSFEService(authService interfaces.AuthService, logger interfaces.Logger) (interfaces.WSFEService, error) {
	return newWSFEService(authService, logger)
}
