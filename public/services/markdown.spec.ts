import * as markdown from "./markdown"
import { fiderAllowedSchemes } from "@fider/hooks"
fiderAllowedSchemes.get = () => "^monero:[48]\n^bitcoin:(1|3|bc1)"

const testCases = [
  {
    input: "Visit [GitHub](https://github.com) to learn more.",
    expectedFull: '<p>Visit <a target="_blank" rel="noopener nofollow" href="https://github.com" class="text-link">GitHub</a> to learn more.</p>',
    expectedPlainText: "Visit GitHub to learn more.",
  },
  {
    input: "My Picture ![](http://demo.dev.fider.io:3000/images/100/28)",
    expectedFull: '<p>My Picture <img alt="" src="http://demo.dev.fider.io:3000/images/100/28"></p>',
    expectedPlainText: "My Picture",
  },
  {
    input: "My Fider Picture ![](fider-image:attachments/zy0hBtqrjQki7M56p26AuAXljRoaNUSwZO6MOky5gnYm2nW1rsMmrp3dwhjGk7ok-aden.jpeg)",
    expectedFull:
      '<p>My Fider Picture <img data-bkey="attachments/zy0hBtqrjQki7M56p26AuAXljRoaNUSwZO6MOky5gnYm2nW1rsMmrp3dwhjGk7ok-aden.jpeg" class="fider-inline-image" alt="" src="/static/images/attachments/zy0hBtqrjQki7M56p26AuAXljRoaNUSwZO6MOky5gnYm2nW1rsMmrp3dwhjGk7ok-aden.jpeg"></p>',
    expectedPlainText: "My Fider Picture",
  },
  {
    input: "# Hello World",
    expectedFull: "<h1>Hello World</h1>",
    expectedSimple: "<p>Hello World</p>",
    expectedPlainText: "Hello World",
  },
  // {
  //   input: "Hello <b>Beautiful</b> World",
  //   expectedFull: "<p>Hello &lt;b&gt;Beautiful&lt;/b&gt; World</p>",
  //   expectedPlainText: "Hello &lt;b&gt;Beautiful&lt;/b&gt; World",
  // },
  {
    input: `[Uh oh...]("onerror="alert('XSS'))`,
    expectedFull: '<p><a target="_blank" rel="noopener nofollow" href="" class="text-link">Uh oh...</a></p>',
    expectedPlainText: "Uh oh...",
  },
  {
    input: "~~Option 3~~",
    expectedFull: "<p><del>Option 3</del></p>",
    expectedPlainText: "Option 3",
  },
  {
    input: "Check this out: `HEEEY`",
    expectedFull: "<p>Check this out: <code>HEEEY</code></p>",
    expectedPlainText: "Check this out: HEEEY",
  },
  {
    input: `# Hello World
How are you?`,
    expectedFull: `<h1>Hello World</h1>
<p>How are you?</p>`,
    expectedPlainText: "Hello World How are you?",
  },
  {
    input: `-123
-456
-789`,
    expectedFull: "<p>-123<br>-456<br>-789</p>",
    expectedPlainText: "-123 -456 -789",
  },
  {
    input:
      "[monero](monero:83zJ2jMbBoxJkhtpaRLk6fQVrvfmGbd8gYUL7FSdLxU91JSpiWXoLUtAMGqmfvfq3qRS5gJUvMY7oLFSx71wxhKRGG6ypMt) [bitcoin](bitcoin:1CgLs6CxXMAY4Pj4edQq5vyaFoP9NdqVKH) [litecoin](litecoin:ltc1qg0elpp0hxguwlsapl68gvklt5ngemj8k8lu0f5)",
    expectedFull:
      '<p><a href="monero:83zJ2jMbBoxJkhtpaRLk6fQVrvfmGbd8gYUL7FSdLxU91JSpiWXoLUtAMGqmfvfq3qRS5gJUvMY7oLFSx71wxhKRGG6ypMt" target="_blank" rel="noopener nofollow" class="text-link">monero</a> <a href="bitcoin:1CgLs6CxXMAY4Pj4edQq5vyaFoP9NdqVKH" target="_blank" rel="noopener nofollow" class="text-link">bitcoin</a> <a target="_blank" rel="noopener nofollow" class="text-link">litecoin</a></p>',
    expectedPlainText: "monero bitcoin litecoin",
  },
  {
    input: "Jane's & Jim's > [Matt](https://example.com)",
    expectedFull: '<p>Jane\'s &amp; Jim\'s &gt; <a target="_blank" rel="noopener nofollow" href="https://example.com" class="text-link">Matt</a></p>',
    expectedPlainText: "Jane's & Jim's > Matt",
  },
]

