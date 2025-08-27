package auth

import (
	"arca_invoice_lib/internal/shared"
	"arca_invoice_lib/pkg/interfaces"
)

// NewAuthService crea un nuevo servicio de autenticación
func NewAuthService(config *shared.InternalConfig, logger interfaces.Logger) interfaces.AuthService {
	return newAuthService(config, logger)
}
