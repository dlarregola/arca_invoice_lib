# ARCA Go Library - Multi-Tenant

[![Go](https://github.com/YOUR_USERNAME/invoiceservice/workflows/Go/badge.svg)](https://github.com/YOUR_USERNAME/invoiceservice/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/YOUR_USERNAME/invoiceservice)](https://goreportcard.com/report/github.com/YOUR_USERNAME/invoiceservice)
[![GoDoc](https://godoc.org/github.com/YOUR_USERNAME/invoiceservice?status.svg)](https://godoc.org/github.com/YOUR_USERNAME/invoiceservice)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Librería en Go para interactuar con los Web Services de facturación electrónica de ARCA (WSFEv1 y WSFEXv1) con soporte para arquitectura multi-tenant.

## Características

- ✅ **Multi-Tenant** - Soporte para múltiples empresas con cache de clientes
- ✅ **WSFEv1** - Factura Electrónica Nacional (Manual v4.0, R.G. N° 4.291)
- ✅ **WSFEXv1** - Factura Electrónica de Exportación (Manual v3.1.0, R.G. N° 2.758)
- ✅ **Autenticación WSAA** - Manejo automático de tickets de acceso
- ✅ **Thread-safe** - Uso concurrente seguro
- ✅ **Cache inteligente** - Cache de clientes con limpieza automática
- ✅ **Factory Pattern** - Creación flexible de managers
- ✅ **Interfaces públicas** - API limpia y extensible
- ✅ **Retry automático** - Reintentos con backoff exponencial
- ✅ **Logging estructurado** - Logs detallados para debugging
- ✅ **Validaciones** - Validación de datos antes del envío

## Arquitectura Multi-Tenant

La librería implementa un patrón multi-tenant donde cada empresa gestiona su propia conexión a ARCA:

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Application   │    │  ARCA Manager    │    │   ARCA Client   │
│                 │    │                  │    │                 │
│ Company A ──────┼───▶│  Cache Manager   │───▶│  Company A      │
│ Company B ──────┼───▶│  Config Provider │───▶│  Company B      │
│ Company C ──────┼───▶│  Factory Pattern │───▶│  Company C      │
└─────────────────┘    └──────────────────┘    └─────────────────┘
```

### Ventajas del Patrón Multi-Tenant

- **Aislamiento**: Cada empresa tiene su propia configuración y conexión
- **Escalabilidad**: Cache inteligente reduce la creación de conexiones
- **Flexibilidad**: Configuración dinámica por empresa
- **Mantenibilidad**: Código limpio con interfaces bien definidas
- **Performance**: Reutilización de conexiones activas

## Instalación

```bash
go get github.com/arca-go
```

## Uso Multi-Tenant

### Ventajas del Nuevo Enfoque

El nuevo patrón de configuración en tiempo de ejecución ofrece varias ventajas importantes:

- **🚀 Mejor Escalabilidad**: No es necesario precargar todas las configuraciones de empresas al inicio
- **⚡ Configuración Dinámica**: Las configuraciones se pueden obtener de bases de datos, APIs, o sistemas externos en tiempo real
- **🔒 Seguridad Mejorada**: Los certificados y claves privadas se cargan solo cuando son necesarios
- **💾 Uso Eficiente de Memoria**: Solo se mantienen en memoria las configuraciones de empresas activas
- **🔄 Flexibilidad**: Fácil integración con sistemas de gestión de configuraciones dinámicas

### 1. Configuración Básica

```go
package main

import (
    "context"
    "time"
    
    "arca_invoice_lib/pkg/factory"
    "arca_invoice_lib/pkg/interfaces"
)

// CompanyConfiguration implementa la interfaz CompanyConfig
type CompanyConfiguration struct {
    CompanyID   string
    CUIT        string
    Certificate []byte
    PrivateKey  []byte
    Environment string
}

func (c *CompanyConfiguration) GetCUIT() string { return c.CUIT }
func (c *CompanyConfiguration) GetCertificate() []byte { return c.Certificate }
func (c *CompanyConfiguration) GetPrivateKey() []byte { return c.PrivateKey }
func (c *CompanyConfiguration) GetEnvironment() string { return c.Environment }
func (c *CompanyConfiguration) GetCompanyID() string { return c.CompanyID }

func main() {
    // 1. Crear factory
    factory := factory.NewClientManagerFactory()

    // 2. Configurar manager
    manager := factory.CreateManager(factory.ManagerConfig{
        CompanyConfigProvider: func(ctx context.Context, companyID string) (interfaces.CompanyConfig, error) {
            // Aquí buscarías en tu DB los datos de la empresa
            return &CompanyConfiguration{
                CompanyID:   companyID,
                CUIT:        "20-12345678-9",
                Certificate: []byte("certificado..."),
                PrivateKey:  []byte("clave privada..."),
                Environment: "testing",
            }, nil
        },
        ClientCacheSize:   100,
        ClientIdleTimeout: 30 * time.Minute,
        HTTPTimeout:       30 * time.Second,
        MaxRetryAttempts:  3,
    })

    // 3. Usar en tu servicio
    ctx := context.Background()
    
    // Crear configuración de empresa (en tiempo de ejecución)
    companyConfig := &models.CompanyConfiguration{
        CompanyID:   "empresa-001",
        CUIT:        "20-12345678-9",
        Certificate: []byte("certificado..."),
        PrivateKey:  []byte("clave privada..."),
        Environment: "testing",
    }
    
    // Obtener cliente específico de la empresa
    client, err := manager.GetClientForCompany(ctx, companyConfig)
    if err != nil {
        log.Fatal(err)
    }

    // Usar servicios...
}
```

### 2. Uso en Servicios

```go
// InvoiceService maneja facturas para múltiples empresas
type InvoiceService struct {
    arcaManager interfaces.ARCAClientManager
}

func (s *InvoiceService) CreateInvoice(ctx context.Context, companyConfig interfaces.CompanyConfig, invoiceData *models.Invoice) error {
    // Obtener cliente específico de la empresa
    client, err := s.arcaManager.GetClientForCompany(ctx, companyConfig)
    if err != nil {
        return fmt.Errorf("failed to get ARCA client: %w", err)
    }
    
    // Usar servicios específicos
    response, err := client.WSFE().AuthorizeInvoice(ctx, invoiceData)
    if err != nil {
        return fmt.Errorf("failed to authorize invoice: %w", err)
    }
    
    // Procesar respuesta...
    return nil
}
```

### 3. Gestión de Cache

```go
// Obtener estadísticas del cache
stats := manager.GetCacheStats()
fmt.Printf("Total Clients: %d\n", stats.TotalClients)
fmt.Printf("Active Clients: %d\n", stats.ActiveClients)

// Limpiar cache de clientes inactivos
manager.CleanupInactiveClients(5 * time.Minute)

// Invalidar cliente específico
manager.InvalidateClient("empresa-001")
```

## Interfaces Principales

### ARCAClientManager

```go
type ARCAClientManager interface {
    GetClientForCompany(ctx context.Context, companyConfig CompanyConfig) (ARCAClient, error)
    ValidateCompanyConfig(config CompanyConfig) error
    CleanupInactiveClients(maxIdleTime time.Duration)
    InvalidateClient(companyID string)
    GetCacheStats() CacheStats
}
```

### ARCAClient

```go
type ARCAClient interface {
    WSFE() WSFEService
    WSFEX() WSFEXService
    GetCompanyInfo() CompanyInfo
    IsHealthy(ctx context.Context) error
    Close() error
}
```

### CompanyConfig

```go
type CompanyConfig interface {
    GetCUIT() string
    GetCertificate() []byte
    GetPrivateKey() []byte
    GetEnvironment() string
    GetCompanyID() string
}
```

## Configuración

### 1. Obtener Certificados

Para usar los Web Services de ARCA necesitas:

1. **CUIT habilitado** en ARCA
2. **Certificado X.509** (.crt)
3. **Clave privada** (.key)

### 2. Configuración por Empresa

Cada empresa debe tener su propia configuración:

```go
type CompanyConfig struct {
    CompanyID   string
    CUIT        string
    Certificate []byte
    PrivateKey  []byte
    Environment string // "testing" o "production"
}
```

## Uso de Servicios

### Facturación Nacional (WSFEv1)

```go
// Autorizar factura
response, err := client.WSFE().AuthorizeInvoice(ctx, invoice)
if err != nil {
    return err
}

// Consultar factura
query := &models.InvoiceQuery{
    InvoiceType:   models.InvoiceTypeA,
    PointOfSale:   1,
    InvoiceNumber: 1001,
    DateFrom:      time.Now().AddDate(0, 0, -30),
    DateTo:        time.Now(),
}

invoice, err := client.WSFE().QueryInvoice(ctx, query)
```

### Facturación Internacional (WSFEXv1)

```go
// Autorizar factura de exportación
exportInvoice := &models.ExportInvoice{
    InvoiceBase: models.InvoiceBase{...},
    Destination:     "Estados Unidos",
    DestinationCode: "US",
    ExportDate:      time.Now(),
    ExportType:      "Definitiva",
}

response, err := client.WSFEX().AuthorizeExportInvoice(ctx, exportInvoice)
```

## Documentación

La documentación está dividida en dos secciones principales:

### 📚 [Documentación Técnica](docs/ARCHITECTURE.md)
Para desarrolladores y arquitectos que necesitan entender la arquitectura interna, decisiones de diseño, y implementación de la librería.

### 🚀 [Guía de Uso](docs/USAGE.md)
Para desarrolladores que quieren usar la librería, con ejemplos prácticos y casos de uso.

### 📖 [Índice de Documentación](docs/README.md)
Punto de entrada principal con navegación a todas las secciones.

## Ejemplos

- [Ejemplo Básico](examples/basic_example.go) - Uso básico de la librería
- [Ejemplo Avanzado](examples/advanced_usage.go) - Uso multi-tenant con múltiples empresas

El nuevo patrón de configuración en tiempo de ejecución se puede usar de la siguiente manera:

```go
// 1. Crear manager
factory := arca.NewClientManagerFactory()
manager := factory.CreateManager(arca.ManagerConfig{
    ClientCacheSize:   100,
    ClientIdleTimeout: 30 * time.Minute,
    HTTPTimeout:       30 * time.Second,
    MaxRetryAttempts:  3,
    Logger:            &MyLogger{},
})

// 2. En tu servicio, obtener configuración de empresa en tiempo de ejecución
func (s *InvoiceService) CreateInvoice(ctx context.Context, companyID string, invoiceData *models.Invoice) error {
    // Obtener configuración de empresa desde tu base de datos/API
    companyConfig, err := s.dbService.GetCompanyConfig(ctx, companyID)
    if err != nil {
        return fmt.Errorf("failed to get company config: %w", err)
    }
    
    // Obtener cliente usando la configuración
    client, err := s.arcaManager.GetClientForCompany(ctx, companyConfig)
    if err != nil {
        return fmt.Errorf("failed to get ARCA client: %w", err)
    }
    
    // Usar servicios
    response, err := client.WSFE().AuthorizeInvoice(ctx, invoiceData)
    if err != nil {
        return fmt.Errorf("failed to authorize invoice: %w", err)
    }
    
    return nil
}
```

## Estructura del Proyecto

```
arca-go/
├── pkg/
│   ├── interfaces/            # Interfaces públicas
│   │   ├── client.go         # ARCAClientManager, ARCAClient
│   │   ├── wsfe.go           # WSFEService interface
│   │   ├── wsfex.go          # WSFEXService interface
│   │   └── auth.go           # AuthService interface
│   ├── factory/              # Factory para crear clientes
│   │   └── client_factory.go
│   ├── models/               # Modelos públicos (solo datos)
│   │   ├── common.go
│   │   ├── invoice.go
│   │   └── errors.go
│   └── errors/               # Errores públicos
│       └── errors.go
├── internal/
│   ├── client/               # Implementaciones privadas
│   │   ├── manager.go        # clientManager (privado)
│   │   └── arca_client.go    # arcaClient (privado)
│   ├── services/             # Servicios privados
│   │   ├── wsfe/
│   │   ├── wsfex/
│   │   └── auth/
│   ├── config/               # Configuración interna
│   └── utils/
├── examples/
└── README.md
```

## Contribuir

1. Fork el proyecto
2. Crea una rama para tu feature (`git checkout -b feature/AmazingFeature`)
3. Commit tus cambios (`git commit -m 'Add some AmazingFeature'`)
4. Push a la rama (`git push origin feature/AmazingFeature`)
5. Abre un Pull Request

## Licencia

Este proyecto está bajo la Licencia MIT - ver el archivo [LICENSE](LICENSE) para detalles.

## Soporte

Para soporte y preguntas:

- 📧 Email: support@arca-go.com
- 📖 Documentación: [docs/](docs/)
- 🐛 Issues: [GitHub Issues](https://github.com/arca-go/issues)
