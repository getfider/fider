import React from "react"
import { renderToStaticMarkup } from "react-dom/server"
import { Fider, FiderContext } from "./services/fider"
import { IconContext } from "react-icons"
import { Header } from "./components"
import { resolveRootComponent, route } from "./router"

import HomePage from "./pages/Home/Home.page"
import ShowPostPage from "./pages/ShowPost/ShowPost.page"
import SignInPage from "./pages/SignIn/SignIn.page"
import SignUpPage from "./pages/SignUp/SignUp.page"
import UIToolkitPage from "./pages/UI/UIToolkit.page"

// Only public routes should be here
// Routes behind authentication are not crawled
const routes = [
  route("", HomePage),
  route("/posts/:number*", ShowPostPage),
  route("/signin", SignInPage, false),
  route("/signup", SignUpPage, false),
  route("/-/ui", UIToolkitPage),
]

function ssrRender(url: string, pathname: string, args: any) {
  const fider = Fider.initialize({ ...args })
  const config = resolveRootComponent(pathname, routes)
  window.location.href = url

  return renderToStaticMarkup(
    <FiderContext.Provider value={fider}>
      <IconContext.Provider value={{ className: "icon" }}>
        {config.showHeader && <Header />}
        {React.createElement(config.component, args.props)}
      </IconContext.Provider>
    </FiderContext.Provider>
  )
}

;(globalThis as any).ssrRender = ssrRender
