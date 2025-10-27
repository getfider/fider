import { formatter } from "@lingui/format-json"

// NOTE: All locale definitions are centralized in app/models/enum/locale.go
// See locale/locales.ts for the complete list of steps when adding a new locale

export default {
  catalogs: [
    {
      path: "<rootDir>/locale/{locale}/client",
      include: ["<rootDir>/public/**/*.{ts,tsx}"],
    },
  ],
  orderBy: "messageId",
  fallbackLocales: {
    default: "en",
  },
  sourceLocale: "en",
  format: formatter({ style: "minimal", explicitIdAsDefault: true, sort: true }),
  locales: ["pt-BR", "es-ES", "nl", "sv-SE", "fr", "de", "en", "pl", "ru", "ja", "sk", "tr", "el", "it", "zh-CN", "ar", "fa"],
}
