import { resolveHref } from "./navigator"
import { Fider } from "./fider"

beforeAll(() => {
  // Fider is usually initialized from a <script id="server-data"> element
  // injected by the SSR renderer. Tests have no such element, so we
  // initialize with an empty props/settings object and mutate baseURL per
  // test via the helper below.
  Fider.initialize({ settings: { baseURL: "" }, props: {} })
})

// Helper: configure Fider's baseURL for the duration of a test.
function withBaseURL(baseURL: string, fn: () => void) {
  const original = Fider.settings.baseURL
  Fider.settings.baseURL = baseURL
  try {
    fn()
  } finally {
    Fider.settings.baseURL = original
  }
}

describe("resolveHref", () => {
  describe("when hosted at the domain root", () => {
    test("returns root-relative paths unchanged", () => {
      withBaseURL("https://example.com", () => {
        expect(resolveHref("/posts/1")).toBe("/posts/1")
        expect(resolveHref("/signin")).toBe("/signin")
        expect(resolveHref("/")).toBe("/")
      })
    })
  })

  describe("when hosted under a sub-path", () => {
    test("prepends the sub-path to root-relative paths", () => {
      withBaseURL("https://example.com/feedback", () => {
        expect(resolveHref("/posts/1")).toBe("/feedback/posts/1")
        expect(resolveHref("/signin")).toBe("/feedback/signin")
      })
    })

    test("is idempotent when called with an already-prefixed href", () => {
      withBaseURL("https://example.com/feedback", () => {
        expect(resolveHref("/feedback/posts/1")).toBe("/feedback/posts/1")
      })
    })

    test("handles a trailing slash in baseURL", () => {
      withBaseURL("https://example.com/feedback/", () => {
        expect(resolveHref("/posts/1")).toBe("/feedback/posts/1")
      })
    })
  })

  describe("passes through non-root-relative hrefs unchanged", () => {
    test.each([
      ["absolute URL", "https://external.example/page"],
      ["fragment", "#section"],
      ["mailto", "mailto:foo@example.com"],
      ["empty string", ""],
      ["relative path", "posts/1"],
    ])("%s", (_label, href) => {
      withBaseURL("https://example.com/feedback", () => {
        expect(resolveHref(href)).toBe(href)
      })
    })
  })

  test("returns the input unchanged when baseURL is malformed", () => {
    withBaseURL("not a url", () => {
      expect(resolveHref("/posts/1")).toBe("/posts/1")
    })
  })
})
