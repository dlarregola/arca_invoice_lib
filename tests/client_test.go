package tests

import (
	"arca_invoice_lib/pkg/client"
	"arca_invoice_lib/pkg/models"
	"context"
	"testing"
	"time"
)

func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		config  client.Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: client.Config{
				Environment:   models.EnvironmentTesting,
				CUIT:          "20-12345678-9",
				Certificate:   []byte("test certificate"),
				PrivateKey:    []byte("test private key"),
				Timeout:       30 * time.Second,
				RetryAttempts: 3,
				AuthCacheTTL:  23 * time.Hour,
			},
			wantErr: false,
		},
		{
			name: "invalid environment",
			config: client.Config{
				Environment:   "invalid",
				CUIT:          "20-12345678-9",
				Certificate:   []byte("test certificate"),
				PrivateKey:    []byte("test private key"),
				Timeout:       30 * time.Second,
				RetryAttempts: 3,
			},
			wantErr: true,
		},
		{
			name: "empty CUIT",
			config: client.Config{
				Environment:   models.EnvironmentTesting,
				CUIT:          "",
				Certificate:   []byte("test certificate"),
				PrivateKey:    []byte("test private key"),
				Timeout:       30 * time.Second,
				RetryAttempts: 3,
			},
			wantErr: true,
		},
		{
			name: "empty certificate",
			config: client.Config{
				Environment:   models.EnvironmentTesting,
				CUIT:          "20-12345678-9",
				Certificate:   []byte{},
				PrivateKey:    []byte("test private key"),
				Timeout:       30 * time.Second,
				RetryAttempts: 3,
			},
			wantErr: true,
		},
		{
			name: "empty private key",
			config: client.Config{
				Environment:   models.EnvironmentTesting,
				CUIT:          "20-12345678-9",
				Certificate:   []byte("test certificate"),
				PrivateKey:    []byte{},
				Timeout:       30 * time.Second,
				RetryAttempts: 3,
			},
			wantErr: true,
		},
		{
			name: "zero timeout",
			config: client.Config{
				Environment:   models.EnvironmentTesting,
				CUIT:          "20-12345678-9",
				Certificate:   []byte("test certificate"),
				PrivateKey:    []byte("test private key"),
				Timeout:       0,
				RetryAttempts: 3,
			},
			wantErr: true,
		},
		{
			name: "negative retry attempts",
			config: client.Config{
				Environment:   models.EnvironmentTesting,
				CUIT:          "20-12345678-9",
				Certificate:   []byte("test certificate"),
				PrivateKey:    []byte("test private key"),
				Timeout:       30 * time.Second,
				RetryAttempts: -1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDefaultConfig(t *testing.T) {
	config := client.DefaultConfig()

	// Verificar valores por defecto
	if config.Environment != models.EnvironmentTesting {
		t.Errorf("Default environment should be testing, got %s", config.Environment)
	}

	if config.Timeout != 30*time.Second {
		t.Errorf("Default timeout should be 30s, got %v", config.Timeout)
	}

	if config.RetryAttempts != 3 {
		t.Errorf("Default retry attempts should be 3, got %d", config.RetryAttempts)
	}

	if config.AuthCacheTTL != 23*time.Hour {
		t.Errorf("Default auth cache TTL should be 23h, got %v", config.AuthCacheTTL)
	}
}

func TestConfigURLs(t *testing.T) {
	testingConfig := client.Config{Environment: models.EnvironmentTesting}
	productionConfig := client.Config{Environment: models.EnvironmentProduction}

	// Verificar URLs de testing
	if testingConfig.GetBaseURL() != "https://wswhomo.afip.gov.ar" {
		t.Errorf("Testing base URL should be https://wswhomo.afip.gov.ar, got %s", testingConfig.GetBaseURL())
	}

	// Verificar URLs de production
	if productionConfig.GetBaseURL() != "https://servicios1.afip.gov.ar" {
		t.Errorf("Production base URL should be https://servicios1.afip.gov.ar, got %s", productionConfig.GetBaseURL())
	}

	// Verificar URLs de servicios
	wsaaURL := testingConfig.GetWSAAURL()
	if wsaaURL != "https://wswhomo.afip.gov.ar/ws/services/LoginCms" {
		t.Errorf("WSAA URL should be correct, got %s", wsaaURL)
	}

	wsfeURL := testingConfig.GetWSFEURL()
	if wsfeURL != "https://wswhomo.afip.gov.ar/wsfev1/service.asmx" {
		t.Errorf("WSFE URL should be correct, got %s", wsfeURL)
	}

	wsfexURL := testingConfig.GetWSFEXURL()
	if wsfexURL != "https://wswhomo.afip.gov.ar/wsfexv1/service.asmx" {
		t.Errorf("WSFEX URL should be correct, got %s", wsfexURL)
	}
}

func TestConfigBuilder(t *testing.T) {
	config := client.DefaultConfig()
	config.Environment = models.EnvironmentProduction
	config.CUIT = "30-98765432-1"
	config.Timeout = 60 * time.Second
	config.RetryAttempts = 5
	config.LogLevel = "debug"

	// Verificar que los valores se configuraron correctamente
	if config.Environment != models.EnvironmentProduction {
		t.Errorf("Environment should be production, got %s", config.Environment)
	}

	if config.CUIT != "30-98765432-1" {
		t.Errorf("CUIT should be 30-98765432-1, got %s", config.CUIT)
	}

	if config.Timeout != 60*time.Second {
		t.Errorf("Timeout should be 60s, got %v", config.Timeout)
	}

	if config.RetryAttempts != 5 {
		t.Errorf("Retry attempts should be 5, got %d", config.RetryAttempts)
	}

	if config.LogLevel != "debug" {
		t.Errorf("Log level should be debug, got %s", config.LogLevel)
	}
}

func TestClientCreation(t *testing.T) {
	// Test con configuración válida
	validConfig := client.Config{
		Environment:   models.EnvironmentTesting,
		CUIT:          "20-12345678-9",
		Certificate:   []byte("test certificate"),
		PrivateKey:    []byte("test private key"),
		Timeout:       30 * time.Second,
		RetryAttempts: 3,
		AuthCacheTTL:  23 * time.Hour,
	}

	afipClient, err := client.NewAFIPClient(validConfig)
	if err != nil {
		t.Errorf("NewAFIPClient() should not return error with valid config: %v", err)
	}

	if afipClient == nil {
		t.Error("NewAFIPClient() should return a client")
	}

	// Verificar que los servicios están disponibles
	// Nota: Los servicios WSFE y WSFEX necesitan ser configurados después de crear el cliente
	// Por ahora verificamos que el cliente se creó correctamente
	if afipClient == nil {
		t.Error("AFIP client should not be nil")
	}

	// Test con configuración inválida
	invalidConfig := client.Config{
		Environment:   "invalid",
		CUIT:          "20-12345678-9",
		Certificate:   []byte("test certificate"),
		PrivateKey:    []byte("test private key"),
		Timeout:       30 * time.Second,
		RetryAttempts: 3,
		AuthCacheTTL:  23 * time.Hour,
	}

	_, err = client.NewAFIPClient(invalidConfig)
	if err == nil {
		t.Error("NewAFIPClient() should return error with invalid config")
	}
}

func TestSystemStatus(t *testing.T) {
	config := client.Config{
		Environment:   models.EnvironmentTesting,
		CUIT:          "20-12345678-9",
		Certificate:   []byte("test certificate"),
		PrivateKey:    []byte("test private key"),
		Timeout:       30 * time.Second,
		RetryAttempts: 3,
		AuthCacheTTL:  23 * time.Hour,
	}

	afipClient, err := client.NewAFIPClient(config)
	if err != nil {
		t.Skipf("Skipping test due to client creation error: %v", err)
	}

	ctx := context.Background()
	status, err := afipClient.GetSystemStatus(ctx)

	// En un ambiente de testing sin certificados válidos, esperamos un error
	// pero el status debería tener información básica
	if status == nil {
		t.Error("GetSystemStatus() should return a status object even on error")
	}

	if status.Timestamp.IsZero() {
		t.Error("Status timestamp should not be zero")
	}
}

func TestAuthCache(t *testing.T) {
	config := client.Config{
		Environment:   models.EnvironmentTesting,
		CUIT:          "20-12345678-9",
		Certificate:   []byte("test certificate"),
		PrivateKey:    []byte("test private key"),
		Timeout:       30 * time.Second,
		RetryAttempts: 3,
		AuthCacheTTL:  23 * time.Hour,
	}

	afipClient, err := client.NewAFIPClient(config)
	if err != nil {
		t.Skipf("Skipping test due to client creation error: %v", err)
	}

	// Verificar tamaño inicial del cache
	initialSize := afipClient.GetAuthCacheSize()
	if initialSize != 0 {
		t.Errorf("Initial cache size should be 0, got %d", initialSize)
	}

	// Limpiar cache
	afipClient.ClearAuthCache()

	// Verificar que el cache está vacío
	size := afipClient.GetAuthCacheSize()
	if size != 0 {
		t.Errorf("Cache size after clear should be 0, got %d", size)
	}
}
