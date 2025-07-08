import { cache } from "@fider/services"

export const CACHE_KEYS = {
  TITLE: "PostInput-Title",
  DESCRIPTION: "PostInput-Description",
  ATTACHMENT: "PostInput-Attachment",
  TAGS: "PostInput-Tags",
  POST_PENDING: "PostInput-PostPending",
} as const

export const getCachedTitle = (): string => {
  return getCachedValue(CACHE_KEYS.TITLE)
}

export const setCachedTitle = (title: string): void => {
  cache.local.set(CACHE_KEYS.TITLE, title)
}

export const getCachedDescription = (): string => {
  return getCachedValue(CACHE_KEYS.DESCRIPTION)
}

export const setCachedDescription = (description: string): void => {
  cache.local.set(CACHE_KEYS.DESCRIPTION, description)
}

const getCachedValue = (key: string): string => {
  return cache.local.get(key) || ""
}

export const getCachedTags = (): string[] => {
  const cacheValue = getCachedValue(CACHE_KEYS.TAGS)
  return cacheValue.split(",")
}

export const setCachedTags = (tags: string[]): void => {
  cache.local.set(CACHE_KEYS.TAGS, tags.join(","))
}

export const setPostPending = (value: boolean): void => {
  cache.local.set(CACHE_KEYS.POST_PENDING, value.toString())
}

export const isPostPending = (): boolean => {
  return getCachedValue(CACHE_KEYS.POST_PENDING) === "true"
}

export const clearCache = () => {
  cache.local.remove(...Object.values(CACHE_KEYS))
}

export const setPostCreated = () => {
  cache.session.set("POST_CREATED_SUCCESS", "true")
}
