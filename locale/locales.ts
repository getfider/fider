interface Locale {
  text: string
  translated?: number
}

const locales: { [key: string]: Locale } = {
  en: {
    text: "🇺🇸 English",
  },
  "pt-BR": {
    text: "🇧🇷 Portuguese (Brazilian)",
    translated: 100,
  },
}

export default locales
