import { ReactRenderer } from "@tiptap/react"
import { actions } from "@fider/services"

import MentionList, { MentionListHandle } from "./MentionList"
import { MentionNodeAttrs } from "@tiptap/extension-mention"
interface MentionListProps {
  items: any[]
  command?: (item: MentionNodeAttrs) => void
}

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
    let containerElement: HTMLElement | null = null
    let scrollListener: EventListener | null = null
    let clickOutsideListener: EventListener | null = null
    let initialPosition: { top: number; left: number } | null = null

    return {
      onStart: (props: { editor: any; clientRect?: (() => DOMRect | null) | null }) => {
        reactRenderer = new ReactRenderer(MentionList, {
          props,
          editor: props.editor,
        })

        if (!props.clientRect) {
          return
        }
        // Add click outside listener
        clickOutsideListener = (event: Event) => {
          if (event instanceof MouseEvent && containerElement && event.target instanceof Node && !containerElement.contains(event.target)) {
            // Click was outside the container, clean up
            if (scrollListener) {
              window.removeEventListener("scroll", scrollListener)
            }
            if (clickOutsideListener) {
              document.removeEventListener("click", clickOutsideListener)
            }
            if (containerElement && containerElement.parentNode) {
              document.body.removeChild(containerElement)
              containerElement = null
            }

            // Important: Return focus to editor
            props.editor.view.focus()
          }
        }
        // Use setTimeout to avoid immediate trigger when creating the popup
        setTimeout(() => {
          document.addEventListener("click", clickOutsideListener as EventListener)
        }, 100)

        // Create container for the suggestion list
        containerElement = document.createElement("div")
        containerElement.style.position = "absolute"
        containerElement.style.zIndex = "1000"
        document.body.appendChild(containerElement)
        containerElement.appendChild(reactRenderer.element)

        // Get initial position
        const rect = props.clientRect()
        if (rect) {
          // Store initial position relative to the document
          initialPosition = {
            top: rect.bottom + window.scrollY,
            left: rect.left + window.scrollX,
          }

          // Set initial position
          containerElement.style.left = `${initialPosition.left}px`
          containerElement.style.top = `${initialPosition.top}px`
        }

        // Add scroll listener to maintain position during scrolling
        scrollListener = () => {
          if (containerElement && initialPosition) {
            containerElement.style.top = `${initialPosition.top}px`
          }
        }

        window.addEventListener("scroll", scrollListener, { passive: true })
      },
      onUpdate(props: { clientRect?: (() => DOMRect | null) | null | undefined }) {
        if (!props.clientRect || !containerElement) {
          return
        }

        const rect = props.clientRect()
        if (!rect) {
          return
        }

        reactRenderer.updateProps(props)

        // Update position
        initialPosition = {
          top: rect.bottom + window.scrollY,
          left: rect.left + window.scrollX,
        }

        containerElement.style.left = `${initialPosition.left}px`
        containerElement.style.top = `${initialPosition.top}px`
      },

      onKeyDown(props: { event: KeyboardEvent }) {
        if (props.event.key === "Escape" && containerElement) {
          // Clean up
          if (scrollListener) {
            window.removeEventListener("scroll", scrollListener)
          }

          document.body.removeChild(containerElement)
          containerElement = null
          return true
        }

        return reactRenderer.ref?.onKeyDown(props) || false
      },

      onExit() {
        if (containerElement) {
          if (scrollListener) {
            window.removeEventListener("scroll", scrollListener)
          }
          if (clickOutsideListener) {
            document.removeEventListener("click", clickOutsideListener)
          }

          document.body.removeChild(containerElement)
          containerElement = null
        }

        reactRenderer.destroy()
      },
    }
  },
}
