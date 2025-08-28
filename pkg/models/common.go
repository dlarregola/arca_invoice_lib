package models

import (
	"time"
)

// Environment representa el ambiente de ARCA
type Environment string

const (
	EnvironmentTesting    Environment = "testing"
	EnvironmentProduction Environment = "production"
)

// DocumentType representa los tipos de documento
type DocumentType int

const (
	DocumentTypeDNI  DocumentType = 1
	DocumentTypeCUIT DocumentType = 11
	DocumentTypeCUIL DocumentType = 12
	DocumentTypeCDI  DocumentType = 13
	DocumentTypeLE   DocumentType = 14
	DocumentTypeLC   DocumentType = 15
	DocumentTypeCI   DocumentType = 16
	DocumentTypePAS  DocumentType = 17
	DocumentTypeDE   DocumentType = 18
	DocumentTypeDI   DocumentType = 19
)

// ConceptType representa los tipos de concepto
type ConceptType int

const (
	ConceptTypeProducts ConceptType = 1
	ConceptTypeServices ConceptType = 2
	ConceptTypeMixed    ConceptType = 3
)

// InvoiceType representa los tipos de comprobante
type InvoiceType int

const (
	InvoiceTypeA InvoiceType = 1
	InvoiceTypeB InvoiceType = 6
	InvoiceTypeC InvoiceType = 11
	InvoiceTypeE InvoiceType = 19
	InvoiceTypeM InvoiceType = 51
	InvoiceTypeT InvoiceType = 60
	InvoiceTypeR InvoiceType = 63
)

// CurrencyType representa los tipos de moneda
type CurrencyType string

const (
	CurrencyTypePES CurrencyType = "PES" // Peso Argentino
	CurrencyTypeUSD CurrencyType = "USD" // Dólar Estadounidense
	CurrencyTypeEUR CurrencyType = "EUR" // Euro
	CurrencyTypeBRL CurrencyType = "BRL" // Real Brasileño
)

// TaxType representa los tipos de impuesto
type TaxType int

const (
	TaxTypeIVA TaxType = 1
	TaxTypeII  TaxType = 2
	TaxTypeIO  TaxType = 3
)

// TaxRate representa las alícuotas de IVA
type TaxRate int

const (
	TaxRate0      TaxRate = 0
	TaxRate105    TaxRate = 105
	TaxRate21     TaxRate = 21
	TaxRate27     TaxRate = 27
	TaxRate25     TaxRate = 25
	TaxRate5      TaxRate = 5
	TaxRateExempt TaxRate = -1
)

// BaseEntity representa una entidad base con campos comunes
type BaseEntity struct {
	ID        string    `json:"id,omitempty" xml:"id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty" xml:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty" xml:"updated_at,omitempty"`
}

// Address representa una dirección
type Address struct {
	Street     string `json:"street,omitempty" xml:"street,omitempty"`
	Number     string `json:"number,omitempty" xml:"number,omitempty"`
	Floor      string `json:"floor,omitempty" xml:"floor,omitempty"`
	Apartment  string `json:"apartment,omitempty" xml:"apartment,omitempty"`
	PostalCode string `json:"postal_code,omitempty" xml:"postal_code,omitempty"`
	City       string `json:"city,omitempty" xml:"city,omitempty"`
	State      string `json:"state,omitempty" xml:"state,omitempty"`
	Country    string `json:"country,omitempty" xml:"country,omitempty"`
}

// Person representa una persona física o jurídica
type Person struct {
	DocType   DocumentType `json:"doc_type" xml:"doc_type"`
	DocNumber string       `json:"doc_number" xml:"doc_number"`
	Name      string       `json:"name,omitempty" xml:"name,omitempty"`
	Address   *Address     `json:"address,omitempty" xml:"address,omitempty"`
}

// Tax representa un impuesto
type Tax struct {
	Type   TaxType `json:"type" xml:"type"`
	Rate   TaxRate `json:"rate" xml:"rate"`
	Base   float64 `json:"base" xml:"base"`
	Amount float64 `json:"amount" xml:"amount"`
}

