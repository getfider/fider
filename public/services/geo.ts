import { cache } from "./cache"

const CACHE_COUNTRYCODE_KEY = "CountryCode"

export const getCountryCode = async (): Promise<string | undefined> => {
  const countryCode = cache.session.get(CACHE_COUNTRYCODE_KEY)
  if (countryCode) {
    return countryCode
  }

  try {
    const response = await fetch("https://ipinfo.io/json")
    if (response.status === 200) {
      const data = await response.json()
      cache.session.set(CACHE_COUNTRYCODE_KEY, data.country)
      return data.country
    }
  } catch (err) {
    console.error(err)
  }
}
