import { Marked, RendererObject } from "marked"
import DOMPurify from "dompurify"
import { fiderAllowedSchemes } from "@fider/hooks"
import { fiderImageBkey, MENTION_GLOBAL } from "./markdownSyntax"

if (DOMPurify.isSupported) {
  DOMPurify.setConfig({
    USE_PROFILES: {
      html: true,
    },
    ADD_ATTR: ["target"],
  })

  let allow: RegExp[] | undefined
  DOMPurify.addHook("uponSanitizeAttribute", (currentNode, hookEvent) => {
    if (allow === undefined)
      allow = fiderAllowedSchemes
        .get()
        .split("\n")
        .filter((s) => s)
        .map((s) => new RegExp(s, "i"))

    if (allow && hookEvent.attrName === "href") {
      const href = currentNode.getAttribute("href")
      if (href !== null && !/^javascript/i.test(href)) hookEvent.forceKeepAttr = allow.some((r) => r.test(href))
    }
  })
}

// marked 5+ no longer escapes token text in the lexer; escaping is the renderer's job, so every
// override that writes token text/attributes into HTML must escape it itself.

// Text and attribute values. marked 5+ escapes text (including ' -> &#39;); we keep apostrophes
// literal to match the historical output, while preserving existing entities (input has already
// been through encodeHTML, so "<" arrives as "&lt;") and escaping & < > ".
const escapeText = (s: string): string =>
  s
    .replace(/&(?![#\w]+;)/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;")

// Code (spans and blocks): escape everything, exactly like marked's default code renderers,
// so plainText and full show code the same way.
const escapeCode = (s: string): string => s.replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;").replace(/"/g, "&quot;").replace(/'/g, "&#39;")

// Expand Fider's proprietary mention syntax @[name] into a styled span.
const renderMentions = (html: string): string => html.replace(MENTION_GLOBAL, (_match, name) => `<span class="mention">@${name}</span>`)

// marked 5+ replaced the positional Renderer API with token objects and requires renderer
// overrides to be registered via .use(). Typing as RendererObject binds `this` to the
// renderer (exposing this.parser) and infers the exact token types. Unoverridden methods
// keep marked's defaults.
const fullRenderer: RendererObject = {
  image({ href, text }) {
    // Fider's proprietary inline-image syntax: ![](fider-image:<bkey>). The bkey is validated
    // (safe characters, no "..") so it can go into the attributes as-is.
    const bkey = fiderImageBkey(href)
    if (bkey) {
      return `<img src="/static/images/${bkey}" alt="${escapeText(text || "")}" class="fider-inline-image" data-bkey="${bkey}" />`
    }
    return false // fall back to marked's default image renderer
  },
  link({ href, title, tokens }) {
    const text = this.parser.parseInline(tokens)
    const titleAttr = title ? ` title="${escapeText(title)}"` : ""
    return `<a class="text-link" href="${escapeText(href)}"${titleAttr} rel="noopener nofollow" target="_blank">${text}</a>`
  },
  text(token) {
    const rendered = "tokens" in token && token.tokens ? this.parser.parseInline(token.tokens) : escapeText(token.text)
    return renderMentions(rendered)
  },
}
const markedFull = new Marked({ gfm: true, breaks: true })
markedFull.use({ renderer: fullRenderer })

// Plain-text renderer: strip all formatting down to readable text.
const plainTextRenderer: RendererObject = {
  link({ tokens }) {
    return this.parser.parseInline(tokens)
  },
  image() {
    return ""
  },
  br() {
    return " "
  },
  strong({ tokens }) {
    return this.parser.parseInline(tokens)
  },
  del({ tokens }) {
    return this.parser.parseInline(tokens)
  },
  heading({ tokens }) {
    return this.parser.parseInline(tokens)
  },
  paragraph({ tokens }) {
    return ` ${this.parser.parseInline(tokens)} `
  },
  code({ text }) {
    return escapeCode(text)
  },
  codespan({ text }) {
    return escapeCode(text)
  },
  html({ text }) {
    return escapeText(text)
  },
  blockquote({ tokens }) {
    return this.parser.parse(tokens)
  },
  list(token) {
    return token.items.map((item) => `${this.parser.parse(item.tokens).trim()} `).join("")
  },
}
const markedPlainText = new Marked({ gfm: true, breaks: true })
markedPlainText.use({ renderer: plainTextRenderer })

const encodeHTML = (s: string) => s.replace(/</g, "&lt;")
const stripTags = (input: string) => input.replace(/<[^>]*>/g, "")
const sanitize = (input: string) => (DOMPurify.isSupported ? DOMPurify.sanitize(input) : stripTags(input))
// Helper function to decode HTML entities back to readable characters
const decodeHtmlEntities = (text: string): string => {
  return text.replace(/&[#\w]+;/g, (entity) => {
    const entities: { [key: string]: string } = {
      "&#39;": "'",
      "&quot;": '"',
      "&amp;": "&",
      "&lt;": "<",
      "&gt;": ">",
      "&nbsp;": " ",
    }
    return entities[entity] || entity
  })
}

export const full = (input: string): string => {
  return sanitize((markedFull.parse(encodeHTML(input)) as string).trim())
}

// HTML with the formatting stripped, safe to insert as HTML (<Markdown style="plainText" />)
export const plainText = (input: string): string => {
  return sanitize((markedPlainText.parse(encodeHTML(input)) as string).trim())
}

// Plain text with entities decoded (e.g. auto-generated titles). Never insert the result as HTML
export const toText = (input: string): string => {
  return decodeHtmlEntities(plainText(input)).trim()
}
