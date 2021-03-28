import marked from "marked"
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
}

const link = (href: string, title: string, text: string) => {
  const titleAttr = title ? ` title=${title}` : ""
  return `<a href="${href}"${titleAttr} rel="noopener" target="_blank">${text}</a>`
}

const fullRenderer = new marked.Renderer()
fullRenderer.image = () => ""
fullRenderer.link = link

const plainTextRenderer = new marked.Renderer()
plainTextRenderer.link = (_href, _title, text) => text
plainTextRenderer.image = (_href, _title, _text) => ""
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
const sanitize = (input: string) => (DOMPurify.isSupported ? DOMPurify.sanitize(input) : input)

export const full = (input: string): string => {
  return sanitize(marked(encodeHTML(input), { renderer: fullRenderer }).trim())
}

export const plainText = (input: string): string => {
  return sanitize(marked(encodeHTML(input), { renderer: plainTextRenderer }).trim())
}
