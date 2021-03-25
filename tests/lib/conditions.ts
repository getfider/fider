import { NewablePage, Page, BrowserTab, WebComponent } from "."

export type WaitCondition = (
  browser: BrowserTab
) => {
  function: (...args: any[]) => boolean
  args: any[]
}

export const elementIsVisible = (target: string | WebComponent): WaitCondition => {
  return () => {
    const selector = typeof target === "string" ? target : target.selector
    return {
      function: (query: string) => {
        const node = document.querySelector(query)
        if (!node) {
          return false
        }
        const style = window.getComputedStyle(node)
        return style && style.display !== "none" && style.visibility !== "hidden" && style.opacity !== "0"
      },
      args: [selector],
    }
  }
}

export const elementIsNotVisible = (target: string | WebComponent): WaitCondition => {
  return () => {
    const selector = typeof target === "string" ? target : target.selector
    return {
      function: (query: string) => {
        const node = document.querySelector(query)
        if (!node) {
          return true
        }
        const style = window.getComputedStyle(node)
        return style && (style.display === "none" || style.visibility === "hidden" || style.opacity === "0")
      },
      args: [selector],
    }
  }
}

export function pageHasLoaded<T extends Page>(page: NewablePage<T>): WaitCondition {
  return (tab: BrowserTab) => {
    return new page(tab).loadCondition()(tab)
  }
}
