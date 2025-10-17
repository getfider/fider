import * as url from "./url"

// Test the URL validation utilities
describe("URL validation utilities", () => {
  test("should always allow HTTP/HTTPS URLs even when no allowed schemes provided", () => {
    // HTTP/HTTPS URLs are always allowed, even when no allowed schemes provided
    expect(url.isValidUrl("https://example.com")).toBe(true)
    expect(url.isValidUrl("http://example.com")).toBe(true)
    expect(url.isValidUrl("example.com")).toBe(true) // gets normalized to https://

    // But crypto URLs are rejected when no allowed schemes provided
    expect(url.isValidUrl("monero:4AdUndXHHZ6cFhR4L5n4eL3p5N1L7cKyK4n8v7f3v2")).toBe(false)
  })

  test("should validate URLs when allowed schemes provided", () => {
    const mockAllowedSchemes = [new RegExp("^monero:[48]", "i"), new RegExp("^bitcoin:(1|3|bc1)", "i")]

    // HTTP/HTTPS URLs should work
    expect(url.isValidUrl("https://example.com", mockAllowedSchemes)).toBe(true)
    expect(url.isValidUrl("http://example.com", mockAllowedSchemes)).toBe(true)
    expect(url.isValidUrl("example.com", mockAllowedSchemes)).toBe(true) // gets normalized to https://

    // Crypto URLs should work
    expect(url.isValidUrl("monero:4AdUndXHHZ6cFhR4L5n4eL3p5N1L7cKyK4n8v7f3v2", mockAllowedSchemes)).toBe(true)
    expect(url.isValidUrl("bitcoin:1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa", mockAllowedSchemes)).toBe(true)

    // Non-allowed crypto should be rejected
    expect(url.isValidUrl("litecoin:LTC1qy2e3r4t5y6u7i8o9p0", mockAllowedSchemes)).toBe(false)
  })

  test("should reject invalid URL schemes", () => {
    const mockAllowedSchemes = [new RegExp("^https?://", "i")]
    expect(url.isValidUrl("http:invalid-address", mockAllowedSchemes)).toBe(false)
  })

  test("should reject empty URLs", () => {
    expect(url.isValidUrl("")).toBe(false)
    expect(url.isValidUrl("   ")).toBe(false)
  })

  test("should normalize URLs correctly", () => {
    const regularUrl = url.normalizeUrl("example.com")
    expect(regularUrl).toBe("https://example.com")

    const cryptoUrl = url.normalizeUrl("monero:4AdUndXHHZ6cFhR4L5n4eL3p5N1L7cKyK4n8v7f3v2")
    expect(cryptoUrl).toBe("monero:4AdUndXHHZ6cFhR4L5n4eL3p5N1L7cKyK4n8v7f3v2")
  })
})
