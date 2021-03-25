import { Fider } from "."

export const delay = (ms: number) => {
  return new Promise((resolve) => setTimeout(resolve, ms))
}

export const classSet = (input?: any): string => {
  let classes = ""
  if (input) {
    for (const key in input) {
      if (key && !!input[key]) {
        classes += ` ${key}`
      }
    }
    return classes.trim()
  }
  return ""
}

const monthNames = ["January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"]

export const currencySymbol = (currencyCode: string): string => {
  currencyCode = currencyCode ? currencyCode.toLowerCase() : ""
  switch (currencyCode) {
    case "eur":
      return "€"
    case "usd":
      return "$"
  }
  throw new Error(`Invalid currency code '${currencyCode}'`)
}

const twoDigits = (value: number): string => {
  return value <= 9 ? `0${value}` : value.toString()
}

type DateFormat = "full" | "short"

export const formatDate = (input: Date | string, format: DateFormat = "full"): string => {
  const date = input instanceof Date ? input : new Date(input)

  const monthIndex = date.getMonth()
  const year = date.getFullYear()

  if (format === "short") {
    return `${monthNames[monthIndex].substring(0, 3)} ${year}`
  }

  const day = date.getDate()
  const hours = twoDigits(date.getHours())
  const minutes = twoDigits(date.getMinutes())
  return `${monthNames[monthIndex]} ${day}, ${year} · ${hours}:${minutes}`
}

const templates: { [key: string]: string } = {
  seconds: "less than a minute",
  minute: "about a minute",
  minutes: "%d minutes",
  hour: "about an hour",
  hours: "about %d hours",
  day: "a day",
  days: "%d days",
  month: "about a month",
  months: "%d months",
  year: "about a year",
  years: "%d years",
}

const template = (t: string, n: number): string => {
  return templates[t] && templates[t].replace(/%d/i, Math.abs(Math.round(n)).toString())
}

export const timeSince = (now: Date, date: Date): string => {
  const seconds = (now.getTime() - date.getTime()) / 1000
  const minutes = seconds / 60
  const hours = minutes / 60
  const days = hours / 24
  const years = days / 365

  return (
    ((seconds < 45 && template("seconds", seconds)) ||
      (seconds < 90 && template("minute", 1)) ||
      (minutes < 45 && template("minutes", minutes)) ||
      (minutes < 90 && template("hour", 1)) ||
      (hours < 24 && template("hours", hours)) ||
      (hours < 42 && template("day", 1)) ||
      (days < 30 && template("days", days)) ||
      (days < 45 && template("month", 1)) ||
      (days < 365 && template("months", days / 30)) ||
      (years < 1.5 && template("year", 1)) ||
      template("years", years)) + " ago"
  )
}

export const fileToBase64 = async (file: File): Promise<string> => {
  return new Promise<string>((resolve, reject) => {
    const reader = new FileReader()
    reader.addEventListener(
      "load",
      () => {
        const parts = (reader.result as string).split("base64,")
        resolve(parts[1])
      },
      false
    )

    reader.addEventListener(
      "error",
      () => {
        reject(reader.error)
      },
      false
    )

    reader.readAsDataURL(file)
  })
}

export const isCookieEnabled = (): boolean => {
  try {
    document.cookie = "cookietest=1"
    const ret = document.cookie.indexOf("cookietest=") !== -1
    document.cookie = "cookietest=1; expires=Thu, 01-Jan-1970 00:00:01 GMT"
    return ret
  } catch (e) {
    return false
  }
}

export const uploadedImageURL = (bkey: string | undefined, size?: number): string | undefined => {
  if (bkey) {
    if (size) {
      return `${Fider.settings.tenantAssetsURL}/images/${bkey}?size=${size}`
    }
    return `${Fider.settings.tenantAssetsURL}/images/${bkey}`
  }
  return undefined
}

export const truncate = (input: string, maxLength: number): string => {
  if (input && input.length > 1000) {
    return `${input.substr(0, maxLength)}...`
  }
  return input
}
