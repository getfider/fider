/* eslint-disable @typescript-eslint/no-var-requires */
import React from "react"
import { renderToStaticMarkup } from "react-dom/server"
import { Fider, FiderContext } from "./services/fider"
import { DevBanner, ReadOnlyNotice } from "./components"

import { activateI18NSync } from "./services"
import { I18nProvider } from "@lingui/react"

// Locale files must be bundled for SSR to work synchronously
const messages: { [key: string]: any } = {
  en: require(`../locale/en/client`),
  "pt-BR": require(`../locale/pt-BR/client`),
  "sv-SE": require(`../locale/sv-SE/client`),
  "es-ES": require(`../locale/es-ES/client`),
  nl: require(`../locale/nl/client`),
  de: require(`../locale/de/client`),
  fr: require(`../locale/fr/client`),
  pl: require(`../locale/pl/client`),
  ru: require(`../locale/ru/client`),
  sk: require(`../locale/sk/client`),
  tr: require(`../locale/tr/client`),
}

// ESBuild doesn't support Dynamic Imports, so we need to map them statically
// But at least only public routes will be here, as routes behind authentication won't be crawled anyway
const pages: { [key: string]: any } = {
  "Home/Home.page": require(`./pages/Home/Home.page`),
  "ShowPost/ShowPost.page": require(`./pages/ShowPost/ShowPost.page`),
  "SignIn/SignIn.page": require(`./pages/SignIn/SignIn.page`),
  "SignUp/SignUp.page": require(`./pages/SignUp/SignUp.page`),
  "SignUp/PendingActivation.page": require(`./pages/SignUp/PendingActivation.page`),
  "Legal/Legal.page": require(`./pages/Legal/Legal.page`),
  "DesignSystem/DesignSystem.page": require(`./pages/DesignSystem/DesignSystem.page`),
  "Error/Maintenance.page": require(`./pages/Error/Maintenance.page`),
  "Error/Error401.page": require(`./pages/Error/Error401.page`),
  "Error/Error403.page": require(`./pages/Error/Error403.page`),
  "Error/Error404.page": require(`./pages/Error/Error404.page`),
  "Error/Error410.page": require(`./pages/Error/Error410.page`),
  "Error/Error500.page": require(`./pages/Error/Error500.page`),
  "Error/NotInvited.page": require(`./pages/Error/NotInvited.page`),
}

function ssrRender(url: string, args: any) {
  const fider = Fider.initialize({ ...args })
  const i18n = activateI18NSync(fider.currentLocale, messages[fider.currentLocale].messages)
  const component = pages[fider.session.page]?.default
  if (!component) {
    throw new Error(`Page not found: ${fider.session.page}`)
  }

  window.location.href = url

  return renderToStaticMarkup(
    <I18nProvider i18n={i18n}>
      <FiderContext.Provider value={fider}>
        <DevBanner />
        <ReadOnlyNotice />
        {React.createElement(component, args.props)}
      </FiderContext.Provider>
    </I18nProvider>
  )
}

;(globalThis as any).ssrRender = ssrRender
