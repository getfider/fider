import React from "react"
import { classSet } from "@fider/services"
import { ValidationContext } from "../"
import { DisplayError, hasError } from "./DisplayError"
import Textarea from "react-textarea-autosize"

import "./TextArea.scss"

interface TextAreaProps {
  children?: React.ReactNode
  label?: string
  field: string
  value?: string
  disabled?: boolean
  minRows?: number
  placeholder?: string
  afterLabel?: JSX.Element
  onChange?: (value: string, selectionStart?: number) => void
  onKeyDown?: (e: React.KeyboardEvent<HTMLTextAreaElement>) => void
  inputRef?: React.MutableRefObject<any>
  onFocus?: React.FocusEventHandler<HTMLTextAreaElement>
  className?: string
}

export const TextArea: React.FunctionComponent<TextAreaProps> = (props) => {
  // Original onChange handler
  const onChange = (e: React.FormEvent<HTMLTextAreaElement>) => {
    if (props.onChange) {
      props.onChange(e.currentTarget.value, e.currentTarget.selectionStart)
    }
  }

  // Original onKeyDown handler
  const onKeyDown = (e: React.KeyboardEvent<HTMLTextAreaElement>) => {
    if (props.onKeyDown) {
      props.onKeyDown(e)
    }
  }

  // Original onFocus handler
  const onFocus = (e: React.FocusEvent<HTMLTextAreaElement>) => {
    if (props.onFocus) {
      props.onFocus(e)
    }
  }

  return (
    <ValidationContext.Consumer>
      {(ctx) => (
        <>
          <div className="c-form-field">
            {!!props.label && (
              <label htmlFor={`input-${props.field}`}>
                {props.label}
                {props.afterLabel}
              </label>
            )}
            <Textarea
              className={classSet({
                "c-textarea": true,
                "c-textarea--error": hasError(props.field, ctx.error),
                [props.className || ""]: props.className,
              })}
              id={`input-${props.field}`}
              disabled={props.disabled}
              onChange={(e) => {
                // Clear error for this field when user interacts with it
                if (ctx.clearError && hasError(props.field, ctx.error)) {
                  ctx.clearError(props.field)
                }

                // Call the original onChange handler
                onChange(e)
              }}
              onKeyDown={onKeyDown}
              value={props.value}
              minRows={props.minRows || 3}
              placeholder={props.placeholder}
              ref={props.inputRef}
              onFocus={onFocus}
            />
            <DisplayError fields={[props.field]} error={ctx.error} />
            {props.children}
          </div>
        </>
      )}
    </ValidationContext.Consumer>
  )
}
