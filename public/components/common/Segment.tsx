import "./Segment.scss"

import React from "react"

interface SegmentProps {
  className?: string
}

export const Segments: React.FunctionComponent<SegmentProps> = (props) => {
  return <div className={`c-segments ${props.className || ""}`}>{props.children}</div>
}

export const Segment: React.FunctionComponent<SegmentProps> = (props) => {
  return <div className={`c-segment ${props.className || ""}`}>{props.children}</div>
}
