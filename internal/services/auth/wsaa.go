package auth

import (
	"arca_invoice_lib/internal/shared"
	"arca_invoice_lib/pkg/interfaces"
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

// wsaaService es la implementación privada del servicio de autenticación
type wsaaService struct {
	config     *shared.InternalConfig
	cache      map[string]*interfaces.AccessToken
	cacheMutex sync.RWMutex
	logger     interfaces.Logger
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

// newAuthService crea un nuevo servicio de autenticación
func newAuthService(config *shared.InternalConfig, logger interfaces.Logger) interfaces.AuthService {
	return &wsaaService{
		config: config,
		cache:  make(map[string]*interfaces.AccessToken),
		logger: logger,
	}
}

// GetToken obtiene un token de autenticación válido
func (s *wsaaService) GetToken(ctx context.Context, service string) (*interfaces.AccessToken, error) {
	// Verificar cache primero
	if token := s.getFromCache(service); token != nil {
		return token, nil
	}

	// Generar nuevo token
	return s.generateAccessToken(ctx, service)
}

// ClearCache limpia el cache de tokens
func (s *wsaaService) ClearCache() {
	s.cacheMutex.Lock()
	defer s.cacheMutex.Unlock()

	s.cache = make(map[string]*interfaces.AccessToken)
	s.logger.Debug("Auth cache cleared")
}

// GetCacheSize retorna el tamaño del cache
func (s *wsaaService) GetCacheSize() int {
	s.cacheMutex.RLock()
	defer s.cacheMutex.RUnlock()

	return len(s.cache)
}

// getFromCache obtiene un token del cache
func (s *wsaaService) getFromCache(service string) *interfaces.AccessToken {
	s.cacheMutex.RLock()
	defer s.cacheMutex.RUnlock()

	token, exists := s.cache[service]
	if !exists {
		return nil
	}

	// Verificar si el token aún es válido (con margen de 5 minutos)
	if time.Now().Add(5 * time.Minute).Before(token.ExpirationTime) {
		return token
	}

	// Token expirado, remover del cache
	delete(s.cache, service)
	return nil
}

// addToCache agrega un token al cache
func (s *wsaaService) addToCache(service string, token *interfaces.AccessToken) {
	s.cacheMutex.Lock()
	defer s.cacheMutex.Unlock()

	s.cache[service] = token
}

// generateAccessToken genera un nuevo token de acceso
func (s *wsaaService) generateAccessToken(ctx context.Context, service string) (*interfaces.AccessToken, error) {
	// Parsear certificado
	cert, err := x509.ParseCertificate(s.config.Certificate)
	if err != nil {
		return nil, fmt.Errorf("error parsing certificate: %v", err)
	}

	// Parsear clave privada
	var privateKey *rsa.PrivateKey
	parsedKey, err := x509.ParsePKCS1PrivateKey(s.config.PrivateKey)
	if err != nil {
		// Intentar con PKCS8
		key, err := x509.ParsePKCS8PrivateKey(s.config.PrivateKey)
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
	request.Header.Source = s.config.CUIT
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
	cms, err := s.createCMS(requestXML, cert, privateKey)
	if err != nil {
		return nil, fmt.Errorf("error creating CMS: %v", err)
	}

	// Realizar request a WSAA
	response, err := s.callWSAA(ctx, cms)
	if err != nil {
		return nil, err
	}

	// Parsear respuesta
	var wsaaResponse WSAAResponse
	if err := xml.Unmarshal([]byte(response), &wsaaResponse); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	// Crear token
	token := &interfaces.AccessToken{
		Token:          wsaaResponse.Credentials.Token,
		Sign:           wsaaResponse.Credentials.Sign,
		GenerationTime: time.Now(),
		ExpirationTime: time.Now().Add(24 * time.Hour),
	}

	// Agregar al cache
	s.addToCache(service, token)

	s.logger.Infof("Generated new access token for service %s", service)
	return token, nil
}

// createCMS crea un mensaje CMS firmado
func (s *wsaaService) createCMS(data []byte, cert *x509.Certificate, privateKey *rsa.PrivateKey) (string, error) {
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
func (s *wsaaService) callWSAA(ctx context.Context, cms string) (string, error) {
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
	req, err := http.NewRequestWithContext(ctx, "POST", s.getWSAAURL(), bytes.NewReader([]byte(requestBody)))
	if err != nil {
		return "", fmt.Errorf("error creating HTTP request: %v", err)
	}

	// Configurar headers
	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	req.Header.Set("SOAPAction", "http://wsaa.view.sua.dvadac.desein.afip.gov/loginCms")
	req.Header.Set("User-Agent", "AFIP-Go-Client/1.0")

	// Realizar request
	client := &http.Client{Timeout: s.config.Timeout}
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

// getWSAAURL retorna la URL del servicio WSAA
func (s *wsaaService) getWSAAURL() string {
	baseURL := s.getBaseURL()
	return baseURL + "/ws/services/LoginCms"
}

// getBaseURL retorna la URL base según el environment
func (s *wsaaService) getBaseURL() string {
	switch s.config.Environment {
	case "testing":
		return "https://wswhomo.afip.gov.ar"
	case "production":
		return "https://servicios1.afip.gov.ar"
	default:
		return "https://wswhomo.afip.gov.ar"
	}
}

// generateUniqueID genera un ID único
func generateUniqueID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", bytes), nil
}
