import { JSONContent } from "@tiptap/core"
import { Markdown, MarkdownExtensionOptions, MarkdownExtensionStorage } from "@tiptap/markdown"
import { Marked, marked } from "marked"

// A marked instance for one editor. @tiptap/markdown defaults to the global `marked` singleton
// and registers extension tokenizers (e.g. mentions) on it for every editor created, so each
// editor gets its own instance instead.
const createEditorMarked = (): Marked => {
  const instance = new Marked({ gfm: true, breaks: true })
  // Equivalent of the old tiptap-markdown `html: false`: raw HTML is never tokenized, so tags
  // typed in markdown mode stay literal text instead of becoming rich-text nodes.
  instance.use({
    tokenizer: {
      html: () => undefined,
      tag: () => undefined,
    },
  })
  return instance
}

// MarkdownManager internals used below (private in its type declarations).
interface MarkdownManagerInternals {
  encodeTextForMarkdown(text: string, node: JSONContent, parentNode?: JSONContent): string
  escapeMarkdownSyntax(text: string): string
}

const FiderMarkdownExtension = Markdown.extend<MarkdownExtensionOptions, MarkdownExtensionStorage>({
  onBeforeCreate(event) {
    this.parent?.(event)

    const manager = this.editor.markdown as unknown as MarkdownManagerInternals | undefined
    if (!manager) return

    // @tiptap/markdown HTML-encodes text when serializing (AT&T -> AT&amp;T, a > b -> a &gt; b),
    // which would change the stored markdown. Raw HTML is never parsed (see above) and text
    // tokens are entity-decoded on load, so store special characters as typed, like the old
    // serializer did. Markdown syntax characters are still backslash-escaped, and a leading ">"
    // is escaped so the text can't turn into a blockquote when loaded again.
    const encode = manager.encodeTextForMarkdown.bind(manager)
    manager.encodeTextForMarkdown = (text, node, parentNode) => {
      // Code text (and text with nothing to escape) comes back unchanged: keep it as is.
      if (encode(text, node, parentNode) === text) return text
      return manager.escapeMarkdownSyntax(text).replace(/^(\s{0,3})>/, "$1\\>")
    }
  },
})

// The Markdown extension as configured for Fider's comment editor. Call once per editor.
export const fiderMarkdown = () => FiderMarkdownExtension.configure({ marked: createEditorMarked() as unknown as typeof marked })