// Item representa un ítem de factura
type Item struct {
	Description string  `json:"description" xml:"description"`
	Quantity    float64 `json:"quantity" xml:"quantity"`
	UnitPrice   float64 `json:"unit_price" xml:"unit_price"`
	TotalPrice  float64 `json:"total_price" xml:"total_price"`
	ProductCode string  `json:"product_code,omitempty" xml:"product_code,omitempty"`
	UnitMeasure string  `json:"unit_measure,omitempty" xml:"unit_measure,omitempty"`
	Discount    float64 `json:"discount,omitempty" xml:"discount,omitempty"`
	Country     string  `json:"country,omitempty" xml:"country,omitempty"`
	Taxes       []Tax   `json:"taxes,omitempty" xml:"taxes,omitempty"`
}

// InvoiceBase representa los campos base de una factura
type InvoiceBase struct {
	BaseEntity
	InvoiceType   InvoiceType  `json:"invoice_type" xml:"invoice_type"`
	PointOfSale   int          `json:"point_of_sale" xml:"point_of_sale"`
	InvoiceNumber int          `json:"invoice_number,omitempty" xml:"invoice_number,omitempty"`
	DateFrom      time.Time    `json:"date_from" xml:"date_from"`
	DateTo        time.Time    `json:"date_to" xml:"date_to"`
	ConceptType   ConceptType  `json:"concept_type" xml:"concept_type"`
	CurrencyType  CurrencyType `json:"currency_type" xml:"currency_type"`
	CurrencyRate  float64      `json:"currency_rate,omitempty" xml:"currency_rate,omitempty"`
	Amount        float64      `json:"amount" xml:"amount"`
	TaxAmount     float64      `json:"tax_amount" xml:"tax_amount"`
	TotalAmount   float64      `json:"total_amount" xml:"total_amount"`
	Items         []Item       `json:"items" xml:"items"`
	Taxes         []Tax        `json:"taxes,omitempty" xml:"taxes,omitempty"`
	Notes         string       `json:"notes,omitempty" xml:"notes,omitempty"`
}

// AuthorizationResult representa el resultado de una autorización
type AuthorizationResult struct {
	CAE               string      `json:"cae" xml:"cae"`
	CAEExpirationDate time.Time   `json:"cae_expiration_date" xml:"cae_expiration_date"`
	InvoiceNumber     int         `json:"invoice_number" xml:"invoice_number"`
	PointOfSale       int         `json:"point_of_sale" xml:"point_of_sale"`
	InvoiceType       InvoiceType `json:"invoice_type" xml:"invoice_type"`
	AuthorizationDate time.Time   `json:"authorization_date" xml:"authorization_date"`
	Status            string      `json:"status" xml:"status"`
	Message           string      `json:"message,omitempty" xml:"message,omitempty"`
}

// Parameters representa los parámetros del sistema
type Parameters struct {
	DocumentTypes []DocumentTypeInfo `json:"document_types" xml:"document_types"`
	InvoiceTypes  []InvoiceTypeInfo  `json:"invoice_types" xml:"invoice_types"`
	CurrencyTypes []CurrencyTypeInfo `json:"currency_types" xml:"currency_types"`
	TaxRates      []TaxRateInfo      `json:"tax_rates" xml:"tax_rates"`
	ConceptTypes  []ConceptTypeInfo  `json:"concept_types" xml:"concept_types"`
	LastUpdate    time.Time          `json:"last_update" xml:"last_update"`
}

// DocumentTypeInfo representa información de un tipo de documento
type DocumentTypeInfo struct {
	ID          DocumentType `json:"id" xml:"id"`
	Description string       `json:"description" xml:"description"`
	Active      bool         `json:"active" xml:"active"`
}

// InvoiceTypeInfo representa información de un tipo de factura
type InvoiceTypeInfo struct {
	ID          InvoiceType `json:"id" xml:"id"`
	Description string      `json:"description" xml:"description"`
	Active      bool        `json:"active" xml:"active"`
}

// CurrencyTypeInfo representa información de un tipo de moneda
type CurrencyTypeInfo struct {
	ID          CurrencyType `json:"id" xml:"id"`
	Description string       `json:"description" xml:"description"`
	Active      bool         `json:"active" xml:"active"`
}

// TaxRateInfo representa información de una alícuota
type TaxRateInfo struct {
	ID          TaxRate `json:"id" xml:"id"`
	Description string  `json:"description" xml:"description"`
	Active      bool    `json:"active" xml:"active"`
}

// ConceptTypeInfo representa información de un tipo de concepto
type ConceptTypeInfo struct {
	ID          ConceptType `json:"id" xml:"id"`
	Description string      `json:"description" xml:"description"`
	Active      bool        `json:"active" xml:"active"`
}
