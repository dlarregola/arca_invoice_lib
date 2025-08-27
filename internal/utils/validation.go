package utils

import (
	"github.com/dlarregola/arca_invoice_lib/pkg/models"
	"fmt"
	"regexp"
	"strconv"
	"time"
)

// ValidateCUIT valida el formato de un CUIT
func ValidateCUIT(cuit string) error {
	if cuit == "" {
		return models.NewValidationError("cuit", "CUIT no puede estar vacío", cuit)
	}

	// Validar formato: XX-XXXXXXXX-X
	re := regexp.MustCompile(`^\d{2}-\d{8}-\d$`)
	if !re.MatchString(cuit) {
		return models.NewValidationError("cuit", "CUIT debe tener formato XX-XXXXXXXX-X", cuit)
	}

	// Validar dígito verificador
	if !validateCUITCheckDigit(cuit) {
		return models.NewValidationError("cuit", "CUIT inválido - dígito verificador incorrecto", cuit)
	}

	return nil
}

// validateCUITCheckDigit valida el dígito verificador del CUIT
func validateCUITCheckDigit(cuit string) bool {
	// Remover guiones
	cuit = regexp.MustCompile(`-`).ReplaceAllString(cuit, "")

	if len(cuit) != 11 {
		return false
	}

	// Factores para el cálculo
	factors := []int{5, 4, 3, 2, 7, 6, 5, 4, 3, 2}

	var sum int
	for i := 0; i < 10; i++ {
		digit, _ := strconv.Atoi(string(cuit[i]))
		sum += digit * factors[i]
	}

	remainder := sum % 11
	var expectedDigit int

	if remainder == 0 {
		expectedDigit = 0
	} else if remainder == 1 {
		// Casos especiales
		firstTwo := cuit[:2]
		if firstTwo == "20" || firstTwo == "23" || firstTwo == "24" || firstTwo == "27" || firstTwo == "30" || firstTwo == "33" || firstTwo == "34" {
			expectedDigit = 9
		} else {
			expectedDigit = 0
		}
	} else {
		expectedDigit = 11 - remainder
	}

	actualDigit, _ := strconv.Atoi(string(cuit[10]))
	return actualDigit == expectedDigit
}

// ValidateDocumentNumber valida un número de documento
func ValidateDocumentNumber(docType models.DocumentType, docNumber string) error {
	if docNumber == "" {
		return models.NewValidationError("doc_number", "Número de documento no puede estar vacío", docNumber)
	}

	switch docType {
	case models.DocumentTypeDNI:
		if len(docNumber) != 7 && len(docNumber) != 8 {
			return models.NewValidationError("doc_number", "DNI debe tener 7 u 8 dígitos", docNumber)
		}
	case models.DocumentTypeCUIT, models.DocumentTypeCUIL:
		return ValidateCUIT(docNumber)
	case models.DocumentTypeCDI:
		if len(docNumber) < 5 || len(docNumber) > 20 {
			return models.NewValidationError("doc_number", "CDI debe tener entre 5 y 20 caracteres", docNumber)
		}
	case models.DocumentTypeLE, models.DocumentTypeLC:
		if len(docNumber) < 5 || len(docNumber) > 20 {
			return models.NewValidationError("doc_number", "LE/LC debe tener entre 5 y 20 caracteres", docNumber)
		}
	case models.DocumentTypeCI:
		if len(docNumber) < 5 || len(docNumber) > 20 {
			return models.NewValidationError("doc_number", "CI debe tener entre 5 y 20 caracteres", docNumber)
		}
	case models.DocumentTypePAS:
		if len(docNumber) < 5 || len(docNumber) > 20 {
			return models.NewValidationError("doc_number", "Pasaporte debe tener entre 5 y 20 caracteres", docNumber)
		}
	default:
		if len(docNumber) < 1 || len(docNumber) > 20 {
			return models.NewValidationError("doc_number", "Número de documento debe tener entre 1 y 20 caracteres", docNumber)
		}
	}

	return nil
}

// ValidateAmount valida un monto
func ValidateAmount(amount float64, fieldName string) error {
	if amount < 0 {
		return models.NewValidationError(fieldName, "Monto no puede ser negativo", amount)
	}

	if amount > 999999999.99 {
		return models.NewValidationError(fieldName, "Monto excede el máximo permitido", amount)
	}

	return nil
}

// ValidateInvoiceNumber valida un número de factura
func ValidateInvoiceNumber(invoiceNumber int) error {
	if invoiceNumber <= 0 {
		return models.NewValidationError("invoice_number", "Número de factura debe ser mayor a 0", invoiceNumber)
	}

	if invoiceNumber > 99999999 {
		return models.NewValidationError("invoice_number", "Número de factura excede el máximo permitido", invoiceNumber)
	}

	return nil
}

