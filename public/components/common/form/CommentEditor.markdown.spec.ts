import { Editor } from "@tiptap/core"
import StarterKit from "@tiptap/starter-kit"
import Link from "@tiptap/extension-link"
import { Markdown } from "@tiptap/markdown"
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
      Markdown.configure({ markedOptions: { breaks: true, gfm: true } }),
      CustomMention.configure({ HTMLAttributes: { class: "mention" } }),
      CustomImage.configure({ allowBase64: true }),
    ],
  })

const roundTrip = (md: string): string => {
  const editor = makeEditor()
  editor.commands.setContent(md, { emitUpdate: false, contentType: "markdown" })
  const out = editor.getMarkdown().trim()
  editor.destroy()
  return out
}

describe("CommentEditor markdown round-trip (@tiptap/markdown)", () => {
  const cases: Array<[string, string]> = [
    ["mention", "@[Jane Doe]"],
    ["mention in text", "Hey @[Jane Doe], welcome aboard"],
    ["two mentions", "@[Jane Doe] and @[John Smith]"],
    ["fider image", "![](fider-image:attachments/zy0hBtqrjQki7M56p26AuAXljRoaNUSwZO6MOky5gnYm2nW1rsMmrp3dwhjGk7ok-aden.jpeg)"],
    ["fider image in text", "look ![](fider-image:attachments/abc-x.jpeg) here"],
    ["plain image", "![](http://demo.dev.fider.io:3000/images/100/28)"],
    ["bold", "**bold**"],
    ["italic", "*italic*"],
    ["strike", "~~struck~~"],
    ["inline code", "`code`"],
    ["mixed inline", "a **bold** and *italic* and `code` and @[Jane Doe]"],
    ["h2", "## Heading two"],
    ["h3", "### Heading three"],
    ["bullet list", "- one\n- two"],
    ["ordered list", "1. one\n2. two"],
    ["blockquote", "> quoted"],
    ["link", "[GitHub](https://github.com)"],
    ["link with mention", "see [GitHub](https://github.com) cc @[Jane Doe]"],
    ["multi paragraph", "First paragraph.\n\nSecond paragraph."],
    // Hard breaks serialize to the standard two-space form. This is render-equivalent: the
    // marked renderer emits identical HTML (<br>) for bare "\n" and "  \n" under breaks:true,
    // so existing stored content (bare \n) is unaffected.
    ["hard breaks", "line one  \nline two  \nline three"],
    ["heading then text", "## Title\n\nSome body text."],
  ]

  cases.forEach(([name, md]) => {
    test(`${name} round-trips byte-identically: ${md}`, () => {
      expect(roundTrip(md)).toEqual(md)
    })
  })
})
