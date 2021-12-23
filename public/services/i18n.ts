import { i18n, I18n } from "@lingui/core"
import { en, pt, fr, de, se, pl, ru, sk, nl, es, tr } from "make-plural/plurals"

export function activateI18NSync(locale: string, messages?: any): I18n {
  i18n.loadLocaleData("en", { plurals: en })
  i18n.loadLocaleData("pt-BR", { plurals: pt })
  i18n.loadLocaleData("sv-SE", { plurals: se })
  i18n.loadLocaleData("es-ES", { plurals: es })
  i18n.loadLocaleData("nl", { plurals: nl })
  i18n.loadLocaleData("de", { plurals: de })
  i18n.loadLocaleData("fr", { plurals: fr })
  i18n.loadLocaleData("pl", { plurals: pl })
  i18n.loadLocaleData("ru", { plurals: ru })
  i18n.loadLocaleData("sk", { plurals: sk })
  i18n.loadLocaleData("tr", { plurals: tr })
  i18n.load(locale, messages)
  i18n.activate(locale)
  return i18n
}

export async function activateI18N(locale: string): Promise<I18n> {
  try {
    const content = await import(
      /* webpackChunkName: "locale-[request]" */
      `@locale/${locale}/client.json`
    )
    return activateI18NSync(locale, content.messages)
  } catch (err) {
    console.error(err)
    return activateI18NSync(locale)
  }
}
