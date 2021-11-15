export const analytics = {
  event: (eventCategory: string, eventAction: string): void => {
    if (window.ga) {
      window.ga("send", "event", {
        eventCategory,
        eventAction,
      })
    }
  },
  error: (err?: Error): void => {
    if (window.ga) {
      window.ga("send", "exception", {
        exDescription: err ? err.stack : "<not available>",
        exFatal: false,
      })
    }
  },
}
