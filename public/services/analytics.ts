// GA4 event names may only contain letters, numbers and underscores, and are limited to 40 characters
export const toEventName = (eventCategory: string, eventAction: string): string =>
  `${eventCategory}_${eventAction}`.replace(/[^A-Za-z0-9_]/g, "_").substring(0, 40)

export const analytics = {
  event: (eventCategory: string, eventAction: string): void => {
    if (window.gtag) {
      window.gtag("event", toEventName(eventCategory, eventAction), {
        event_category: eventCategory,
      })
    }
  },
  error: (err?: Error): void => {
    if (window.gtag) {
      // GA4 truncates parameter values to 100 characters, so a stack trace would be cut off after its first line
      window.gtag("event", "exception", {
        description: err ? `${err.name}: ${err.message}` : "<not available>",
        fatal: false,
      })
    }
  },
}
