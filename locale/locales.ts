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
  "es-ES": {
    text: "🇪🇸 Spanish",
    translated: 97,
  },
  de: {
    text: "🇩🇪 German",
    translated: 97,
  },
  fr: {
    text: "🇫🇷 French",
    translated: 97,
  },
  "sv-SE": {
    text: "🇸🇪 Swedish",
    translated: 97,
  },
  nl: {
    text: "🇳🇱 Dutch",
    translated: 97,
  },
  pl: {
    text: "🇵🇱 Polish",
    translated: 97,
  },
  ru: {
    text: "🇷🇺 Russian",
    translated: 97,
  },
  sk: {
    text: "🇸🇰 Slovak",
    translated: 97,
  },
  tr: {
    text: "🇹🇷 Turkish",
    translated: 97,
  },
}

export default locales
