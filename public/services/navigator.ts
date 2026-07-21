import { Fider } from "@fider/services"

export function basePath(): string {
  try {
    const p = new URL(Fider.settings.baseURL).pathname.replace(/\/$/, "")
    return p === "/" ? "" : p
  } catch {
    return ""
  }
}

// resolveHref prepends basePath() to root-relative hrefs. Non-root-relative
// hrefs (absolute URLs, fragments, mailto:, etc.) and already-prefixed hrefs
// are returned unchanged, making this safe to apply idempotently.
export function resolveHref(href: string): string {
  if (!href.startsWith("/")) return href
  const bp = basePath()
  if (bp && !href.startsWith(bp)) {
    return bp + href
  }
  return href
}

const navigator = {
  url: () => {
    return window.location.href
  },
  goHome: () => {
    window.location.href = Fider.settings.baseURL
  },
  goTo: (url: string) => {
    url = resolveHref(url)
    const isEqual = window.location.href === url || window.location.pathname === url
    if (!isEqual) {
      window.location.href = url
    }
  },
  replaceState: (path: string): void => {
    if (history.replaceState !== undefined) {
      const newURL = Fider.settings.baseURL + path
      window.history.replaceState({ path: newURL }, "", newURL)
    }
  },
}

export default navigator
