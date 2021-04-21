import { cache } from "./cache"

test("cache starts empty", () => {
  const value = cache.local.get("my-key")
  expect(value).toBeNull()

  expect(cache.local.has("my-key")).toBeFalsy()
})

test("can set, remove and get from cache", () => {
  cache.local.set("my-key", "Hello World")
  const value = cache.local.get("my-key")
  expect(value).toBe("Hello World")
  expect(cache.local.has("my-key")).toBeTruthy()

  cache.local.remove("my-key")
  const newValue = cache.local.get("my-key")
  expect(newValue).toBeNull()
})
