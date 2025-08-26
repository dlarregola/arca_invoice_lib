package wsfex

import (
	"arca_invoice_lib/pkg/models"
	"time"
)

// ExportInvoice representa una factura de exportación
type ExportInvoice struct {
	models.InvoiceBase
	DocType       models.DocumentType `json:"doc_type" xml:"doc_type"`
	DocNumber     string              `json:"doc_number" xml:"doc_number"`
	DocTypeFrom   models.DocumentType `json:"doc_type_from" xml:"doc_type_from"`
	DocNumberFrom string              `json:"doc_number_from" xml:"doc_number_from"`
	NameFrom      string              `json:"name_from,omitempty" xml:"name_from,omitempty"`
	AddressFrom   *models.Address     `json:"address_from,omitempty" xml:"address_from,omitempty"`
	CountryFrom   string              `json:"country_from,omitempty" xml:"country_from,omitempty"`
	ServiceFrom   string              `json:"service_from,omitempty" xml:"service_from,omitempty"`
	CAE           string              `json:"cae,omitempty" xml:"cae,omitempty"`
	CAEDueDate    time.Time           `json:"cae_due_date,omitempty" xml:"cae_due_date,omitempty"`
}

// ExportInvoiceItem representa un ítem de factura de exportación
type ExportInvoiceItem struct {
	models.Item
	ProductCode string  `json:"product_code,omitempty" xml:"product_code,omitempty"`
	UnitMeasure string  `json:"unit_measure,omitempty" xml:"unit_measure,omitempty"`
	Discount    float64 `json:"discount,omitempty" xml:"discount,omitempty"`
	Country     string  `json:"country,omitempty" xml:"country,omitempty"`
}

// ExportAuthorizationRequest representa el request de autorización de exportación
type ExportAuthorizationRequest struct {
	Auth struct {
		Token string `xml:"token"`
		Sign  string `xml:"sign"`
		CUIT  string `xml:"cuit"`
	} `xml:"Auth"`
	Request struct {
		InvoiceType   int       `xml:"FeCabReq"`
		PointOfSale   int       `xml:"FeCabReq"`
		InvoiceNumber int       `xml:"FeCabReq"`
		DateFrom      time.Time `xml:"FeCabReq"`
		DateTo        time.Time `xml:"FeCabReq"`
		ServiceFrom   string    `xml:"FeCabReq"`
		Amount        float64   `xml:"FeCabReq"`
		TaxAmount     float64   `xml:"FeCabReq"`
		TotalAmount   float64   `xml:"FeCabReq"`
		CurrencyType  string    `xml:"FeCabReq"`
		CurrencyRate  float64   `xml:"FeCabReq"`
		ConceptType   int       `xml:"FeCabReq"`
		DocType       int       `xml:"FeDetReq"`
		DocNumber     string    `xml:"FeDetReq"`
		DocTypeFrom   int       `xml:"FeDetReq"`
		DocNumberFrom string    `xml:"FeDetReq"`
		NameFrom      string    `xml:"FeDetReq"`
		CountryFrom   string    `xml:"FeDetReq"`
		Items         []struct {
			Description string  `xml:"Concepto"`
			Quantity    float64 `xml:"Cantidad"`
			UnitPrice   float64 `xml:"PrecioUnit"`
			TotalPrice  float64 `xml:"Importe"`
			ProductCode string  `xml:"CodProd"`
			UnitMeasure string  `xml:"UnidadMedida"`
			Discount    float64 `xml:"Descuento"`
			Country     string  `xml:"PaisDestino"`
		} `xml:"FeDetReq"`
	} `xml:"FEXAuthorize"`
}

// ExportAuthorizationResponse representa la respuesta de autorización de exportación
type ExportAuthorizationResponse struct {
	Result struct {
		CAE               string    `xml:"CAE"`
		CAEDueDate        time.Time `xml:"CAEFchVto"`
		InvoiceNumber     int       `xml:"CbteDesde"`
		PointOfSale       int       `xml:"PuntoVta"`
		InvoiceType       int       `xml:"CbteTipo"`
		AuthorizationDate time.Time `xml:"FchProceso"`
		Status            string    `xml:"Resultado"`
		Message           string    `xml:"Observaciones"`
	} `xml:"FEXResultAuth"`
	Errors []struct {
		Code    string `xml:"Code"`
		Message string `xml:"Msg"`
	} `xml:"Errors"`
}

// ExportQueryRequest representa el request de consulta de exportación
type ExportQueryRequest struct {
	Auth struct {
		Token string `xml:"token"`
		Sign  string `xml:"sign"`
		CUIT  string `xml:"cuit"`
	} `xml:"Auth"`
	Request struct {
		InvoiceType   int `xml:"FEXGetCMP"`
		PointOfSale   int `xml:"FEXGetCMP"`
		InvoiceNumber int `xml:"FEXGetCMP"`
	} `xml:"FEXGetCMP"`
}

