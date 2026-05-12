import { MentionNodeAttrs } from "@tiptap/extension-mention"
import "./MentionList.scss"

import React, { forwardRef, useEffect, useImperativeHandle, useState } from "react"

interface Props {
  items: MentionNodeAttrs[]
  command?: (item: MentionNodeAttrs) => void
}

export interface MentionListHandle {
  onKeyDown: (args: { event: KeyboardEvent }) => boolean
}

const MentionList = forwardRef<MentionListHandle, Props>((props, ref) => {
  const [selectedIndex, setSelectedIndex] = useState(0)

  const selectItem = (index: number) => {
    const item = props.items?.[index]

    if (item) {
      props.command?.(item)
    }
  }
  const upHandler = () => {
    setSelectedIndex((selectedIndex + props.items.length - 1) % props.items.length)
  }

  const downHandler = () => {
    setSelectedIndex((selectedIndex + 1) % props.items.length)
  }

  const enterHandler = () => {
    selectItem(selectedIndex)
  }

  useEffect(() => setSelectedIndex(0), [props.items])

  useImperativeHandle(ref, () => ({
    onKeyDown: ({ event }: { event: KeyboardEvent }): boolean => {
      if (event.key === "ArrowUp") {
        upHandler()
        return true
      }

      if (event.key === "ArrowDown") {
        downHandler()
        return true
      }

      if (event.key === "Enter") {
        enterHandler()
        return true
      }

      return false
    },
  }))

  return (
    <div className="dropdown-menu">
      {props.items.length ? (
        props.items.map((item, index) => (
          <button className={`${index === selectedIndex ? "is-selected" : ""}`} key={index} onClick={() => selectItem(index)}>
            {item.label}
          </button>
        ))
      ) : (
        <div className="item">No result</div>
      )}
    </div>
  )
})

MentionList.displayName = "MentionList"

export default MentionList
