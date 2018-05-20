import * as React from "react";
import { Failure, classSet } from "@fider/services";
import { ValidationContext } from "./Form2";
import { DisplayError, hasError } from "./DisplayError";
import { Button } from "@fider/components";

interface InputProps {
  field: string;
  label?: string;
  autoFocus?: boolean;
  afterLabel?: JSX.Element;
  icon?: string;
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

export class Input extends React.Component<InputProps, {}> {
  constructor(props: InputProps) {
    super(props);
  }

  private onChange = (e: React.FormEvent<HTMLInputElement>) => {
    if (this.props.onChange) {
      this.props.onChange(e.currentTarget.value);
    }
  };

  private onKeyDown = (event: React.KeyboardEvent<HTMLInputElement>): void => {
    if (event.keyCode === 13 && this.props.onSubmit) {
      this.props.onSubmit();
      event.preventDefault();
    }
  };

  public render() {
    const suffix =
      typeof this.props.suffix === "string" ? (
        <span className="c-form-input-suffix">{this.props.suffix}</span>
      ) : (
        this.props.suffix
      );

    return (
      <ValidationContext.Consumer>
        {ctx => (
          <>
            <div
              className={classSet({
                "c-form-field": true,
                "m-suffix": this.props.suffix,
                "m-error": hasError(this.props.field, ctx.error),
                "m-icon": this.props.icon
              })}
            >
              {!!this.props.label && (
                <label htmlFor={`input-${this.props.field}`}>
                  {this.props.label}
                  {this.props.afterLabel}
                </label>
              )}
              <div className="c-form-field-wrapper">
                <input
                  id={`input-${this.props.field}`}
                  type="text"
                  ref={this.props.inputRef}
                  autoFocus={this.props.autoFocus}
                  onFocus={this.props.onFocus}
                  maxLength={this.props.maxLength}
                  disabled={this.props.disabled}
                  value={this.props.value}
                  placeholder={this.props.placeholder}
                  onKeyDown={this.props.onSubmit ? this.onKeyDown : undefined}
                  onChange={this.onChange}
                />
                {!!this.props.icon && (
                  <i
                    onClick={this.props.onIconClick}
                    className={classSet({
                      icon: true,
                      [this.props.icon]: true,
                      link: this.props.onIconClick
                    })}
                  />
                )}
                {suffix}
              </div>
              <DisplayError fields={[this.props.field]} error={ctx.error} />
              {this.props.children}
            </div>
          </>
        )}
      </ValidationContext.Consumer>
    );
  }
}
