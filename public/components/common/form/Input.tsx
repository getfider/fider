import React from "react";
import { classSet } from "@fider/services";
import { ValidationContext } from "./Form";
import { DisplayError, hasError } from "./DisplayError";
import { IconType } from "react-icons";

interface InputProps {
  field: string;
  label?: string;
  className?: string;
  autoFocus?: boolean;
  afterLabel?: JSX.Element;
  icon?: IconType;
  maxLength?: number;
  value?: string;
  disabled?: boolean;
  suffix?: string | JSX.Element;
  placeholder?: string;
  onIconClick?: () => void;
  onSubmit?: () => void;
  onFocus?: () => void;
  inputRef?: (node: HTMLInputElement) => void;
  onChange?: (value: string) => void;
}

export const Input: React.StatelessComponent<InputProps> = props => {
  const onChange = (e: React.FormEvent<HTMLInputElement>) => {
    if (props.onChange) {
      props.onChange(e.currentTarget.value);
    }
  };

  const onKeyDown = (event: React.KeyboardEvent<HTMLInputElement>): void => {
    if (event.keyCode === 13 && props.onSubmit) {
      props.onSubmit();
      event.preventDefault();
    }
  };

  const suffix =
    typeof props.suffix === "string" ? <span className="c-form-input-suffix">{props.suffix}</span> : props.suffix;

  const icon = !!props.icon
    ? React.createElement(props.icon, {
        onClick: props.onIconClick,
        className: classSet({ link: !!props.onIconClick })
      })
    : undefined;

  return (
    <ValidationContext.Consumer>
      {ctx => (
        <div
          className={classSet({
            "c-form-field": true,
            "m-suffix": props.suffix,
            "m-error": hasError(props.field, ctx.error),
            "m-icon": !!props.icon,
            [`${props.className}`]: props.className
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
              ref={props.inputRef}
              autoFocus={props.autoFocus}
              onFocus={props.onFocus}
              maxLength={props.maxLength}
              disabled={props.disabled}
              value={props.value}
              placeholder={props.placeholder}
              onKeyDown={props.onSubmit ? onKeyDown : undefined}
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
  );
};
