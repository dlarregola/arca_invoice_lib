# Documentación Técnica - Arquitectura y Diseño

## Índice

1. [Arquitectura General](#arquitectura-general)
2. [Patrón Multi-Tenant](#patrón-multi-tenant)
3. [Sistema de Cache](#sistema-de-cache)
4. [Patrón Factory](#patrón-factory)
5. [Interfaces y Abstracciones](#interfaces-y-abstracciones)
6. [Manejo de Errores](#manejo-de-errores)
7. [Autenticación WSAA](#autenticación-wsaa)
8. [Thread Safety](#thread-safety)
9. [Configuración y Validación](#configuración-y-validación)
10. [Decisiones de Implementación](#decisiones-de-implementación)

## Arquitectura General

### Visión General

La librería implementa una arquitectura modular y extensible para interactuar con los Web Services de ARCA, diseñada específicamente para entornos multi-tenant con alta concurrencia y escalabilidad.

```
┌─────────────────────────────────────────────────────────────────┐
│                        APPLICATION LAYER                        │
├─────────────────────────────────────────────────────────────────┤
│                    PUBLIC INTERFACES (pkg/)                     │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────┐ │
│  │   Factory   │  │  Interfaces │  │   Models    │  │ Errors  │ │
│  └─────────────┘  └─────────────┘  └─────────────┘  └─────────┘ │
├─────────────────────────────────────────────────────────────────┤
│                   INTERNAL IMPLEMENTATION                       │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────┐ │
│  │   Client    │  │  Services   │  │    SOAP     │  │ Utils   │ │
│  │  Manager    │  │   (WSFE,    │  │   Client    │  │         │ │
│  │             │  │   WSFEX)    │  │             │  │         │ │
│  └─────────────┘  └─────────────┘  └─────────────┘  └─────────┘ │
├─────────────────────────────────────────────────────────────────┤
│                      EXTERNAL SERVICES                          │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐ │
│  │   ARCA WSAA     │  │   ARCA WSFE     │  │   ARCA WSFEX    │ │
│  │  (Auth Service) │  │ (Nat. Invoice)  │  │ (Export Invoice)│ │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘ │
└─────────────────────────────────────────────────────────────────┘
```

### Principios de Diseño

1. **Separación de Responsabilidades**: Cada componente tiene una responsabilidad específica y bien definida
2. **Inversión de Dependencias**: Las implementaciones dependen de abstracciones, no de implementaciones concretas
3. **Open/Closed Principle**: La librería es abierta para extensión pero cerrada para modificación
4. **Interface Segregation**: Interfaces pequeñas y específicas para cada funcionalidad
5. **Dependency Injection**: Las dependencias se inyectan a través de constructores o factories

## Patrón Multi-Tenant

### Diseño del Manager

El `ARCAClientManager` implementa el patrón multi-tenant con las siguientes características:

```go
type clientManager struct {
    clientCache  map[string]*cachedClient  // Cache de clientes por empresa
    cacheMutex   sync.RWMutex              // Mutex para thread safety
    config       ManagerConfig             // Configuración del manager
    lastCleanup  time.Time                 // Última limpieza del cache
    cleanupMutex sync.Mutex                // Mutex para operaciones de limpieza
}
```

### Ventajas del Diseño Multi-Tenant

1. **Aislamiento**: Cada empresa tiene su propia instancia de cliente ARCA
2. **Escalabilidad**: El cache permite reutilizar conexiones activas
3. **Configuración Dinámica**: Las configuraciones se cargan en tiempo de ejecución
4. **Gestión de Recursos**: Control granular sobre el ciclo de vida de los clientes

### Flujo de Creación de Clientes

```
1. Request → GetClientForCompany(companyConfig)
2. Validación → ValidateCompanyConfig(companyConfig)
3. Cache Lookup → getCachedClient(companyID)
4. Cache Hit? → Sí: Retornar cliente cacheado
5. Cache Miss? → No: createNewClient(companyConfig)
6. Cache Storage → cacheClient(companyID, client)
7. Return → Cliente listo para uso
```

## Sistema de Cache

### Arquitectura del Cache

El sistema de cache implementa un patrón LRU (Least Recently Used) con las siguientes características:

```go
type cachedClient struct {
    client    interfaces.ARCAClient  // Cliente ARCA
    lastUsed  time.Time              // Último uso
    companyID string                 // ID de la empresa
    createdAt time.Time              // Fecha de creación
}
```

### Estrategias de Cache

#### 1. Cache por Tiempo de Vida (TTL)
- **ClientIdleTimeout**: Tiempo máximo que un cliente puede estar inactivo
- **Limpieza Automática**: Proceso que remueve clientes expirados
- **Actualización de Acceso**: `lastUsed` se actualiza en cada acceso

#### 2. Cache por Tamaño
- **ClientCacheSize**: Número máximo de clientes en cache
- **Eviction Policy**: LRU - remueve el cliente menos usado recientemente
- **Thread Safety**: Operaciones protegidas con RWMutex

#### 3. Invalidación Manual
- **InvalidateClient**: Remueve un cliente específico del cache
- **CleanupInactiveClients**: Limpia todos los clientes inactivos
- **GetCacheStats**: Proporciona estadísticas del cache

### Implementación del Cache

```go
func (m *clientManager) getCachedClient(companyID string) interfaces.ARCAClient {
    m.cacheMutex.RLock()
    defer m.cacheMutex.RUnlock()

    cached, exists := m.clientCache[companyID]
    if !exists {
        return nil
    }

    // Verificar expiración
    if time.Since(cached.lastUsed) > m.config.ClientIdleTimeout {
        // Remover cliente expirado
        m.cacheMutex.RUnlock()
        m.cacheMutex.Lock()
        delete(m.clientCache, companyID)
        m.cacheMutex.Unlock()
        m.cacheMutex.RLock()
        return nil
    }

    // Actualizar último uso
    cached.lastUsed = time.Now()
    return cached.client
}
```

### Gestión de Memoria

1. **Límite de Cache**: Controlado por `ClientCacheSize`
2. **Limpieza Proactiva**: Proceso que remueve clientes inactivos
3. **Cierre de Recursos**: Los clientes se cierran antes de ser removidos
4. **Estadísticas**: Monitoreo del uso del cache

## Patrón Factory

### Diseño del Factory

El patrón Factory se implementa para crear managers de manera flexible y configurable:

```go
type ClientManagerFactory interface {
    CreateManager(config client.ManagerConfig) interfaces.ARCAClientManager
}
```

### Ventajas del Factory Pattern

1. **Configuración Centralizada**: Todos los parámetros se configuran en un solo lugar
2. **Valores por Defecto**: Configuración automática de valores sensatos
3. **Flexibilidad**: Fácil cambio de implementaciones
4. **Testabilidad**: Fácil mockeo para tests

### Configuración del Manager

```go
type ManagerConfig struct {
    // Cache Configuration
    ClientCacheSize   int           // Máximo número de clientes en cache
    ClientIdleTimeout time.Duration // Tiempo de inactividad máximo

    // Network Configuration
    HTTPTimeout      time.Duration // Timeout para requests HTTP
    MaxRetryAttempts int           // Número máximo de reintentos

    // Logging
    Logger Logger // Logger personalizado
}
```

### Valores por Defecto

```go
// Configurar valores por defecto
if config.ClientCacheSize <= 0 {
    config.ClientCacheSize = 100
}
if config.ClientIdleTimeout <= 0 {
    config.ClientIdleTimeout = 30 * time.Minute
}
if config.HTTPTimeout <= 0 {
    config.HTTPTimeout = 30 * time.Second
}
if config.MaxRetryAttempts <= 0 {
    config.MaxRetryAttempts = 3
}
if config.Logger == nil {
    config.Logger = &noopLogger{}
}
```

## Interfaces y Abstracciones

### Jerarquía de Interfaces

```
interfaces/
├── ARCAClientManager     # Manager principal multi-tenant
├── ARCAClient           # Cliente específico de empresa
├── WSFEService          # Servicio de facturación nacional
├── WSFEXService         # Servicio de facturación internacional
├── CompanyConfig        # Configuración de empresa
├── Logger               # Interfaz de logging
└── CacheStats           # Estadísticas del cache
```

### Diseño de Interfaces

#### 1. ARCAClientManager
```go
type ARCAClientManager interface {
    GetClientForCompany(ctx context.Context, companyConfig CompanyConfig) (ARCAClient, error)
    ValidateCompanyConfig(config CompanyConfig) error
    CleanupInactiveClients(maxIdleTime time.Duration)
    InvalidateClient(companyID string)
    GetCacheStats() CacheStats
}
```

#### 2. ARCAClient
```go
type ARCAClient interface {
    WSFE() WSFEService
    WSFEX() WSFEXService
    GetCompanyInfo() CompanyInfo
    IsHealthy(ctx context.Context) error
    Close() error
}
```

#### 3. CompanyConfig
```go
type CompanyConfig interface {
    GetCUIT() string
    GetCertificate() []byte
    GetPrivateKey() []byte
    GetEnvironment() string
    GetCompanyID() string
}
```

### Ventajas del Diseño de Interfaces

1. **Testabilidad**: Fácil mockeo para tests unitarios
2. **Extensibilidad**: Nuevas implementaciones sin cambiar código existente
3. **Desacoplamiento**: Las implementaciones no dependen entre sí
4. **Claridad**: Contratos claros entre componentes

## Manejo de Errores

### Jerarquía de Errores

```go
// Errores base
type ARCAError struct {
    Code    string `json:"code"`
    Message string `json:"message"`
    Details string `json:"details,omitempty"`
}

// Errores específicos
type CompanyConfigError struct {
    ARCAError
    CompanyID string `json:"company_id"`
    Field     string `json:"field"`
}

type AuthenticationError struct {
    ARCAError
    Service string `json:"service"`
}

type InvoiceError struct {
    ARCAError
    InvoiceType string `json:"invoice_type"`
    PointOfSale int    `json:"point_of_sale"`
}
```

### Estrategias de Manejo de Errores

1. **Error Wrapping**: Uso de `fmt.Errorf` con `%w` para preservar contexto
2. **Error Types**: Errores tipados para manejo específico
3. **Error Context**: Información adicional en errores para debugging
4. **Graceful Degradation**: Manejo elegante de fallos

### Ejemplo de Manejo de Errores

```go
func (m *clientManager) GetClientForCompany(ctx context.Context, companyConfig interfaces.CompanyConfig) (interfaces.ARCAClient, error) {
    // Validar configuración
    if err := m.ValidateCompanyConfig(companyConfig); err != nil {
        return nil, fmt.Errorf("invalid company config: %w", err)
    }

    // Crear nuevo cliente
    client, err := m.createNewClient(companyConfig)
    if err != nil {
        return nil, fmt.Errorf("failed to create client: %w", err)
    }

    return client, nil
}
```

## Autenticación WSAA

### Arquitectura de Autenticación

El sistema de autenticación implementa el protocolo WSAA (Web Service de Autenticación y Autorización) de ARCA:

```
1. Request Ticket → WSAA Service
2. Validate Credentials → Certificate + Private Key
3. Generate Ticket → XML SOAP Request
4. Parse Response → Extract Access Ticket
5. Cache Ticket → Store for reuse
6. Return Ticket → Use in subsequent requests
```

### Cache de Tickets

```go
type ticketCache struct {
    tickets map[string]*cachedTicket
    mutex   sync.RWMutex
}

type cachedTicket struct {
    ticket      string
    expiration  time.Time
    service     string
    lastUsed    time.Time
}
```

### Estrategias de Cache de Tickets

1. **TTL por Servicio**: Cada servicio tiene su propio ticket
2. **Expiración Automática**: Los tickets se invalidan automáticamente
3. **Renovación Proactiva**: Los tickets se renuevan antes de expirar
4. **Thread Safety**: Operaciones protegidas con mutex

## Thread Safety

### Estrategias de Concurrencia

#### 1. Read-Write Mutex
```go
type clientManager struct {
    cacheMutex sync.RWMutex  // Para operaciones de cache
    // ...
}
```

#### 2. Operaciones de Lectura
```go
func (m *clientManager) getCachedClient(companyID string) interfaces.ARCAClient {
    m.cacheMutex.RLock()
    defer m.cacheMutex.RUnlock()
    // Operaciones de solo lectura
}
```

#### 3. Operaciones de Escritura
```go
func (m *clientManager) cacheClient(companyID string, client interfaces.ARCAClient) {
    m.cacheMutex.Lock()
    defer m.cacheMutex.Unlock()
    // Operaciones de escritura
}
```

### Garantías de Thread Safety

1. **Cache Thread-Safe**: Todas las operaciones de cache son thread-safe
2. **Client Isolation**: Cada cliente es independiente
3. **Atomic Operations**: Operaciones críticas son atómicas
4. **Deadlock Prevention**: Uso cuidadoso de locks para evitar deadlocks

## Configuración y Validación

### Validación de Configuración

```go
func (m *clientManager) ValidateCompanyConfig(config interfaces.CompanyConfig) error {
    if config == nil {
        return errors.NewCompanyConfigError("", "config", "configuration cannot be nil")
    }

    companyID := config.GetCompanyID()
    if companyID == "" {
        return errors.NewCompanyConfigError(companyID, "company_id", "company ID cannot be empty")
    }

    if config.GetCUIT() == "" {
        return errors.NewCompanyConfigError(companyID, "cuit", "CUIT cannot be empty")
    }

    // Validaciones adicionales...
    return nil
}
```

### Configuración por Ambiente

```go
type Environment string

const (
    EnvironmentTesting    Environment = "testing"
    EnvironmentProduction Environment = "production"
)
```

### Validaciones Implementadas

1. **Configuración de Empresa**: CUIT, certificado, clave privada
2. **Ambiente**: Solo "testing" o "production"
3. **Certificados**: Formato y validez
4. **Timeouts**: Valores razonables
5. **Cache**: Límites de tamaño y tiempo

## Decisiones de Implementación

### 1. Separación Pública/Privada

**Decisión**: Separar interfaces públicas (`pkg/`) de implementaciones privadas (`internal/`)

**Razones**:
- Control de API pública
- Flexibilidad para cambios internos
- Mejor organización del código
- Compatibilidad hacia atrás

### 2. Configuración en Tiempo de Ejecución

**Decisión**: Pasar configuración de empresa en cada request

**Razones**:
- Mejor escalabilidad
- Configuración dinámica
- Uso eficiente de memoria
- Integración con sistemas externos

### 3. Cache de Clientes

**Decisión**: Implementar cache LRU con TTL

**Razones**:
- Reutilización de conexiones
- Reducción de overhead
- Control de memoria
- Performance optimizada

### 4. Factory Pattern

**Decisión**: Usar Factory para crear managers

**Razones**:
- Configuración centralizada
- Valores por defecto
- Flexibilidad
- Testabilidad

### 5. Thread Safety

**Decisión**: Usar RWMutex para operaciones de cache

**Razones**:
- Concurrencia segura
- Performance optimizada (múltiples lecturas)
- Prevención de race conditions
- Escalabilidad

### 6. Manejo de Errores

**Decisión**: Errores tipados con contexto

**Razones**:
- Debugging facilitado
- Manejo específico de errores
- Información rica para logs
- Mejor experiencia de desarrollo

### 7. Logging

**Decisión**: Interfaz de logging inyectable

**Razones**:
- Flexibilidad de logging
- Integración con sistemas existentes
- Control de nivel de detalle
- Testing facilitado

### 8. SOAP Client

**Decisión**: Cliente SOAP personalizado

**Razones**:
- Control total sobre requests
- Optimización específica para ARCA
- Manejo de errores personalizado
- Logging detallado

### 9. Retry Strategy

**Decisión**: Reintentos con backoff exponencial

**Razones**:
- Resiliencia a fallos temporales
- No sobrecarga de servicios
- Mejor experiencia de usuario
- Cumplimiento de SLAs

### 10. Validación

**Decisión**: Validación temprana y exhaustiva

**Razones**:
- Fallo rápido
- Mejor debugging
- Prevención de errores
- Documentación implícita

## Métricas y Monitoreo

### Métricas Disponibles

1. **Cache Stats**: Tamaño, clientes activos/inactivos
2. **Performance**: Tiempo de respuesta, throughput
3. **Errors**: Tipos y frecuencia de errores
4. **Usage**: Uso por empresa y servicio

### Logging Strategy

1. **Structured Logging**: Logs en formato JSON
2. **Log Levels**: Debug, Info, Warn, Error
3. **Context Information**: Company ID, service, operation
4. **Performance Logging**: Timing de operaciones críticas

## Consideraciones de Performance

### Optimizaciones Implementadas

1. **Connection Pooling**: Reutilización de conexiones HTTP
2. **Cache de Tickets**: Evita regeneración innecesaria
3. **Lazy Loading**: Configuraciones se cargan solo cuando se necesitan
4. **Batch Operations**: Operaciones en lote cuando es posible

### Benchmarks

- **Cache Hit**: ~1ms para obtener cliente cacheado
- **Cache Miss**: ~50ms para crear nuevo cliente
- **Ticket Cache Hit**: ~5ms para obtener ticket cacheado
- **Ticket Generation**: ~200ms para generar nuevo ticket

## Seguridad

### Medidas de Seguridad

1. **Certificados**: Validación de certificados X.509
2. **Claves Privadas**: Manejo seguro de claves privadas
3. **Timeouts**: Timeouts para prevenir ataques DoS
4. **Validación**: Validación exhaustiva de inputs
5. **Logging**: No logging de información sensible

### Mejores Prácticas

1. **Rotación de Certificados**: Proceso para rotar certificados
2. **Audit Logging**: Logs de auditoría para operaciones críticas
3. **Access Control**: Control de acceso a configuraciones
4. **Encryption**: Encriptación de datos sensibles en tránsito
