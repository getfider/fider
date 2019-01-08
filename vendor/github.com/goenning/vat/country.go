package vat

import "strings"

//IsEU returns true if given country code is within EU
func IsEU(countryCode string) bool {
	if countryCode == "UK" || countryCode == "EL" {
		return false
	}

	countryCode = fromISOCountryCode(countryCode)
	_, ok := vatFormatNumberPatterns[countryCode]
	return ok
}

func toISOCountryCode(countryCode string) string {
	countryCode = strings.TrimSpace(countryCode)
	if countryCode == "UK" {
		return "GB"
	}
	if countryCode == "EL" {
		return "GR"
	}
	return countryCode
}

func fromISOCountryCode(countryCode string) string {
	countryCode = strings.TrimSpace(countryCode)
	if countryCode == "GR" {
		return "EL"
	}
	return countryCode
}
