import { formatter } from "@lingui/format-json"

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
