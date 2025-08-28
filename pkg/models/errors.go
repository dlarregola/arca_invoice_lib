package models

import (
	"fmt"
	"strings"
)

// ARCAError representa un error específico de ARCA
type ARCAError struct {
	Code    string `json:"code" xml:"code"`
	Message string `json:"message" xml:"message"`
	Details string `json:"details,omitempty" xml:"details,omitempty"`
}

// Error implementa la interfaz error
func (e *ARCAError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("ARCA Error %s: %s - %s", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("ARCA Error %s: %s", e.Code, e.Message)
}

// IsARCAError verifica si un error es un ARCAError
func IsARCAError(err error) bool {
	_, ok := err.(*ARCAError)
	return ok
}

// GetARCAError extrae un ARCAError de un error
func GetARCAError(err error) *ARCAError {
	if arcaErr, ok := err.(*ARCAError); ok {
		return arcaErr
	}
	return nil
}

// ValidationError representa un error de validación
type ValidationError struct {
	Field   string      `json:"field" xml:"field"`
	Message string      `json:"message" xml:"message"`
	Value   interface{} `json:"value,omitempty" xml:"value,omitempty"`
}

// Error implementa la interfaz error
func (e *ValidationError) Error() string {
	return fmt.Sprintf("Validation Error in field '%s': %s", e.Field, e.Message)
}

// ValidationErrors representa múltiples errores de validación
type ValidationErrors []ValidationError

// Error implementa la interfaz error
func (e ValidationErrors) Error() string {
	if len(e) == 0 {
		return "No validation errors"
	}

	var messages []string
	for _, err := range e {
		messages = append(messages, err.Error())
	}
	return strings.Join(messages, "; ")
}

// HasErrors verifica si hay errores de validación
func (e ValidationErrors) HasErrors() bool {
	return len(e) > 0
}

// Add agrega un error de validación
func (e *ValidationErrors) Add(field, message string, value interface{}) {
	*e = append(*e, ValidationError{
		Field:   field,
		Message: message,
		Value:   value,
	})
}

// AuthenticationError representa un error de autenticación
type AuthenticationError struct {
	Message string `json:"message" xml:"message"`
	Code    string `json:"code,omitempty" xml:"code,omitempty"`
}

// Error implementa la interfaz error
func (e *AuthenticationError) Error() string {
	if e.Code != "" {
		return fmt.Sprintf("Authentication Error %s: %s", e.Code, e.Message)
	}
	return fmt.Sprintf("Authentication Error: %s", e.Message)
}

// NetworkError representa un error de red
type NetworkError struct {
	Message string `json:"message" xml:"message"`
	URL     string `json:"url,omitempty" xml:"url,omitempty"`
	Status  int    `json:"status,omitempty" xml:"status,omitempty"`
}

// Error implementa la interfaz error
func (e *NetworkError) Error() string {
	if e.Status != 0 {
		return fmt.Sprintf("Network Error %d: %s", e.Status, e.Message)
	}
	return fmt.Sprintf("Network Error: %s", e.Message)
}

// Códigos de error comunes de ARCA
const (
	// Errores de autenticación
	ErrorCodeCUITNotEnabled     = "10015"
	ErrorCodeInvalidCertificate = "10016"
	ErrorCodeExpiredCertificate = "10017"
	ErrorCodeInvalidToken       = "10018"
	ErrorCodeTokenExpired       = "10019"

	// Errores de facturación
	ErrorCodeInvalidInvoiceType    = "20001"
	ErrorCodeInvalidPointOfSale    = "20002"
	ErrorCodeInvalidInvoiceNumber  = "20003"
	ErrorCodeInvalidAmount         = "20004"
	ErrorCodeInvalidTaxAmount      = "20005"
	ErrorCodeInvalidTotalAmount    = "20006"
	ErrorCodeInvalidDate           = "20007"
	ErrorCodeInvalidCurrency       = "20008"
	ErrorCodeInvalidConceptType    = "20009"
	ErrorCodeInvalidDocumentType   = "20010"
	ErrorCodeInvalidDocumentNumber = "20011"

	// Errores de sistema
	ErrorCodeServiceUnavailable = "30001"
	ErrorCodeTimeout            = "30002"
	ErrorCodeInvalidResponse    = "30003"
	ErrorCodeRateLimitExceeded  = "30004"
)

// ErrorMessages mapea códigos de error a mensajes descriptivos
var ErrorMessages = map[string]string{
	ErrorCodeCUITNotEnabled:     "CUIT no habilitado para facturación electrónica",
	ErrorCodeInvalidCertificate: "Certificado inválido o no encontrado",
	ErrorCodeExpiredCertificate: "Certificado expirado",
	ErrorCodeInvalidToken:       "Token de acceso inválido",
	ErrorCodeTokenExpired:       "Token de acceso expirado",

	ErrorCodeInvalidInvoiceType:    "Tipo de comprobante inválido",
	ErrorCodeInvalidPointOfSale:    "Punto de venta inválido",
	ErrorCodeInvalidInvoiceNumber:  "Número de comprobante inválido",
	ErrorCodeInvalidAmount:         "Monto inválido",
	ErrorCodeInvalidTaxAmount:      "Monto de impuestos inválido",
	ErrorCodeInvalidTotalAmount:    "Monto total inválido",
	ErrorCodeInvalidDate:           "Fecha inválida",
	ErrorCodeInvalidCurrency:       "Moneda inválida",
	ErrorCodeInvalidConceptType:    "Tipo de concepto inválido",
	ErrorCodeInvalidDocumentType:   "Tipo de documento inválido",
	ErrorCodeInvalidDocumentNumber: "Número de documento inválido",

	ErrorCodeServiceUnavailable: "Servicio no disponible",
	ErrorCodeTimeout:            "Timeout en la comunicación",
	ErrorCodeInvalidResponse:    "Respuesta inválida del servidor",
	ErrorCodeRateLimitExceeded:  "Límite de requests excedido",
}

// GetErrorMessage obtiene el mensaje descriptivo para un código de error
func GetErrorMessage(code string) string {
	if message, exists := ErrorMessages[code]; exists {
		return message
	}
	return "Error desconocido"
}

// NewARCAError crea un nuevo error de ARCA
func NewARCAError(code, details string) *ARCAError {
	return &ARCAError{
		Code:    code,
		Message: GetErrorMessage(code),
		Details: details,
	}
}

// NewValidationError crea un nuevo error de validación
func NewValidationError(field, message string, value interface{}) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: message,
		Value:   value,
	}
}

// NewAuthenticationError crea un nuevo error de autenticación
func NewAuthenticationError(message, code string) *AuthenticationError {
	return &AuthenticationError{
		Message: message,
		Code:    code,
	}
}

// NewNetworkError crea un nuevo error de red
func NewNetworkError(message, url string, status int) *NetworkError {
	return &NetworkError{
		Message: message,
		URL:     url,
		Status:  status,
	}
}
