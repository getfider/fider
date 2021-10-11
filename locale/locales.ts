interface Locale {
  text: string
  translated: number
}

const locales: { [key: string]: Locale } = {
  en: {
    text: "🇺🇸 English",
    translated: 100,
  },
  "pt-BR": {
    text: "🇧🇷 Portuguese (Brazilian)",
    translated: 100,
  },
  de: {
    text: "🇩🇪 German",
    translated: 100,
  },
  fr: {
    text: "🇫🇷 French",
    translated: 100,
  },
  "sv-SE": {
    text: "🇸🇪 Swedish",
    translated: 100,
  },
  nl: {
    text: "🇳🇱 Dutch",
    translated: 100,
  },
  ru: {
    text: "🇷🇺 Russian",
    translated: 100,
  },
  sk: {
    text: "🇸🇰 Slovak",
    translated: 100,
  },
}

export default locales
