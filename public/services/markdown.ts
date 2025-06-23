import { marked } from "marked"
import DOMPurify from "dompurify"

marked.setOptions({
  headerIds: false,
  xhtml: true,
  smartLists: true,
  gfm: true,
  breaks: true,
})

if (DOMPurify.isSupported) {
  DOMPurify.setConfig({
    USE_PROFILES: {
      html: true,
    },
    ADD_ATTR: ["target"],
  })

  let allow: RegExp[] | undefined
  DOMPurify.addHook("uponSanitizeAttribute", (currentNode, hookEvent) => {
    if (allow === undefined && window.MARKDOWN_ALLOW !== undefined)
      allow = window.MARKDOWN_ALLOW.split("\n")
        .filter((s) => s)
        .map((s) => new RegExp(s, "i"))

    if (allow && hookEvent.attrName === "href") {
      const href = currentNode.getAttribute("href")
      if (href !== null && !href.startsWith("javascript")) hookEvent.forceKeepAttr = allow.some((r) => href.match(r))
    }
  })
}

const link = (href: string, title: string, text: string) => {
  const titleAttr = title ? ` title=${title}` : ""
  return `<a class="text-link" href="${href}"${titleAttr} rel="noopener nofollow" target="_blank">${text}</a>`
}

const fullRenderer = new marked.Renderer()
fullRenderer.image = (href, title, alt) => {
  // Check if this is our special fider-image syntax
  if (href && href.startsWith("fider-image:")) {
    const bkey = href.substring("fider-image:".length)
    return `<img src="/static/images/${bkey}" alt="${alt || ""}" class="fider-inline-image" data-bkey="${bkey}" />`
  }
  return "" // Ignore other images
}
fullRenderer.link = link
fullRenderer.text = (text: string) => {
  // Handling mention links (they're in the format @[name])
  return text.replace(/@\[(.*?)\]/g, (match, name) => {
    return `<span class="text-blue-600">@${name}</span>`
  })
}

const plainTextRenderer = new marked.Renderer()
plainTextRenderer.link = (_href, _title, text) => text
plainTextRenderer.image = () => ""
plainTextRenderer.br = () => " "
plainTextRenderer.strong = (text) => text
plainTextRenderer.list = (body) => body
plainTextRenderer.listitem = (text) => `${text} `
plainTextRenderer.heading = (text) => text
plainTextRenderer.paragraph = (text) => ` ${text} `
plainTextRenderer.code = (code) => code
plainTextRenderer.codespan = (code) => code
plainTextRenderer.html = (html) => html
plainTextRenderer.del = (text) => text

const entities: { [key: string]: string } = {
  // "<": "&lt;",
  // ">": "&gt;",
}

const encodeHTML = (s: string) => s.replace(/[<>]/g, (tag) => entities[tag] || tag)
const sanitize = (input: string) => (DOMPurify.isSupported ? DOMPurify.sanitize(input) : input)

export const full = (input: string): string => {
  return sanitize(marked(encodeHTML(input), { renderer: fullRenderer }).trim())
}

export const plainText = (input: string): string => {
  return sanitize(marked(encodeHTML(input), { renderer: plainTextRenderer }).trim())
}
