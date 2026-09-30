export const analytics = {
  event: (eventCategory: string, eventAction: string): void => {
    if (window.gtag) {
      // GA4 event names may only contain letters, numbers and underscores
      window.gtag("event", `${eventCategory}_${eventAction}`.replace(/-/g, "_"), {
        event_category: eventCategory,
      })
    }
  },
  error: (err?: Error): void => {
    if (window.gtag) {
      window.gtag("event", "exception", {
        description: err ? err.stack : "<not available>",
        fatal: false,
      })
    }
  },
}
