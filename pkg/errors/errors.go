package errors

import (
	"fmt"
)

// AFIPError representa un error específico de AFIP
type AFIPError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Company string `json:"company,omitempty"`
}

func (e *AFIPError) Error() string {
	if e.Company != "" {
		return fmt.Sprintf("[%s] AFIP Error %s: %s", e.Company, e.Code, e.Message)
	}
	return fmt.Sprintf("AFIP Error %s: %s", e.Code, e.Message)
}

// NewAFIPError crea un nuevo error de AFIP
func NewAFIPError(code, message string) *AFIPError {
	return &AFIPError{
		Code:    code,
		Message: message,
	}
}

// NewAFIPErrorWithCompany crea un nuevo error de AFIP con información de empresa
func NewAFIPErrorWithCompany(code, message, company string) *AFIPError {
	return &AFIPError{
		Code:    code,
		Message: message,
		Company: company,
	}
}

// CompanyConfigError representa un error de configuración de empresa
type CompanyConfigError struct {
	CompanyID string `json:"company_id"`
	Field     string `json:"field"`
	Message   string `json:"message"`
}

func (e *CompanyConfigError) Error() string {
	return fmt.Sprintf("Company %s config error in field %s: %s", e.CompanyID, e.Field, e.Message)
}

// NewCompanyConfigError crea un nuevo error de configuración de empresa
func NewCompanyConfigError(companyID, field, message string) *CompanyConfigError {
	return &CompanyConfigError{
		CompanyID: companyID,
		Field:     field,
		Message:   message,
	}
}

// ClientCacheError representa un error del cache de clientes
type ClientCacheError struct {
	CompanyID string `json:"company_id"`
	Operation string `json:"operation"`
	Message   string `json:"message"`
}

func (e *ClientCacheError) Error() string {
	return fmt.Sprintf("Client cache error for company %s during %s: %s", e.CompanyID, e.Operation, e.Message)
}

// NewClientCacheError crea un nuevo error del cache de clientes
func NewClientCacheError(companyID, operation, message string) *ClientCacheError {
	return &ClientCacheError{
		CompanyID: companyID,
		Operation: operation,
		Message:   message,
	}
}

// AuthenticationError representa un error de autenticación
type AuthenticationError struct {
	CompanyID string `json:"company_id"`
	Service   string `json:"service"`
	Message   string `json:"message"`
}

func (e *AuthenticationError) Error() string {
	return fmt.Sprintf("Authentication error for company %s service %s: %s", e.CompanyID, e.Service, e.Message)
}

// NewAuthenticationError crea un nuevo error de autenticación
func NewAuthenticationError(companyID, service, message string) *AuthenticationError {
	return &AuthenticationError{
		CompanyID: companyID,
		Service:   service,
		Message:   message,
	}
}
