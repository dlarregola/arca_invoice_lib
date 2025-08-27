package wsfe

import (
	"github.com/dlarregola/arca_invoice_lib/internal/utils"
	"github.com/dlarregola/arca_invoice_lib/pkg/client"
	"github.com/dlarregola/arca_invoice_lib/pkg/models"
	"context"
	"fmt"
)

// Service representa el servicio WSFEv1
type Service struct {
	config *client.Config
	auth   *client.WSAAAuth
	logger interface{}
}

// NewService crea un nuevo servicio WSFEv1
func NewService(config *client.Config, auth *client.WSAAAuth, logger interface{}) *Service {
	return &Service{
		config: config,
		auth:   auth,
		logger: logger,
	}
}

// AuthorizeInvoice autoriza una factura
func (s *Service) AuthorizeInvoice(ctx context.Context, invoice *Invoice) (*models.AuthorizationResult, error) {
	// Validar factura
	if err := s.validateInvoice(invoice); err != nil {
		return nil, err
	}

	// Obtener ticket de acceso
	ticket, err := s.auth.GetAccessTicket(ctx, "wsfe")
	if err != nil {
		return nil, fmt.Errorf("error getting access ticket: %w", err)
	}

	// Crear request
	request := &AuthorizationRequest{}
	request.Auth.Token = ticket.Token
	request.Auth.Sign = ticket.Sign
	request.Auth.CUIT = s.config.CUIT

	// Configurar datos de la factura
	request.Request.InvoiceType = int(invoice.InvoiceType)
	request.Request.PointOfSale = invoice.PointOfSale
	request.Request.InvoiceNumber = invoice.InvoiceNumber
	request.Request.DateFrom = invoice.DateFrom
	request.Request.DateTo = invoice.DateTo
	request.Request.ServiceFrom = invoice.ServiceFrom
	request.Request.Amount = invoice.Amount
	request.Request.TaxAmount = invoice.TaxAmount
	request.Request.TotalAmount = invoice.TotalAmount
	request.Request.CurrencyType = string(invoice.CurrencyType)
	request.Request.CurrencyRate = invoice.CurrencyRate
	request.Request.ConceptType = int(invoice.ConceptType)
	request.Request.DocType = int(invoice.DocType)
	request.Request.DocNumber = invoice.DocNumber
	request.Request.DocTypeFrom = int(invoice.DocTypeFrom)
	request.Request.DocNumberFrom = invoice.DocNumberFrom
	request.Request.NameFrom = invoice.NameFrom

	// Configurar ítems
	for _, item := range invoice.Items {
		requestItem := struct {
			Description string  `xml:"Concepto"`
			Quantity    float64 `xml:"Cantidad"`
			UnitPrice   float64 `xml:"PrecioUnit"`
			TotalPrice  float64 `xml:"Importe"`
			ProductCode string  `xml:"CodProd"`
			UnitMeasure string  `xml:"UnidadMedida"`
			Discount    float64 `xml:"Descuento"`
		}{
			Description: item.Description,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
			TotalPrice:  item.TotalPrice,
			ProductCode: item.ProductCode,
			UnitMeasure: item.UnitMeasure,
			Discount:    item.Discount,
		}
		request.Request.Items = append(request.Request.Items, requestItem)
	}

	// Realizar llamada SOAP
	var response AuthorizationResponse
	if err := s.callSOAP(ctx, "FECAESolicitar", request, &response); err != nil {
		return nil, err
	}

	// Verificar errores
	if len(response.Errors) > 0 {
		error := response.Errors[0]
		return nil, models.NewAFIPError(error.Code, error.Message)
	}

	// Crear resultado
	result := &models.AuthorizationResult{
		CAE:               response.Result.CAE,
		CAEExpirationDate: response.Result.CAEDueDate,
		InvoiceNumber:     response.Result.InvoiceNumber,
		PointOfSale:       response.Result.PointOfSale,
		InvoiceType:       models.InvoiceType(response.Result.InvoiceType),
		AuthorizationDate: response.Result.AuthorizationDate,
		Status:            response.Result.Status,
		Message:           response.Result.Message,
	}

	return result, nil
}

