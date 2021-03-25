import React from "react"

interface RadioButtonOption {
  value: string
  label: string
}

interface RadioButtonProps {
  label: string
  field: string
  defaultOption: RadioButtonOption
  options: RadioButtonOption[]
  onSelect?: (value: RadioButtonOption) => void
}

interface RadioButtonState {
  selected: RadioButtonOption
}

export class RadioButton extends React.Component<RadioButtonProps, RadioButtonState> {
  constructor(props: RadioButtonProps) {
    super(props)
    this.state = {
      selected: props.defaultOption || props.options[0],
    }
  }

  private onChange = (selected: RadioButtonOption) => {
    this.setState({ selected }, () => {
      if (this.props.onSelect) {
        this.props.onSelect(this.state.selected)
      }
    })
  }

  public render() {
    const inputs = this.props.options.map((option) => {
      return (
        <div key={option.value} className="c-form-radio-option">
          <input
            id={`visibility-${option.value}`}
            type="radio"
            name={`input-${this.props.field}`}
            checked={this.state.selected === option}
            onChange={this.onChange.bind(this, option)}
          />
          <label htmlFor={`visibility-${option.value}`}>{option.label}</label>
        </div>
      )
    })

    return (
      <div className="c-form-field">
        <label htmlFor={`input-${this.props.field}`}>{this.props.label}</label>
        {inputs}
      </div>
    )
  }
}
