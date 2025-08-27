package wsfe

import (
	"github.com/dlarregola/arca_invoice_lib/pkg/models"
	"time"
)

// Invoice representa una factura nacional
type Invoice struct {
	models.InvoiceBase
	DocType       models.DocumentType `json:"doc_type" xml:"doc_type"`
	DocNumber     string              `json:"doc_number" xml:"doc_number"`
	DocTypeFrom   models.DocumentType `json:"doc_type_from" xml:"doc_type_from"`
	DocNumberFrom string              `json:"doc_number_from" xml:"doc_number_from"`
	NameFrom      string              `json:"name_from,omitempty" xml:"name_from,omitempty"`
	AddressFrom   *models.Address     `json:"address_from,omitempty" xml:"address_from,omitempty"`
	ServiceFrom   string              `json:"service_from,omitempty" xml:"service_from,omitempty"`
	CAE           string              `json:"cae,omitempty" xml:"cae,omitempty"`
	CAEDueDate    time.Time           `json:"cae_due_date,omitempty" xml:"cae_due_date,omitempty"`
}

// InvoiceItem representa un ítem de factura nacional
type InvoiceItem struct {
	models.Item
	ProductCode string  `json:"product_code,omitempty" xml:"product_code,omitempty"`
	UnitMeasure string  `json:"unit_measure,omitempty" xml:"unit_measure,omitempty"`
	Discount    float64 `json:"discount,omitempty" xml:"discount,omitempty"`
}

// AuthorizationRequest representa el request de autorización
type AuthorizationRequest struct {
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
		Items         []struct {
			Description string  `xml:"Concepto"`
			Quantity    float64 `xml:"Cantidad"`
			UnitPrice   float64 `xml:"PrecioUnit"`
			TotalPrice  float64 `xml:"Importe"`
			ProductCode string  `xml:"CodProd"`
			UnitMeasure string  `xml:"UnidadMedida"`
			Discount    float64 `xml:"Descuento"`
		} `xml:"FeDetReq"`
	} `xml:"FeCAEReq"`
}

// AuthorizationResponse representa la respuesta de autorización
type AuthorizationResponse struct {
	Result struct {
		CAE               string    `xml:"CAE"`
		CAEDueDate        time.Time `xml:"CAEFchVto"`
		InvoiceNumber     int       `xml:"CbteDesde"`
		PointOfSale       int       `xml:"PuntoVta"`
		InvoiceType       int       `xml:"CbteTipo"`
		AuthorizationDate time.Time `xml:"FchProceso"`
		Status            string    `xml:"Resultado"`
		Message           string    `xml:"Observaciones"`
	} `xml:"FeCabResp"`
	Errors []struct {
		Code    string `xml:"Code"`
		Message string `xml:"Msg"`
	} `xml:"Errors"`
}

// QueryRequest representa el request de consulta
type QueryRequest struct {
	Auth struct {
		Token string `xml:"token"`
		Sign  string `xml:"sign"`
		CUIT  string `xml:"cuit"`
	} `xml:"Auth"`
	Request struct {
		InvoiceType   int `xml:"FeCompConsReq"`
		PointOfSale   int `xml:"FeCompConsReq"`
		InvoiceNumber int `xml:"FeCompConsReq"`
	} `xml:"FeCompConsReq"`
}

// QueryResponse representa la respuesta de consulta
type QueryResponse struct {
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
	} `xml:"FeCompConsResult"`
	Errors []struct {
		Code    string `xml:"Code"`
		Message string `xml:"Msg"`
	} `xml:"Errors"`
}

// LastAuthorizedRequest representa el request para obtener el último autorizado
type LastAuthorizedRequest struct {
	Auth struct {
		Token string `xml:"token"`
		Sign  string `xml:"sign"`
		CUIT  string `xml:"cuit"`
	} `xml:"Auth"`
	Request struct {
		InvoiceType int `xml:"FeCompUltimoAutorizadoReq"`
		PointOfSale int `xml:"FeCompUltimoAutorizadoReq"`
	} `xml:"FeCompUltimoAutorizadoReq"`
}

// LastAuthorizedResponse representa la respuesta del último autorizado
type LastAuthorizedResponse struct {
	Result struct {
		InvoiceType   int       `xml:"CbteTipo"`
		PointOfSale   int       `xml:"PuntoVta"`
		InvoiceNumber int       `xml:"CbteNro"`
		DateFrom      time.Time `xml:"CbteFch"`
		Amount        float64   `xml:"ImpTotal"`
		CurrencyType  string    `xml:"MonId"`
		CurrencyRate  float64   `xml:"MonCotIz"`
	} `xml:"FeCompUltimoAutorizadoResult"`
	Errors []struct {
		Code    string `xml:"Code"`
		Message string `xml:"Msg"`
	} `xml:"Errors"`
}

// ParametersRequest representa el request de parámetros
type ParametersRequest struct {
	Auth struct {
		Token string `xml:"token"`
		Sign  string `xml:"sign"`
		CUIT  string `xml:"cuit"`
	} `xml:"Auth"`
}

// ParametersResponse representa la respuesta de parámetros
type ParametersResponse struct {
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
	TaxRates []struct {
		ID          int    `xml:"Id"`
		Description string `xml:"Desc"`
		Active      bool   `xml:"FchDesde"`
	} `xml:"IvaTipo"`
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

// CAEARequest representa el request de CAEA
type CAEARequest struct {
	Auth struct {
		Token string `xml:"token"`
		Sign  string `xml:"sign"`
		CUIT  string `xml:"cuit"`
	} `xml:"Auth"`
	Request struct {
		Period     int `xml:"CAEAReq"`
		Order      int `xml:"CAEAReq"`
		FiscalYear int `xml:"CAEAReq"`
	} `xml:"CAEAReq"`
}

// CAEAResponse representa la respuesta de CAEA
type CAEAResponse struct {
	Result struct {
		CAEA       string    `xml:"CAEA"`
		Period     int       `xml:"Periodo"`
		Order      int       `xml:"Orden"`
		FiscalYear int       `xml:"FchVigDesde"`
		DueDate    time.Time `xml:"FchVigHasta"`
		MaxAmount  float64   `xml:"MaximoImporte"`
		Status     string    `xml:"Resultado"`
		Message    string    `xml:"Observaciones"`
	} `xml:"CAEAResult"`
	Errors []struct {
		Code    string `xml:"Code"`
		Message string `xml:"Msg"`
	} `xml:"Errors"`
}
