import { ImageUpload } from "@fider/models"
import Image from "@tiptap/extension-image"
import * as MarkdownIt from "markdown-it"

export interface CustomImageOptions {
  HTMLAttributes?: Record<string, any>
  allowBase64?: boolean
  onImageUpload?: (upload: ImageUpload) => void
  onImageRemove?: (bkey: string) => void
  onGetImageSrc?: (bkey: string) => string
}

export const CustomImage = Image.extend<CustomImageOptions>({
  name: "customImage",

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

  addStorage() {
    return {
      images: {},
      markdown: {
        serialize: (state: any, node: any) => {
          // When serializing to markdown, we use a special syntax: ![](fider-image:bkey)
          // Use bkey if available, otherwise fall back to id
          const imageId = node.attrs.bkey || node.attrs.id || ""
          state.write(`![](fider-image:${imageId})`)
        },
        parse: {
          setup: (markdownit: MarkdownIt) => {
            // Custom rule to parse our special image syntax
            markdownit.inline.ruler.before("image", "fider-image", (state: MarkdownIt.StateInline, silent: boolean) => {
              const match = state.src.slice(state.pos).match(/^!\[\]\(fider-image:([a-zA-Z0-9_/.-]+)\)/)
              if (!match) return false

              if (!silent) {
                const imageId = match[1]
                const token = state.push("image", "img", 0)

                // Initialize attrs as an empty array first
                token.attrs = []

                let imageSrc = this.options.onGetImageSrc ? this.options.onGetImageSrc(imageId) : ""
                if (imageSrc.length === 0) {
                  imageSrc = `/static/images/${imageId}`
                }

                token.attrSet("src", imageSrc)
                token.attrSet("alt", "")
                token.attrSet("data-id", imageId)
                token.attrSet("data-bkey", imageId)

                token.children = []
                token.content = ""
              }

              state.pos += match[0].length
              return true
            })
          },
        },
      },
    }
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
