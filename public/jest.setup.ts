import { configure } from "enzyme"
import Adapter from "@wojtekmaj/enzyme-adapter-react-17"

configure({ adapter: new Adapter() })

let localStorageCache: {
  [key: string]: string | undefined
}

beforeEach(() => {
  localStorageCache = {}
})
;(window as any).localStorage = {
  getItem: (key: string) => {
    const value = localStorageCache[key]
    return typeof value === "undefined" ? null : value
  },
  setItem: (key: string, value: string) => {
    localStorageCache[key] = value
  },
  removeItem: (key: string) => {
    return delete localStorageCache[key]
  },
}
