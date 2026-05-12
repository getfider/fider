# Adding a New Language to Fider 💪

This guide explains how to add a new language/locale to Fider.

## Overview

All locale definitions are centralized in **`app/models/enum/locale.go`**. This is the single source of truth for backend locale metadata including locale codes, display names, pluralization rules, PostgreSQL text search configurations, and language detection mappings.

## Step 1: Create Translation Files

1. Create a new directory in `/locale` with the locale code (e.g., `ko` for Korean, `pt-BR` for Brazilian Portuguese)

2. Create two JSON files in the new directory:
   - `client.json` - Frontend translations (UI strings)
   - `server.json` - Backend translations (emails, server-side messages)

The easiest way to do this is copy the english files, and remove all the translations, replacing them with ones in your language.

## Step 2: Add Locale to Backend (Go)

Edit **`app/models/enum/locale.go`**:

1. Add a new locale constant:

   ```go
   LocaleKorean = Locale{
       Code:              "ko",                 // Locale code
       Name:              "Korean",             // Display name
       MessageFormatCode: "ko",                 // Pluralization library culture code
       PostgresConfig:    "simple",             // PostgreSQL text search config
       LinguaLanguage:    lingua.Korean,        // Lingua-go language enum
       IsRTL:             false,                // Right-to-left text direction
   }
   ```

2. Add the new locale to the `AllLocales` slice:
   ```go
   AllLocales = []Locale{
       LocaleEnglish,
       LocalePortugueseBR,
       // ... other locales
       LocaleKorean,  // Add your new locale here
   }
   ```

### Field Explanations

- **Code**: The locale code (e.g., `"en"`, `"pt-BR"`, `"ko"`)
- **Name**: Human-readable name shown in language selector
- **MessageFormatCode**: Culture code for the [messageformat](https://github.com/gotnospirit/messageformat) library (used for pluralization)
- **PostgresConfig**: PostgreSQL text search configuration name (see [PostgreSQL docs](https://www.postgresql.org/docs/current/textsearch-dictionaries.html))
  - Available configs: `"english"`, `"german"`, `"french"`, `"spanish"`, `"portuguese"`, `"italian"`, `"dutch"`, `"russian"`, `"swedish"`, `"turkish"`, `"arabic"`
  - Use `"simple"` if PostgreSQL doesn't have native support for your language
- **LinguaLanguage**: The [lingua-go](https://github.com/pemistahl/lingua-go) enum for automatic language detection
  - See lingua-go documentation for available languages
- **IsRTL**: Set to `true` for right-to-left languages (Arabic, Hebrew, Persian, etc.)

## Step 3: Add Locale to Frontend

1. Update `locale/locales.ts`

Add your locale to the `locales` object:

```typescript
const locales: { [key: string]: Locale } = {
  en: {
    text: "English",
  },
  // ... other locales
  ko: {
    text: "Korean",
  },
}
```

2. Update `lingui.config.js`

Add your locale code to the `locales` array:

```javascript
locales: ["pt-BR", "es-ES", "nl", /* ... other locales, */ "ko"],
```

3. Update `public/ssr.tsx`

Add your locale to the `messages` object for server-side rendering:

```typescript
const messages: { [key: string]: any } = {
  en: require(`../locale/en/client`),
  // ... other locales
  ko: require(`../locale/ko/client`),
}
```

## Step 4: Generate and Compile Translations

After adding the locale to all configuration files:

1. Check that everything works..

   ```bash
   make build
   ```

2. Run the server:

   ```bash
   make run
   ```

3. Verify:

   - The new language appears in the language selector
   - Translations display correctly
   - Right-to-left rendering works (if applicable)
   - Search functionality works with the correct PostgreSQL configuration

4. Open a pull request

## When New Translation Strings Are Added

Fider uses [LinguiJS](https://lingui.dev) for i18n management. When new features are added:

1. New translation keys are automatically extracted and added to all locale files (`client.json` / `server.json`)
2. A [GitHub Action](/.github/workflows/locale.yml) uses Google Cloud Translate to automatically translate missing keys
3. Machine translations are submitted as a pull request for review

While the automated translations work well for most cases, you may want to refine them for accuracy and cultural context. Feel free to submit pull requests with improved translations at any time.
