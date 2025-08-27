package interfaces

import (
	"context"
	"time"
)

// AuthService es la interfaz para el servicio de autenticación
type AuthService interface {
	// GetToken obtiene un token de autenticación válido
	GetToken(ctx context.Context, service string) (*AccessToken, error)

	// ClearCache limpia el cache de tokens
	ClearCache()

	// GetCacheSize retorna el tamaño del cache
	GetCacheSize() int
}

// AccessToken representa un token de acceso
type AccessToken struct {
	Token          string    `json:"token"`
	Sign           string    `json:"sign"`
	ExpirationTime time.Time `json:"expiration_time"`
	GenerationTime time.Time `json:"generation_time"`
}
