import React from "react"
import { classSet } from "@fider/services"
import { ValidationContext } from "./Form"
import { DisplayError, hasError } from "./DisplayError"
import { IconType } from "react-icons"

interface InputProps {
  field: string
  label?: string
  className?: string
  autoComplete?: string
  autoFocus?: boolean
  noTabFocus?: boolean
  afterLabel?: JSX.Element
  icon?: IconType
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
  const onChange = (e: React.FormEvent<HTMLInputElement>) => {
    if (props.onChange) {
      props.onChange(e.currentTarget.value)
    }
  }

  const suffix = typeof props.suffix === "string" ? <span className="c-form-input-suffix">{props.suffix}</span> : props.suffix

  const icon = props.icon
    ? React.createElement(props.icon, {
        onClick: props.onIconClick,
        className: classSet({ link: !!props.onIconClick }),
      })
    : undefined

  return (
    <ValidationContext.Consumer>
      {(ctx) => (
        <div
          className={classSet({
            "c-form-field": true,
            "m-suffix": props.suffix,
            "m-error": hasError(props.field, ctx.error),
            "m-icon": !!props.icon,
            [`${props.className}`]: props.className,
          })}
        >
          {!!props.label && (
            <label htmlFor={`input-${props.field}`}>
              {props.label}
              {props.afterLabel}
            </label>
          )}
          <div className="c-form-field-wrapper">
            <input
              id={`input-${props.field}`}
              type="text"
              autoComplete={props.autoComplete}
              tabIndex={props.noTabFocus ? -1 : undefined}
              ref={props.inputRef}
              autoFocus={props.autoFocus}
              onFocus={props.onFocus}
              maxLength={props.maxLength}
              disabled={props.disabled}
              value={props.value}
              placeholder={props.placeholder}
              onChange={onChange}
            />
            {icon}
            {suffix}
          </div>
          <DisplayError fields={[props.field]} error={ctx.error} />
          {props.children}
        </div>
      )}
    </ValidationContext.Consumer>
  )
}
