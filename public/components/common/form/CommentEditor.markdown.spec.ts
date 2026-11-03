import { Editor, JSONContent } from "@tiptap/core"
import Link from "@tiptap/extension-link"
import { CustomImage } from "./CustomImage"
import { CustomMention } from "./CustomMention"
import { fiderMarkdown } from "./FiderMarkdown"
import { RawMarkdown } from "./RawMarkdown"
import { fiderStarterKit } from "./FiderStarterKit"
import * as markdown from "@fider/services/markdown"
import { fiderAllowedSchemes } from "@fider/hooks"

fiderAllowedSchemes.get = () => ""

// Headless editor with the same markdown-relevant extensions CommentEditor uses (minus UI-only
// bits like suggestions and placeholder). Locks the @tiptap/markdown round-trip
// (markdown -> doc -> markdown): stored content must come back byte-for-byte, including
// content the schema can't represent, raw HTML and special characters.
const makeEditor = (imageOptions: Parameters<typeof CustomImage.configure>[0] = {}) =>
  new Editor({
    extensions: [
      ...fiderStarterKit(),
      Link.configure({ openOnClick: true, autolink: true, defaultProtocol: "https" }),
      fiderMarkdown(),
      RawMarkdown,
      CustomMention.configure({ HTMLAttributes: { class: "mention" } }),
      CustomImage.configure({ allowBase64: true, ...imageOptions }),
    ],
  })

const load = (md: string, editor = makeEditor()): Editor => {
  editor.commands.setContent(md, { emitUpdate: false, contentType: "markdown" })
  return editor
}

const roundTrip = (md: string): string => {
  const editor = load(md)
  const out = editor.getMarkdown().trim()
  editor.destroy()
  return out
}

// Serialize a document built directly (as if typed in rich-text mode)
const serialize = (content: JSONContent[]): string => {
  const editor = makeEditor()
  editor.commands.setContent({ type: "doc", content }, { emitUpdate: false })
  const out = editor.getMarkdown().trim()
  editor.destroy()
  return out
}

const paragraph = (text: string): JSONContent => ({ type: "paragraph", content: [{ type: "text", text }] })

const nodesOfType = (editor: Editor, type: string): JSONContent[] => {
  const found: JSONContent[] = []
  const walk = (node: JSONContent) => {
    if (node.type === type) found.push(node)
    ;(node.content || []).forEach(walk)
  }
  walk(editor.getJSON())
  return found
}

// Parse editor HTML the way the browser would and list the elements it creates.
const elementsIn = (html: string): string[] => {
  const container = document.createElement("div")
  container.innerHTML = html
  return Array.from(container.querySelectorAll("*")).map((e) => e.tagName.toLowerCase())
}

describe("CommentEditor markdown round-trip (@tiptap/markdown)", () => {
  const cases: Array<[string, string]> = [
    ["mention", "@[Jane Doe]"],
    ["mention in text", "Hey @[Jane Doe], welcome aboard"],
    ["two mentions", "@[Jane Doe] and @[John Smith]"],
    ["fider image", "![](fider-image:attachments/zy0hBtqrjQki7M56p26AuAXljRoaNUSwZO6MOky5gnYm2nW1rsMmrp3dwhjGk7ok-aden.jpeg)"],
    ["fider image in text", "look ![](fider-image:attachments/abc-x.jpeg) here"],
    ["plain image", "![](http://demo.dev.fider.io:3000/images/100/28)"],
    ["plain image with alt and title", '![alt](http://demo.dev.fider.io:3000/images/100/28 "title")'],
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
    // Content the schema has no node for must survive untouched (it used to be deleted)
    ["gfm table", "| a | b |\n| --- | --- |\n| 1 | 2 |"],
    ["table between paragraphs", "Before\n\n| a |\n| - |\n| 1 |\n\nAfter"],
    ["task list", "- [ ] task\n- [x] done"],
    ["mixed task list", "- a\n- [ ] b"],
    ["nested task list", "- a\n  - [ ] nested\n- c"],
    // Raw HTML is literal text, never markup
    ["inline html", "<b>hi</b> <img src=x onerror=alert(1)>"],
    ["block html", "<div>hi</div>"],
    ["html comment", "<!-- note -->"],
    // Special characters are stored as typed, not as HTML entities
    ["ampersand and angle brackets", "AT&T > x < y"],
    ["code keeps entity-like text", "`&lt;b&gt; & <i>`"],
    // Only safe fider-image references become inline images; anything else is kept as-is
    ["fider image with caption", "![caption](fider-image:attachments/abc.png)"],
    ["fider image path traversal", "![caption](fider-image:../../api/v1/admin/x?y)"],
  ]

  cases.forEach(([name, md]) => {
    test(`${name} round-trips byte-identically: ${md}`, () => {
      expect(roundTrip(md)).toEqual(md)
    })
  })
})

