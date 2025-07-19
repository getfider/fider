import React from "react"
import { classSet } from "@fider/services"
import { ValidationContext } from "./Form"
import { DisplayError, hasError } from "./DisplayError"
import { Icon } from "@fider/components"

import "./Input.scss"
import { HStack } from "@fider/components/layout"

interface InputProps {
  children?: React.ReactNode
  field: string
  label?: string
  className?: string
  autoComplete?: string
  autoFocus?: boolean
  noTabFocus?: boolean
  afterLabel?: JSX.Element
  icon?: SpriteSymbol
  maxLength?: number
  value?: string
  disabled?: boolean
  suffix?: string | JSX.Element
  placeholder?: string
  onIconClick?: () => void
  onFocus?: () => void
  inputRef?: React.MutableRefObject<any>
  onChange?: (value: string) => void
}

export const Input: React.FunctionComponent<InputProps> = (props) => {
  // Original onChange handler
  const onChange = (e: React.FormEvent<HTMLInputElement>) => {
    if (props.onChange) {
      props.onChange(e.currentTarget.value)
    }
  }

  // Original onFocus handler
  const onFocus = () => {
    if (props.onFocus) {
      props.onFocus()
    }
  }

  const suffix = typeof props.suffix === "string" ? <span className="c-input__suffix">{props.suffix}</span> : props.suffix

  const icon = props.icon ? <Icon sprite={props.icon} onClick={props.onIconClick} className={classSet({ clickable: !!props.onIconClick })} /> : undefined

  return (
    <ValidationContext.Consumer>
      {(ctx) => (
        <div
          className={classSet({
            "c-form-field": true,
            [`${props.className}`]: props.className,
          })}
        >
          {!!props.label && (
            <label htmlFor={`input-${props.field}`}>
              {props.label}
              {props.afterLabel}
            </label>
          )}
          <HStack spacing={0} align={props.icon ? "center" : undefined} className="relative">
            <input
              className={classSet({
                "c-input": true,
                "c-input--icon": !!props.icon,
                "c-input--error": hasError(props.field, ctx.error),
                "c-input--suffixed": !!suffix,
              })}
              id={`input-${props.field}`}
              type="text"
              autoComplete={props.autoComplete}
              tabIndex={props.noTabFocus ? -1 : undefined}
              ref={props.inputRef}
              autoFocus={props.autoFocus}
              onFocus={onFocus}
              maxLength={props.maxLength}
              disabled={props.disabled}
              value={props.value}
              placeholder={props.placeholder}
              onChange={(e) => {
                // Clear error for this field when user interacts with it
                if (ctx.clearError && hasError(props.field, ctx.error)) {
                  ctx.clearError(props.field)
                }

                // Call the original onChange handler
                onChange(e)
              }}
            />
            {icon}
            {suffix}
          </HStack>
          <DisplayError fields={[props.field]} error={ctx.error} />
          {props.children}
        </div>
      )}
    </ValidationContext.Consumer>
  )
}
