import { I18n, setupI18n } from "@lingui/core"
import { en, pt } from "make-plural/plurals"

let instance: I18n

export function reset() {
  instance = setupI18n()
}

async function getMessages(locale: string): Promise<any> {
  const content = await import(
    /* webpackChunkName: "locale-[request]" */
    `@locale/${locale}.json`
  )
  return content.messages
}

export async function activate(locale: string, messages?: any) {
  locale = locale || "en"
  try {
    if (!messages) {
      messages = { ...(await getMessages("en")), ...(await getMessages(locale)) }
    }

    instance = setupI18n({ missing: (_, key) => `⚠️ Missing Translation: ${key}` })
    instance.loadLocaleData("en", { plurals: en })
    instance.loadLocaleData("pt-BR", { plurals: pt })
    instance.load(locale, messages)
    instance.activate(locale)
  } catch (err) {
    console.error(err)
    throw err
  }
}

export function t(message: string, values?: Record<string, any>): string {
  return instance._(message, values)
}