// ExportQueryResponse representa la respuesta de consulta de exportación
type ExportQueryResponse struct {
	Result struct {
		InvoiceType   int       `xml:"CbteTipo"`
		PointOfSale   int       `xml:"PuntoVta"`
		InvoiceNumber int       `xml:"CbteNro"`
		DateFrom      time.Time `xml:"CbteFch"`
		Amount        float64   `xml:"ImpTotal"`
		CurrencyType  string    `xml:"MonId"`
		CurrencyRate  float64   `xml:"MonCotIz"`
		Status        string    `xml:"Resultado"`
		Message       string    `xml:"Observaciones"`
	} `xml:"FEXResultGet"`
	Errors []struct {
		Code    string `xml:"Code"`
		Message string `xml:"Msg"`
	} `xml:"Errors"`
}

// ExportLastAuthorizedRequest representa el request para obtener el último autorizado de exportación
type ExportLastAuthorizedRequest struct {
	Auth struct {
		Token string `xml:"token"`
		Sign  string `xml:"sign"`
		CUIT  string `xml:"cuit"`
	} `xml:"Auth"`
	Request struct {
		InvoiceType int `xml:"FEXGetLast_CMP"`
		PointOfSale int `xml:"FEXGetLast_CMP"`
	} `xml:"FEXGetLast_CMP"`
}

// ExportLastAuthorizedResponse representa la respuesta del último autorizado de exportación
type ExportLastAuthorizedResponse struct {
	Result struct {
		InvoiceType   int       `xml:"CbteTipo"`
		PointOfSale   int       `xml:"PuntoVta"`
		InvoiceNumber int       `xml:"CbteNro"`
		DateFrom      time.Time `xml:"CbteFch"`
		Amount        float64   `xml:"ImpTotal"`
		CurrencyType  string    `xml:"MonId"`
		CurrencyRate  float64   `xml:"MonCotIz"`
	} `xml:"FEXResultLast_CMP"`
	Errors []struct {
		Code    string `xml:"Code"`
		Message string `xml:"Msg"`
	} `xml:"Errors"`
}

// ExportParametersRequest representa el request de parámetros de exportación
type ExportParametersRequest struct {
	Auth struct {
		Token string `xml:"token"`
		Sign  string `xml:"sign"`
		CUIT  string `xml:"cuit"`
	} `xml:"Auth"`
}

// ExportParametersResponse representa la respuesta de parámetros de exportación
type ExportParametersResponse struct {
	DocumentTypes []struct {
		ID          int    `xml:"Id"`
		Description string `xml:"Desc"`
		Active      bool   `xml:"FchDesde"`
	} `xml:"DocTipo"`
	InvoiceTypes []struct {
		ID          int    `xml:"Id"`
		Description string `xml:"Desc"`
		Active      bool   `xml:"FchDesde"`
	} `xml:"CbteTipo"`
	CurrencyTypes []struct {
		ID          string `xml:"Id"`
		Description string `xml:"Desc"`
		Active      bool   `xml:"FchDesde"`
	} `xml:"MonId"`
	Countries []struct {
		ID          string `xml:"Id"`
		Description string `xml:"Desc"`
		Active      bool   `xml:"FchDesde"`
	} `xml:"Pais"`
	ConceptTypes []struct {
		ID          int    `xml:"Id"`
		Description string `xml:"Desc"`
		Active      bool   `xml:"FchDesde"`
	} `xml:"ConceptoTipo"`
	LastUpdate time.Time `xml:"FchServDesde"`
	Errors     []struct {
		Code    string `xml:"Code"`
		Message string `xml:"Msg"`
	} `xml:"Errors"`
}

// ExportCAEARequest representa el request de CAEA para exportación
type ExportCAEARequest struct {
	Auth struct {
		Token string `xml:"token"`
		Sign  string `xml:"sign"`
		CUIT  string `xml:"cuit"`
	} `xml:"Auth"`
	Request struct {
		Period     int `xml:"FEXGetCAEA"`
		Order      int `xml:"FEXGetCAEA"`
		FiscalYear int `xml:"FEXGetCAEA"`
	} `xml:"FEXGetCAEA"`
}

// ExportCAEAResponse representa la respuesta de CAEA para exportación
type ExportCAEAResponse struct {
	Result struct {
		CAEA       string    `xml:"CAEA"`
		Period     int       `xml:"Periodo"`
		Order      int       `xml:"Orden"`
		FiscalYear int       `xml:"FchVigDesde"`
		DueDate    time.Time `xml:"FchVigHasta"`
		MaxAmount  float64   `xml:"MaximoImporte"`
		Status     string    `xml:"Resultado"`
		Message    string    `xml:"Observaciones"`
	} `xml:"FEXResultGetCAEA"`
	Errors []struct {
		Code    string `xml:"Code"`
		Message string `xml:"Msg"`
	} `xml:"Errors"`
}
