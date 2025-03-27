import Mention from "@tiptap/extension-mention"
import * as MarkdownIt from "markdown-it"

export const CustomMention = Mention.extend({
  name: "mention",
  addStorage() {
    return {
      markdown: {
        serialize: (state: any, node: { attrs: { id: string; label: string } }) => {
          state.write(`@[${node.attrs.label}]`)
        },
        parse: {
          setup(markdownit: MarkdownIt) {
            markdownit.renderer.rules.mention_open = (tokens, idx) => {
              const token = tokens[idx]
              const id = token.attrGet("id")
              const label = token.attrGet("label")
              return `<span data-type="mention" data-id="${id}" data-label="${label}" class="mention">`
            }

            markdownit.renderer.rules.mention_close = () => {
              return "</span>"
            }

            markdownit.inline.ruler.before("text", "mention", (state: MarkdownIt.StateInline, silent: boolean) => {
              const match = state.src.slice(state.pos).match(/^@\[(.+?)\]/)
              if (!match) return false
              if (!silent) {
                const label = match[1]
                const token = state.push("mention_open", "span", 1)
                token.attrs = [["label", label]]

                // Add the text content token
                const contentToken = state.push("text", "", 0)
                contentToken.content = `@${label}`

                // Close the mention span
                state.push("mention_close", "span", -1)
              }

              state.pos += match[0].length
              return true
            })
          },
        },
      },
    }
  },
})
