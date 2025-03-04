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
    USE_PROFILES: { html: true },
    ADD_TAGS: ["iframe"],
    ADD_ATTR: [
      "allow",
      "allowfullscreen",
      "frameborder",
      "sandbox",
      "src",
      "width",
      "height",
      "title",
      "target",
      "href",
    ]
  })
}

const defaultLink = (href: string, title: string, text: string) => {
  const titleAttr = title ? ` title="${title}"` : ""
  return `<a class="text-link" href="${href}"${titleAttr} rel="noopener nofollow" target="_blank">${text}</a>`
}

const fullRenderer = new marked.Renderer()
fullRenderer.image = () => ""

// Only auto-embed if the href is a valid YouTube URL AND the link text equals the URL (i.e. a naked link).
fullRenderer.link = (href: string | null, title: string, text: string) => {
  if (!href) return text
  const youtubeRegex = /^(?:https?:\/\/)?(?:www\.)?(?:youtube\.com\/watch\?v=|youtu\.be\/)([A-Za-z0-9_-]{11})$/
  const match = href.match(youtubeRegex)
  if (match && text.trim() === href.trim()) {
    const videoId = match[1]
    const embedUrl = `https://www.youtube.com/embed/${videoId}`
    return `<iframe style="width: 100%; height: auto; aspect-ratio: 16/9;" src="${embedUrl}" frameborder="0" allow="accelerometer; clipboard-write; encrypted-media; gyroscope; picture-in-picture" allowfullscreen sandbox="allow-same-origin allow-scripts allow-presentation" title="YouTube video"></iframe>`
  }
  return defaultLink(href, title, text)
}

fullRenderer.text = (text: string) => {
  // Handling mention links (they're in the format @{id:1234, name:'John Doe'})
  return text.replace(/@{([^}]+)}/g, (match) => {
    try {
      const json = match.substring(1).replace(/&quot;/g, '"')
      const mention = JSON.parse(json)
      return `<span class="text-blue-600">@${mention.name}</span>`
    } catch {
      return match
    }
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
  "<": "&lt;",
  ">": "&gt;",
}

const encodeHTML = (s: string) => s.replace(/[<>]/g, (tag) => entities[tag] || tag)
const sanitize = (input: string) => DOMPurify.isSupported ? DOMPurify.sanitize(input) : input

export const full = (input: string): string => {
  return sanitize(marked(encodeHTML(input), { renderer: fullRenderer }).trim())
}

export const plainText = (input: string): string => {
  return sanitize(marked(encodeHTML(input), { renderer: plainTextRenderer }).trim())
}
