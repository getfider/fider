import React from "react"
import { classSet } from "@fider/services"
import { ValidationContext } from "./Form"
import { DisplayError, hasError } from "./DisplayError"

interface FieldProps {
  className?: string
  label?: string
  field?: string
  afterLabel?: JSX.Element
}

export const Field: React.FunctionComponent<FieldProps> = (props) => {
  const fields = props.field ? [props.field] : undefined
  return (
    <ValidationContext.Consumer>
      {(ctx) => (
        <div
          className={classSet({
            "c-form-field": true,
            "m-error": hasError(props.field, ctx.error),
            [props.className || ""]: props.className,
          })}
        >
          {!!props.label && (
            <label>
              {props.label}
              {props.afterLabel}
            </label>
          )}
          {props.children}
          <DisplayError fields={fields} error={ctx.error} />
        </div>
      )}
    </ValidationContext.Consumer>
  )
}
