import { ReactRenderer } from "@tiptap/react"
import tippy from "tippy.js"

import MentionList from "./MentionList"

export default {
  items: ({ query }: { query: string }) => {
    return [
      "Lea Thompson",
      "Cyndi Lauper",
      "Tom Cruise",
      "Madonna",
      "Jerry Hall",
      "Joan Collins",
      "Winona Ryder",
      "Christina Applegate",
      "Alyssa Milano",
      "Molly Ringwald",
      "Ally Sheedy",
      "Debbie Harry",
      "Olivia Newton-John",
      "Elton John",
      "Michael J. Fox",
      "Axl Rose",
      "Emilio Estevez",
      "Ralph Macchio",
      "Rob Lowe",
      "Jennifer Grey",
      "Mickey Rourke",
      "John Cusack",
      "Matthew Broderick",
      "Justine Bateman",
      "Lisa Bonet",
    ]
      .filter((item) => item.toLowerCase().startsWith(query.toLowerCase()))
      .slice(0, 5)
  },
  render: () => {
    let component: ReactRenderer
    let popup: any

    return {
      onStart: (props: { editor: any; clientRect?: (() => DOMRect | null) | null }) => {
        component = new ReactRenderer(MentionList, {
          props,
          editor: props.editor,
        })

        if (!props.clientRect) {
          return
        }

        popup = tippy(document.body, {
          getReferenceClientRect: props.clientRect as () => DOMRect,
          appendTo: () => document.body,
          content: component.element,
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

        component.updateProps(props)

        popup.setProps({
          getReferenceClientRect: props.clientRect,
        })
      },

      onKeyDown(props: { event: KeyboardEvent }) {
        if (props.event.key === "Escape") {
          popup.hide()

          return true
        }

        // Todo: need to call onKeyDown on the MentionList, but not sure how to get it.
        return false
      },

      onExit() {
        popup.destroy()
        component.destroy()
      },
    }
  },
}
