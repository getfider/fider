package enum

import (
	"github.com/pemistahl/lingua-go"
)

// Locale represents a supported language/locale in Fider
type Locale struct {
	Code              string          // Locale code (e.g., "pt-BR", "en")
	Name              string          // Display name (e.g., "Portuguese (Brazilian)")
	MessageFormatCode string          // Culture code for messageformat pluralization (e.g., "pt")
	PostgresConfig    string          // PostgreSQL text search configuration (e.g., "portuguese", "simple")
	LinguaLanguage    lingua.Language // Lingua-go language enum for detection
	IsRTL             bool            // Right-to-left text direction
}

var (
	// LocaleEnglish represents English
	LocaleEnglish = Locale{
		Code:              "en",
		Name:              "English",
		MessageFormatCode: "en",
		PostgresConfig:    "english",
		LinguaLanguage:    lingua.English,
		IsRTL:             false,
	}

	// LocalePortugueseBR represents Portuguese (Brazilian)
	LocalePortugueseBR = Locale{
		Code:              "pt-BR",
		Name:              "Portuguese (Brazilian)",
		MessageFormatCode: "pt",
		PostgresConfig:    "portuguese",
		LinguaLanguage:    lingua.Portuguese,
		IsRTL:             false,
	}

	// LocaleSpanishES represents Spanish
	LocaleSpanishES = Locale{
		Code:              "es-ES",
		Name:              "Spanish",
		MessageFormatCode: "es",
		PostgresConfig:    "spanish",
		LinguaLanguage:    lingua.Spanish,
		IsRTL:             false,
	}

	// LocaleGerman represents German
	LocaleGerman = Locale{
		Code:              "de",
		Name:              "German",
		MessageFormatCode: "de",
		PostgresConfig:    "german",
		LinguaLanguage:    lingua.German,
		IsRTL:             false,
	}

	// LocaleFrench represents French
	LocaleFrench = Locale{
		Code:              "fr",
		Name:              "French",
		MessageFormatCode: "fr",
		PostgresConfig:    "french",
		LinguaLanguage:    lingua.French,
		IsRTL:             false,
	}

	// LocaleSwedishSE represents Swedish
	LocaleSwedishSE = Locale{
		Code:              "sv-SE",
		Name:              "Swedish",
		MessageFormatCode: "se",
		PostgresConfig:    "swedish",
		LinguaLanguage:    lingua.Swedish,
		IsRTL:             false,
	}

	// LocaleItalian represents Italian
	LocaleItalian = Locale{
		Code:              "it",
		Name:              "Italian",
		MessageFormatCode: "it",
		PostgresConfig:    "italian",
		LinguaLanguage:    lingua.Italian,
		IsRTL:             false,
	}

	// LocaleJapanese represents Japanese
	LocaleJapanese = Locale{
		Code:              "ja",
		Name:              "Japanese",
		MessageFormatCode: "ja",
		PostgresConfig:    "simple",
		LinguaLanguage:    lingua.Japanese,
		IsRTL:             false,
	}

	// LocaleDutch represents Dutch
	LocaleDutch = Locale{
		Code:              "nl",
		Name:              "Dutch",
		MessageFormatCode: "nl",
		PostgresConfig:    "dutch",
		LinguaLanguage:    lingua.Dutch,
		IsRTL:             false,
	}

	// LocalePolish represents Polish
	LocalePolish = Locale{
		Code:              "pl",
		Name:              "Polish",
		MessageFormatCode: "pl",
		PostgresConfig:    "simple",
		LinguaLanguage:    lingua.Polish,
		IsRTL:             false,
	}

	// LocaleRussian represents Russian
	LocaleRussian = Locale{
		Code:              "ru",
		Name:              "Russian",
		MessageFormatCode: "ru",
		PostgresConfig:    "russian",
		LinguaLanguage:    lingua.Russian,
		IsRTL:             false,
	}

	// LocaleSlovak represents Slovak
	LocaleSlovak = Locale{
		Code:              "sk",
		Name:              "Slovak",
		MessageFormatCode: "sk",
		PostgresConfig:    "simple",
		LinguaLanguage:    lingua.Slovak,
		IsRTL:             false,
	}

	// LocaleTurkish represents Turkish
	LocaleTurkish = Locale{
		Code:              "tr",
		Name:              "Turkish",
		MessageFormatCode: "tr",
		PostgresConfig:    "turkish",
		LinguaLanguage:    lingua.Turkish,
		IsRTL:             false,
	}

	// LocaleGreek represents Greek
	LocaleGreek = Locale{
		Code:              "el",
		Name:              "Greek",
		MessageFormatCode: "el",
		PostgresConfig:    "simple",
		LinguaLanguage:    lingua.Greek,
		IsRTL:             false,
	}

	// LocaleArabic represents Arabic
	LocaleArabic = Locale{
		Code:              "ar",
		Name:              "Arabic",
		MessageFormatCode: "ar",
		PostgresConfig:    "arabic",
		LinguaLanguage:    lingua.Arabic,
		IsRTL:             true,
	}

	// LocaleChineseCN represents Chinese (Simplified)
	LocaleChineseCN = Locale{
		Code:              "zh-CN",
		Name:              "Chinese (Simplified)",
		MessageFormatCode: "zh",
		PostgresConfig:    "simple",
		LinguaLanguage:    lingua.Chinese,
		IsRTL:             false,
	}

	// LocalePersian represents Persian (Farsi)
	LocalePersian = Locale{
		Code:              "fa",
		Name:              "Persian (پارسی)",
		MessageFormatCode: "fa",
		PostgresConfig:    "simple",
		LinguaLanguage:    lingua.Persian,
		IsRTL:             false,
	}

	// AllLocales contains all supported locales
	AllLocales = []Locale{
		LocaleEnglish,
		LocalePortugueseBR,
		LocaleSpanishES,
		LocaleGerman,
		LocaleFrench,
		LocaleSwedishSE,
		LocaleItalian,
		LocaleJapanese,
		LocaleDutch,
		LocalePolish,
		LocaleRussian,
		LocaleSlovak,
		LocaleTurkish,
		LocaleGreek,
		LocaleArabic,
		LocaleChineseCN,
		LocalePersian,
	}
)

