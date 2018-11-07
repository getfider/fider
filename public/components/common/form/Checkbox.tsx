import React from "react";
import { classSet } from "@fider/services";
import { DisplayError, ValidationContext, hasError } from "../";

interface CheckboxProps {
  field: string;
  checked?: boolean;
  onChange: (checked: boolean) => void;
}

interface CheckboxState {
  checked: boolean;
}

export class Checkbox extends React.Component<CheckboxProps, CheckboxState> {
  constructor(props: CheckboxProps) {
    super(props);
    this.state = {
      checked: props.checked || false
    };
  }

  private onChange = (e: React.FormEvent<HTMLInputElement>) => {
    this.setState({ checked: e.currentTarget.checked }, () => {
      if (this.props.onChange) {
        this.props.onChange(this.state.checked);
      }
    });
  };

  public render() {
    return (
      <ValidationContext.Consumer>
        {ctx => (
          <div
            className={classSet({
              "c-form-field m-checkbox": true,
              "m-error": hasError(this.props.field, ctx.error)
            })}
          >
            <label htmlFor={`input-${this.props.field}`}>
              <input
                id={`input-${this.props.field}`}
                type="checkbox"
                defaultChecked={this.props.checked}
                onChange={this.onChange}
              />
              {this.props.children}
            </label>
            <DisplayError fields={[this.props.field]} error={ctx.error} />
          </div>
        )}
      </ValidationContext.Consumer>
    );
  }
}
