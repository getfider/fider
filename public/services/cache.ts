const set = (storage: Storage, key: string, value: string): void => {
  if (storage) {
    storage.setItem(key, value)
  }
}

const get = (storage: Storage, key: string): string | null => {
  if (window.localStorage) {
    return storage.getItem(key)
  }
  return null
}

const has = (storage: Storage, key: string): boolean => {
  if (storage) {
    return !!storage.getItem(key)
  }
  return false
}

const remove = (storage: Storage, ...keys: string[]): void => {
  if (storage && keys) {
    for (const key of keys) {
      storage.removeItem(key)
    }
  }
}

export const cache = {
  local: {
    set: (key: string, value: string): void => {
      set(window.localStorage, key, value)
    },
    get: (key: string): string | null => {
      return get(window.localStorage, key)
    },
    has: (key: string): boolean => {
      return has(window.localStorage, key)
    },
    remove: (...keys: string[]): void => {
      remove(window.localStorage, ...keys)
    },
  },
  session: {
    set: (key: string, value: string): void => {
      set(window.sessionStorage, key, value)
    },
    get: (key: string): string | null => {
      return get(window.sessionStorage, key)
    },
    has: (key: string): boolean => {
      return has(window.sessionStorage, key)
    },
    remove: (...keys: string[]): void => {
      remove(window.sessionStorage, ...keys)
    },
  },
}
