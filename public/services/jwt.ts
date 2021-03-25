export const jwt = {
  decode: (token: string): any => {
    if (token) {
      const segments = token.split(".")
      try {
        return JSON.parse(window.atob(segments[1]))
      } catch {
        return undefined
      }
    }
  },
}
