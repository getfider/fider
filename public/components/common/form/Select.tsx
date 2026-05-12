import React from "react"
import { classSet } from "@fider/services"
import { ValidationContext } from "./Form"
import { DisplayError, hasError } from "./DisplayError"

import "./Select.scss"

export interface SelectOption {
  value: string
  label: string
}

interface SelectProps {
  children?: React.ReactNode
  field: string
  label?: string
  maxLength?: number
  defaultValue?: string
  options: SelectOption[]
  onChange?: (option?: SelectOption) => void
}

export const Select: React.FunctionComponent<SelectProps> = (props) => {
  const getOption = (value?: string) => {
    if (value && props.options) {
      const filtered = props.options.filter((x) => x.value === value)
      if (filtered && filtered.length > 0) {
        return filtered[0]
      }
    }
  }
  const [selected, setSelected] = React.useState<SelectOption | undefined>(getOption(props.defaultValue))
  // Original onChange handler
  const handleChange = (e: React.FormEvent<HTMLSelectElement>) => {
    let selected: SelectOption | undefined
    if (e.currentTarget.value) {
      const options = props.options.filter((o) => o.value === e.currentTarget.value)
      if (options && options.length > 0) {
        selected = options[0]
      }
    }

    setSelected(selected)
    if (props.onChange) {
      props.onChange(selected)
    }
  }

  return (
    <ValidationContext.Consumer>
      {(ctx) => (
        <>
          <div className="c-form-field">
            {!!props.label && <label htmlFor={`input-${props.field}`}>{props.label}</label>}
            <select
              className={classSet({
                "c-select": true,
                "c-select--error": hasError(props.field, ctx.error),
              })}
              value={selected?.value}
              id={`input-${props.field}`}
              defaultValue={props.defaultValue}
              onChange={(e) => {
                // Clear error for this field when user interacts with it
                if (ctx.clearError && hasError(props.field, ctx.error)) {
                  ctx.clearError(props.field)
                }

                // Call the original onChange handler
                handleChange(e)
              }}
            >
              {props.options.map((option) => (
                <option key={option.value} value={option.value}>
                  {option.label}
                </option>
              ))}
            </select>
            <DisplayError fields={[props.field]} error={ctx.error} />
            {props.children}
          </div>
        </>
      )}
    </ValidationContext.Consumer>
  )
}