// GetInvoice consulta una factura específica
func (s *Service) GetInvoice(ctx context.Context, pointOfSale, invoiceType, invoiceNumber int) (*Invoice, error) {
	// Validar parámetros
	if err := utils.ValidatePointOfSale(pointOfSale); err != nil {
		return nil, err
	}
	if err := utils.ValidateInvoiceType(models.InvoiceType(invoiceType)); err != nil {
		return nil, err
	}
	if err := utils.ValidateInvoiceNumber(invoiceNumber); err != nil {
		return nil, err
	}

	// Obtener ticket de acceso
	ticket, err := s.auth.GetAccessTicket(ctx, "wsfe")
	if err != nil {
		return nil, fmt.Errorf("error getting access ticket: %w", err)
	}

	// Crear request
	request := &QueryRequest{}
	request.Auth.Token = ticket.Token
	request.Auth.Sign = ticket.Sign
	request.Auth.CUIT = s.config.CUIT
	request.Request.InvoiceType = invoiceType
	request.Request.PointOfSale = pointOfSale
	request.Request.InvoiceNumber = invoiceNumber

	// Realizar llamada SOAP
	var response QueryResponse
	if err := s.callSOAP(ctx, "FECompConsultar", request, &response); err != nil {
		return nil, err
	}

	// Verificar errores
	if len(response.Errors) > 0 {
		error := response.Errors[0]
		return nil, models.NewAFIPError(error.Code, error.Message)
	}

	// Crear factura
	invoice := &Invoice{
		InvoiceBase: models.InvoiceBase{
			InvoiceType:   models.InvoiceType(response.Result.InvoiceType),
			PointOfSale:   response.Result.PointOfSale,
			InvoiceNumber: response.Result.InvoiceNumber,
			DateFrom:      response.Result.DateFrom,
			Amount:        response.Result.Amount,
			CurrencyType:  models.CurrencyType(response.Result.CurrencyType),
			CurrencyRate:  response.Result.CurrencyRate,
		},
	}

	return invoice, nil
}

// GetLastAuthorizedInvoice obtiene el último comprobante autorizado
func (s *Service) GetLastAuthorizedInvoice(ctx context.Context, pointOfSale, invoiceType int) (*models.AuthorizationResult, error) {
	// Validar parámetros
	if err := utils.ValidatePointOfSale(pointOfSale); err != nil {
		return nil, err
	}
	if err := utils.ValidateInvoiceType(models.InvoiceType(invoiceType)); err != nil {
		return nil, err
	}

	// Obtener ticket de acceso
	ticket, err := s.auth.GetAccessTicket(ctx, "wsfe")
	if err != nil {
		return nil, fmt.Errorf("error getting access ticket: %w", err)
	}

	// Crear request
	request := &LastAuthorizedRequest{}
	request.Auth.Token = ticket.Token
	request.Auth.Sign = ticket.Sign
	request.Auth.CUIT = s.config.CUIT
	request.Request.InvoiceType = invoiceType
	request.Request.PointOfSale = pointOfSale

	// Realizar llamada SOAP
	var response LastAuthorizedResponse
	if err := s.callSOAP(ctx, "FECompUltimoAutorizado", request, &response); err != nil {
		return nil, err
	}

	// Verificar errores
	if len(response.Errors) > 0 {
		error := response.Errors[0]
		return nil, models.NewAFIPError(error.Code, error.Message)
	}

	// Crear resultado
	result := &models.AuthorizationResult{
		InvoiceNumber:     response.Result.InvoiceNumber,
		PointOfSale:       response.Result.PointOfSale,
		InvoiceType:       models.InvoiceType(response.Result.InvoiceType),
		AuthorizationDate: response.Result.DateFrom,
		Status:            "A",
	}

	return result, nil
}

// GetParameters obtiene los parámetros del sistema
func (s *Service) GetParameters(ctx context.Context) (*models.Parameters, error) {
	// Obtener ticket de acceso
	ticket, err := s.auth.GetAccessTicket(ctx, "wsfe")
	if err != nil {
		return nil, fmt.Errorf("error getting access ticket: %w", err)
	}

	// Crear request
	request := &ParametersRequest{}
	request.Auth.Token = ticket.Token
	request.Auth.Sign = ticket.Sign
	request.Auth.CUIT = s.config.CUIT

	// Realizar llamada SOAP
	var response ParametersResponse
	if err := s.callSOAP(ctx, "FEParamGetTiposConcepto", request, &response); err != nil {
		return nil, err
	}

	// Verificar errores
	if len(response.Errors) > 0 {
		error := response.Errors[0]
		return nil, models.NewAFIPError(error.Code, error.Message)
	}

	// Crear parámetros
	params := &models.Parameters{
		LastUpdate: response.LastUpdate,
	}

	// Convertir tipos de documento
	for _, dt := range response.DocumentTypes {
		params.DocumentTypes = append(params.DocumentTypes, models.DocumentTypeInfo{
			ID:          models.DocumentType(dt.ID),
			Description: dt.Description,
			Active:      dt.Active,
		})
	}

	// Convertir tipos de factura
	for _, it := range response.InvoiceTypes {
		params.InvoiceTypes = append(params.InvoiceTypes, models.InvoiceTypeInfo{
			ID:          models.InvoiceType(it.ID),
			Description: it.Description,
			Active:      it.Active,
		})
	}

	// Convertir tipos de moneda
	for _, ct := range response.CurrencyTypes {
		params.CurrencyTypes = append(params.CurrencyTypes, models.CurrencyTypeInfo{
			ID:          models.CurrencyType(ct.ID),
			Description: ct.Description,
			Active:      ct.Active,
		})
	}

	// Convertir alícuotas
	for _, tr := range response.TaxRates {
		params.TaxRates = append(params.TaxRates, models.TaxRateInfo{
			ID:          models.TaxRate(tr.ID),
			Description: tr.Description,
			Active:      tr.Active,
		})
	}

	// Convertir tipos de concepto
	for _, ct := range response.ConceptTypes {
		params.ConceptTypes = append(params.ConceptTypes, models.ConceptTypeInfo{
			ID:          models.ConceptType(ct.ID),
			Description: ct.Description,
			Active:      ct.Active,
		})
	}

	return params, nil
}

