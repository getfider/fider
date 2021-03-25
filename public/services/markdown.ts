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

const simpleRenderer = new marked.Renderer()
simpleRenderer.heading = (_text, _level, raw) => `<p>${raw}</p>`
simpleRenderer.image = () => ""
simpleRenderer.link = link

const fullRenderer = new marked.Renderer()
fullRenderer.link = link

const entities: { [key: string]: string } = {
  "<": "&lt;",
  ">": "&gt;",
}

const encodeHTML = (s: string) => s.replace(/[<>]/g, (tag) => entities[tag] || tag)
const sanitize = (input: string) => (DOMPurify.isSupported ? DOMPurify.sanitize(input) : input)

export const full = (input: string): string => {
  return sanitize(marked(encodeHTML(input), { renderer: fullRenderer }).trim())
}

export const simple = (input: string): string => {
  return sanitize(marked(encodeHTML(input), { renderer: simpleRenderer }).trim())
}
