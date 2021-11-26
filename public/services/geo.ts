import { cache } from "./cache"

const CACHE_CURRENCYCODE_KEY = "CurrencyCode"

export const getCurrencyCode = async (): Promise<string> => {
  const currencyCode = cache.session.get(CACHE_CURRENCYCODE_KEY)
  if (currencyCode) {
    return currencyCode
  }

  try {
    const response = await fetch("https://mygeo.vercel.app")
    if (response.status === 200) {
      const data = await response.json()
      cache.session.set(CACHE_CURRENCYCODE_KEY, data.currencyCode)
      return data.currencyCode
    }
  } catch (err) {
    console.error(err)
  }

  return "USD"
}
