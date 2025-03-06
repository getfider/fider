import { ReactRenderer } from "@tiptap/react"
import tippy from "tippy.js"
import { actions } from "@fider/services"

import MentionList, { MentionListHandle } from "./MentionList"
import { MentionNodeAttrs } from "@tiptap/extension-mention"
interface MentionListProps {
  items: any[]
  command?: (item: MentionNodeAttrs) => void
}
// interface ReactRendererOptions {
//   props: { editor: any; clientRect?: (() => DOMRect | null) | null }
//   editor: any
// }

// Cache for storing users
let cachedUsers: MentionNodeAttrs[] = []

export default {
  items: async ({ query }: { query: string }) => {
    // If we don't have cached users yet, fetch them
    if (cachedUsers.length === 0) {
      const result = await actions.getTaggableUsers("")
      if (result.ok) {
        cachedUsers = result.data.map((user) => ({ id: user.id.toString(), label: user.name }))
      }
    }

    // Filter the cached users based on the query
    return cachedUsers.filter((item) => item.label?.toLowerCase().startsWith(query.toLowerCase())).slice(0, 5)
  },
  render: () => {
    let reactRenderer: ReactRenderer<MentionListHandle, MentionListProps>
    let popup: any

    return {
      onStart: (props: { editor: any; clientRect?: (() => DOMRect | null) | null }) => {
        reactRenderer = new ReactRenderer(MentionList, {
          props,
          editor: props.editor,
        })

        if (!props.clientRect) {
          return
        }

        popup = tippy(document.body, {
          getReferenceClientRect: props.clientRect as () => DOMRect,
          appendTo: () => document.body,
          content: reactRenderer.element,
          showOnCreate: true,
          interactive: true,
          trigger: "manual",
          placement: "bottom-start",
        })
      },
      onUpdate(props: { clientRect?: (() => DOMRect | null) | null | undefined }) {
        if (!props.clientRect) {
          return
        }

        const rect = props.clientRect()
        if (!rect) {
          return
        }

        reactRenderer.updateProps(props)

        popup.setProps({
          getReferenceClientRect: props.clientRect,
        })
      },

      onKeyDown(props: { event: KeyboardEvent }) {
        if (props.event.key === "Escape") {
          popup.hide()

          return true
        }

        return reactRenderer.ref?.onKeyDown(props) || false
      },

      onExit() {
        popup.destroy()
        reactRenderer.destroy()
      },
    }
  },
}
