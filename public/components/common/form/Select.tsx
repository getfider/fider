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
  field: string
  label?: string
  maxLength?: number
  defaultValue?: string
  options: SelectOption[]
  onChange?: (option?: SelectOption) => void
}

interface SelectState {
  selected?: SelectOption
}

export class Select extends React.Component<SelectProps, SelectState> {
  constructor(props: SelectProps) {
    super(props)
    this.state = {
      selected: this.getOption(props.defaultValue),
    }
  }

  private getOption(value?: string): SelectOption | undefined {
    if (value && this.props.options) {
      const filtered = this.props.options.filter((x) => x.value === value)
      if (filtered && filtered.length > 0) {
        return filtered[0]
      }
    }
  }

  private onChange = (e: React.FormEvent<HTMLSelectElement>) => {
    let selected: SelectOption | undefined
    if (e.currentTarget.value) {
      const options = this.props.options.filter((o) => o.value === e.currentTarget.value)
      if (options && options.length > 0) {
        selected = options[0]
      }
    }

    this.setState({ selected }, () => {
      if (this.props.onChange) {
        this.props.onChange(this.state.selected)
      }
    })
  }

  public render() {
    const options = this.props.options.map((option) => {
      return (
        <option key={option.value} value={option.value}>
          {option.label}
        </option>
      )
    })

    return (
      <ValidationContext.Consumer>
        {(ctx) => (
          <>
            <div className="c-form-field">
              {!!this.props.label && <label htmlFor={`input-${this.props.field}`}>{this.props.label}</label>}
              <select
                className={classSet({
                  "c-select": true,
                  "c-select--error": hasError(this.props.field, ctx.error),
                })}
                id={`input-${this.props.field}`}
                defaultValue={this.props.defaultValue}
                onChange={this.onChange}
              >
                {options}
              </select>
              <DisplayError fields={[this.props.field]} error={ctx.error} />
              {this.props.children}
            </div>
          </>
        )}
      </ValidationContext.Consumer>
    )
  }
}
