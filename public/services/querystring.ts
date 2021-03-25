import navigator from "./navigator"

export const getNumber = (name: string): number | undefined => {
  return parseInt(get(name), 10) || undefined
}

export const set = (name: string, value: any): string => {
  const uri = navigator.url()
  const re = new RegExp("([?&])" + name + "=.*?(&|$)", "i")
  if (uri.match(re)) {
    return uri.replace(re, "$1" + name + "=" + value + "$2")
  } else {
    const separator = uri.indexOf("?") !== -1 ? "&" : "?"
    return uri + separator + name + "=" + value
  }
}

export const get = (name: string): string => {
  name = name.replace(/[[\]]/g, "\\$&")
  const regex = new RegExp("[?&]" + name + "(=([^&#]*)|&|#|$)")
  const results = regex.exec(navigator.url())

  if (!results || !results[2]) {
    return ""
  }

  return decodeURIComponent(results[2].replace(/\+/g, " "))
}

export const getArray = (name: string): string[] => {
  const qs = get(name)
  if (qs) {
    return qs.split(",").filter((i) => i)
  }

  return []
}

export interface QueryString {
  [key: string]: string | string[] | number | undefined
}

export const stringify = (object: QueryString | undefined): string => {
  if (!object) {
    return ""
  }

  let qs = ""

  for (const key of Object.keys(object)) {
    const symbol = qs ? "&" : "?"
    const value = object[key]
    if (value instanceof Array) {
      if (value.length > 0) {
        qs += `${symbol}${key}=${value.join(",")}`
      }
    } else if (value) {
      qs += `${symbol}${key}=${encodeURIComponent(value.toString()).replace(/%20/g, "+")}`
    }
  }

  return qs
}
