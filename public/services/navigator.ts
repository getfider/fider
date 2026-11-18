import { Fider } from "@fider/services"

export interface PageData {
  page: string
  title: string
  props: { [key: string]: any }
}

type PageListener = (data: PageData, pop: boolean) => void
let pageListener: PageListener | undefined

// Routes that are not React pages, or have side effects on GET: always full navigation
const nonPageRoute = /^\/(api|_api|feed|static|assets|oauth|signin\/|signout|attachments|images|admin\/export\/)/

const fullLoad = (url: string) => {
  window.location.href = url
}

const navigator = {
  url: () => {
    return window.location.href
  },
  goHome: () => {
    navigator.goTo("/")
  },
  goTo: (url: string) => {
    const isEqual = window.location.href === url || window.location.pathname === url
    if (!isEqual) {
      navigator.navigate(url)
    }
  },
  // Re-fetch the current page payload and remount the page; the shell (header) stays put.
  // Drop-in replacement for location.reload() after a mutation.
  reload: (): Promise<void> => {
    return navigator.navigate(window.location.href, { replace: true })
  },
  replaceState: (path: string): void => {
    if (history.replaceState !== undefined) {
      const newURL = Fider.settings.baseURL + path
      window.history.replaceState({ path: newURL }, "", newURL)
    }
  },
  onPageChange: (fn: PageListener) => {
    pageListener = fn
  },
  isClientNavigable: (url: URL): boolean => {
    return url.origin === window.location.origin && !nonPageRoute.test(url.pathname)
  },
  // Client-side navigation: fetch the page payload as JSON and swap the page component.
  // Anything that doesn't answer with a page payload falls back to a full load.
  navigate: async (url: string, opts: { replace?: boolean; pop?: boolean } = {}): Promise<void> => {
    const target = new URL(url, window.location.href)
    if (!pageListener || !navigator.isClientNavigable(target)) {
      return fullLoad(url)
    }

    let response: Response
    try {
      response = await fetch(target.href, { headers: { Accept: "application/json", "X-Fider-Navigation": "1" }, credentials: "same-origin" })
    } catch {
      return fullLoad(url)
    }

    if (!(response.headers.get("Content-Type") || "").includes("application/json")) {
      return fullLoad(url)
    }

    const data = await response.json()
    if (!data || !data.page) {
      return fullLoad(url)
    }

    // response.url is the final URL after server-side redirects (e.g. /admin -> /signin)
    const finalURL = new URL(response.url || target.href)
    const path = finalURL.pathname + finalURL.search + finalURL.hash
    if (opts.replace) {
      window.history.replaceState({ spa: true }, "", path)
    } else {
      window.history.pushState({ spa: true }, "", path)
    }

    document.title = data.title
    Fider.initialize(data)
    pageListener(data, !!opts.pop)

    if (window.gtag) {
      window.gtag("event", "page_view", { page_location: finalURL.href, page_title: data.title })
    }
  },
}

export default navigator
