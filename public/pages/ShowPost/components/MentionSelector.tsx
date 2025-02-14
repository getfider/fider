import React from "react"

import "./MentionSelector.scss"

const MentionSelector: React.FC<{ names: string[]; cursorPosition: { top: number; left: number } }> = ({ names, cursorPosition }) => {
  return (
    <div
      className="c-mention-selector"
      style={{
        top: cursorPosition.top,
        left: cursorPosition.left,
      }}
    >
      {names.map((name, index) => (
        <div key={index} className="c-mention-selector--item clickable hover p-2">
          {name}
        </div>
      ))}
    </div>
  )
}
export default MentionSelector
