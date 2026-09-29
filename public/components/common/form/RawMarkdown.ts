import { Extension, JSONContent, MarkdownToken, mergeAttributes, Node } from "@tiptap/core"

// Markdown the editor schema has no node for (GFM tables and task lists) would otherwise be
// dropped when loaded, losing it on the next save. It is kept verbatim in a plain-text block
// instead: shown (and editable) as the original markdown, and written back unchanged.

const rawMarkdownNode = (raw: string): JSONContent => {
  const text = raw.replace(/\s+$/, "")
  return { type: "rawMarkdown", content: text ? [{ type: "text", text }] : [] }
}

// marked tokenizes "- [ ] task" as a list token whose items are flagged as tasks.
const RawMarkdownTaskList = Extension.create({
  name: "rawMarkdownTaskList",
  // Runs before the bullet/ordered list handlers for the same "list" token.
  priority: 500,

  markdownTokenName: "list",

  parseMarkdown(token: MarkdownToken) {
    const items = (token.items as Array<{ task?: boolean }> | undefined) || []
    if (!items.some((item) => item.task)) {
      return [] // a normal list: an empty result passes it on to BulletList / OrderedList
    }
    return rawMarkdownNode(token.raw || "")
  },
})

export const RawMarkdown = Node.create({
  name: "rawMarkdown",
  group: "block",
  content: "text*",
  marks: "",
  code: true,
  defining: true,
  whitespace: "pre",

  addExtensions() {
    return [RawMarkdownTaskList]
  },

  parseHTML() {
    return [{ tag: 'div[data-type="raw-markdown"]', preserveWhitespace: "full" as const }]
  },

  renderHTML({ HTMLAttributes }) {
    return ["div", mergeAttributes(HTMLAttributes, { "data-type": "raw-markdown", class: "c-editor-raw-markdown" }), 0]
  },

  // Note: `this` inside parseMarkdown is NOT the extension, so the node type is a literal.
  markdownTokenName: "table",

  parseMarkdown(token: MarkdownToken) {
    return rawMarkdownNode(token.raw || "")
  },

  renderMarkdown(node: JSONContent): string {
    return (node.content || []).map((child) => child.text || "").join("")
  },
})
