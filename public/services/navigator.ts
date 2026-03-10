import { Fider } from "@fider/services"

export function basePath(): string {
  try {
    const p = new URL(Fider.settings.baseURL).pathname.replace(/\/$/, "")
    return p === "/" ? "" : p
  } catch {
    return ""
  }
}

const navigator = {
  url: () => {
    return window.location.href
  },
  goHome: () => {
    window.location.href = Fider.settings.baseURL
  },
  goTo: (url: string) => {
    if (url.startsWith("/")) {
      const bp = basePath()
      if (bp && !url.startsWith(bp)) {
        url = bp + url
      }
    }
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
