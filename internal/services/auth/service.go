package auth

import (
	"github.com/dlarregola/arca_invoice_lib/internal/shared"
	"github.com/dlarregola/arca_invoice_lib/pkg/interfaces"
)

// NewAuthService crea un nuevo servicio de autenticación
func NewAuthService(config *shared.InternalConfig, logger interfaces.Logger) interfaces.AuthService {
	return newAuthService(config, logger)
}
