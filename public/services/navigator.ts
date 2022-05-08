import { Fider } from "@fider/services"

const navigator = {
  url: () => {
    return window.location.href
  },
  goHome: () => {
    window.location.href = "/"
  },
  goTo: (url: string) => {
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