describe("CommentEditor raw HTML (html: false)", () => {
  test("typed HTML is loaded as literal text", () => {
    const editor = load("<b>hi</b> <img src=x onerror=alert(1)>\n\n<div>block</div>")
    expect(elementsIn(editor.getHTML())).toEqual(["p", "p"])
    expect(editor.getText()).toEqual("<b>hi</b> <img src=x onerror=alert(1)>\n\n<div>block</div>")
    editor.destroy()
  })
})

describe("CommentEditor unsupported markdown", () => {
  test("a table is kept as literal text in the editor", () => {
    const editor = load("| a | b |\n| --- | --- |\n| 1 | 2 |")
    expect(editor.getText()).toContain("| 1 | 2 |")
    editor.destroy()
  })

  test("task-list checkbox state is kept", () => {
    const editor = load("- [ ] task\n- [x] done")
    expect(editor.getText()).toContain("[x] done")
    editor.destroy()
  })
})

describe("CommentEditor markdown serialization", () => {
  test("special characters typed in rich-text mode are stored as typed", () => {
    expect(serialize([paragraph("AT&T <b>bold?</b> a > b")])).toEqual("AT&T <b>bold?</b> a > b")
  })

  test("text that looks like a blockquote stays a paragraph", () => {
    const md = serialize([paragraph("> not a quote")])
    const editor = load(md)
    expect(nodesOfType(editor, "blockquote")).toEqual([])
    expect(editor.getText()).toEqual("> not a quote")
    editor.destroy()
  })
})

describe("CommentEditor fider images", () => {
  const unsafe = ["![](fider-image:../../api/v1/admin/x)", "![caption](fider-image:../../api/v1/admin/x?y)", "![](fider-image:attachments/../x.png)"]

  unsafe.forEach((md) => {
    test(`does not become a fider image: ${md}`, () => {
      const editor = load(md)
      const images = nodesOfType(editor, "customImage")
      images.forEach((img) => {
        expect(img.attrs?.bkey).toBeNull()
        expect(img.attrs?.src).not.toContain("/static/images/")
      })
      editor.destroy()
    })
  })

  test("a pending (unsaved) upload is shown from its cached data, not /static/images", () => {
    const bkey = "attachments/pending-image.png"
    const onGetImageSrc = (key: string) => (key === bkey ? "data:image/png;base64,AAAA" : "")
    const editor = load(`look ![](fider-image:${bkey}) and ![](fider-image:attachments/saved.png)`, makeEditor({ onGetImageSrc }))

    const html = editor.getHTML()
    expect(html).toContain('src="data:image/png;base64,AAAA"')
    expect(html).toContain('src="/static/images/attachments/saved.png"')
    expect(editor.getMarkdown().trim()).toEqual(`look ![](fider-image:${bkey}) and ![](fider-image:attachments/saved.png)`)
    editor.destroy()
  })
})

describe("CommentEditor and renderer agree on mentions", () => {
  const inputs = ["@[Jane Doe]", "@[] and @[Jane Doe]", "@[a]b] c", "x @[Jane] y @[John Smith] z", "@[unclosed and @[Jane]"]

  inputs.forEach((md) => {
    test(`same mentions: ${md}`, () => {
      const editor = load(md)
      const editorLabels = nodesOfType(editor, "mention").map((m) => `@${m.attrs?.label}`)
      editor.destroy()

      const container = document.createElement("div")
      container.innerHTML = markdown.full(md)
      const rendererLabels = Array.from(container.querySelectorAll(".mention")).map((e) => e.textContent)

      expect(editorLabels).toEqual(rendererLabels)
    })
  })
})

describe("block toggles over the whole document", () => {
  // tiptap v3's trailing empty paragraph used to be swept into Select All toggles, leaving empty
  // list items, quotes and code blocks behind.
  const commands = ["toggleBulletList", "toggleOrderedList", "toggleBlockquote", "toggleCodeBlock"] as const

  commands.forEach((command) => {
    test(`Select All + ${command} twice leaves no empty blocks`, () => {
      const editor = load("sample text")
      for (let i = 0; i < 2; i++) {
        editor.commands.selectAll()
        editor.commands[command]()
      }

      const emptyTextblocks: string[] = []
      editor.state.doc.descendants((node) => {
        if (node.isTextblock && node.content.size === 0) emptyTextblocks.push(node.type.name)
      })
      expect(emptyTextblocks).toEqual([])
      expect(editor.state.doc.textContent).toEqual("sample text")
      editor.destroy()
    })
  })

  test("a raw-markdown block at the end still gets a paragraph after it to type into", () => {
    const editor = load("| a | b |\n| --- | --- |\n| 1 | 2 |")
    expect(editor.state.doc.lastChild?.type.name).toEqual("paragraph")
    expect(editor.getMarkdown().trim()).toEqual("| a | b |\n| --- | --- |\n| 1 | 2 |")
    editor.destroy()
  })
})