// GetLocaleByCode returns a Locale by its code, or false if not found
func GetLocaleByCode(code string) (Locale, bool) {
	for _, locale := range AllLocales {
		if locale.Code == code {
			return locale, true
		}
	}
	return Locale{}, false
}

// IsValidLocale returns true if the given locale code is valid
func IsValidLocale(code string) bool {
	_, valid := GetLocaleByCode(code)
	return valid
}

// GetLinguaLanguages returns all lingua.Language enums for language detection
func GetLinguaLanguages() []lingua.Language {
	languages := make([]lingua.Language, len(AllLocales))
	for i, locale := range AllLocales {
		languages[i] = locale.LinguaLanguage
	}
	return languages
}

// GetLocaleByLinguaLanguage returns a Locale by its lingua.Language enum
func GetLocaleByLinguaLanguage(lang lingua.Language) (Locale, bool) {
	for _, locale := range AllLocales {
		if locale.LinguaLanguage == lang {
			return locale, true
		}
	}
	return Locale{}, false
}

// MapLocaleToTSConfig maps a locale code to a PostgreSQL text search configuration
// Returns "simple" if the locale is not found or doesn't have native Postgres support
func MapLocaleToTSConfig(code string) string {
	locale, found := GetLocaleByCode(code)
	if !found {
		return "simple"
	}
	return locale.PostgresConfig
}

// GetLocaleCodes returns all locale codes
func GetLocaleCodes() []string {
	codes := make([]string, len(AllLocales))
	for i, locale := range AllLocales {
		codes[i] = locale.Code
	}
	return codes
}
