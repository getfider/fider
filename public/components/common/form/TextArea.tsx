import React from "react"
import { classSet } from "@fider/services"
import { ValidationContext } from "../"
import { DisplayError, hasError } from "./DisplayError"
import Textarea from "react-textarea-autosize"

interface TextAreaProps {
  label?: string
  field: string
  value?: string
  disabled?: boolean
  minRows?: number
  placeholder?: string
  onChange?: (value: string) => void
  inputRef?: React.MutableRefObject<any>
  onFocus?: React.FocusEventHandler<HTMLTextAreaElement>
}

export const TextArea: React.FunctionComponent<TextAreaProps> = (props) => {
  const onChange = (e: React.FormEvent<HTMLTextAreaElement>) => {
    if (props.onChange) {
      props.onChange(e.currentTarget.value)
    }
  }

  return (
    <ValidationContext.Consumer>
      {(ctx) => (
        <>
          <div
            className={classSet({
              "c-form-field": true,
              "m-error": hasError(props.field, ctx.error),
            })}
          >
            {!!props.label && <label htmlFor={`input-${props.field}`}>{props.label}</label>}
            <div className="c-form-field-wrapper">
              <Textarea
                id={`input-${props.field}`}
                disabled={props.disabled}
                onChange={onChange}
                value={props.value}
                minRows={props.minRows || 3}
                placeholder={props.placeholder}
                ref={props.inputRef}
                onFocus={props.onFocus}
              />
            </div>
            <DisplayError fields={[props.field]} error={ctx.error} />
            {props.children}
          </div>
        </>
      )}
    </ValidationContext.Consumer>
  )
}
