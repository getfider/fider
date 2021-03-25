import React, { useState } from "react"
import { classSet } from "@fider/services"
import { DisplayError, ValidationContext, hasError } from "../"

interface CheckboxProps {
  field: string
  checked?: boolean
  onChange: (checked: boolean) => void
}

export const Checkbox: React.FC<CheckboxProps> = (props) => {
  const [checked, setChecked] = useState<boolean>(props.checked || false)

  const onChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const isChecked: boolean = e.currentTarget.checked

    setChecked(isChecked)
    props.onChange(isChecked)
  }

  return (
    <ValidationContext.Consumer>
      {(ctx) => (
        <div
          className={classSet({
            "c-form-field m-checkbox": true,
            "m-error": hasError(props.field, ctx.error),
          })}
        >
          <label htmlFor={`input-${props.field}`}>
            <input id={`input-${props.field}`} type="checkbox" checked={checked} onChange={onChange} />
            {props.children}
          </label>
          <DisplayError fields={[props.field]} error={ctx.error} />
        </div>
      )}
    </ValidationContext.Consumer>
  )
}
