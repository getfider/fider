const global = (1, eval)("this")
global.global = global
global.globalThis = global
global.frames = global
global.self = global

const document = {
  documentElement: {},
  getElementById: () => undefined,
}

const window = {
  document,
  location: {
    href: "",
  },
}

const navigator = {}
global.navigator = navigator
global.window = window
global.document = document

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
