import { Editor } from "@tiptap/core"
import StarterKit from "@tiptap/starter-kit"
import Link from "@tiptap/extension-link"
import { Markdown, MarkdownStorage } from "tiptap-markdown"
import { CustomImage } from "./CustomImage"
import { CustomMention } from "./CustomMention"

// Headless editor with the same extension set CommentEditor uses (minus UI-only bits).
// This locks the markdown round-trip (markdown -> doc -> markdown) so the later engine
// swap to @tiptap/markdown can be verified byte-for-byte against it.
const makeEditor = () =>
  new Editor({
    extensions: [
      StarterKit.configure({ link: false, underline: false }),
      Link.configure({ openOnClick: true, autolink: true, defaultProtocol: "https" }),
      Markdown.configure({ html: false, breaks: true }),
      CustomMention.configure({ HTMLAttributes: { class: "mention" } }),
      CustomImage.configure({ allowBase64: true }),
    ],
  })

const roundTrip = (md: string): string => {
  const editor = makeEditor()
  editor.commands.setContent(md, { emitUpdate: false })
  const out = (editor.storage as unknown as { markdown: MarkdownStorage }).markdown.getMarkdown()
  editor.destroy()
  return out
}

describe("CommentEditor markdown round-trip (tiptap-markdown 0.9)", () => {
  const cases: Array<[string, string]> = [
    ["mention", "@[Jane Doe]"],
    ["fider image", "![](fider-image:attachments/zy0hBtqrjQki7M56p26AuAXljRoaNUSwZO6MOky5gnYm2nW1rsMmrp3dwhjGk7ok-aden.jpeg)"],
    ["plain image", "![](http://demo.dev.fider.io:3000/images/100/28)"],
    ["bold", "**bold**"],
    ["italic", "*italic*"],
    ["strike", "~~struck~~"],
    ["inline code", "`code`"],
    ["h2", "## Heading two"],
    ["h3", "### Heading three"],
    ["bullet list", "- one\n- two"],
    ["ordered list", "1. one\n2. two"],
    ["blockquote", "> quoted"],
    ["link", "[GitHub](https://github.com)"],
  ]

  cases.forEach(([name, md]) => {
    test(`${name} round-trips byte-identically: ${md}`, () => {
      expect(roundTrip(md)).toEqual(md)
    })
  })
})
