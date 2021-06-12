import { I18n, setupI18n } from "@lingui/core"
import { en, pt } from "make-plural/plurals"

export function activateI18NSync(locale: string, messages: any): I18n {
  const i18n = setupI18n({ missing: (_, key) => `⚠️ Missing Translation: ${key}` })
  locale = locale || "en"
  i18n.loadLocaleData("en", { plurals: en })
  i18n.loadLocaleData("pt-BR", { plurals: pt })
  i18n.load(locale, messages)
  i18n.activate(locale)
  return i18n
}

export async function activateI18N(locale: string): Promise<I18n> {
  const i18n = setupI18n({ missing: (_, key) => `⚠️ Missing Translation: ${key}` })
  locale = locale || "en"
  try {
    const content = await import(
      /* webpackChunkName: "locale-[request]" */
      `@locale/client/${locale}.po`
    )

    i18n.loadLocaleData("en", { plurals: en })
    i18n.loadLocaleData("pt-BR", { plurals: pt })
    i18n.load(locale, content.messages)
    i18n.activate(locale)
  } catch (err) {
    console.error(err)
  }
  return i18n
}
