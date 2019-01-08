package vat

import (
	"fmt"
	"regexp"
	"strings"
)

var vatFormatNumberPatterns = map[string]string{
	"AT": `U[A-Z0-9]{8}`,
	"BE": `(0[0-9]{9}|[0-9]{10})`,
	"BG": `[0-9]{9,10}`,
	"CY": `[0-9]{8}[A-Z]`,
	"CZ": `[0-9]{8,10}`,
	"DE": `[0-9]{9}`,
	"DK": `[0-9]{8}`,
	"EE": `[0-9]{9}`,
	"EL": `[0-9]{9}`,
	"ES": `[A-Z][0-9]{7}[A-Z]|[0-9]{8}[A-Z]|[A-Z][0-9]{8}`,
	"FI": `[0-9]{8}`,
	"FR": `([A-Z]{2}|[0-9]{2})[0-9]{9}`,
	"GB": `[0-9]{9}|[0-9]{12}|(GD|HA)[0-9]{3}`,
	"HR": `[0-9]{11}`,
	"HU": `[0-9]{8}`,
	"IE": `[A-Z0-9]{7}[A-Z]|[A-Z0-9]{7}[A-W][A-I]`,
	"IT": `[0-9]{11}`,
	"LT": `([0-9]{9}|[0-9]{12})`,
	"LU": `[0-9]{8}`,
	"LV": `[0-9]{11}`,
	"MT": `[0-9]{8}`,
	"NL": `[0-9]{9}B[0-9]{2}`,
	"PL": `[0-9]{10}`,
	"PT": `[0-9]{9}`,
	"RO": `[0-9]{2,10}`,
	"SE": `[0-9]{12}`,
	"SI": `[0-9]{8}`,
	"SK": `[0-9]{10}`,
}

// ValidateNumberFormat returns true if given
func ValidateNumberFormat(vatNumber string) (bool, string) {
	vatNumber = sanitizeVATNumber(vatNumber)

	if len(vatNumber) < 3 {
		return false, ""
	}

	vatNumber = strings.ToUpper(vatNumber)
	euCoutryCode := vatNumber[0:2]
	pattern, ok := vatFormatNumberPatterns[euCoutryCode]
	if !ok {
		return false, ""
	}

	matched, err := regexp.MatchString(pattern, vatNumber[2:])
	if err != nil {
		panic(fmt.Sprintf("Failed to compile expression '%s'.", pattern))
	}
	return matched, toISOCountryCode(euCoutryCode)
}

func formatAddress(address string) string {
	address = strings.TrimSpace(address)
	address = strings.Replace(address, "\n\n", ", ", -1)
	address = strings.Replace(address, "\n", ", ", -1)
	return address
}

func sanitizeVATNumber(vatNumber string) string {
	vatNumber = strings.TrimSpace(vatNumber)
	return regexp.MustCompile(" ").ReplaceAllString(vatNumber, "")
}
