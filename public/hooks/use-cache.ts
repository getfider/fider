import { Dispatch, useEffect, useState } from "react"

const isClient = typeof window !== "undefined"

export function useCache(key: string, defaultValue: string): [string, Dispatch<string>] {
  const [value, setValue] = useState(defaultValue)

  const setCachedValue = (newValue: string) => {
    if (isClient && window.sessionStorage) {
      window.sessionStorage.setItem(key, newValue)
    }
    setValue(newValue)
  }

  useEffect(() => {
    if (isClient) {
      const cachedValue = window.sessionStorage?.getItem(key)
      if (cachedValue) {
        setValue(cachedValue)
      }
    }
  }, [key, setValue])

  return [value, setCachedValue]
}
