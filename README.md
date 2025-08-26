# AFIP Go Library

[![Go](https://github.com/YOUR_USERNAME/invoiceservice/workflows/Go/badge.svg)](https://github.com/YOUR_USERNAME/invoiceservice/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/YOUR_USERNAME/invoiceservice)](https://goreportcard.com/report/github.com/YOUR_USERNAME/invoiceservice)
[![GoDoc](https://godoc.org/github.com/YOUR_USERNAME/invoiceservice?status.svg)](https://godoc.org/github.com/YOUR_USERNAME/invoiceservice)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Librería en Go para interactuar con los Web Services de facturación electrónica de AFIP (WSFEv1 y WSFEXv1).

## Características

- ✅ **WSFEv1** - Factura Electrónica Nacional (Manual v4.0, R.G. N° 4.291)
- ✅ **WSFEXv1** - Factura Electrónica de Exportación (Manual v3.1.0, R.G. N° 2.758)
- ✅ **Autenticación WSAA** - Manejo automático de tickets de acceso
- ✅ **Thread-safe** - Uso concurrente seguro
- ✅ **Retry automático** - Reintentos con backoff exponencial
- ✅ **Logging estructurado** - Logs detallados para debugging
- ✅ **Validaciones** - Validación de datos antes del envío

## Instalación

```bash
go get github.com/afip-go
```

## Configuración

### 1. Obtener Certificados

Para usar los Web Services de AFIP necesitas:

1. **CUIT habilitado** en AFIP
2. **Certificado X.509** (.crt)
3. **Clave privada** (.key)

### 2. Configurar el Cliente

```go
package main

import (
    "context"
    "log"
    "time"
    
    "github.com/afip-go/pkg/client"
    "github.com/afip-go/pkg/wsfe"
)

func main() {
    // Cargar certificado y clave privada
    cert, err := os.ReadFile("certificate.crt")
    if err != nil {
        log.Fatal(err)
    }
    
    privateKey, err := os.ReadFile("private.key")
    if err != nil {
        log.Fatal(err)
    }
    
    // Configurar cliente
    config := client.Config{
        Environment:   "testing", // "testing" o "production"
        CUIT:         "20123456789",
        Certificate:  cert,
        PrivateKey:   privateKey,
        Timeout:      30 * time.Second,
        RetryAttempts: 3,
    }
    
    // Crear cliente
    afipClient, err := client.NewAFIPClient(config)
    if err != nil {
        log.Fatal(err)
    }
    
    // Usar servicios...
}
```

## Uso

### Facturación Nacional (WSFEv1)

```go
// Crear factura
invoice := &wsfe.Invoice{
    DocType:      11, // CUIT
    DocNumber:    "20123456789",
    ConceptType:  1,  // Productos
    DocTypeFrom:  1,  // DNI
    DocNumberFrom: "12345678",
    Amount:       1000.00,
    TaxAmount:    210.00,
    TotalAmount:  1210.00,
    CurrencyType: "PES",
    DateFrom:     time.Now(),
    DateTo:       time.Now(),
    ServiceFrom:  "Servicio de ejemplo",
    Items: []wsfe.InvoiceItem{
        {
            Description: "Producto 1",
            Quantity:    1,
            UnitPrice:   1000.00,
            TotalPrice:  1000.00,
        },
    },
}

// Autorizar factura
ctx := context.Background()
result, err := afipClient.WSFE.AuthorizeInvoice(ctx, invoice)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("CAE: %s\n", result.CAE)
fmt.Printf("Vencimiento CAE: %s\n", result.CAEExpirationDate)
```

### Facturación Internacional (WSFEXv1)

```go
// Crear factura de exportación
exportInvoice := &wsfex.ExportInvoice{
    DocType:      11,
    DocNumber:    "20123456789",
    ConceptType:  1,
    DocTypeFrom:  1,
    DocNumberFrom: "12345678",
    Amount:       1000.00,
    CurrencyType: "USD",
    DateFrom:     time.Now(),
    DateTo:       time.Now(),
    ServiceFrom:  "Servicio de exportación",
    Items: []wsfex.ExportInvoiceItem{
        {
            Description: "Producto exportación",
            Quantity:    1,
            UnitPrice:   1000.00,
            TotalPrice:  1000.00,
        },
    },
}

// Autorizar factura de exportación
result, err := afipClient.WSFEX.AuthorizeExportInvoice(ctx, exportInvoice)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("CAE: %s\n", result.CAE)
```

### Consultas

```go
// Consultar último comprobante autorizado
lastAuth, err := afipClient.WSFE.GetLastAuthorizedInvoice(ctx, 1, 6) // Punto de venta 1, tipo 6 (Factura B)
if err != nil {
    log.Fatal(err)
}

// Consultar comprobante específico
invoice, err := afipClient.WSFE.GetInvoice(ctx, 1, 6, 1) // Punto de venta 1, tipo 6, número 1
if err != nil {
    log.Fatal(err)
}

// Obtener parámetros
params, err := afipClient.WSFE.GetParameters(ctx)
if err != nil {
    log.Fatal(err)
}
```

## Estructura del Proyecto

```
afip-go/
├── pkg/
│   ├── client/          # Cliente principal y autenticación
│   ├── wsfe/           # Facturación Nacional
│   ├── wsfex/          # Facturación Internacional
│   └── models/         # Modelos compartidos
├── internal/
│   ├── soap/           # Cliente SOAP interno
│   └── utils/          # Utilidades internas
├── examples/           # Ejemplos de uso
└── tests/              # Tests unitarios
```

## Ambientes

- **Testing**: `https://wswhomo.afip.gov.ar/`
- **Production**: `https://servicios1.afip.gov.ar/`

## Manejo de Errores

La librería incluye errores específicos de AFIP:

```go
if err != nil {
    if afipErr, ok := err.(*models.AFIPError); ok {
        switch afipErr.Code {
        case "10015":
            fmt.Println("Error: CUIT no habilitado")
        case "10016":
            fmt.Println("Error: Certificado inválido")
        default:
            fmt.Printf("Error AFIP: %s - %s\n", afipErr.Code, afipErr.Message)
        }
    }
}
```

## Logging

```go
// Habilitar logs detallados
afipClient.SetLogLevel(logrus.DebugLevel)

// Los logs incluyen:
// - Requests/responses SOAP
// - Errores de autenticación
// - Métricas de performance
```

## Testing

```bash
# Ejecutar tests unitarios
go test ./...

# Ejecutar tests con coverage
go test -cover ./...

# Ejecutar tests de integración (requiere certificados válidos)
go test -tags=integration ./tests/
```

## Contribuir

1. Fork el proyecto
2. Crea una rama para tu feature (`git checkout -b feature/AmazingFeature`)
3. Commit tus cambios (`git commit -m 'Add some AmazingFeature'`)
4. Push a la rama (`git push origin feature/AmazingFeature`)
5. Abre un Pull Request

## Licencia

Este proyecto está bajo la Licencia MIT. Ver el archivo `LICENSE` para más detalles.

## Referencias

- [Manual WSFEv1 v4.0](https://www.afip.gob.ar/ws/documentacion/ws-factura-electronica.asp)
- [Manual WSFEXv1 v3.1.0](https://www.afip.gob.ar/ws/documentacion/ws-factura-electronica-exportacion.asp)
- [Portal AFIP WebServices](https://www.afip.gob.ar/ws)
