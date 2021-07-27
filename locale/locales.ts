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
    translated: 98,
  },
  de: {
    text: "🇩🇪 German",
    translated: 98,
  },
  fr: {
    text: "🇫🇷 French",
    translated: 98,
  },
}

export default locales
