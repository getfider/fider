const global = (1, eval)("this")
global.global = global
global.globalThis = global
global.frames = global
global.self = global

const document = {
  documentElement: {
    style: {},
  },
  getElementById: () => undefined,
}

const navigator = {
  platform: "win32",
}

const window = {
  document,
  location: {
    href: "",
  },
  navigator,
}

global.navigator = navigator
global.window = window
global.document = document

// Polyfill URLSearchParams (which is a constructor). Just add a dummy "get" method that returns an empty string
global.URLSearchParams = class {
  get() {
    return ""
  }
}

// Intl polyfill is required until v8go supports Intl
class NoopFormat {
  format(arg0) {
    return arg0 ? arg0.toString() : ""
  }
}

global.Intl = {
  NumberFormat: NoopFormat,
  DateTimeFormat: NoopFormat,
}

class TextEncoder {
  encode(str) {
    const arr = new Uint8Array(str.length)
    for (let i = 0; i < str.length; i++) {
      arr[i] = str.charCodeAt(i)
    }
    return arr
  }
}
global.TextEncoder = TextEncoder
