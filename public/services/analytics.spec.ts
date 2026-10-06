import { analytics, toEventName } from "./analytics"

test("toEventName joins category and action with an underscore", () => {
  expect(toEventName("post", "vote")).toBe("post_vote")
  expect(toEventName("comment", "approve-and-verify")).toBe("comment_approve_and_verify")
})

test("toEventName replaces characters GA4 does not allow", () => {
  expect(toEventName("post", "approve and verify")).toBe("post_approve_and_verify")
  expect(toEventName("tag", "import.csv")).toBe("tag_import_csv")
})

test("toEventName truncates to 40 characters", () => {
  const name = toEventName("category", "a-very-long-action-name-that-keeps-on-going")
  expect(name).toHaveLength(40)
  expect(name).toMatch(/^[A-Za-z0-9_]+$/)
})

describe("analytics", () => {
  afterEach(() => {
    delete window.gtag
  })

  test("does nothing when gtag is not loaded", () => {
    expect(() => analytics.event("post", "vote")).not.toThrow()
    expect(() => analytics.error(new Error("boom"))).not.toThrow()
  })

  test("sends events to gtag", () => {
    window.gtag = jest.fn()
    analytics.event("post", "toggle-vote")
    expect(window.gtag).toHaveBeenCalledWith("event", "post_toggle_vote", { event_category: "post" })
  })

  test("sends errors as exception events with the error message", () => {
    window.gtag = jest.fn()
    analytics.error(new TypeError("Cannot read properties of undefined"))
    expect(window.gtag).toHaveBeenCalledWith("event", "exception", {
      description: "TypeError: Cannot read properties of undefined",
      fatal: false,
    })
  })
})
