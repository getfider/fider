import * as plurals from "make-plural/plurals"

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

// Lingui formats plural messages with Intl.PluralRules, so back it with the CLDR rules from make-plural
class PluralRules {
  constructor(locales, options) {
    const locale = [].concat(locales || [])[0] || "en"
    this.rule = plurals[locale.replace("-", "_")] || plurals[locale.split("-")[0]] || plurals.en
    this.ordinal = !!options && options.type === "ordinal"
  }

  select(n) {
    return this.rule(n, this.ordinal)
  }
}

global.Intl = {
  NumberFormat: NoopFormat,
  DateTimeFormat: NoopFormat,
  PluralRules,
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
