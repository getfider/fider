import { i18n } from "@lingui/core"

export function activateI18NSync(locale: string, messages?: any) {
  i18n.load(locale, messages)
  i18n.activate(locale)
  return i18n
}

export async function activateI18N(locale: string) {
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
