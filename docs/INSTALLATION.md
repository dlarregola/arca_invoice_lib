# Guía de Instalación y Configuración

## Requisitos Previos

### 1. CUIT Habilitado
Para usar los Web Services de AFIP necesitas:

- **CUIT habilitado** para facturación electrónica
- **Certificado X.509** (.crt)
- **Clave privada** (.key)

### 2. Obtener Certificados

#### Paso 1: Solicitar Certificado
1. Ingresar a [AFIP](https://www.afip.gob.ar)
2. Ir a **Mi AFIP** > **Web Services**
3. Solicitar certificado para **WSFEv1** y **WSFEXv1**

#### Paso 2: Descargar Certificados
1. Descargar el certificado (.crt)
2. Descargar la clave privada (.key)
3. Guardar en ubicación segura

### 3. Instalación de la Librería

```bash
# Instalar la librería
go get github.com/afip-go

# O clonar el repositorio
git clone https://github.com/afip-go/afip-go.git
cd afip-go
go mod tidy
```

## Configuración

### 1. Configuración Básica

```go
package main

import (
    "github.com/afip-go/pkg/client"
    "github.com/afip-go/pkg/models"
    "time"
)

func main() {
    // Cargar certificados
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
        Environment:   models.EnvironmentTesting, // Usar testing primero
        CUIT:         "20-12345678-9",           // Tu CUIT
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
}
```

### 2. Configuración con Builder Pattern

```go
config := client.DefaultConfig().
    WithEnvironment(models.EnvironmentTesting).
    WithCUIT("20-12345678-9").
    WithCertificate(cert).
    WithPrivateKey(privateKey).
    WithTimeout(30 * time.Second).
    WithRetryAttempts(3).
    WithLogLevel("debug")
```

### 3. Configuración desde Archivo

```go
// Cargar configuración desde YAML
configData, err := os.ReadFile("config.yaml")
if err != nil {
    log.Fatal(err)
}

var config client.Config
err = yaml.Unmarshal(configData, &config)
if err != nil {
    log.Fatal(err)
}
```

## Ambientes

### Testing (Homologación)
- **URL**: `https://wswhomo.afip.gov.ar`
- **Uso**: Para pruebas y desarrollo
- **Certificados**: Especiales para testing

### Production
- **URL**: `https://servicios1.afip.gov.ar`
- **Uso**: Para uso en producción
- **Certificados**: Certificados oficiales

## Validación de Configuración

```go
// Validar configuración antes de usar
if err := config.Validate(); err != nil {
    log.Fatal("Configuración inválida:", err)
}

// Probar conexión
ctx := context.Background()
if err := afipClient.TestConnection(ctx); err != nil {
    log.Fatal("Error de conexión:", err)
}
```

## Logging

### Configurar Logging Básico

```go
import "github.com/sirupsen/logrus"

logger := logrus.New()
logger.SetLevel(logrus.DebugLevel)

afipClient.SetLogger(logger)
```

### Configurar Logging Detallado

```go
config := client.DefaultConfig().
    WithLogLevel("debug").
    WithLogRequests(true).
    WithLogResponses(true)
```

## Manejo de Errores

### Errores Comunes

```go
if err != nil {
    if afipErr, ok := err.(*models.AFIPError); ok {
        switch afipErr.Code {
        case "10015":
            log.Println("CUIT no habilitado")
        case "10016":
            log.Println("Certificado inválido")
        case "10017":
            log.Println("Certificado expirado")
        default:
            log.Printf("Error AFIP: %s - %s", afipErr.Code, afipErr.Message)
        }
    }
}
```

### Errores de Validación

```go
if err != nil {
    if validationErrs, ok := err.(models.ValidationErrors); ok {
        for _, validationErr := range validationErrs {
            log.Printf("Error en campo %s: %s", validationErr.Field, validationErr.Message)
        }
    }
}
```

## Cache de Autenticación

La librería maneja automáticamente el cache de tickets de acceso:

```go
// Verificar tamaño del cache
size := afipClient.GetAuthCacheSize()
fmt.Printf("Tickets en cache: %d\n", size)

// Limpiar cache manualmente
afipClient.ClearAuthCache()
```

## Configuración de Seguridad

### 1. Almacenamiento Seguro de Certificados

```go
// Usar variables de entorno
cert := []byte(os.Getenv("AFIP_CERTIFICATE"))
privateKey := []byte(os.Getenv("AFIP_PRIVATE_KEY"))

// O usar un gestor de secretos
cert := getSecret("afip-certificate")
privateKey := getSecret("afip-private-key")
```

### 2. Configuración de Timeouts

```go
config := client.DefaultConfig().
    WithTimeout(30 * time.Second).
    WithRetryAttempts(3).
    WithRetryDelay(1 * time.Second)
```

### 3. Configuración de Reintentos

```go
config := client.DefaultConfig().
    WithRetryAttempts(5).
    WithRetryDelay(2 * time.Second)
```

## Verificación de Instalación

### Test Básico

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/afip-go/pkg/client"
    "github.com/afip-go/pkg/models"
)

func main() {
    // Configuración mínima para testing
    config := client.Config{
        Environment:   models.EnvironmentTesting,
        CUIT:         "20-12345678-9",
        Certificate:  []byte("test"),
        PrivateKey:   []byte("test"),
        Timeout:      30 * time.Second,
        RetryAttempts: 3,
    }
    
    // Crear cliente
    afipClient, err := client.NewAFIPClient(config)
    if err != nil {
        log.Fatal("Error creando cliente:", err)
    }
    
    // Probar servicios
    ctx := context.Background()
    
    if afipClient.WSFE() != nil {
        fmt.Println("✓ Servicio WSFE disponible")
    }
    
    if afipClient.WSFEX() != nil {
        fmt.Println("✓ Servicio WSFEX disponible")
    }
    
    fmt.Println("✓ Instalación exitosa")
}
```

## Próximos Pasos

1. **Configurar certificados reales**
2. **Probar en ambiente de testing**
3. **Implementar manejo de errores**
4. **Configurar logging apropiado**
5. **Migrar a producción cuando esté listo**

## Soporte

- **Documentación**: [README.md](../README.md)
- **Ejemplos**: [examples/](../examples/)
- **Tests**: [tests/](../tests/)
- **Issues**: [GitHub Issues](https://github.com/afip-go/afip-go/issues)