// ValidatePointOfSale valida un punto de venta
func ValidatePointOfSale(pointOfSale int) error {
	if pointOfSale <= 0 {
		return models.NewValidationError("point_of_sale", "Punto de venta debe ser mayor a 0", pointOfSale)
	}

	if pointOfSale > 9999 {
		return models.NewValidationError("point_of_sale", "Punto de venta excede el máximo permitido", pointOfSale)
	}

	return nil
}

// ValidateDate valida una fecha
func ValidateDate(date time.Time, fieldName string) error {
	if date.IsZero() {
		return models.NewValidationError(fieldName, "Fecha no puede estar vacía", date)
	}

	// Validar que la fecha no sea futura (con tolerancia de 1 día)
	if date.After(time.Now().AddDate(0, 0, 1)) {
		return models.NewValidationError(fieldName, "Fecha no puede ser futura", date)
	}

	// Validar que la fecha no sea muy antigua (más de 1 año)
	if date.Before(time.Now().AddDate(-1, 0, 0)) {
		return models.NewValidationError(fieldName, "Fecha no puede ser anterior a 1 año", date)
	}

	return nil
}

// ValidateCurrencyType valida un tipo de moneda
func ValidateCurrencyType(currency models.CurrencyType) error {
	switch currency {
	case models.CurrencyTypePES, models.CurrencyTypeUSD, models.CurrencyTypeEUR, models.CurrencyTypeBRL:
		return nil
	default:
		return models.NewValidationError("currency_type", "Tipo de moneda no válido", currency)
	}
}

// ValidateInvoiceType valida un tipo de factura
func ValidateInvoiceType(invoiceType models.InvoiceType) error {
	switch invoiceType {
	case models.InvoiceTypeA, models.InvoiceTypeB, models.InvoiceTypeC, models.InvoiceTypeE, models.InvoiceTypeM, models.InvoiceTypeT, models.InvoiceTypeR:
		return nil
	default:
		return models.NewValidationError("invoice_type", "Tipo de factura no válido", invoiceType)
	}
}

// ValidateConceptType valida un tipo de concepto
func ValidateConceptType(conceptType models.ConceptType) error {
	switch conceptType {
	case models.ConceptTypeProducts, models.ConceptTypeServices, models.ConceptTypeMixed:
		return nil
	default:
		return models.NewValidationError("concept_type", "Tipo de concepto no válido", conceptType)
	}
}

// ValidateDocumentType valida un tipo de documento
func ValidateDocumentType(docType models.DocumentType) error {
	switch docType {
	case models.DocumentTypeDNI, models.DocumentTypeCUIT, models.DocumentTypeCUIL, models.DocumentTypeCDI, models.DocumentTypeLE, models.DocumentTypeLC, models.DocumentTypeCI, models.DocumentTypePAS, models.DocumentTypeDE, models.DocumentTypeDI:
		return nil
	default:
		return models.NewValidationError("doc_type", "Tipo de documento no válido", docType)
	}
}

// ValidateTaxRate valida una alícuota de impuesto
func ValidateTaxRate(taxRate models.TaxRate) error {
	switch taxRate {
	case models.TaxRate0, models.TaxRate105, models.TaxRate21, models.TaxRate27, models.TaxRate25, models.TaxRate5, models.TaxRateExempt:
		return nil
	default:
		return models.NewValidationError("tax_rate", "Alícuota de impuesto no válida", taxRate)
	}
}

// ValidateItems valida los ítems de una factura
func ValidateItems(items []models.Item) error {
	if len(items) == 0 {
		return models.NewValidationError("items", "La factura debe tener al menos un ítem", items)
	}

	if len(items) > 1000 {
		return models.NewValidationError("items", "La factura no puede tener más de 1000 ítems", len(items))
	}

	for i, item := range items {
		if err := ValidateItem(item, fmt.Sprintf("items[%d]", i)); err != nil {
			return err
		}
	}

	return nil
}

// ValidateItem valida un ítem individual
func ValidateItem(item models.Item, fieldPrefix string) error {
	if item.Description == "" {
		return models.NewValidationError(fieldPrefix+".description", "Descripción del ítem no puede estar vacía", item.Description)
	}

	if len(item.Description) > 200 {
		return models.NewValidationError(fieldPrefix+".description", "Descripción del ítem no puede exceder 200 caracteres", item.Description)
	}

	if err := ValidateAmount(item.Quantity, fieldPrefix+".quantity"); err != nil {
		return err
	}

	if err := ValidateAmount(item.UnitPrice, fieldPrefix+".unit_price"); err != nil {
		return err
	}

	if err := ValidateAmount(item.TotalPrice, fieldPrefix+".total_price"); err != nil {
		return err
	}

	// Validar que el total sea consistente
	expectedTotal := item.Quantity * item.UnitPrice
	if abs(item.TotalPrice-expectedTotal) > 0.01 {
		return models.NewValidationError(fieldPrefix+".total_price", "Total del ítem no coincide con cantidad * precio unitario", item.TotalPrice)
	}

	return nil
}

// abs retorna el valor absoluto de un float64
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
