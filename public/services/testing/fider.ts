import { Fider } from "@fider/services"

export const fiderMock = {
  notAuthenticated: () => {
    Fider.initialize()
    Object.defineProperty(Fider.session, "isAuthenticated", {
      get() {
        return false
      },
    })
  },
  authenticated: () => {
    Fider.initialize()
    Object.defineProperty(Fider.session, "isAuthenticated", {
      get() {
        return true
      },
    })
  },
}
