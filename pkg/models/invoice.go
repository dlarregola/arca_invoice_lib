package models

import (
	"time"
)

// Invoice representa una factura
type Invoice struct {
	InvoiceBase
	DocType       DocumentType `json:"doc_type" xml:"doc_type"`
	DocNumber     string       `json:"doc_number" xml:"doc_number"`
	DocTypeFrom   DocumentType `json:"doc_type_from" xml:"doc_type_from"`
	DocNumberFrom string       `json:"doc_number_from" xml:"doc_number_from"`
	NameFrom      string       `json:"name_from" xml:"name_from"`
	ServiceFrom   time.Time    `json:"service_from" xml:"service_from"`
}

// ExportInvoice representa una factura de exportación
type ExportInvoice struct {
	InvoiceBase
	Destination     string    `json:"destination" xml:"destination"`
	DestinationCode string    `json:"destination_code" xml:"destination_code"`
	ExportDate      time.Time `json:"export_date" xml:"export_date"`
	ExportType      string    `json:"export_type" xml:"export_type"`
}

// InvoiceQuery representa una consulta de factura
type InvoiceQuery struct {
	InvoiceType   InvoiceType `json:"invoice_type" xml:"invoice_type"`
	PointOfSale   int         `json:"point_of_sale" xml:"point_of_sale"`
	InvoiceNumber int         `json:"invoice_number" xml:"invoice_number"`
	DateFrom      time.Time   `json:"date_from" xml:"date_from"`
	DateTo        time.Time   `json:"date_to" xml:"date_to"`
}

// ExportInvoiceQuery representa una consulta de factura de exportación
type ExportInvoiceQuery struct {
	InvoiceType   InvoiceType `json:"invoice_type" xml:"invoice_type"`
	PointOfSale   int         `json:"point_of_sale" xml:"point_of_sale"`
	InvoiceNumber int         `json:"invoice_number" xml:"invoice_number"`
	DateFrom      time.Time   `json:"date_from" xml:"date_from"`
	DateTo        time.Time   `json:"date_to" xml:"date_to"`
}

// AuthorizationResponse representa la respuesta de autorización
type AuthorizationResponse struct {
	CAE               string      `json:"cae" xml:"cae"`
	CAEExpirationDate time.Time   `json:"cae_expiration_date" xml:"cae_expiration_date"`
	InvoiceNumber     int         `json:"invoice_number" xml:"invoice_number"`
	PointOfSale       int         `json:"point_of_sale" xml:"point_of_sale"`
	InvoiceType       InvoiceType `json:"invoice_type" xml:"invoice_type"`
	AuthorizationDate time.Time   `json:"authorization_date" xml:"authorization_date"`
	Status            string      `json:"status" xml:"status"`
	Message           string      `json:"message,omitempty" xml:"message,omitempty"`
}

// ExportAuthResponse representa la respuesta de autorización de exportación
type ExportAuthResponse struct {
	AuthorizationResponse
	ExportType string `json:"export_type" xml:"export_type"`
}

// LastInvoiceResponse representa la respuesta del último comprobante
type LastInvoiceResponse struct {
	InvoiceType   InvoiceType `json:"invoice_type" xml:"invoice_type"`
	PointOfSale   int         `json:"point_of_sale" xml:"point_of_sale"`
	InvoiceNumber int         `json:"invoice_number" xml:"invoice_number"`
	Date          time.Time   `json:"date" xml:"date"`
}

// CAEAResponse representa la respuesta de consulta CAEA
type CAEAResponse struct {
	CAEA           string    `json:"caea" xml:"caea"`
	ExpirationDate time.Time `json:"expiration_date" xml:"expiration_date"`
	Status         string    `json:"status" xml:"status"`
	Message        string    `json:"message,omitempty" xml:"message,omitempty"`
}

// Currency representa una moneda
type Currency struct {
	ID          string `json:"id" xml:"id"`
	Description string `json:"description" xml:"description"`
	Active      bool   `json:"active" xml:"active"`
}

// Destination representa un destino de exportación
type Destination struct {
	ID          string `json:"id" xml:"id"`
	Description string `json:"description" xml:"description"`
	Active      bool   `json:"active" xml:"active"`
}

// UnitType representa un tipo de unidad
type UnitType struct {
	ID          string `json:"id" xml:"id"`
	Description string `json:"description" xml:"description"`
	Active      bool   `json:"active" xml:"active"`
}
