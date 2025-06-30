import { Fider } from "@fider/services"

/**
 * Handles direct post URLs by checking if the current URL is a post URL
 * and returns the post number if it is.
 *
 * @returns The post number if the current URL is a post URL, otherwise null
 */
export const getPostNumberFromURL = (): number | null => {
  const path = window.location.pathname
  const postRegex = /\/posts\/(\d+)(?:\/.*)?$/
  const match = path.match(postRegex)

  if (match && match[1]) {
    return parseInt(match[1], 10)
  }

  return null
}

/**
 * Checks if the current page is the home page
 *
 * @returns True if the current page is the home page
 */
export const isHomePage = (): boolean => {
  return Fider.session.page === "Home/Home.page"
}

/**
 * Sets up browser history handling for post modals
 *
 * @param onPopState Function to call when the user navigates back/forward
 */
export const setupHistoryHandling = (onPopState: (postNumber: number | null) => void): (() => void) => {
  const handlePopState = (event: PopStateEvent) => {
    const postNumber = event.state?.postNumber || getPostNumberFromURL()
    onPopState(postNumber)
  }

  window.addEventListener("popstate", handlePopState)

  // Return a cleanup function
  return () => {
    window.removeEventListener("popstate", handlePopState)
  }
}
