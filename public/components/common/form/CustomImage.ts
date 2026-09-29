import Image from "@tiptap/extension-image"
import Paragraph from "@tiptap/extension-paragraph"
import { JSONContent, mergeAttributes } from "@tiptap/core"
import { fiderImageBkey, fiderImageMarkdown } from "@fider/services/markdownSyntax"

export interface CustomImageOptions {
  HTMLAttributes?: Record<string, any>
  allowBase64?: boolean
  // Returns the src (e.g. a data: URL) of an image that was attached but isn't saved yet, or ""
  // when the bkey refers to a stored image (served from /static/images/<bkey>).
  onGetImageSrc?: (bkey: string) => string
}

// marked image token shape (the parts we read)
type ImageToken = { href?: string; text?: string; title?: string | null }

// tiptap's Paragraph unwraps a paragraph holding only an image, expecting a block image node.
// Fider's image is inline, so an unwrapped image would sit directly in the document: invalid
// content that makes every later edit fail. Posts often end with such paragraphs (the server
// appends unreferenced attachments as "![](fider-image:<bkey>)"), so keep them as paragraphs.
// Use in place of StarterKit's paragraph.
const parseParagraph = Paragraph.config.parseMarkdown
export const ImageParagraph = Paragraph.extend({
  parseMarkdown: (token, helpers) => {
    const tokens = token.tokens || []
    if (tokens.length === 1 && tokens[0].type === "image") {
      return helpers.createNode("paragraph", undefined, helpers.parseInline(tokens))
    }
    return parseParagraph ? parseParagraph(token, helpers) : helpers.createNode("paragraph", undefined, helpers.parseInline(tokens))
  },
})

export const CustomImage = Image.extend<CustomImageOptions>({
  name: "customImage",

  // marked tokenizes images as INLINE tokens, so the node must be inline to sit inside
  // paragraph content (a block image would be dropped during markdown parse). This also
  // matches Fider's `fider-inline-image` rendering.
  inline: true,
  group: "inline",

  addOptions() {
    return {
      ...this.parent?.(),
      HTMLAttributes: {},
      allowBase64: true,
      onGetImageSrc: undefined,
    }
  },

  addAttributes() {
    return {
      ...this.parent?.(),
      id: {
        default: null,
        parseHTML: (element) => element.getAttribute("data-id"),
        renderHTML: (attributes) => {
          if (!attributes.id) {
            return {}
          }
          return {
            "data-id": attributes.id,
          }
        },
      },
      bkey: {
        default: null,
        parseHTML: (element) => element.getAttribute("data-bkey"),
        renderHTML: (attributes) => {
          if (!attributes.bkey) {
            return {}
          }
          return {
            "data-bkey": attributes.bkey,
          }
        },
      },
    }
  },

  // Unsaved uploads are referenced by bkey in the markdown but don't exist on the server yet,
  // so resolve their src when rendering. (parseMarkdown can't do this: `this` there is not the
  // extension, so options aren't available.) Also show safe fider-image references that aren't
  // editor images (e.g. ![caption](fider-image:...), see parseMarkdown) the way the renderer does.
  renderHTML({ node, HTMLAttributes }) {
    const bkey: string | undefined = node.attrs.bkey || fiderImageBkey(node.attrs.src)
    const pendingSrc = bkey && this.options.onGetImageSrc ? this.options.onGetImageSrc(bkey) : ""
    let src = HTMLAttributes.src
    if (pendingSrc) {
      src = pendingSrc
    } else if (bkey && !node.attrs.bkey) {
      src = `/static/images/${bkey}`
    }
    return ["img", mergeAttributes(this.options.HTMLAttributes || {}, HTMLAttributes, { src })]
  },

  // --- @tiptap/markdown integration (marked engine) ---
  // Handle the standard marked "image" token, detecting Fider's ![](fider-image:<bkey>) syntax.
  markdownTokenName: "image",

  parseMarkdown(token: ImageToken): JSONContent {
    // Note: `this` inside parseMarkdown is NOT the extension (no this.name / this.options),
    // so use a literal node type.
    // Only the exact form the editor writes, ![](fider-image:<safe bkey>), becomes an editor
    // image (tracked as an attachment and written back without alt text). Anything else stays
    // a plain image with its original src/alt so it round-trips unchanged.
    const bkey = token.text ? undefined : fiderImageBkey(token.href)
    if (bkey) {
      return {
        type: "customImage",
        attrs: { src: `/static/images/${bkey}`, alt: "", id: bkey, bkey },
      }
    }
    return {
      type: "customImage",
      attrs: { src: token.href || "", alt: token.text || "", title: token.title || null, id: null, bkey: null },
    }
  },

  renderMarkdown(node: JSONContent): string {
    const attrs = node.attrs || {}
    if (attrs.bkey || attrs.id) {
      // Fider inline-image syntax; use bkey if available, otherwise id.
      return fiderImageMarkdown(attrs.bkey || attrs.id)
    }
    const title = attrs.title ? ` "${String(attrs.title).replace(/"/g, '\\"')}"` : ""
    return `![${attrs.alt || ""}](${attrs.src || ""}${title})`
  },

  // Override the addImage command to include our custom attributes and handle uploads
  addCommands() {
    return {
      setImage:
        (options: any) =>
        ({ tr, dispatch }: { tr: any; dispatch: any }) => {
          const { src, alt, title, id, bkey } = options

          // Create a node with our custom attributes
          const node = this.type.create({
            src,
            alt,
            title,
            id,
            bkey,
          })

          if (dispatch) {
            tr.replaceSelectionWith(node)
          }

          return true
        },
    }
  },

  // Add event handlers to the editor
  onSelectionUpdate() {
    // When an image is selected, we can add a delete button or handle keyboard events
    // This is a placeholder for future implementation
  },
})
