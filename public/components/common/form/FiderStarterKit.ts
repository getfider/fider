import StarterKit from "@tiptap/starter-kit"
import { ImageParagraph } from "./CustomImage"

// StarterKit as configured for Fider's comment editor, shared with the markdown specs so they
// exercise the same schema. Link is configured by CommentEditor, underline has no markdown
// form, and paragraph is replaced by ImageParagraph (see CustomImage).
//
// TrailingNode (new in tiptap v3) keeps an empty paragraph after the last block. Select All then
// sweeps it into block toggles, leaving stray empty list items, quotes or code blocks in the
// markdown. Headings, lists, quotes and code blocks can already be left with Enter or the arrow
// keys (as in tiptap v2), so only keep the trailing paragraph after a raw-markdown block.
export const fiderStarterKit = () => [
  StarterKit.configure({
    link: false,
    underline: false,
    paragraph: false,
    trailingNode: { notAfter: ["heading", "bulletList", "orderedList", "blockquote", "codeBlock", "horizontalRule"] },
  }),
  ImageParagraph,
]
