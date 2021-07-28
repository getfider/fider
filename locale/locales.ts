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
    translated: 96,
  },
  fr: {
    text: "🇫🇷 French",
    translated: 96,
  },
}

export default locales
