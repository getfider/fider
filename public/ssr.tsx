import React from "react"
import { renderToStaticMarkup } from "react-dom/server"
import { Fider, FiderContext } from "./services/fider"
import { Header } from "./components"
import { resolveRootComponent, route } from "./router"

import HomePage from "./pages/Home/Home.page"
import ShowPostPage from "./pages/ShowPost/ShowPost.page"
import SignInPage from "./pages/SignIn/SignIn.page"
import SignUpPage from "./pages/SignUp/SignUp.page"
import DesignSystemPage from "./pages/DesignSystem/DesignSystem.page"
import LegalPage from "./pages/Legal/Legal.page"
import { i18n } from "./services"

// Locale files must be bundled for SSR to work synchronously
const messages: { [key: string]: any } = {
  en: require(`../locale/en`),
  "pt-BR": require(`../locale/pt-BR`),
}

// Only public routes should be here
// Routes behind authentication are not crawled
const routes = [
  route("", HomePage),
  route("/posts/:number*", ShowPostPage),
  route("/signin", SignInPage, false),
  route("/signup", SignUpPage, false),
  route("/terms", LegalPage, false),
  route("/privacy", LegalPage, false),
  route("/_design", DesignSystemPage),
]

function ssrRender(url: string, pathname: string, args: any) {
  const fider = Fider.initialize({ ...args })
  i18n.activate(fider.settings.locale, messages[fider.settings.locale].messages)
  const config = resolveRootComponent(pathname, routes)
  window.location.href = url

  return renderToStaticMarkup(
    <FiderContext.Provider value={fider}>
      {config.showHeader && <Header />}
      {React.createElement(config.component, args.props)}
    </FiderContext.Provider>
  )
}

;(globalThis as any).ssrRender = ssrRender
