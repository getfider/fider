import React, { useState } from "react"
import { classSet } from "@fider/services"
import { DisplayError, ValidationContext, hasError } from "../"
import { HStack } from "@fider/components/layout"

import "./Checkbox.scss"

interface CheckboxProps {
  field: string
  checked?: boolean
  onChange?: (checked: boolean) => void
}

export const Checkbox: React.FC<CheckboxProps> = (props) => {
  const [checked, setChecked] = useState<boolean>(props.checked || false)

  const onChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const isChecked: boolean = e.currentTarget.checked

    setChecked(isChecked)
    if (props.onChange) {
      props.onChange(isChecked)
    }
  }

  return (
    <ValidationContext.Consumer>
      {(ctx) => (
        <div
          className={classSet({
            "c-form-field": true,
            "m-error": hasError(props.field, ctx.error),
          })}
        >
          <div className="c-checkbox">
            <HStack>
              <input id={`input-${props.field}`} type="checkbox" checked={checked} onChange={onChange} />
              <label htmlFor={`input-${props.field}`} className="text-sm">
                {props.children}
              </label>
            </HStack>
            <DisplayError fields={[props.field]} error={ctx.error} />
          </div>
        </div>
      )}
    </ValidationContext.Consumer>
  )
}
