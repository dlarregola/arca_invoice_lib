package soap

import (
	"arca_invoice_lib/pkg/models"
	"bytes"
	"context"
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

// Client representa un cliente SOAP
type Client struct {
	httpClient *http.Client
	logger     *logrus.Logger
	baseURL    string
}

// NewClient crea un nuevo cliente SOAP
func NewClient(baseURL string, timeout time.Duration, logger *logrus.Logger) *Client {
	httpClient := &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: false,
			},
		},
	}

	return &Client{
		httpClient: httpClient,
		logger:     logger,
		baseURL:    baseURL,
	}
}

// Call realiza una llamada SOAP
func (c *Client) Call(ctx context.Context, action string, request interface{}, response interface{}) error {
	// Serializar request a XML
	requestXML, err := xml.MarshalIndent(request, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling request: %w", err)
	}

	// Crear envelope SOAP
	envelope := &SOAPEnvelope{
		XMLName: xml.Name{Space: "http://schemas.xmlsoap.org/soap/envelope/", Local: "Envelope"},
		Header:  &SOAPHeader{},
		Body: SOAPBody{
			Content: requestXML,
		},
	}

	// Serializar envelope
	envelopeXML, err := xml.MarshalIndent(envelope, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling envelope: %w", err)
	}

	// Log request si está habilitado
	if c.logger.GetLevel() >= logrus.DebugLevel {
		c.logger.WithFields(logrus.Fields{
			"action": action,
			"url":    c.baseURL,
		}).Debug("SOAP Request")
		c.logger.Debug(string(envelopeXML))
	}

	// Crear request HTTP
	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL, bytes.NewReader(envelopeXML))
	if err != nil {
		return fmt.Errorf("error creating HTTP request: %w", err)
	}

	// Configurar headers
	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	req.Header.Set("SOAPAction", action)
	req.Header.Set("User-Agent", "AFIP-Go-Client/1.0")

	// Realizar request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return models.NewNetworkError(fmt.Sprintf("error making HTTP request: %v", err), c.baseURL, 0)
	}
	defer resp.Body.Close()

	// Leer response
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.NewNetworkError(fmt.Sprintf("error reading response body: %v", err), c.baseURL, resp.StatusCode)
	}

	// Log response si está habilitado
	if c.logger.GetLevel() >= logrus.DebugLevel {
		c.logger.WithFields(logrus.Fields{
			"status_code": resp.StatusCode,
			"action":      action,
		}).Debug("SOAP Response")
		c.logger.Debug(string(responseBody))
	}

	// Verificar status code
	if resp.StatusCode != http.StatusOK {
		return models.NewNetworkError(fmt.Sprintf("HTTP error: %s", resp.Status), c.baseURL, resp.StatusCode)
	}

	// Parsear response SOAP
	var responseEnvelope SOAPEnvelope
	if err := xml.Unmarshal(responseBody, &responseEnvelope); err != nil {
		return models.NewAFIPError(models.ErrorCodeInvalidResponse, fmt.Sprintf("error unmarshaling SOAP response: %v", err))
	}

	// Verificar si hay error SOAP
	if responseEnvelope.Body.Fault != nil {
		fault := responseEnvelope.Body.Fault
		return models.NewAFIPError(fault.FaultCode, fault.FaultString)
	}

	// Parsear contenido de respuesta
	if err := xml.Unmarshal(responseEnvelope.Body.Content, response); err != nil {
		return models.NewAFIPError(models.ErrorCodeInvalidResponse, fmt.Sprintf("error unmarshaling response content: %v", err))
	}

	return nil
}

// SOAPEnvelope representa un envelope SOAP
type SOAPEnvelope struct {
	XMLName xml.Name    `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
	Header  *SOAPHeader `xml:"Header,omitempty"`
	Body    SOAPBody    `xml:"Body"`
}

// SOAPHeader representa el header SOAP
type SOAPHeader struct {
	XMLName xml.Name `xml:"Header"`
}

// SOAPBody representa el body SOAP
type SOAPBody struct {
	XMLName xml.Name   `xml:"Body"`
	Content []byte     `xml:",innerxml"`
	Fault   *SOAPFault `xml:"Fault,omitempty"`
}

// SOAPFault representa un fault SOAP
type SOAPFault struct {
	XMLName     xml.Name `xml:"Fault"`
	FaultCode   string   `xml:"faultcode"`
	FaultString string   `xml:"faultstring"`
	Detail      string   `xml:"detail,omitempty"`
}

// SetTimeout actualiza el timeout del cliente HTTP
func (c *Client) SetTimeout(timeout time.Duration) {
	c.httpClient.Timeout = timeout
}

// SetLogger actualiza el logger del cliente
func (c *Client) SetLogger(logger *logrus.Logger) {
	c.logger = logger
}
