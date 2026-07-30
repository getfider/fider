import { ImageUpload } from "@fider/models"
import Image from "@tiptap/extension-image"
import { JSONContent } from "@tiptap/core"

export interface CustomImageOptions {
  HTMLAttributes?: Record<string, any>
  allowBase64?: boolean
  onImageUpload?: (upload: ImageUpload) => void
  onImageRemove?: (bkey: string) => void
  onGetImageSrc?: (bkey: string) => string
}

// marked image token shape (the parts we read)
type ImageToken = { href?: string; text?: string; title?: string | null }

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
      onImageUpload: undefined,
      onImageRemove: undefined,
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

  // --- @tiptap/markdown integration (marked engine) ---
  // Handle the standard marked "image" token, detecting Fider's ![](fider-image:<bkey>) syntax.
  markdownTokenName: "image",

  parseMarkdown(token: ImageToken): JSONContent {
    // Note: `this` inside parseMarkdown is NOT the extension (no this.name / this.options),
    // so use a literal node type. Parsed content is already-stored images, which resolve
    // via the static path (the same fallback the editor used before); live base64 uploads
    // arrive through the setImage command, not markdown parsing.
    const href = token.href || ""
    if (href.startsWith("fider-image:")) {
      const imageId = href.substring("fider-image:".length)
      return {
        type: "customImage",
        attrs: { src: `/static/images/${imageId}`, alt: "", id: imageId, bkey: imageId },
      }
    }
    return {
      type: "customImage",
      attrs: { src: href, alt: token.text || "", id: null, bkey: null },
    }
  },

  renderMarkdown(node: JSONContent): string {
    const attrs = node.attrs || {}
    if (attrs.bkey || attrs.id) {
      // Fider inline-image syntax; use bkey if available, otherwise id.
      return `![](fider-image:${attrs.bkey || attrs.id})`
    }
    return `![${attrs.alt || ""}](${attrs.src || ""})`
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
