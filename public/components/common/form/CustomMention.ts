import Mention from "@tiptap/extension-mention"
import { JSONContent, MarkdownToken } from "@tiptap/core"

export const CustomMention = Mention.extend({
  name: "mention",

  // --- @tiptap/markdown integration (marked engine) ---
  // @[name] is non-standard markdown, so register a custom marked tokenizer for it.
  markdownTokenizer: {
    name: "mention",
    level: "inline" as const,
    start(src: string) {
      return src.indexOf("@[")
    },
    tokenize(src: string): MarkdownToken | undefined {
      const match = /^@\[(.+?)\]/.exec(src)
      if (!match) {
        return undefined
      }
      return { type: "mention", raw: match[0], label: match[1] }
    },
  },

  markdownTokenName: "mention",

  parseMarkdown(token: MarkdownToken): JSONContent {
    // Note: `this` inside parseMarkdown is NOT the extension (this.name is undefined),
    // so the node type must be a literal.
    const label = (token.label as string) || ""
    return { type: "mention", attrs: { id: label, label } }
  },

  renderMarkdown(node: JSONContent): string {
    return `@[${node.attrs?.label ?? ""}]`
  },
})
