package client

import (
	"bytes"
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

// WSAAAuth maneja la autenticación con el Web Service de Autenticación y Autorización
type WSAAAuth struct {
	config     *Config
	cache      map[string]*AccessTicket
	cacheMutex sync.RWMutex
	logger     interface{}
}

// AccessTicket representa un ticket de acceso de ARCA
type AccessTicket struct {
	Token          string    `xml:"token"`
	Sign           string    `xml:"sign"`
	ExpirationTime time.Time `xml:"expirationTime"`
	GenerationTime time.Time `xml:"generationTime"`
}

// WSAARequest representa el request para WSAA
type WSAARequest struct {
	XMLName xml.Name `xml:"loginTicketRequest"`
	Version string   `xml:"version"`
	Header  struct {
		Source         string `xml:"source"`
		Destination    string `xml:"destination"`
		UniqueID       string `xml:"uniqueId"`
		GenerationTime string `xml:"generationTime"`
		ExpirationTime string `xml:"expirationTime"`
	} `xml:"header"`
	Service string `xml:"service"`
}

// WSAAResponse representa la respuesta de WSAA
type WSAAResponse struct {
	XMLName xml.Name `xml:"loginTicketResponse"`
	Header  struct {
		Source         string `xml:"source"`
		Destination    string `xml:"destination"`
		UniqueID       string `xml:"uniqueId"`
		GenerationTime string `xml:"generationTime"`
		ExpirationTime string `xml:"expirationTime"`
	} `xml:"header"`
	Credentials struct {
		Token string `xml:"token"`
		Sign  string `xml:"sign"`
	} `xml:"credentials"`
}

// NewWSAAAuth crea un nuevo autenticador WSAA
func NewWSAAAuth(config *Config, logger interface{}) *WSAAAuth {
	return &WSAAAuth{
		config: config,
		cache:  make(map[string]*AccessTicket),
		logger: logger,
	}
}

// GetAccessTicket obtiene un ticket de acceso válido
func (a *WSAAAuth) GetAccessTicket(ctx context.Context, service string) (*AccessTicket, error) {
	// Verificar cache primero
	if ticket := a.getFromCache(service); ticket != nil {
		return ticket, nil
	}

	// Generar nuevo ticket
	return a.generateAccessTicket(ctx, service)
}

// getFromCache obtiene un ticket del cache
func (a *WSAAAuth) getFromCache(service string) *AccessTicket {
	a.cacheMutex.RLock()
	defer a.cacheMutex.RUnlock()

	ticket, exists := a.cache[service]
	if !exists {
		return nil
	}

	// Verificar si el ticket aún es válido (con margen de 5 minutos)
	if time.Now().Add(5 * time.Minute).Before(ticket.ExpirationTime) {
		return ticket
	}

	// Ticket expirado, remover del cache
	delete(a.cache, service)
	return nil
}

// addToCache agrega un ticket al cache
func (a *WSAAAuth) addToCache(service string, ticket *AccessTicket) {
	a.cacheMutex.Lock()
	defer a.cacheMutex.Unlock()

	a.cache[service] = ticket
}

// generateAccessTicket genera un nuevo ticket de acceso
func (a *WSAAAuth) generateAccessTicket(ctx context.Context, service string) (*AccessTicket, error) {
	// Parsear certificado
	cert, err := x509.ParseCertificate(a.config.Certificate)
	if err != nil {
		return nil, fmt.Errorf("error parsing certificate: %v", err)
	}

	// Parsear clave privada
	var privateKey *rsa.PrivateKey
	parsedKey, err := x509.ParsePKCS1PrivateKey(a.config.PrivateKey)
	if err != nil {
		// Intentar con PKCS8
		key, err := x509.ParsePKCS8PrivateKey(a.config.PrivateKey)
		if err != nil {
			return nil, fmt.Errorf("error parsing private key: %v", err)
		}
		parsedKey, ok := key.(*rsa.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("private key is not RSA")
		}
		privateKey = parsedKey
	} else {
		privateKey = parsedKey
	}

	// Generar unique ID
	uniqueID, err := generateUniqueID()
	if err != nil {
		return nil, fmt.Errorf("error generating unique ID: %v", err)
	}

	// Crear request
	request := &WSAARequest{
		Version: "1.0",
		Service: service,
	}
	request.Header.Source = a.config.CUIT
	request.Header.Destination = "cn=wsaahomo,o=afip,c=ar,serialNumber=CUIT 33693450239"
	request.Header.UniqueID = uniqueID
	request.Header.GenerationTime = time.Now().UTC().Format("2006-01-02T15:04:05.000-07:00")
	request.Header.ExpirationTime = time.Now().Add(24 * time.Hour).UTC().Format("2006-01-02T15:04:05.000-07:00")

	// Serializar request
	requestXML, err := xml.MarshalIndent(request, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}

	// Crear CMS (Cryptographic Message Syntax)
	cms, err := a.createCMS(requestXML, cert, privateKey)
	if err != nil {
		return nil, fmt.Errorf("error creating CMS: %v", err)
	}

	// Realizar request a WSAA
	response, err := a.callWSAA(ctx, cms)
	if err != nil {
		return nil, err
	}

	// Parsear respuesta
	var wsaaResponse WSAAResponse
	if err := xml.Unmarshal([]byte(response), &wsaaResponse); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	// Crear ticket
	ticket := &AccessTicket{
		Token:          wsaaResponse.Credentials.Token,
		Sign:           wsaaResponse.Credentials.Sign,
		GenerationTime: time.Now(),
		ExpirationTime: time.Now().Add(24 * time.Hour),
	}

	// Agregar al cache
	a.addToCache(service, ticket)

	return ticket, nil
}

// createCMS crea un mensaje CMS firmado
func (a *WSAAAuth) createCMS(data []byte, cert *x509.Certificate, privateKey *rsa.PrivateKey) (string, error) {
	// Crear hash SHA1 del data
	hash := sha1.Sum(data)

	// Firmar el hash
	_, err := rsa.SignPKCS1v15(nil, privateKey, crypto.SHA1, hash[:])
	if err != nil {
		return "", err
	}

	// Crear estructura CMS simplificada
	cms := fmt.Sprintf(`<cm:loginTicketRequest xmlns:cm="http://wsaa.view.sua.dvadac.desein.afip.gov">
%s
</cm:loginTicketRequest>`, string(data))

	// Codificar en base64
	return base64.StdEncoding.EncodeToString([]byte(cms)), nil
}

// callWSAA realiza la llamada al servicio WSAA
func (a *WSAAAuth) callWSAA(ctx context.Context, cms string) (string, error) {
	// Crear request HTTP
	requestBody := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:wsaa="http://wsaa.view.sua.dvadac.desein.afip.gov">
   <soapenv:Header/>
   <soapenv:Body>
      <wsaa:loginCms>
         <wsaa:in0>%s</wsaa:in0>
      </wsaa:loginCms>
   </soapenv:Body>
</soapenv:Envelope>`, cms)

	// Crear request HTTP
	req, err := http.NewRequestWithContext(ctx, "POST", a.config.GetWSAAURL(), bytes.NewReader([]byte(requestBody)))
	if err != nil {
		return "", fmt.Errorf("error creating HTTP request: %v", err)
	}

	// Configurar headers
	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	req.Header.Set("SOAPAction", "http://wsaa.view.sua.dvadac.desein.afip.gov/loginCms")
	req.Header.Set("User-Agent", "ARCA-Go-Client/1.0")

	// Realizar request
	client := &http.Client{Timeout: a.config.Timeout}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// Leer response
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	// Verificar status code
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP error: %s", resp.Status)
	}

	// Parsear respuesta SOAP
	var soapResponse struct {
		XMLName xml.Name `xml:"Envelope"`
		Body    struct {
			LoginCmsResponse struct {
				LoginCmsReturn string `xml:"loginCmsReturn"`
			} `xml:"loginCmsResponse"`
		} `xml:"Body"`
	}

	if err := xml.Unmarshal(responseBody, &soapResponse); err != nil {
		return "", fmt.Errorf("error unmarshaling SOAP response: %v", err)
	}

	return soapResponse.Body.LoginCmsResponse.LoginCmsReturn, nil
}

// generateUniqueID genera un ID único
func generateUniqueID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", bytes), nil
}

// ClearCache limpia el cache de tickets
func (a *WSAAAuth) ClearCache() {
	a.cacheMutex.Lock()
	defer a.cacheMutex.Unlock()

	a.cache = make(map[string]*AccessTicket)
}

// GetCacheSize retorna el tamaño del cache
func (a *WSAAAuth) GetCacheSize() int {
	a.cacheMutex.RLock()
	defer a.cacheMutex.RUnlock()

	return len(a.cache)
}
