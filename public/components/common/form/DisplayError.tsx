import React from "react"
import { Failure } from "@fider/services"

const arrayToTag = (items: string[]) => {
  return items.map((m) => <li key={m}>{m}</li>)
}

interface DisplayErrorProps {
  error?: Failure
  fields?: string[]
}

export const hasError = (field?: string, error?: Failure): boolean => {
  if (field && error && error.errors) {
    for (const err of error.errors) {
      if (err.field === field) {
        return true
      }
    }
  }
  return false
}

export const DisplayError = (props: DisplayErrorProps) => {
  if (!props.error || !props.error.errors) {
    return null
  }

  const dict = props.error.errors.reduce((result, err) => {
    result[err.field || ""] = result[err.field || ""] || []
    result[err.field || ""].push(err.message)
    return result
  }, {} as { [key: string]: string[] })

  let items: JSX.Element[] = []

  if (dict[""] && !props.fields) {
    items = arrayToTag(dict[""])
  } else if (props.fields) {
    for (const field of props.fields || Object.keys(dict)) {
      if (Object.prototype.hasOwnProperty.call(dict, field)) {
        const tags = arrayToTag(dict[field])
        tags.forEach((t) => items.push(t))
      }
    }
  }

  return items.length > 0 ? (
    <div className={`c-form-field-error`}>
      <ul>{items}</ul>
    </div>
  ) : null
}
