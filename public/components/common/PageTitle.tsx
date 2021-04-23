import React from "react"
import { classSet } from "@fider/services"

interface PageTitleLogo {
  title: string
  subtitle?: string
  className?: string
}

export const PageTitle = (props: PageTitleLogo) => {
  const className = classSet({
    "mb-4": true,
    [`${props.className}`]: props.className,
  })

  return (
    <div className={className}>
      <div className="text-display2 mb-1">{props.title}</div>
      <div className="text-gray-700">{props.subtitle}</div>
    </div>
  )
}
