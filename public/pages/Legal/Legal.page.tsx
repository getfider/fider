import { Markdown } from "@fider/components"
import React from "react"
import "./Legal.page.scss"

export interface LegalPageProps {
  content: string
}

const LegalPage = (props: LegalPageProps) => {
  return (
    <div id="p-legal" className="page container">
      <Markdown text={props.content} style="full" />
    </div>
  )
}

export default LegalPage
