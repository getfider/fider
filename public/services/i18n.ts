import { i18n } from "@lingui/core"
import { en, pt } from "make-plural/plurals"

i18n.loadLocaleData("en", { plurals: en })
i18n.loadLocaleData("pt-BR", { plurals: pt })

export async function activate(locale: string, messages?: any) {
  try {
    if (!messages) {
      const content = await import(
        /* webpackChunkName: "locale-[request]" */
        `@fider/../locale/${locale}.js`
      )
      messages = content.messages
    }

    i18n.load(locale, messages)
    i18n.activate(locale)
  } catch (err) {
    console.error(err)
    throw err
  }
}

export function t(message: string, values?: Record<string, any>): string {
  return i18n._(message, values)
}
