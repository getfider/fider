export const cache = {
  set: (key: string, value: string): void => {
    if (window.sessionStorage) {
      window.sessionStorage.setItem(key, value);
    }
  },
  get: (key: string): string | null => {
    if (window.sessionStorage) {
      return window.sessionStorage.getItem(key);
    }
    return null;
  },
  remove: (...keys: string[]): void => {
    if (window.sessionStorage && keys) {
      for (const key of keys) {
        window.sessionStorage.removeItem(key);
      }
    }
  }
};