// GetCAEA obtiene un CAEA
func (s *Service) GetCAEA(ctx context.Context, period, order, fiscalYear int) (*CAEAResponse, error) {
	// Obtener ticket de acceso
	ticket, err := s.auth.GetAccessTicket(ctx, "wsfe")
	if err != nil {
		return nil, fmt.Errorf("error getting access ticket: %w", err)
	}

	// Crear request
	request := &CAEARequest{}
	request.Auth.Token = ticket.Token
	request.Auth.Sign = ticket.Sign
	request.Auth.CUIT = s.config.CUIT
	request.Request.Period = period
	request.Request.Order = order
	request.Request.FiscalYear = fiscalYear

	// Realizar llamada SOAP
	var response CAEAResponse
	if err := s.callSOAP(ctx, "FECAEASolicitar", request, &response); err != nil {
		return nil, err
	}

	// Verificar errores
	if len(response.Errors) > 0 {
		error := response.Errors[0]
		return nil, models.NewAFIPError(error.Code, error.Message)
	}

	return &response, nil
}

// validateInvoice valida una factura
func (s *Service) validateInvoice(invoice *Invoice) error {
	var errors models.ValidationErrors

	// Validar campos básicos
	if err := utils.ValidateInvoiceType(invoice.InvoiceType); err != nil {
		errors.Add("invoice_type", err.Error(), invoice.InvoiceType)
	}

	if err := utils.ValidatePointOfSale(invoice.PointOfSale); err != nil {
		errors.Add("point_of_sale", err.Error(), invoice.PointOfSale)
	}

	if err := utils.ValidateInvoiceNumber(invoice.InvoiceNumber); err != nil {
		errors.Add("invoice_number", err.Error(), invoice.InvoiceNumber)
	}

	if err := utils.ValidateDate(invoice.DateFrom, "date_from"); err != nil {
		errors.Add("date_from", err.Error(), invoice.DateFrom)
	}

	if err := utils.ValidateDate(invoice.DateTo, "date_to"); err != nil {
		errors.Add("date_to", err.Error(), invoice.DateTo)
	}

	if err := utils.ValidateConceptType(invoice.ConceptType); err != nil {
		errors.Add("concept_type", err.Error(), invoice.ConceptType)
	}

	if err := utils.ValidateCurrencyType(invoice.CurrencyType); err != nil {
		errors.Add("currency_type", err.Error(), invoice.CurrencyType)
	}

	if err := utils.ValidateAmount(invoice.Amount, "amount"); err != nil {
		errors.Add("amount", err.Error(), invoice.Amount)
	}

	if err := utils.ValidateAmount(invoice.TaxAmount, "tax_amount"); err != nil {
		errors.Add("tax_amount", err.Error(), invoice.TaxAmount)
	}

	if err := utils.ValidateAmount(invoice.TotalAmount, "total_amount"); err != nil {
		errors.Add("total_amount", err.Error(), invoice.TotalAmount)
	}

	// Validar documento
	if err := utils.ValidateDocumentType(invoice.DocType); err != nil {
		errors.Add("doc_type", err.Error(), invoice.DocType)
	}

	if err := utils.ValidateDocumentNumber(invoice.DocType, invoice.DocNumber); err != nil {
		errors.Add("doc_number", err.Error(), invoice.DocNumber)
	}

	// Validar documento del cliente
	if err := utils.ValidateDocumentType(invoice.DocTypeFrom); err != nil {
		errors.Add("doc_type_from", err.Error(), invoice.DocTypeFrom)
	}

	if err := utils.ValidateDocumentNumber(invoice.DocTypeFrom, invoice.DocNumberFrom); err != nil {
		errors.Add("doc_number_from", err.Error(), invoice.DocNumberFrom)
	}

	// Validar ítems
	if err := utils.ValidateItems(invoice.Items); err != nil {
		errors.Add("items", err.Error(), invoice.Items)
	}

	if errors.HasErrors() {
		return errors
	}

	return nil
}

// callSOAP realiza una llamada SOAP
func (s *Service) callSOAP(ctx context.Context, action string, request interface{}, response interface{}) error {
	// Esta es una implementación simplificada
	// En una implementación real, usarías el cliente SOAP interno
	return fmt.Errorf("SOAP call not implemented yet")
}
