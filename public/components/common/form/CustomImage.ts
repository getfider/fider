import Image from "@tiptap/extension-image"
import * as MarkdownIt from "markdown-it"
import { ImageUpload } from "@fider/models"

// This is a unique ID generator for images
const generateImageId = () => `img_${Math.random().toString(36).substring(2, 15)}`

export interface CustomImageOptions {
  HTMLAttributes?: Record<string, any>
  allowBase64?: boolean
  onImageUpload?: (upload: ImageUpload) => void
  onImageRemove?: (id: string) => void
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
          const imageId = node.attrs.bkey || node.attrs.id || generateImageId()
          state.write(`![](fider-image:${imageId})`)
        },
        parse: {
          setup(markdownit: MarkdownIt) {
            // Custom rule to parse our special image syntax
            markdownit.inline.ruler.before("text", "fider-image", (state: MarkdownIt.StateInline, silent: boolean) => {
              const match = state.src.slice(state.pos).match(/^!\[\]\(fider-image:([a-zA-Z0-9_/.-]+)\)/)
              if (!match) return false

              if (!silent) {
                const id = match[1]
                const token = state.push("image", "img", 0)
                token.attrs = [
                  ["src", ""], // Empty src, will be filled by the renderer
                  ["data-id", id],
                  ["data-bkey", id], // Store as both id and bkey
                  ["alt", ""],
                ]
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
        ({ tr, dispatch }) => {
          const { src, alt, title, id } = options

          // Generate a unique ID for this image if not provided
          const imageId = id || generateImageId()

          // Create a node with our custom attributes
          const node = this.type.create({
            src,
            alt,
            title,
            id: imageId,
          })

          if (dispatch) {
            tr.replaceSelectionWith(node)
          }

          return true
        },
    }
  },
})
