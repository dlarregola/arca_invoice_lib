package shared

import (
	"time"
)

// InternalConfig representa la configuración interna del cliente
type InternalConfig struct {
	CUIT          string
	Certificate   []byte
	PrivateKey    []byte
	Environment   string
	Timeout       time.Duration
	RetryAttempts int
}

// GetBaseURL retorna la URL base según el environment
func (c *InternalConfig) GetBaseURL() string {
	switch c.Environment {
	case "testing":
		return "https://wswhomo.afip.gov.ar"
	case "production":
		return "https://servicios1.afip.gov.ar"
	default:
		return "https://wswhomo.afip.gov.ar"
	}
}

// GetWSAAURL retorna la URL del servicio WSAA
func (c *InternalConfig) GetWSAAURL() string {
	return c.GetBaseURL() + "/ws/services/LoginCms"
}

// GetWSFEURL retorna la URL del servicio WSFEv1
func (c *InternalConfig) GetWSFEURL() string {
	return c.GetBaseURL() + "/wsfev1/service.asmx"
}

// GetWSFEXURL retorna la URL del servicio WSFEXv1
func (c *InternalConfig) GetWSFEXURL() string {
	return c.GetBaseURL() + "/wsfexv1/service.asmx"
}
