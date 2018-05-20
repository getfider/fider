import * as React from "react";
import { Failure, classSet } from "@fider/services";
import { ValidationContext } from "./Form";
import { DisplayError, hasError } from "./DisplayError";

export interface SelectOption {
  value: string;
  label: string;
}

interface SelectProps {
  field: string;
  label?: string;
  maxLength?: number;
  defaultOption?: SelectOption;
  options: SelectOption[];
  onChange?: (option?: SelectOption) => void;
}

interface SelectState {
  selected?: SelectOption;
}

export class Select extends React.Component<SelectProps, SelectState> {
  constructor(props: SelectProps) {
    super(props);
    this.state = {
      selected: props.defaultOption
    };
  }

  private onChange = (e: React.FormEvent<HTMLSelectElement>) => {
    let selected: SelectOption | undefined;
    if (e.currentTarget.value) {
      const options = this.props.options.filter(o => o.value === e.currentTarget.value);
      if (options && options.length > 0) {
        selected = options[0];
      }
    }

    this.setState({ selected }, () => {
      if (this.props.onChange) {
        this.props.onChange(this.state.selected);
      }
    });
  };

  public render() {
    const options = this.props.options.map(option => {
      return (
        <option key={option.value} value={option.value}>
          {option.label}
        </option>
      );
    });

    return (
      <ValidationContext.Consumer>
        {ctx => (
          <>
            <div
              className={classSet({
                "c-form-field": true,
                "m-error": hasError(this.props.field, ctx.error)
              })}
            >
              {!!this.props.label && <label htmlFor={`input-${this.props.field}`}>{this.props.label}</label>}
              <div className="c-form-field-wrapper">
                <select
                  id={`input-${this.props.field}`}
                  defaultValue={this.props.defaultOption && this.props.defaultOption.value}
                  onChange={this.onChange}
                >
                  {options}
                </select>
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
