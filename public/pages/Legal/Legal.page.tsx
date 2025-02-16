import { Header, Markdown } from "@fider/components"
import React from "react"
import "./Legal.page.scss"

export interface LegalPageProps {
  content: string
}

const LegalPage = (props: LegalPageProps) => {
  return (
    <>
      <Header />
      <div id="p-legal" className="page container w-max-10xl">
        <Markdown text={props.content} style="full" />
      </div>
    </>
  )
}

export default LegalPage
