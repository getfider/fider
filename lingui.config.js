import { formatter } from "@lingui/format-json"

module.exports = {
  catalogs: [
    {
      path: "<rootDir>/locale/{locale}/client",
      include: ["<rootDir>"],
    },
  ],
  fallbackLocales: {
    default: "en",
  },
  sourceLocale: "en",
  format: formatter({ style: "minimal" }),
  locales: ["pt-BR", "es-ES", "nl", "sv-SE", "fr", "de", "en", "pl", "ru", "sk", "tr", "el"],
}
