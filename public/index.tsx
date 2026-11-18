import "@fider/assets/styles/index.scss"

import React, { Suspense, startTransition, useEffect, useState } from "react"
import { createRoot } from "react-dom/client"
import { AppShell, ErrorBoundary, Loader, ReadOnlyNotice, DevBanner } from "@fider/components"
import { classSet, Fider, FiderContext, actions, activateI18N, navigator } from "@fider/services"

import { I18n } from "@lingui/core"
import { I18nProvider } from "@lingui/react"
import { AsyncPage } from "./AsyncPages"

const Loading = () => (
  <div className="page">
    <Loader />
  </div>
)

const logProductionError = (err: Error) => {
  if (Fider.isProduction()) {
    console.error(err)
    actions.logError(`react.ErrorBoundary: ${err.message}`, err)
  }
}

window.addEventListener("unhandledrejection", (evt: PromiseRejectionEvent) => {
  if (evt.reason instanceof Error) {
    actions.logError(`window.unhandledrejection: ${evt.reason.message}`, evt.reason)
  } else if (evt.reason) {
    actions.logError(`window.unhandledrejection: ${evt.reason.toString()}`)
  }
})

window.addEventListener("error", (evt: ErrorEvent) => {
  if (evt.error && evt.colno > 0 && evt.lineno > 0) {
    actions.logError(`window.error: ${evt.message}`, evt.error)
  }
})

interface AppState {
  page: string
  props: { [key: string]: any }
  nav: number
  pop: boolean
}

const App = () => {
  const [state, setState] = useState<AppState>({ page: fider.session.page, props: fider.session.props, nav: 0, pop: false })

  useEffect(() => {
    navigator.onPageChange((data, pop) => {
      // transition: keep the current page on screen until the next page's chunk is loaded
      startTransition(() => setState((s) => ({ page: data.page, props: data.props, nav: s.nav + 1, pop })))
    })
  }, [])

  useEffect(() => {
    if (state.nav > 0 && !state.pop) {
      window.scrollTo(0, 0)
    }
  }, [state.nav])

  const Component = AsyncPage(state.page)
  return (
    <AppShell page={state.page}>
      <Suspense fallback={<Loading />}>
        <Component key={`${state.page}:${state.nav}`} {...state.props} />
      </Suspense>
    </AppShell>
  )
}

// Same-origin links navigate client-side; anything that isn't a page falls back to a full load
document.addEventListener("click", (e) => {
  if (e.defaultPrevented || e.button !== 0 || e.metaKey || e.ctrlKey || e.shiftKey || e.altKey) {
    return
  }
  const anchor = (e.target as Element).closest?.("a[href]") as HTMLAnchorElement | null
  if (!anchor || anchor.target || anchor.hasAttribute("download")) {
    return
  }
  const url = new URL(anchor.href)
  if (!navigator.isClientNavigable(url)) {
    return
  }
  if (url.hash && url.pathname === window.location.pathname && url.search === window.location.search) {
    return
  }
  e.preventDefault()
  navigator.navigate(anchor.href)
})

window.addEventListener("popstate", (e) => {
  // History entries owned by the post overlay of the page currently on screen are handled by usePostOverlay
  if (e.state?.overlayPage && e.state.overlayPage === fider.session.page) {
    return
  }
  navigator.navigate(window.location.href, { replace: true, pop: true })
})

const bootstrapApp = (i18n: I18n) => {
  document.body.className = classSet({
    "is-authenticated": fider.session.isAuthenticated,
    "is-staff": fider.session.isAuthenticated && fider.session.user.isCollaborator,
  })

  const rootElement = document.getElementById("root")
  if (rootElement) {
    const root = createRoot(rootElement)
    root.render(
      <React.StrictMode>
        <ErrorBoundary onError={logProductionError}>
          <I18nProvider i18n={i18n}>
            <FiderContext.Provider value={fider}>
              <DevBanner />
              <ReadOnlyNotice />
              <App />
            </FiderContext.Provider>
          </I18nProvider>
        </ErrorBoundary>
      </React.StrictMode>
    )
  }
}
const fider = Fider.initialize()
__webpack_nonce__ = fider.session.contextID
__webpack_public_path__ = `${fider.settings.assetsURL}/assets/`
activateI18N(fider.currentLocale).then(bootstrapApp).catch(bootstrapApp)
