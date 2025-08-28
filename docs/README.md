# Documentación ARCA Go Library

## Bienvenido a la Documentación

Esta documentación está dividida en dos secciones principales para facilitar la comprensión y uso de la librería ARCA Go:

## 📚 Secciones de Documentación

### 1. [Documentación Técnica](ARCHITECTURE.md)
**Para desarrolladores y arquitectos**

Esta sección contiene información detallada sobre:
- 🏗️ **Arquitectura General** - Diseño y estructura de la librería
- 🏢 **Patrón Multi-Tenant** - Implementación multi-empresa
- 💾 **Sistema de Cache** - Estrategias de cache y gestión de memoria
- 🏭 **Patrón Factory** - Creación y configuración de managers
- 🔌 **Interfaces y Abstracciones** - Diseño de APIs públicas
- ⚠️ **Manejo de Errores** - Jerarquía y estrategias de errores
- 🔐 **Autenticación WSAA** - Sistema de autenticación ARCA
- 🔒 **Thread Safety** - Concurrencia y sincronización
- ✅ **Configuración y Validación** - Validaciones y configuraciones
- 🤔 **Decisiones de Implementación** - Justificación de decisiones técnicas

### 2. [Guía de Uso](USAGE.md)
**Para desarrolladores que quieren usar la librería**

Esta sección contiene información práctica sobre:
- 📦 **Instalación** - Cómo instalar y configurar la librería
- ⚙️ **Configuración Inicial** - Primeros pasos y certificados ARCA
- 🚀 **Uso Básico** - Ejemplos básicos de uso
- 🏢 **Patrón Multi-Tenant** - Cómo usar con múltiples empresas
- 🔧 **Servicios Disponibles** - WSFE y WSFEX con ejemplos
- ⚠️ **Manejo de Errores** - Cómo manejar errores en la práctica
- ⚙️ **Configuración Avanzada** - Configuraciones avanzadas
- 💡 **Ejemplos Prácticos** - Casos de uso reales
- ✅ **Mejores Prácticas** - Recomendaciones de uso
- 🔧 **Troubleshooting** - Solución de problemas comunes

## 🎯 ¿Qué Documentación Leer?

### Si eres un **Desarrollador que va a usar la librería**:
1. Comienza con la [Guía de Uso](USAGE.md)
2. Revisa los [Ejemplos Prácticos](USAGE.md#ejemplos-prácticos)
3. Consulta la [Documentación Técnica](ARCHITECTURE.md) solo si necesitas entender detalles internos

### Si eres un **Arquitecto o Desarrollador Senior**:
1. Comienza con la [Documentación Técnica](ARCHITECTURE.md)
2. Revisa las [Decisiones de Implementación](ARCHITECTURE.md#decisiones-de-implementación)
3. Consulta la [Guía de Uso](USAGE.md) para ejemplos prácticos

### Si eres un **DevOps o SRE**:
1. Revisa la [Documentación Técnica](ARCHITECTURE.md) para entender la arquitectura
2. Consulta [Configuración Avanzada](USAGE.md#configuración-avanzada) para configuraciones de producción
3. Revisa [Troubleshooting](USAGE.md#troubleshooting) para monitoreo

## 🚀 Inicio Rápido

### Instalación
```bash
go get github.com/your-org/arca_invoice_lib
```

### Uso Básico
```go
package main

import (
    "context"
    "log"
    "time"
    
    "arca_invoice_lib/pkg/factory"
    "arca_invoice_lib/pkg/interfaces"
)

func main() {
    // 1. Crear factory
    factory := factory.NewClientManagerFactory()

    // 2. Configurar manager
    manager := factory.CreateManager(factory.ManagerConfig{
        ClientCacheSize:   100,
        ClientIdleTimeout: 30 * time.Minute,
        HTTPTimeout:       30 * time.Second,
        MaxRetryAttempts:  3,
    })

    // 3. Crear configuración de empresa
    companyConfig := &CompanyConfiguration{
        CompanyID:   "empresa-001",
        CUIT:        "20-12345678-9",
        Certificate: certData,  // []byte del certificado
        PrivateKey:  keyData,   // []byte de la clave privada
        Environment: "testing",
    }

    // 4. Obtener cliente y usar servicios
    ctx := context.Background()
    client, err := manager.GetClientForCompany(ctx, companyConfig)
    if err != nil {
        log.Fatal("Failed to get client:", err)
    }

    // 5. Usar servicios (ver ejemplos en USAGE.md)
    // client.WSFE().AuthorizeInvoice(...)
    // client.WSFEX().AuthorizeExportInvoice(...)
}
```

## 📋 Requisitos Previos

- **Go 1.19+** - Versión mínima de Go
- **Certificado ARCA** - Certificado X.509 de ARCA
- **Clave Privada** - Clave privada correspondiente al certificado
- **CUIT Habilitado** - CUIT habilitado en ARCA para los servicios

## 🔗 Enlaces Útiles

- [ARCA Web Services](https://www.afip.gob.ar/ws) - Documentación oficial de ARCA
- [Go Modules](https://golang.org/ref/mod) - Documentación de módulos de Go
- [GitHub Repository](https://github.com/your-org/arca_invoice_lib) - Código fuente

## 🤝 Contribuir

Para contribuir a la documentación:

1. **Reportar Errores**: Usa el sistema de issues de GitHub
2. **Sugerir Mejoras**: Abre un issue con la etiqueta `documentation`
3. **Contribuir Código**: Sigue las guías de contribución del proyecto

## 📞 Soporte

- 📧 **Email**: support@arca-go.com
- 🐛 **Issues**: [GitHub Issues](https://github.com/your-org/arca_invoice_lib/issues)
- 📖 **Documentación**: Esta documentación

---

**¿Necesitas ayuda?** Comienza con la [Guía de Uso](USAGE.md) para ejemplos prácticos, o consulta la [Documentación Técnica](ARCHITECTURE.md) para entender la arquitectura interna.
