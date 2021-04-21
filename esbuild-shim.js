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