testCases.forEach((x) => {
  test(`Can parse markdown ${x.input} to ${x.expectedFull} (full mode)`, () => {
    const result = markdown.full(x.input)
    expect(result).toEqual(x.expectedFull)
  })
})

testCases.forEach((x) => {
  test(`Can parse markdown ${x.input} to ${x.expectedPlainText} (plain text)`, () => {
    const result = markdown.plainText(x.input)
    expect(result).toEqual(x.expectedPlainText)
  })
})

// Test the new URL validation utilities
describe("URL validation utilities", () => {
  test("should always allow HTTP/HTTPS URLs even when no allowed schemes provided", () => {
    // HTTP/HTTPS URLs are always allowed, even when no allowed schemes provided
    expect(markdown.isValidUrl("https://example.com")).toBe(true)
    expect(markdown.isValidUrl("http://example.com")).toBe(true)
    expect(markdown.isValidUrl("example.com")).toBe(true) // gets normalized to https://

    // But crypto URLs are rejected when no allowed schemes provided
    expect(markdown.isValidUrl("monero:4AdUndXHHZ6cFhR4L5n4eL3p5N1L7cKyK4n8v7f3v2")).toBe(false)
  })

  test("should validate URLs when allowed schemes provided", () => {
    const mockAllowedSchemes = [new RegExp("^monero:[48]", "i"), new RegExp("^bitcoin:(1|3|bc1)", "i")]

    // HTTP/HTTPS URLs should work
    expect(markdown.isValidUrl("https://example.com", mockAllowedSchemes)).toBe(true)
    expect(markdown.isValidUrl("http://example.com", mockAllowedSchemes)).toBe(true)
    expect(markdown.isValidUrl("example.com", mockAllowedSchemes)).toBe(true) // gets normalized to https://

    // Crypto URLs should work
    expect(markdown.isValidUrl("monero:4AdUndXHHZ6cFhR4L5n4eL3p5N1L7cKyK4n8v7f3v2", mockAllowedSchemes)).toBe(true)
    expect(markdown.isValidUrl("bitcoin:1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa", mockAllowedSchemes)).toBe(true)

    // Non-allowed crypto should be rejected
    expect(markdown.isValidUrl("litecoin:LTC1qy2e3r4t5y6u7i8o9p0", mockAllowedSchemes)).toBe(false)
  })

  test("should reject invalid URL schemes", () => {
    const mockAllowedSchemes = [new RegExp("^https?://", "i")]
    expect(markdown.isValidUrl("http:invalid-address", mockAllowedSchemes)).toBe(false)
  })

  test("should reject empty URLs", () => {
    expect(markdown.isValidUrl("")).toBe(false)
    expect(markdown.isValidUrl("   ")).toBe(false)
  })

  test("should normalize URLs correctly", () => {
    const regularUrl = markdown.normalizeUrl("example.com")
    expect(regularUrl).toBe("https://example.com")

    const cryptoUrl = markdown.normalizeUrl("monero:4AdUndXHHZ6cFhR4L5n4eL3p5N1L7cKyK4n8v7f3v2")
    expect(cryptoUrl).toBe("monero:4AdUndXHHZ6cFhR4L5n4eL3p5N1L7cKyK4n8v7f3v2")
  })
})
