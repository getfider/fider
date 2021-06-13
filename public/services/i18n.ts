import { i18n, I18n } from "@lingui/core"
import { en, pt } from "make-plural/plurals"

const DEFAULT_LOCALE = "en"

export function activateI18NSync(locale: string, messages?: any): I18n {
  i18n.loadLocaleData("en", { plurals: en })
  i18n.loadLocaleData("pt-BR", { plurals: pt })
  if (messages) {
    i18n._missing = (_, key) => `⚠️ Missing Translation: ${key}`
    i18n.load(locale, messages)
  }
  i18n.activate(locale)
  return i18n
}

export async function activateI18N(locale: string): Promise<I18n> {
  if (locale !== DEFAULT_LOCALE) {
    try {
      const content = await import(
        /* webpackChunkName: "locale-[request]" */
        `@locale/client/${locale}.json`
      )
      return activateI18NSync(locale, content.messages)
    } catch (err) {
      console.error(err)
    }
  }

  return activateI18NSync(locale)
}
