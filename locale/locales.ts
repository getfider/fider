interface Locale {
  text: string;
}

const locales: { [key: string]: Locale; } = {
  en: {
    text: "English",
  },
  "pt-BR": {
    text: "Portuguese (Brazilian)",
  },
  "es-ES": {
    text: "Spanish",
  },
  de: {
    text: "German",
  },
  fr: {
    text: "French",
  },
  "sv-SE": {
    text: "Swedish",
  },
  it: {
    text: "Italian",
  },
  ja: {
    text: "Japanese",
  },
  nl: {
    text: "Dutch",
  },
  pl: {
    text: "Polish",
  },
  ru: {
    text: "Russian",
  },
  sk: {
    text: "Slovak",
  },
  tr: {
    text: "Turkish",
  },
  el: {
    text: "Greek",
  },
  ar: {
    text: "Arabic",
  },
  "zh-CN": {
    text: "Chinese (Simplified)",
  },
  fa: {
    text: "Persian (پارسی)"
  }
};

export default locales;
