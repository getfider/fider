import * as React from "react";
import { Failure, classSet } from "@fider/services";
import { ValidationContext } from "./Form2";
import { DisplayError, hasError } from "./DisplayError";

interface InputProps {
  field: string;
  label?: string;
  icon?: string;
  maxLength?: number;
  defaultValue?: string;
  disabled?: boolean;
  suffix?: string;
  placeholder?: string;
  onChange?: (value: string) => void;
}

interface InputState {
  value: string;
}

export class Input extends React.Component<InputProps, InputState> {
  constructor(props: InputProps) {
    super(props);
    this.state = {
      value: props.defaultValue || ""
    };
  }

  private onChange = (e: React.FormEvent<HTMLInputElement>) => {
    this.setState({ value: e.currentTarget.value }, () => {
      if (this.props.onChange) {
        this.props.onChange(this.state.value);
      }
    });
  };

  public render() {
    const suffix = this.props.suffix ? <span className="c-form-input-suffix">{this.props.suffix}</span> : undefined;
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
              {!!this.props.label && <label htmlFor={`input-${this.props.field}`}>{this.props.label}</label>}
              <div className="c-form-field-wrapper">
                <input
                  id={`input-${this.props.field}`}
                  type="text"
                  maxLength={this.props.maxLength}
                  disabled={this.props.disabled}
                  value={this.state.value}
                  placeholder={this.props.placeholder}
                  onChange={this.onChange}
                />
                {!!this.props.icon && <i className={`${this.props.icon} icon`} />}
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
